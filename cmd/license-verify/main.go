// Copyright 2018 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/drone/go-license/license"
	"github.com/drone/go-license/license/licenseutil"
)

var (
	pub  = flag.String("pub", "id_ed25519.pub", "")
	file = flag.String("file", "license.txt", "")
)

func main() {
	flag.Parse()

	publicKey, err := licenseutil.ReadPublicKey(*pub)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	l, err := license.DecodeFile(*file, publicKey)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	enc.Encode(l)
}
