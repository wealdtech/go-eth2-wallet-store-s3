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

package s3

import (
	"fmt"

	"github.com/google/uuid"
)

func (s *Store) walletPath(walletID uuid.UUID) string {
	return walletID.String()
}

func (s *Store) walletHeaderPath(walletID uuid.UUID) string {
	return fmt.Sprintf("%s/%s", s.walletPath(walletID), s.walletPath(walletID))
}

func (s *Store) accountPath(walletID uuid.UUID, accountID uuid.UUID) string {
	return fmt.Sprintf("%s/%s", s.walletPath(walletID), accountID.String())
}

func (s *Store) walletIndexPath(walletID uuid.UUID) string {
	return fmt.Sprintf("%s/index", s.walletPath(walletID))
}
