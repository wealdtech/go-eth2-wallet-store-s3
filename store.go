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
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	session "github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/pkg/errors"
	util "github.com/wealdtech/go-eth2-util"
	wtypes "github.com/wealdtech/go-eth2-wallet-types/v2"
)

// options are the options for the S3 store
type options struct {
	id         []byte
	region     string
	passphrase []byte
}

// Option gives options to New
type Option interface {
	apply(*options)
}

type optionFunc func(*options)

func (f optionFunc) apply(o *options) {
	f(o)
}

// WithPassphrase sets the passphrase for the store.
func WithPassphrase(passphrase []byte) Option {
	return optionFunc(func(o *options) {
		o.passphrase = passphrase
	})
}

// WithID sets the ID for the store
func WithID(t []byte) Option {
	return optionFunc(func(o *options) {
		o.id = t
	})
}

// WithRegion sets the AWS region for the store
func WithRegion(t string) Option {
	return optionFunc(func(o *options) {
		o.region = t
	})
}

// Store is the store for the wallet held encrypted on Amazon S3.
type Store struct {
	session    *session.Session
	id         []byte
	region     string
	bucket     string
	passphrase []byte
}

// New creates a new Amazon S3 store.
// This takes the following options:
//  - region: a string specifying the Amazon S3 region, defaults to "us-east-1", set with WithRegion()
//  - id: a byte array specifying an identifying key for the store, defaults to nil, set with WithID()
// This expects the access credentials to be in a standard place, e.g. ~/.aws/credentials
func New(opts ...Option) (wtypes.Store, error) {
	options := options{
		region: "us-east-1",
	}
	for _, o := range opts {
		o.apply(&options)
	}

	session, err := session.NewSession(&aws.Config{Region: aws.String(options.region)})
	if err != nil {
		return nil, err
	}

	creds, err := session.Config.Credentials.Get()
	if err != nil {
		return nil, err
	}
	cryptKeyCopy := make([]byte, len(creds.AccessKeyID))
	copy(cryptKeyCopy, creds.AccessKeyID)

	// Generate a bucket name from the cryptKey.  This will be the SHA256 hash of a string unique to the account, as a hex string
	// of 63 charaters (as S3 only allows bucket names up to 63 characters in length).
	hash := util.SHA256([]byte(fmt.Sprintf("Ethereum 2 wallet:%s", creds.AccessKeyID)), options.id)
	bucket := hex.EncodeToString(hash)[:63]

	// Check the bucket exists; if not create it
	conn := s3.New(session)
	_, err = conn.GetBucketAcl(&s3.GetBucketAclInput{Bucket: &bucket})
	if err != nil {
		if !strings.Contains(err.Error(), "NoSuchBucket") {
			return nil, errors.Wrap(err, "unable to access bucket")
		}
		// Create the bucket
		_, err = conn.CreateBucket(&s3.CreateBucketInput{Bucket: aws.String(bucket)})
		if err != nil {
			return nil, errors.Wrap(err, "unable to create bucket")
		}
		err = conn.WaitUntilBucketExists(&s3.HeadBucketInput{Bucket: aws.String(bucket)})
		if err != nil {
			return nil, errors.Wrap(err, "failed to confirm bucket creation")
		}
	}

	return &Store{
		session:    session,
		region:     options.region,
		id:         options.id,
		bucket:     bucket,
		passphrase: options.passphrase,
	}, nil
}

// Name returns the name of this store.
func (s *Store) Name() string {
	return "s3"
}

// Location returns the location of this store.
func (s *Store) Location() string {
	return s.bucket
}
