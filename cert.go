package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"sync"
)

type Number struct {
	Number       string `json:"number"`
	sync.RWMutex        //synchronized
}

type Blanks struct {
	Blanks       []Number `json:"blanks"`
	sync.RWMutex          //synchronized
}

var emptyBlank Blanks
var sliceEmptyBlank []Number

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
	emptyBlank.RLock()
	emptyBlank = cert
	// sliceEmptyBlank = emptyBlank.Blanks[0:]
	emptyBlank.RUnlock()
	fmt.Fprintf(w, "%+v", cert)
}

//return number certificate
func numberEmptyCert(w http.ResponseWriter, r *http.Request) {
	emptyBlank.RLock()
	fmt.Fprintf(w, emptyBlank.Blanks[0].Number)
	emptyBlank.Blanks = emptyBlank.Blanks[1:]
	emptyBlank.RUnlock()
}
