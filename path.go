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

package s3

import (
	"encoding/hex"
	"fmt"

	util "github.com/wealdtech/go-eth2-util"
)

func (s *Store) walletPath(walletName string) string {
	// We hash the wallet name for two reasons.  Firstly, to avoid illegal characters in S3 file names.  Secondly, to make it
	// slightly more difficult for wallets with incorrect permissions to be easily searchable
	return hex.EncodeToString(util.SHA256([]byte(walletName)))[:63]
}

func (s *Store) walletHeaderPath(walletName string) string {
	encodedHeaderName := hex.EncodeToString(util.SHA256([]byte("_header.json")))[:63]
	return fmt.Sprintf("%s/%s", s.walletPath(walletName), encodedHeaderName)
}

func (s *Store) accountPath(walletName string, accountName string) string {
	// We hash the key name to avoid illegal characters in S3 file names
	return s.walletPath(walletName) + "/" + hex.EncodeToString(util.SHA256([]byte(accountName)))[:63]
}
