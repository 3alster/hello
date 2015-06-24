package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

type Number struct {
	Number string `json:"number"`
}

type Blanks struct {
	Blanks []Number `json:"blanks"`
}

var emptyBlank Blanks

//get empty blanks certificates and save its 
func getCertEmpty(w http.ResponseWriter, r *http.Request) {

	var cert Blanks

	var hah []byte
	hah, _ = ioutil.ReadAll(r.Body)

	s := string(hah)

	u, _ := url.QueryUnescape(s)

	err := json.Unmarshal([]byte(u), &cert)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

	}
	fmt.Fprintf(w, "%+v", cert)
	emptyBlank = cert
}
