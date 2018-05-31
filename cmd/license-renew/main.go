// Copyright 2018 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/drone/go-license/license"
	"github.com/drone/go-license/license/licenseutil"
)

var (
	pub  = flag.String("pub", "id_ed25519.pub", "")
	file = flag.String("file", "license.txt", "")
	addr = flag.String("addr", "http://localhost:9000", "")
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

	buf := new(bytes.Buffer)
	json.NewEncoder(buf).Encode(l)

	res, err := http.Post(*addr, "application/json", buf)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer res.Body.Close()
	io.Copy(os.Stdout, res.Body)
}
