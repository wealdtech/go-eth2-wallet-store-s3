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

func TestStoreRetrieveAccount(t *testing.T) {
	//nolint:gosec
	rand := rand.New(rand.NewSource(time.Now().Unix()))
	id := fmt.Sprintf("%s-%d", t.Name(), rand.Int31())
	store, err := s3.New(
		s3.WithID([]byte(id)),
		s3.WithCredentialsID(os.Getenv("S3_CREDENTIALS_ID")),
		s3.WithCredentialsSecret(os.Getenv("S3_CREDENTIALS_SECRET")),
		s3.WithBucket(os.Getenv("S3_BUCKET")),
	)
	if err != nil {
		t.Skip("unable to access S3; skipping test")
	}

	walletID := uuid.New()
	walletName := fmt.Sprintf("test wallet for TestStoreRetrieveAccount %d", time.Now().UnixNano())
	walletData := []byte(fmt.Sprintf(`{"name":%q,"uuid":%q}`, walletName, walletID.String()))
	accountID := uuid.New()
	accountName := "test account"
	accountData := []byte(fmt.Sprintf(`{"name":%q,"uuid":%q}`, accountName, accountID.String()))

	err = store.StoreWallet(walletID, walletName, walletData)
	require.Nil(t, err)
	err = store.StoreAccount(walletID, accountID, accountData)
	require.Nil(t, err)
	retData, err := store.RetrieveAccount(walletID, accountID)
	require.Nil(t, err)
	assert.Equal(t, accountData, retData)

	store.RetrieveWallets()

	_, err = store.RetrieveAccount(walletID, uuid.New())
	assert.NotNil(t, err)
}

func TestDuplicateAccounts(t *testing.T) {
	//nolint:gosec
	rand := rand.New(rand.NewSource(time.Now().Unix()))
	id := fmt.Sprintf("%s-%d", t.Name(), rand.Int31())
	store, err := s3.New(s3.WithID([]byte(id)),
		s3.WithCredentialsID(os.Getenv("S3_CREDENTIALS_ID")),
		s3.WithCredentialsSecret(os.Getenv("S3_CREDENTIALS_SECRET")),
		s3.WithBucket(os.Getenv("S3_BUCKET")),
	)
	if err != nil {
		t.Skip("unable to access S3; skipping test")
	}

	walletID := uuid.New()
	walletName := fmt.Sprintf("test wallet for TestDuplicateAccounts %d", time.Now().UnixNano())
	walletData := []byte(fmt.Sprintf(`{"name":%q,"uuid":%q}`, walletName, walletID.String()))
	accountID := uuid.New()
	accountName := "test account"
	accountData := []byte(fmt.Sprintf(`{"name":%q,"uuid":%q}`, accountName, accountID.String()))

	err = store.StoreWallet(walletID, walletName, walletData)
	require.Nil(t, err)
	err = store.StoreAccount(walletID, accountID, accountData)
	require.Nil(t, err)

	// Try to store account with the same name and same ID; should succeed
	err = store.StoreAccount(walletID, accountID, accountData)
	assert.Nil(t, err)
}

func TestRetrieveNonExistentAccount(t *testing.T) {
	//nolint:gosec
	rand := rand.New(rand.NewSource(time.Now().Unix()))
	id := fmt.Sprintf("%s-%d", t.Name(), rand.Int31())
	store, err := s3.New(s3.WithID([]byte(id)),
		s3.WithCredentialsID(os.Getenv("S3_CREDENTIALS_ID")),
		s3.WithCredentialsSecret(os.Getenv("S3_CREDENTIALS_SECRET")),
		s3.WithBucket(os.Getenv("S3_BUCKET")),
	)
	if err != nil {
		t.Skip("unable to access S3; skipping test")
	}

	walletID := uuid.New()

	_, err = store.RetrieveAccount(walletID, uuid.New())
	assert.NotNil(t, err)
}

func TestStoreNonExistentAccount(t *testing.T) {
	//nolint:gosec
	rand := rand.New(rand.NewSource(time.Now().Unix()))
	id := fmt.Sprintf("%s-%d", t.Name(), rand.Int31())
	store, err := s3.New(s3.WithID([]byte(id)),
		s3.WithCredentialsID(os.Getenv("S3_CREDENTIALS_ID")),
		s3.WithCredentialsSecret(os.Getenv("S3_CREDENTIALS_SECRET")),
		s3.WithBucket(os.Getenv("S3_BUCKET")),
	)
	if err != nil {
		t.Skip("unable to access S3; skipping test")
	}

	walletID := uuid.New()
	accountID := uuid.New()
	accountName := "test account"
	data := []byte(fmt.Sprintf(`"uuid":%q,"name":%q}`, accountID, accountName))

	err = store.StoreAccount(walletID, accountID, data)
	assert.NotNil(t, err)
}

func TestStoreAccount(t *testing.T) {
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
				s3.WithBucket(fmt.Sprintf("teststoreaccount-specificbucket-%d", time.Now().UnixNano())),
			},
		},
		{
			name: "SpecificPath",
			opts: []s3.Option{
				s3.WithID([]byte(fmt.Sprintf("%d", rand.Int31()))),
				s3.WithCredentialsID(os.Getenv("S3_CREDENTIALS_ID")),
				s3.WithCredentialsSecret(os.Getenv("S3_CREDENTIALS_SECRET")),
				s3.WithBucket(fmt.Sprintf("teststoreaccount-specificpath-%d", time.Now().UnixNano())),
				s3.WithPath("a/b/c"),
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			store, err := s3.New(test.opts...)
			require.NoError(t, err)

			walletID := uuid.New()
			walletName := fmt.Sprintf("test wallet for TestStoreAccount/%s %d", test.name, time.Now().UnixNano())
			walletData := []byte(fmt.Sprintf(`{"name":%q,"uuid":%q}`, walletName, walletID))
			require.NoError(t, store.StoreWallet(walletID, walletName, walletData))
			retData, err := store.RetrieveWallet(walletName)
			require.NoError(t, err)
			assert.Equal(t, walletData, retData)

			accountID := uuid.New()
			accountName := "test account"
			accountData := []byte(fmt.Sprintf(`{"name":%q,"uuid":%q}`, accountName, accountID.String()))
			require.NoError(t, store.StoreAccount(walletID, accountID, accountData))
			retData, err = store.RetrieveAccount(walletID, accountID)
			require.NoError(t, err)
			assert.Equal(t, accountData, retData)
		})
	}
}
