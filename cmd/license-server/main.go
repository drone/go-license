// Copyright 2018 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"encoding/json"
	"flag"
	"net/http"
	"time"

	"github.com/drone/go-license/license"
	"github.com/drone/go-license/license/licenseutil"
)

var monthly = time.Hour * 24 * 30

var (
	addr    = flag.String("addr", ":9000", "")
	renewal = flag.Duration("dur", monthly, "")
	signer  = flag.String("signer", "id_ed25519", "")
)

func main() {
	flag.Parse()

	http.HandleFunc("/", handler)
	http.ListenAndServe(*addr, nil)
}

func handler(w http.ResponseWriter, r *http.Request) {
	l := new(license.License)
	err := json.NewDecoder(r.Body).Decode(l)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// TODO verify the License.Cus (customer id) is valid
	// and in good standing prior to issuing a new license.

	l.Exp = time.Now().UTC().Add(*renewal)

	privateKey, err := licenseutil.ReadPrivateKey(*signer)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	encoded, err := l.Encode(privateKey)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(encoded)
}
