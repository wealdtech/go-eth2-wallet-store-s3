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

package s3_test

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	s3 "github.com/wealdtech/go-eth2-wallet-store-s3"
	e2wtypes "github.com/wealdtech/go-eth2-wallet-types/v2"
)

func TestStoreRetrieveBatch(t *testing.T) {
	ctx := context.Background()

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
	walletName := fmt.Sprintf("test wallet for TestStoreRetrieveBatch %d", time.Now().UnixNano())
	data := []byte(fmt.Sprintf(`{"name":%q,"uuid":%q}`, walletName, walletID))
	require.Nil(t, store.StoreWallet(walletID, walletName, data))

	batchData := []byte(`{"test":true}`)
	require.NoError(t, store.(e2wtypes.BatchStorer).StoreBatch(ctx, walletID, walletName, batchData))

	retrievedBatchData, err := store.(e2wtypes.BatchRetriever).RetrieveBatch(ctx, walletID)
	require.NoError(t, err)
	require.Equal(t, batchData, retrievedBatchData)
}

func TestStoreBatchNonExistentWallet(t *testing.T) {
	ctx := context.Background()

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
	walletName := fmt.Sprintf("test wallet for TestStoreBatchNonExistentWallet %d", time.Now().UnixNano())

	batchData := []byte(`{"test":true}`)
	require.ErrorContains(t, store.(e2wtypes.BatchStorer).StoreBatch(ctx, walletID, walletName, batchData), "wallet not found")
}

func TestRetrieveBatchNonExistentWallet(t *testing.T) {
	ctx := context.Background()

	//nolint:gosec
	rand := rand.New(rand.NewSource(time.Now().Unix()))
	id := fmt.Sprintf("%s-%d", t.Name(), rand.Int31())
	store, err := s3.New(
		s3.WithID([]byte(id)),
		s3.WithCredentialsID(os.Getenv("S3_CREDENTIALS_ID")),
		s3.WithCredentialsSecret(os.Getenv("S3_CREDENTIALS_SECRET")),
	)
	if err != nil {
		t.Skip("unable to access S3; skipping test")
	}

	walletID := uuid.New()

	_, err = store.(e2wtypes.BatchRetriever).RetrieveBatch(ctx, walletID)
	require.ErrorContains(t, err, "wallet not found")
}

func TestRetrieveNonExistentBatch(t *testing.T) {
	ctx := context.Background()

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
	walletName := fmt.Sprintf("test wallet for TestRetrieveNonExistentBatch %d", time.Now().UnixNano())
	data := []byte(fmt.Sprintf(`{"name":%q,"uuid":%q}`, walletName, walletID))
	require.Nil(t, store.StoreWallet(walletID, walletName, data))

	_, err = store.(e2wtypes.BatchRetriever).RetrieveBatch(ctx, walletID)
	require.ErrorContains(t, err, "The specified key does not exist")
}
