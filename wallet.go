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

// StoreWallet stores wallet-level data.  It will fail if it cannot store the data.
// Note that this will overwrite any existing data; it is up to higher-level functions to check for the presence of a wallet with
// the wallet name and handle clashes accordingly.
func (s *Store) StoreWallet(id uuid.UUID, _ string, data []byte) error {
	path := s.walletHeaderPath(id)
	var err error
	data, err = s.encryptIfRequired(data)
	if err != nil {
		return errors.Wrap(err, "failed to encrypt wallet")
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
	for data := range s.RetrieveWallets() {
		info := &struct {
			Name string `json:"name"`
		}{}
		err := json.Unmarshal(data, info)
		if err == nil && info.Name == walletName {
			return data, nil
		}
	}

	return nil, errors.New("wallet not found")
}

// RetrieveWalletByID retrieves wallet-level data.  It will fail if it cannot retrieve the data.
func (s *Store) RetrieveWalletByID(walletID uuid.UUID) ([]byte, error) {
	for data := range s.RetrieveWallets() {
		info := &struct {
			ID uuid.UUID `json:"uuid"`
		}{}
		err := json.Unmarshal(data, info)
		if err == nil && info.ID == walletID {
			return data, nil
		}
	}

	return nil, errors.New("wallet not found")
}

// RetrieveWallets retrieves wallet-level data for all wallets.
func (s *Store) RetrieveWallets() <-chan []byte {
	ch := make(chan []byte, elementCapacity)
	go func() {
		conn := s3.New(s.session)

		contents := make([]*s3.Object, 0, elementCapacity)
		var continuationToken *string
		for finished := false; !finished; {
			resp, err := conn.ListObjectsV2(&s3.ListObjectsV2Input{
				Bucket:            aws.String(s.bucket),
				Prefix:            aws.String(s.path),
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
				// Directory.
				continue
			}
			// This is only a wallet if the last two components of the path are the same.
			components := strings.Split(*content.Key, "/")
			if components[len(components)-1] != components[len(components)-2] {
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
