// Copyright 2019 - 2023 Weald Technology Trading.
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

package s3_test

import (
	"fmt"
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	s3 "github.com/wealdtech/go-eth2-wallet-store-s3"
)

func TestStoreWallet(t *testing.T) {
	if os.Getenv("S3_CREDENTIALS_ID") == "" ||
		os.Getenv("S3_CREDENTIALS_SECRET") == "" {
		t.Skip("unable to access S3; skipping test")
	}

	tests := []struct {
		name string
		opts []s3.Option
		err  string
	}{
		{
			name: "Defaults",
			opts: []s3.Option{
				s3.WithID([]byte(fmt.Sprintf("%d", rand.Int31()))),
				s3.WithCredentialsID(os.Getenv("S3_CREDENTIALS_ID")),
				s3.WithCredentialsSecret(os.Getenv("S3_CREDENTIALS_SECRET")),
			},
		},
		{
			name: "SpecificBucket",
			opts: []s3.Option{
				s3.WithID([]byte(fmt.Sprintf("%d", rand.Int31()))),
				s3.WithCredentialsID(os.Getenv("S3_CREDENTIALS_ID")),
				s3.WithCredentialsSecret(os.Getenv("S3_CREDENTIALS_SECRET")),
				s3.WithBucket(fmt.Sprintf("teststorewallet-specificbucket-%d", time.Now().UnixNano())),
			},
		},
		{
			name: "SpecificPath",
			opts: []s3.Option{
				s3.WithID([]byte(fmt.Sprintf("%d", rand.Int31()))),
				s3.WithCredentialsID(os.Getenv("S3_CREDENTIALS_ID")),
				s3.WithCredentialsSecret(os.Getenv("S3_CREDENTIALS_SECRET")),
				s3.WithBucket(os.Getenv("S3_BUCKET")),
				s3.WithPath("a/b/c"),
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			store, err := s3.New(test.opts...)
			require.NoError(t, err)

			walletID := uuid.New()
			walletName := fmt.Sprintf("test wallet for TestStoreWallet/%s %d", test.name, time.Now().UnixNano())
			data := []byte(fmt.Sprintf(`{"name":%q,"uuid":%q}`, walletName, walletID))
			err = store.StoreWallet(walletID, walletName, data)
			require.Nil(t, err)
			retData, err := store.RetrieveWallet(walletName)
			require.Nil(t, err)
			assert.Equal(t, data, retData)

			for wallet := range store.RetrieveWallets() {
				require.Equal(t, data, wallet)
			}
		})
	}
}
