// Copyright 2019 - 2022 Weald Technology Trading
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
	"github.com/google/uuid"
)

func (s *Store) walletPath(walletID uuid.UUID) string {
	return join(s.path, walletID.String())
}

func (s *Store) walletHeaderPath(walletID uuid.UUID) string {
	return join(s.walletPath(walletID), walletID.String())
}

func (s *Store) accountPath(walletID uuid.UUID, accountID uuid.UUID) string {
	return join(s.walletPath(walletID), accountID.String())
}

func (s *Store) walletIndexPath(walletID uuid.UUID) string {
	return join(s.walletPath(walletID), "index")
}

// join joins multiple segments of a path.
func join(elem ...string) string {
	res := ""
	for _, e := range elem {
		if e == "" {
			continue
		}
		if res != "" {
			res += "/"
		}
		res += e
	}

	return res
}
