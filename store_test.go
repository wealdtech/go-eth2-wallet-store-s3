// Copyright 2019, 2020 Weald Technology Trading
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
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	s3 "github.com/wealdtech/go-eth2-wallet-store-s3"
	wtypes "github.com/wealdtech/go-eth2-wallet-types/v2"
)

func TestNew(t *testing.T) {
	store, err := s3.New(
		s3.WithCredentialsID(os.Getenv("S3_CREDENTIALS_ID")),
		s3.WithCredentialsSecret(os.Getenv("S3_CREDENTIALS_SECRET")),
	)
	if err != nil {
		t.Skip("unable to access S3; skipping test")
	}
	assert.Equal(t, "s3", store.Name())
	store, err = s3.New(
		s3.WithCredentialsID(os.Getenv("S3_CREDENTIALS_ID")),
		s3.WithCredentialsSecret(os.Getenv("S3_CREDENTIALS_SECRET")),
		s3.WithRegion("us-west-1"),
		s3.WithID([]byte("west")),
	)
	require.Nil(t, err)
	assert.Equal(t, "s3", store.Name())
	store, err = s3.New(
		s3.WithCredentialsID(os.Getenv("S3_CREDENTIALS_ID")),
		s3.WithCredentialsSecret(os.Getenv("S3_CREDENTIALS_SECRET")),
		s3.WithRegion("us-west-1"),
		s3.WithID([]byte("west")),
		s3.WithPassphrase([]byte("secret")),
	)
	require.Nil(t, err)
	assert.Equal(t, "s3", store.Name())

	storeLocationProvider, ok := store.(wtypes.StoreLocationProvider)
	assert.True(t, ok)
	assert.Equal(t, "67038ae26ce874153859c347eebba98cebd31639b3a959c42d6b47e0452b185", storeLocationProvider.Location())
}

func TestNewOptions(t *testing.T) {
	ts := time.Now().UnixNano()
	tests := []struct {
		name     string
		opts     []s3.Option
		err      string
		location string
	}{
		{
			name: "Defaults",
			opts: []s3.Option{
				s3.WithCredentialsID(os.Getenv("S3_CREDENTIALS_ID")),
				s3.WithCredentialsSecret(os.Getenv("S3_CREDENTIALS_SECRET")),
			},
			location: "f264781ef3f6e9b903723a3f6909c167d1d4667403e56a900d7ba02ddb97f44",
		},
		{
			name: "SpecificBucket",
			opts: []s3.Option{
				s3.WithCredentialsID(os.Getenv("S3_CREDENTIALS_ID")),
				s3.WithCredentialsSecret(os.Getenv("S3_CREDENTIALS_SECRET")),
				s3.WithBucket(fmt.Sprintf("testnewoptions-specificbucket-%d", ts)),
			},
			location: fmt.Sprintf("testnewoptions-specificbucket-%d", ts),
		},
		{
			name: "SpecificPath",
			opts: []s3.Option{
				s3.WithCredentialsID(os.Getenv("S3_CREDENTIALS_ID")),
				s3.WithCredentialsSecret(os.Getenv("S3_CREDENTIALS_SECRET")),
				s3.WithBucket(fmt.Sprintf("testnewoptions-specificpath-%d", ts)),
				s3.WithPath("a/b/c"),
			},
			location: fmt.Sprintf("testnewoptions-specificpath-%d/a/b/c", ts),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			store, err := s3.New(test.opts...)
			if test.err != "" {
				require.EqualError(t, err, test.err)
			} else {
				require.NoError(t, err)
				require.Equal(t, "s3", store.Name())
				require.Equal(t, test.location, store.(wtypes.StoreLocationProvider).Location())
			}
		})
	}
}
