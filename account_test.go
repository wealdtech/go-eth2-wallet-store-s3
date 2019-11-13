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

package s3_test

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	s3 "github.com/wealdtech/go-eth2-wallet-store-s3"
)

func TestStoreRetrieveAccount(t *testing.T) {
	rand.Seed(time.Now().Unix())
	id := fmt.Sprintf("%s-%d", t.Name(), rand.Int31())
	store, err := s3.New(s3.WithID([]byte(id)))
	if err != nil {
		t.Skip("unable to access S3; skipping test")
	}

	walletID := uuid.New()
	walletName := "test wallet"
	accountID := uuid.New()
	accountName := "test account"
	data := []byte(fmt.Sprintf(`{"name":%q,"uuid":%q}`, accountName, accountID.String()))

	err = store.StoreWallet(walletID, walletName, data)
	require.Nil(t, err)
	err = store.StoreAccount(walletID, walletName, accountID, accountName, data)
	require.Nil(t, err)
	retData, err := store.RetrieveAccount(walletID, walletName, accountName)
	require.Nil(t, err)
	assert.Equal(t, data, retData)

	store.RetrieveWallets()

	_, err = store.RetrieveAccount(walletID, walletName, "not present")
	assert.NotNil(t, err)
}

func TestDuplicateAccounts(t *testing.T) {
	rand.Seed(time.Now().Unix())
	id := fmt.Sprintf("%s-%d", t.Name(), rand.Int31())
	store, err := s3.New(s3.WithID([]byte(id)))
	if err != nil {
		t.Skip("unable to access S3; skipping test")
	}

	walletID := uuid.New()
	walletName := "test wallet"
	accountID := uuid.New()
	accountName := "test account"
	data := []byte(fmt.Sprintf(`{"name":%q,"uuid":%q}`, accountName, accountID.String()))

	err = store.StoreWallet(walletID, walletName, data)
	require.Nil(t, err)
	err = store.StoreAccount(walletID, walletName, accountID, accountName, data)
	require.Nil(t, err)

	// Try to store account with the same name but different ID; should fail
	err = store.StoreAccount(walletID, walletName, uuid.New(), accountName, data)
	assert.NotNil(t, err)

	// Try to store account with the same name and same ID; should succeed
	err = store.StoreAccount(walletID, walletName, accountID, accountName, data)
	assert.Nil(t, err)
}

func TestRetrieveNonExistentAccount(t *testing.T) {
	rand.Seed(time.Now().Unix())
	id := fmt.Sprintf("%s-%d", t.Name(), rand.Int31())
	store, err := s3.New(s3.WithID([]byte(id)))
	if err != nil {
		t.Skip("unable to access S3; skipping test")
	}

	walletID := uuid.New()
	walletName := "test wallet"
	accountName := "test account"

	_, err = store.RetrieveAccount(walletID, walletName, accountName)
	assert.NotNil(t, err)
}

func TestStoreNonExistentAccount(t *testing.T) {
	rand.Seed(time.Now().Unix())
	id := fmt.Sprintf("%s-%d", t.Name(), rand.Int31())
	store, err := s3.New(s3.WithID([]byte(id)))
	if err != nil {
		t.Skip("unable to access S3; skipping test")
	}

	walletID := uuid.New()
	walletName := "test wallet"
	accountID := uuid.New()
	accountName := "test account"
	data := []byte(fmt.Sprintf(`"uuid":%q,"name":%q}`, accountID, accountName))

	err = store.StoreAccount(walletID, walletName, accountID, accountName, data)
	assert.NotNil(t, err)
}
