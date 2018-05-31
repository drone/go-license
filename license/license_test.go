// Copyright 2018 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package license

import (
	"os"
	"testing"
	"time"

	"github.com/drone/go-license/license/licenseutil"
)

func TestExpired(t *testing.T) {
	license := &License{}
	if license.Expired() {
		t.Errorf("Expect zero value expiration to never expire")
	}

	license = &License{}
	license.Exp = time.Now().Add(time.Hour)
	if license.Expired() == true {
		t.Errorf("Expect license is not expired")
	}

	license = &License{}
	license.Exp = time.Now().Add(time.Hour * -1)
	if license.Expired() == false {
		t.Errorf("Expect license is expired")
	}
}

func TestEncodeDecode(t *testing.T) {
	privateKey, err := licenseutil.ReadPrivateKey("licenseutil/testdata/id_ed25519")
	if err != nil {
		t.Error(err)
	}
	publicKey, err := licenseutil.ReadPublicKey("licenseutil/testdata/id_ed25519.pub")
	if err != nil {
		t.Error(err)
	}
	license := &License{
		Iss: "Acme, Inc",
		Cus: "cus_CxoyqaC4p4Hjl0",
		Sub: "sub_A4l9XkCxyZPcS2",
		Typ: "trial",
		Lim: 50,
		Iat: time.Now().UTC(),
		Exp: time.Now().Add(time.Hour).UTC(),
	}
	encoded, err := license.Encode(privateKey)
	if err != nil {
		t.Error(err)
	}
	decoded, err := Decode(encoded, publicKey)
	if err != nil {
		t.Error(err)
	}
	if got, want := decoded.Cus, license.Cus; got != want {
		t.Errorf("Want license Cus %v, got %v", want, got)
	}
	if got, want := decoded.Exp, license.Exp; got != want {
		t.Errorf("Want license Exp %v, got %v", want, got)
	}
	if got, want := decoded.Iat, license.Iat; got != want {
		t.Errorf("Want license Iat %v, got %v", want, got)
	}
	if got, want := decoded.Iss, license.Iss; got != want {
		t.Errorf("Want license Iss %v, got %v", want, got)
	}
	if got, want := decoded.Lim, license.Lim; got != want {
		t.Errorf("Want license Lim %v, got %v", want, got)
	}
	if got, want := decoded.Sub, license.Sub; got != want {
		t.Errorf("Want license Sub %v, got %v", want, got)
	}
	if got, want := decoded.Typ, license.Typ; got != want {
		t.Errorf("Want license Sub %v, got %v", want, got)
	}
}

func TestDecodeFile(t *testing.T) {
	publicKey, err := licenseutil.ReadPublicKey("licenseutil/testdata/id_ed25519.pub")
	if err != nil {
		t.Error(err)
	}
	decoded, err := DecodeFile("licenseutil/testdata/license.txt", publicKey)
	if err != nil {
		t.Error(err)
	}
	if got, want := decoded.Cus, "cus_CxoyqaC4p4Hjl0"; got != want {
		t.Errorf("Want license Cus %v, got %v", want, got)
	}
	if got, want := decoded.Sub, "sub_A4l9XkCxyZPcS2"; got != want {
		t.Errorf("Want license Exp %v, got %v", want, got)
	}
}

func TestDecodeFile_InvalidSignature(t *testing.T) {
	publicKey, err := licenseutil.ReadPublicKey("licenseutil/testdata/id_ed25519.pub")
	if err != nil {
		t.Error(err)
	}
	_, err = DecodeFile("licenseutil/testdata/license_invalid_signature.txt", publicKey)
	if err != ErrInvalidSignature {
		t.Errorf("Expected invalid signature error")
	}
}

func TestDecodeFile_InvalidEncoding(t *testing.T) {
	publicKey, err := licenseutil.ReadPublicKey("licenseutil/testdata/id_ed25519.pub")
	if err != nil {
		t.Error(err)
	}
	_, err = DecodeFile("licenseutil/testdata/license_invalid_encoding.txt", publicKey)
	if err != ErrMalformedLicense {
		t.Errorf("Expected invalid signature error")
	}
}

func TestDecodeFile_InvalidJson(t *testing.T) {
	publicKey, err := licenseutil.ReadPublicKey("licenseutil/testdata/id_ed25519.pub")
	if err != nil {
		t.Error(err)
	}
	_, err = DecodeFile("licenseutil/testdata/license_invalid_json.txt", publicKey)
	if err != ErrMalformedLicense {
		t.Errorf("Expected invalid signature error")
	}
}

func TestDecodeFile_PathError(t *testing.T) {
	publicKey, err := licenseutil.ReadPublicKey("licenseutil/testdata/id_ed25519.pub")
	if err != nil {
		t.Error(err)
	}
	_, err = DecodeFile("path/does/not/exist", publicKey)
	if _, ok := err.(*os.PathError); !ok {
		t.Errorf("Expect path error")
	}
}
