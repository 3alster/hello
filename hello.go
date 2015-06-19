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
type citInfo struct {
	Snils int `json:"snils"`
}
type mzmkApplication struct {
	IncomingNum string  `json:"incomingNum"`
	CitizenInfo citInfo `json:"citizenInfo"`
}

type MZMK struct {
	ApplicationList []mzmkApplication `json:"applicationList"`
}

const basePath = "/Users/artur/Yandex.Disk/docs/pflb_prj/'15/05.22_Rstyle/mq/"

var certNum = 0

func mzmkMQh(w http.ResponseWriter, r *http.Request) {

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
var newMZMKs MZMK
var newMZRKs MZMK

func mzrkMQh(w http.ResponseWriter, r *http.Request) {
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
			str = fmt.Sprintf(str, num, cert.СтраховойНомер, cert.СерияСертификата, cert.НомерСертификата, ident)
			fmt.Fprintf(w, str)
		} else {
			http.Error(w, "Не загружены связи МЗМК <-> Сертификат", http.StatusPreconditionFailed)

		}

	}
}

func mzmkLoadedh(w http.ResponseWriter, r *http.Request) {

	var mzmk MZMK
	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)
	b := buf.Bytes()
	err := json.Unmarshal(b, &mzmk)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

	}
	fmt.Fprintf(w, "%+v", mzmk)
	newMZMKs = mzmk

}

func mzrkLoadedh(w http.ResponseWriter, r *http.Request) {

	var mzmk MZMK
	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)
	b := buf.Bytes()
	err := json.Unmarshal(b, &mzmk)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

	}
	fmt.Fprintf(w, "%+v", mzmk)
	newMZRKs = mzmk

}

func Testh(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")

	fmt.Fprintf(w, "RawQuery["+string(r.URL.RawQuery)+"]\n")
	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)
	b := buf.Bytes()
	s := *(*string)(unsafe.Pointer(&b))
	fmt.Fprintf(w, "Body["+s+"]")
}

func newMZMKh(w http.ResponseWriter, r *http.Request) {
	if len(newMZMKs.ApplicationList) > 0 {
		m := newMZMKs.ApplicationList[0].CitizenInfo.Snils
		newMZMKs.ApplicationList = append(newMZMKs.ApplicationList[:0],
			newMZMKs.ApplicationList[1:]...)
		fmt.Fprintf(w,"%d",m)

	} else {
		http.Error(w, "Не загружены заявления МЗМК", http.StatusPreconditionFailed)
		return
	}

}

func newMZRKh(w http.ResponseWriter, r *http.Request) {
	if len(newMZRKs.ApplicationList) > 0 {
		m := newMZRKs.ApplicationList[0].CitizenInfo.Snils
		newMZRKs.ApplicationList = append(newMZRKs.ApplicationList[:0],
			newMZRKs.ApplicationList[1:]...)
		fmt.Fprintf(w,"%d",m)
	} else {
		http.Error(w, "Не загружены заявления МЗРК", http.StatusPreconditionFailed)
		return
	}

}

type stat struct {
	MzmksNewPollSize  int
	MzrksNewPoolSize  int
	MzrkReadyPoolSize int
	LastCert          GivenCert
}

func statush(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	var m stat
	m.MzmksNewPollSize = len(newMZMKs.ApplicationList)
	m.MzrksNewPoolSize = len(newMZRKs.ApplicationList)
	m.MzrkReadyPoolSize = len(givenCerts)
	if m.MzrkReadyPoolSize > 0 {
		m.LastCert = givenCerts[m.MzrkReadyPoolSize-1]
	}
	b, err := json.Marshal(m)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, string(b))
}

type MQEndpoint struct {
	Host     string
	Port     int
	QManager string
	Channel  string
	QName    string
}

type HTTPEndpoint struct {
	Host  string
	Port  int
	Url   string
	Users []struct {
		Login string
		Pass  string
		Type  string
	}
}

type MSKConf struct {
	MQ   MQEndpoint
	HTTP HTTPEndpoint
}
type testConf struct {
	MSK MSKConf
}

func confh(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json;charset=UTF-8")

	var m testConf
	filename := path.Join(basePath, "conf.json")
	content, err := ioutil.ReadFile(filename)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = json.Unmarshal(content, &m)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

	}

	b, err := json.Marshal(m)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, string(b))

}
func linkMSKCerth(w http.ResponseWriter, r *http.Request) {
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

	}
}

func MSKCerth(w http.ResponseWriter, r *http.Request) {

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
	http.HandleFunc("/mzmkMQ", mzmkMQh)
	http.HandleFunc("/mzrkMQ", mzrkMQh)
	http.HandleFunc("/newCert", MSKCerth)
	http.HandleFunc("/linkCert", linkMSKCerth)

	http.HandleFunc("/saveMzmk", mzmkLoadedh)
	http.HandleFunc("/saveMzrk", mzrkLoadedh)
	http.HandleFunc("/newMzmk", newMZMKh)
	http.HandleFunc("/newMzrk", newMZRKh)
	http.HandleFunc("/config", confh)
	http.HandleFunc("/status", statush)
	http.HandleFunc("/test", Testh)

	http.ListenAndServe(":8080", nil)
}
