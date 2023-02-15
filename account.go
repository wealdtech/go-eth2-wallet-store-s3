// Copyright 2019 - 2023 Weald Technology Trading
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
	"strings"
	"sync"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

// StoreAccount stores an account.  It will fail if it cannot store the data.
// Note this will overwrite an existing account with the same ID.  It will not, however, allow multiple accounts with the same
// name to co-exist in the same wallet.
func (s *Store) StoreAccount(walletID uuid.UUID, accountID uuid.UUID, data []byte) error {
	// Ensure the wallet exists
	_, err := s.RetrieveWalletByID(walletID)
	if err != nil {
		return errors.New("unknown wallet")
	}

	// See if an account with this name already exists
	existingAccount, err := s.RetrieveAccount(walletID, accountID)
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
			return errors.New("account already exists")
		}
	}

	data, err = s.encryptIfRequired(data)
	if err != nil {
		return err
	}

	path := s.accountPath(walletID, accountID)
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
func (s *Store) RetrieveAccount(walletID uuid.UUID, accountID uuid.UUID) ([]byte, error) {
	path := s.accountPath(walletID, accountID)
	buf := aws.NewWriteAtBuffer([]byte{})
	downloader := s3manager.NewDownloader(s.session)
	if _, err := downloader.Download(buf,
		&s3.GetObjectInput{
			Bucket: aws.String(s.bucket),
			Key:    aws.String(path),
		}); err != nil {
		return nil, err
	}
	data, err := s.decryptIfRequired(buf.Bytes())
	if err != nil {
		return nil, err
	}
	return data, nil
}

// RetrieveAccounts retrieves all account-level data for a wallet.
func (s *Store) RetrieveAccounts(walletID uuid.UUID) <-chan []byte {
	path := s.walletPath(walletID)
	ch := make(chan []byte, elementCapacity)
	go func() {
		conn := s3.New(s.session)

		contents := make([]*s3.Object, 0, elementCapacity)
		var continuationToken *string
		for finished := false; !finished; {
			resp, err := conn.ListObjectsV2(&s3.ListObjectsV2Input{
				Bucket:            aws.String(s.bucket),
				Prefix:            aws.String(path + "/"),
				ContinuationToken: continuationToken,
			})
			if err != nil {
				close(ch)
				return
			}
			contents = append(contents, resp.Contents...)
			if resp.IsTruncated != nil && (*resp.IsTruncated) {
				continuationToken = resp.NextContinuationToken
			} else {
				finished = true
			}
		}

		// Download items concurrently (up to concurrency limit).
		wg := sync.WaitGroup{}
		downloader := s3manager.NewDownloader(s.session, func(d *s3manager.Downloader) {
			d.Concurrency = downloadConcurrency
		})
		for _, content := range contents {
			if strings.HasSuffix(*content.Key, "/") {
				// Directory
				continue
			}
			if strings.HasSuffix(*content.Key, walletID.String()) {
				// Wallet
				continue
			}
			wg.Add(1)
			go func(content *s3.Object) {
				defer wg.Done()
				buf := aws.NewWriteAtBuffer(make([]byte, 0, itemCapacity))
				_, err := downloader.Download(buf,
					&s3.GetObjectInput{
						Bucket: aws.String(s.bucket),
						Key:    aws.String(*content.Key),
					})
				if err != nil {
					return
				}
				data, err := s.decryptIfRequired(buf.Bytes())
				if err != nil {
					return
				}
				ch <- data
			}(content)
		}
		wg.Wait()
		close(ch)
	}()
	return ch
}
