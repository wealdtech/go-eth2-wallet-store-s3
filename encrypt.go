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

package s3

import (
	"errors"

	"github.com/wealdtech/go-ecodec"
)

// encryptIfRequired encrypts data if required.
func (s *Store) encryptIfRequired(data []byte) ([]byte, error) {
	if len(data) == 0 {
		// No data means nothing to encrypt.
		return data, nil
	}

	if len(s.passphrase) == 0 {
		// No passphrase means nothing to encrypt with.
		return data, nil
	}

	if len(data) < 16 {
		return nil, errors.New("data must be at least 16 bytes")
	}

	var err error
	if data, err = ecodec.Encrypt(data, s.passphrase); err != nil {
		return nil, err
	}

	return data, nil
}

// decryptIfRequired decrypts data if required.
func (s *Store) decryptIfRequired(data []byte) ([]byte, error) {
	if len(data) == 0 {
		// No data means nothing to decrypt.
		return data, nil
	}

	if len(s.passphrase) == 0 {
		// No passphrase means nothing to decrypt with.
		return data, nil
	}

	if len(data) < 16 {
		return nil, errors.New("data must be at least 16 bytes")
	}

	var err error
	if data, err = ecodec.Decrypt(data, s.passphrase); err != nil {
		return nil, err
	}

	return data, nil
}
