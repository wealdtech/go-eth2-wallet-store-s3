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
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	s3 "github.com/wealdtech/go-eth2-wallet-store-s3"
)

func TestNew(t *testing.T) {
	store, err := s3.New()
	require.Nil(t, err)
	assert.Equal(t, "s3", store.Name())
	store, err = s3.New(s3.WithRegion("us-west-1"), s3.WithID([]byte("west")))
	require.Nil(t, err)
	assert.Equal(t, "s3", store.Name())
	store, err = s3.New(s3.WithRegion("us-west-1"), s3.WithID([]byte("west")), s3.WithPassphrase([]byte("secret")))
	require.Nil(t, err)
	assert.Equal(t, "s3", store.Name())
}
