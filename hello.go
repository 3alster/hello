package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	//"log"
	"bytes"
	"net/http"
	"path"
	"time"
	"unsafe"
)

type GivenCert struct {
	СерияСертификата string
	НомерСертификата int
	СтраховойНомер   string
}

type MZMK struct {
	Total  int `json:"total"`
	ApplicationList []struct{
		Id int `json:"id"`
		} `json:"applicationList"`
}


const basePath = "/Users/artur/Yandex.Disk/docs/pflb_prj/'15/05.22_Rstyle/mq/"

var certNum = 0

func mzmkh(w http.ResponseWriter, r *http.Request) {

	filename := path.Join(basePath, "мзмкTempl-1.xml")
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return

	} else {
		//		w.Header().Set("Content-Type", "text/xml;charset=Windows-1251")

		t := time.Now()
		str := string(content)
		//ident := fmt.Sprintf("%02d%02d/%02d", t.Day(), t.Month(), t.Year()%100)  //ddmm/yy
		num := fmt.Sprintf("%02d%02d%02d/%02d", t.Hour(), t.Minute(), t.Second(), t.Day()) //mmss/hh
		str = fmt.Sprintf(str, newSnils(), num, num)

		fmt.Fprintf(w, str)
	}

}

var givenCerts []GivenCert

func mzrkh(w http.ResponseWriter, r *http.Request) {
	filename := path.Join(basePath, "мзркTempl-4.xml")
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return

	} else {
		t := time.Now()
		const format = "%02d%02d/%02d"
		str := string(content)
		num := fmt.Sprintf(format, t.Day(), t.Month(), t.Year()%100)   //ddmm/yy
		ident := fmt.Sprintf(format, t.Minute(), t.Second(), t.Hour()) //mmss/hh

		var cert GivenCert

		if len(givenCerts) > 0 {
			cert = givenCerts[0]
			givenCerts = append(givenCerts[:0], givenCerts[1:]...)
		}
		str = fmt.Sprintf(str, num, cert.СтраховойНомер, cert.СерияСертификата, cert.НомерСертификата, ident)
		fmt.Fprintf(w, str)
	}
}


func mzmkLoadedh(w http.ResponseWriter, r *http.Request) {

	var mzmk MZMK
	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)
	b := buf.Bytes()
	//fmt.Fprintf(w, "Body %s\n", b)

	err := json.Unmarshal(b, &mzmk)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

	}
	fmt.Fprintf(w, "%+v", mzmk)
}

func Testh(w http.ResponseWriter, r *http.Request) {

	fmt.Fprintf(w, "RawQuery["+string(r.URL.RawQuery)+"]\n")
	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)
	b := buf.Bytes()
	s := *(*string)(unsafe.Pointer(&b))
	fmt.Fprintf(w, "Body["+s+"]")
}

func MSKCerth(w http.ResponseWriter, r *http.Request) {

	if r.Method == "POST" {

		var cert GivenCert
		buf := new(bytes.Buffer)
		buf.ReadFrom(r.Body)
		b := buf.Bytes()
		err := json.Unmarshal(b, &cert)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)

		}
		fmt.Fprintf(w, "%+v", cert)
		givenCerts = append(givenCerts, cert)
		fmt.Fprintf(w, "len=%d", len(givenCerts))

	} else {
		w.Header().Set("Content-Type", "application/json;charset=UTF-8")

		type MSKCert struct {
			СерияСертификата string
			НомерСертификата int
		}
		cert := MSKCert{СерияСертификата: "МК-5",
			НомерСертификата: 9999000 + certNum,
		}
		certNum++
		b, err := json.Marshal(cert)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		fmt.Fprintf(w, string(b))
	}

}

var snilsCnt = 0

func newSnils() string {
	t := time.Now()
	snilsCnt++
	snilsCnt = snilsCnt % 100
	snils := fmt.Sprintf("%02d%07d", snilsCnt, t.Unix()%1e+7)
	check := 0
	for i := 0; i < 9; i++ {
		check += int(snils[i]-'0') * (9 - i)
	}

	if check == 100 || check == 101 {
		check = 0
	} else {
		check %= 101
	}

	snils = fmt.Sprintf("%s%02d", snils, check)
	snils = fmt.Sprintf("%s-%s-%s %s", snils[0:3], snils[3:6], snils[6:9], snils[9:11])
	return snils
}

func main() {
	http.HandleFunc("/mzmk", mzmkh)
	http.HandleFunc("/mzrk", mzrkh)
	http.HandleFunc("/cert", MSKCerth)
	http.HandleFunc("/test", Testh)
	http.HandleFunc("/mzmkLoaded", mzmkLoadedh)

	http.ListenAndServe(":8080", nil)
}
