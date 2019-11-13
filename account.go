// Copyright © 2019 Weald Technology Trading
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package s3

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	ecodec "github.com/wealdtech/go-ecodec"
)

// StoreAccount stores an account.  It will fail if it cannot store the data.
// Note this will overwrite an existing account with the same ID.  It will not, however, allow multiple accounts with the same
// name to co-exist in the same wallet.
func (s *Store) StoreAccount(walletID uuid.UUID, walletName string, accountID uuid.UUID, accountName string, data []byte) error {
	// Ensure the wallet exists
	if _, err := s.RetrieveWallet(walletName); err != nil {
		return errors.Wrapf(err, "no wallet %q", walletName)
	}

	// See if an account with this name already exists
	existingAccount, err := s.RetrieveAccount(walletID, walletName, accountName)
	if err == nil {
		// It does; they need to have the same ID for us to overwrite it
		info := &struct {
			ID string `json:"uuid"`
		}{}
		err := json.Unmarshal(existingAccount, info)
		if err != nil {
			return err
		}
		if info.ID != accountID.String() {
			return fmt.Errorf("account %q already exists", accountName)
		}
	}

	if len(s.passphrase) > 0 {
		data, err = ecodec.Encrypt(data, s.passphrase)
		if err != nil {
			return err
		}
	}

	path := s.accountPath(walletName, accountName)
	uploader := s3manager.NewUploader(s.session)
	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(path),
		Body:   bytes.NewReader(data),
	})
	if err != nil {
		return errors.Wrap(err, "failed to store key")
	}
	return nil
}

// RetrieveAccount retrieves account-level data.  It will fail if it cannot retrieve the data.
func (s *Store) RetrieveAccount(walletID uuid.UUID, walletName string, accountName string) ([]byte, error) {
	path := s.accountPath(walletName, accountName)
	downloader := s3manager.NewDownloader(s.session)
	buf := aws.NewWriteAtBuffer([]byte{})
	_, err := downloader.Download(buf,
		&s3.GetObjectInput{
			Bucket: aws.String(s.bucket),
			Key:    aws.String(path),
		})
	data := buf.Bytes()
	if len(s.passphrase) > 0 {
		data, err = ecodec.Decrypt(data, s.passphrase)
		if err != nil {
			return nil, fmt.Errorf("unable to decrypt account %q", accountName)
		}
	}
	return data, err
}

// RetrieveAccounts retrieves all account-level data for a wallet.
func (s *Store) RetrieveAccounts(id uuid.UUID, name string) <-chan []byte {
	path := s.walletPath(name)
	ch := make(chan []byte, 1024)
	go func() {
		conn := s3.New(s.session)
		resp, err := conn.ListObjectsV2(&s3.ListObjectsV2Input{
			Bucket: aws.String(s.bucket),
			Prefix: aws.String(path + "/"),
		})
		if err == nil {
			headerPath := s.walletHeaderPath(name)
			for _, item := range resp.Contents {
				if strings.HasSuffix(*item.Key, "/") {
					// Directory
					continue
				}
				if *item.Key == headerPath {
					// Header
					continue
				}
				buf := aws.NewWriteAtBuffer([]byte{})
				downloader := s3manager.NewDownloader(s.session)
				_, err := downloader.Download(buf,
					&s3.GetObjectInput{
						Bucket: aws.String(s.bucket),
						Key:    aws.String(*item.Key),
					})
				if err == nil {
					data := buf.Bytes()
					if len(s.passphrase) > 0 {
						data, err = ecodec.Decrypt(data, s.passphrase)
						if err != nil {
							continue
						}
					}
					ch <- data
				}
			}
		}
		close(ch)
	}()
	return ch
}
