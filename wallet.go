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
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/pkg/errors"
	ecodec "github.com/wealdtech/go-ecodec"
	types "github.com/wealdtech/go-eth2-wallet-types"
)

// StoreWallet stores wallet-level data.  It will fail if it cannot store the data.
// Note that this will overwrite any existing data; it is up to higher-level functions to check for the presence of a wallet with
// the wallet name and handle clashes accordingly.
func (s *Store) StoreWallet(wallet types.Wallet, data []byte) error {
	path := s.walletHeaderPath(wallet.Name())
	var err error
	if len(s.passphrase) > 0 {
		data, err = ecodec.Encrypt(data, s.passphrase)
		if err != nil {
			return errors.Wrap(err, "failed to encrypt wallet")
		}
	}
	uploader := s3manager.NewUploader(s.session)
	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(path),
		Body:   bytes.NewReader(data),
	})
	if err != nil {
		return errors.Wrap(err, "failed to store wallet")
	}
	return nil
}

// RetrieveWallet retrieves wallet-level data.  It will fail if it cannot retrieve the data.
func (s *Store) RetrieveWallet(walletName string) ([]byte, error) {
	path := s.walletHeaderPath(walletName)
	downloader := s3manager.NewDownloader(s.session)
	buf := aws.NewWriteAtBuffer([]byte{})
	_, err := downloader.Download(buf,
		&s3.GetObjectInput{
			Bucket: aws.String(s.bucket),
			Key:    aws.String(path),
		})
	if err != nil {
		return nil, err
	}
	data := buf.Bytes()
	if len(s.passphrase) > 0 {
		data, err = ecodec.Decrypt(data, s.passphrase)
		if err != nil {
			return nil, errors.Wrap(err, "failed to decrypt wallet")
		}
	}
	return data, nil
}

// RetrieveWallets retrieves wallet-level data for all wallets.
func (s *Store) RetrieveWallets() <-chan []byte {
	ch := make(chan []byte, 1024)
	go func() {
		// We don't know the wallet name but need the last component of the encoded path for the header
		walletHeaderKey := strings.Split(s.walletHeaderPath(""), "/")[1]

		conn := s3.New(s.session)
		resp, err := conn.ListObjectsV2(&s3.ListObjectsV2Input{Bucket: aws.String(s.bucket)})
		if err == nil {
			for _, item := range resp.Contents {
				if !strings.HasSuffix(*item.Key, walletHeaderKey) {
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
