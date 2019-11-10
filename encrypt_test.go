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
	"encoding/json"
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	keystorev4 "github.com/wealdtech/go-eth2-wallet-encryptor-keystorev4"
	nd "github.com/wealdtech/go-eth2-wallet-nd"
	s3 "github.com/wealdtech/go-eth2-wallet-store-s3"
)

func TestStoreRetrieveEncryptedWallet(t *testing.T) {
	tests := []struct {
		name string
		data []byte
		err  error
	}{
		{
			name: "WithData",
			data: []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
		},
	}

	rand.Seed(time.Now().Unix())
	id := fmt.Sprintf("%s-%d", t.Name(), rand.Int31())
	store, err := s3.New(s3.WithID([]byte(id)), s3.WithPassphrase([]byte("test")))
	require.Nil(t, err)
	encryptor := keystorev4.New()
	wallet, err := nd.CreateWallet("test", store, encryptor)
	require.Nil(t, err)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := store.StoreWallet(wallet, test.data)
			if test.err != nil {
				require.NotNil(t, err)
				assert.Equal(t, test.err.Error(), err.Error())
			} else {
				require.Nil(t, err)
				data, err := store.RetrieveWallet("test")
				require.Nil(t, err)
				assert.Equal(t, test.data, data)

				wallets := false
				for range store.RetrieveWallets() {
					wallets = true
				}
				assert.True(t, wallets)

			}
		})
	}

	store.RetrieveWallets()
}

func TestStoreRetrieveEncryptedAccount(t *testing.T) {
	tests := []struct {
		name string
		data []byte
		err  error
	}{
		{
			name: "WithData",
			data: []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
		},
	}

	rand.Seed(time.Now().Unix())
	id := fmt.Sprintf("%s-%d", t.Name(), rand.Int31())
	store, err := s3.New(s3.WithID([]byte(id)), s3.WithPassphrase([]byte("test")))
	require.Nil(t, err)
	encryptor := keystorev4.New()
	wallet, err := nd.CreateWallet("test", store, encryptor)
	require.Nil(t, err)
	wallet.Unlock(nil)
	require.Nil(t, err)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			account, err := wallet.CreateAccount(test.name, nil)
			require.Nil(t, err)
			data, err := store.RetrieveAccount(wallet, test.name)
			require.Nil(t, err)
			accData, err := json.Marshal(account)
			require.Nil(t, err)
			assert.Equal(t, accData, data)

			accounts := false
			for range store.RetrieveAccounts(wallet) {
				accounts = true
			}
			assert.True(t, accounts)
		})
	}

	store.RetrieveWallets()
}

func TestBadWalletKey(t *testing.T) {
	rand.Seed(time.Now().Unix())
	id := fmt.Sprintf("%s-%d", t.Name(), rand.Int31())
	store, err := s3.New(s3.WithID([]byte(id)), s3.WithPassphrase([]byte("test")))
	require.Nil(t, err)
	encryptor := keystorev4.New()
	wallet, err := nd.CreateWallet("test", store, encryptor)
	require.Nil(t, err)

	data, err := json.Marshal(wallet)
	require.Nil(t, err)

	err = store.StoreWallet(wallet, data)
	require.Nil(t, err)

	// Open wallet with store with different key; should fail
	store, err = s3.New(s3.WithID([]byte(id)), s3.WithPassphrase([]byte("badkey")))
	require.Nil(t, err)
	_, err = nd.OpenWallet("test", store, encryptor)
	require.NotNil(t, err)
}
