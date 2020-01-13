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

package s3

import (
	"github.com/wealdtech/go-ecodec"
)

// encryptIfRequired encrypts data if required.
func (s *Store) encryptIfRequired(data []byte) ([]byte, error) {
	var err error
	if len(s.passphrase) > 0 {
		data, err = ecodec.Encrypt(data, s.passphrase)
	}
	return data, err
}

// decryptIfRequired decrypts data if required.
func (s *Store) decryptIfRequired(data []byte) ([]byte, error) {
	var err error
	if len(s.passphrase) > 0 {
		data, err = ecodec.Decrypt(data, s.passphrase)
	}
	return data, err
}
