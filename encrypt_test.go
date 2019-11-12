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

func TestStoreRetrieveEncryptedWallet(t *testing.T) {
	rand.Seed(time.Now().Unix())
	id := fmt.Sprintf("%s-%d", t.Name(), rand.Int31())
	store, err := s3.New(s3.WithID([]byte(id)), s3.WithPassphrase([]byte("test")))
	if err != nil {
		t.Skip("unable to access S3; skipping test")
	}

	walletID := uuid.New()
	walletName := "test"
	data := []byte(fmt.Sprintf(`{"id":%q,"name":%q}`, walletID, walletName))

	err = store.StoreWallet(walletID, walletName, data)
	require.Nil(t, err)
	retData, err := store.RetrieveWallet(walletName)
	require.Nil(t, err)
	assert.Equal(t, data, retData)

	wallets := false
	for range store.RetrieveWallets() {
		wallets = true
	}
	assert.True(t, wallets)

	store.RetrieveWallets()
}

func TestStoreRetrieveEncryptedAccount(t *testing.T) {
	rand.Seed(time.Now().Unix())
	id := fmt.Sprintf("%s-%d", t.Name(), rand.Int31())
	store, err := s3.New(s3.WithID([]byte(id)), s3.WithPassphrase([]byte("test")))
	if err != nil {
		t.Skip("unable to access S3; skipping test")
	}

	walletID := uuid.New()
	walletName := "test wallet"
	accountID := uuid.New()
	accountName := "test account"
	data := []byte(fmt.Sprintf(`{"name":"%s","id":"%s"}`, accountName, accountID.String()))

	err = store.StoreWallet(walletID, walletName, data)
	require.Nil(t, err)

	err = store.StoreAccount(walletID, walletName, accountID, accountName, data)
	require.Nil(t, err)
	retData, err := store.RetrieveAccount(walletID, walletName, accountName)
	require.Nil(t, err)
	require.Equal(t, data, retData)

	accounts := false
	for range store.RetrieveAccounts(walletID, walletName) {
		accounts = true
	}
	assert.True(t, accounts)
}

func TestBadWalletKey(t *testing.T) {
	rand.Seed(time.Now().Unix())
	id := fmt.Sprintf("%s-%d", t.Name(), rand.Int31())
	store, err := s3.New(s3.WithID([]byte(id)), s3.WithPassphrase([]byte("test")))
	if err != nil {
		t.Skip("unable to access S3; skipping test")
	}

	walletID := uuid.New()
	walletName := "test wallet"
	data := []byte(fmt.Sprintf(`{"id":%q,"name":%q}`, walletID, walletName))

	err = store.StoreWallet(walletID, walletName, data)
	require.Nil(t, err)

	// Open wallet with store with different key; should fail
	store, err = s3.New(s3.WithID([]byte(id)), s3.WithPassphrase([]byte("badkey")))
	require.Nil(t, err)
	_, err = store.RetrieveWallet(walletName)
	require.NotNil(t, err)
}
