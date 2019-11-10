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
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	keystorev4 "github.com/wealdtech/go-eth2-wallet-encryptor-keystorev4"
	nd "github.com/wealdtech/go-eth2-wallet-nd"
	s3 "github.com/wealdtech/go-eth2-wallet-store-s3"
)

func TestStoreRetrieveWallet(t *testing.T) {
	tests := []struct {
		name string
		data []byte
		err  error
	}{
		{
			name: "Empty",
			data: []byte{},
		},
	}

	rand.Seed(time.Now().Unix())
	path := filepath.Join(os.TempDir(), fmt.Sprintf("TestStoreRetrieveWallet-%d", rand.Int31()))
	defer os.RemoveAll(path)
	store, err := s3.New()
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
			}
		})
	}

	for range store.RetrieveWallets() {
	}
}
