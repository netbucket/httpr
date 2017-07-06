// Copyright © 2017 Igor Bondarenko <ibondare@protonmail.com>
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

package tls

import (
	"crypto/ecdsa"
	"crypto/rsa"
	"testing"
)

func TestKeyLength(t *testing.T) {
	cert, err := generateSelfSignedCert()

	if err != nil {
		t.Fatalf("Unexpected error: %v\n", err)
	}

	var keyLength int
	switch key := cert.PrivateKey.(type) {
	case *rsa.PrivateKey:
		keyLength = key.N.BitLen()
	case *ecdsa.PrivateKey:
		keyLength = key.Curve.Params().BitSize
	default:
		t.Fatal("unsupported private key")
	}

	if keyLength < 2048 {
		t.Errorf("Private key length is too small:  %d\n", keyLength)
	}
}
