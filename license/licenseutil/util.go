// Copyright 2018 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package licenseutil

//go:generate license-keygen -out testdata/id_ed25519
//go:generate license-create -in testdata/id_ed25519 -out testdata/license.txt -cus cus_CxoyqaC4p4Hjl0 -sub sub_A4l9XkCxyZPcS2

import (
	"encoding/base64"
	"io/ioutil"

	"golang.org/x/crypto/ed25519"
)

// ReadPublicKey reads a base64-encoded public key file.
func ReadPublicKey(path string) (ed25519.PublicKey, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return DecodePublicKey(data)
}

// ReadPrivateKey reads a base64-encoded private key file.
func ReadPrivateKey(path string) (ed25519.PrivateKey, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return DecodePrivateKey(data)
}

// DecodePublicKey decodes a base64-encoded private key.
func DecodePublicKey(data []byte) (ed25519.PublicKey, error) {
	decoded, err := decode(data)
	if err != nil {
		return nil, err
	}
	return ed25519.PublicKey(decoded), nil
}

// DecodePrivateKey decodes a base64-encoded private key.
func DecodePrivateKey(data []byte) (ed25519.PrivateKey, error) {
	decoded, err := decode(data)
	if err != nil {
		return nil, err
	}
	return ed25519.PrivateKey(decoded), nil
}

func decode(b []byte) ([]byte, error) {
	enc := base64.StdEncoding
	buf := make([]byte, enc.DecodedLen(len(b)))
	n, err := enc.Decode(buf, b)
	return buf[:n], err
}
