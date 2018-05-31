// Copyright 2018 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"crypto/rand"
	"encoding/base64"
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"golang.org/x/crypto/ed25519"
)

var out = flag.String("out", "id_ed25519", "")

func main() {
	flag.Parse()

	publicKey, privateKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	publicKeyHex, privateKeyHex := base64.StdEncoding.EncodeToString(publicKey),
		base64.StdEncoding.EncodeToString(privateKey)

	err = ioutil.WriteFile(*out, []byte(privateKeyHex), 0644)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	err = ioutil.WriteFile(*out+".pub", []byte(publicKeyHex), 0644)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
