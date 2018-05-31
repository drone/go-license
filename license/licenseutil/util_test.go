// Copyright 2018 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package licenseutil

import (
	"os"
	"testing"

	"golang.org/x/crypto/ed25519"
)

func TestReadKeyPair(t *testing.T) {
	publicKey, err := ReadPublicKey("testdata/id_ed25519.pub")
	if err != nil {
		t.Error(err)
	}
	privateKey, err := ReadPrivateKey("testdata/id_ed25519")
	if err != nil {
		t.Error(err)
	}
	msg := []byte("hello world")
	sig := ed25519.Sign(privateKey, msg)
	if !ed25519.Verify(publicKey, msg, sig) {
		t.Errorf("Cannot sign and verify. Are keys malformed?")
	}
}

func TestReadPublicKey_NotFound(t *testing.T) {
	_, err := ReadPublicKey("does/not/exist")
	if _, ok := err.(*os.PathError); !ok {
		t.Errorf("Expect path error")
	}
}

func TestReadPrivateKey_NotFound(t *testing.T) {
	_, err := ReadPrivateKey("does/not/exist")
	if _, ok := err.(*os.PathError); !ok {
		t.Errorf("Expect path error")
	}
}
