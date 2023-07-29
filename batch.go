// Copyright 2023 Weald Technology Trading.
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
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

// StoreBatch stores wallet batch data.  It will fail if it cannot store the data.
func (s *Store) StoreBatch(_ context.Context, walletID uuid.UUID, _ string, data []byte) error {
	// Ensure wallet exists.
	_, err := s.RetrieveWalletByID(walletID)
	if err != nil {
		return err
	}

	path := s.walletBatchPath(walletID)
	data, err = s.encryptIfRequired(data)
	if err != nil {
		return errors.Wrap(err, "failed to encrypt batch")
	}
	uploader := s3manager.NewUploader(s.session)
	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(path),
		Body:   bytes.NewReader(data),
	})
	if err != nil {
		return errors.Wrap(err, "failed to store batch")
	}

	return nil
}

// RetrieveBatch retrieves the batch of accounts for a given wallet.
func (s *Store) RetrieveBatch(_ context.Context, walletID uuid.UUID) ([]byte, error) {
	// Ensure wallet exists.
	_, err := s.RetrieveWalletByID(walletID)
	if err != nil {
		return nil, err
	}

	path := s.walletBatchPath(walletID)

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
