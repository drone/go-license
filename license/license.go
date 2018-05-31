// Copyright 2018 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package license

import (
	"bytes"
	"encoding/json"
	"encoding/pem"
	"errors"
	"io/ioutil"
	"time"

	"golang.org/x/crypto/ed25519"
)

var (
	// ErrInvalidSignature indicates the license signature is invalid.
	ErrInvalidSignature = errors.New("Invalid signature")

	// ErrMalformedLicense indicates the license file is malformed.
	ErrMalformedLicense = errors.New("Malformed License")
)

// License defines a software license key.
type License struct {
	Iss string          `json:"iss,omitempty"` // Issued By
	Cus string          `json:"cus,omitempty"` // Customer ID
	Sub string          `json:"sub,omitempty"` // Subscriber ID
	Typ string          `json:"typ,omitempty"` // License Type
	Lim int             `json:"lim,omitempty"` // License Limit (e.g. Seats)
	Iat time.Time       `json:"iat,omitempty"` // Issued At
	Exp time.Time       `json:"exp,omitempty"` // Expires At
	Dat json.RawMessage `json:"dat,omitempty"` // Metadata
}

// Expired returns true if the license is expired.
func (l *License) Expired() bool {
	return l.Exp.IsZero() == false && time.Now().After(l.Exp)
}

// Encode generates returns a PEM encoded license key that is
// signed with the ed25519 private key.
func (l *License) Encode(privateKey ed25519.PrivateKey) ([]byte, error) {
	msg, err := json.Marshal(l)
	if err != nil {
		return nil, err
	}

	sig := ed25519.Sign(privateKey, msg)
	buf := new(bytes.Buffer)
	buf.Write(sig)
	buf.Write(msg)

	block := &pem.Block{
		Type:  "LICENSE KEY",
		Bytes: buf.Bytes(),
	}
	return pem.EncodeToMemory(block), nil
}

// Decode decodes the PEM encoded license key and verifies
// the content signature using the ed25519 public key.
func Decode(data []byte, publicKey ed25519.PublicKey) (*License, error) {
	block, _ := pem.Decode(data)
	if block == nil || len(block.Bytes) < ed25519.SignatureSize {
		return nil, ErrMalformedLicense
	}

	sig := block.Bytes[:ed25519.SignatureSize]
	msg := block.Bytes[ed25519.SignatureSize:]

	verified := ed25519.Verify(publicKey, msg, sig)
	if !verified {
		return nil, ErrInvalidSignature
	}
	out := new(License)
	err := json.Unmarshal(msg, out)
	return out, err
}

// DecodeFile decodes the PEM encoded license file and verifies
// the content signature using the ed25519 public key.
func DecodeFile(path string, publicKey ed25519.PublicKey) (*License, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return Decode([]byte(data), publicKey)
}
