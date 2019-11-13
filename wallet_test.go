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

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	s3 "github.com/wealdtech/go-eth2-wallet-store-s3"
)

func TestStoreRetrieveWallet(t *testing.T) {
	rand.Seed(time.Now().Unix())
	path := filepath.Join(os.TempDir(), fmt.Sprintf("TestStoreRetrieveWallet-%d", rand.Int31()))
	defer os.RemoveAll(path)
	store, err := s3.New()
	if err != nil {
		t.Skip("unable to access S3; skipping test")
	}

	walletID := uuid.New()
	walletName := "test wallet"
	data := []byte(fmt.Sprintf(`{"uuid":%q,"name":%q}`, walletID, walletName))

	err = store.StoreWallet(walletID, walletName, data)
	require.Nil(t, err)
	retData, err := store.RetrieveWallet(walletName)
	require.Nil(t, err)
	assert.Equal(t, data, retData)

	for range store.RetrieveWallets() {
	}
}
