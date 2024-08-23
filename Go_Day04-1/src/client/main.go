package main

//#include "../cow/cow.c"
import "C"

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/lizrice/secure-connections/utils"
)

type OrderJSON struct {
	Money      int64  `json:"money"`
	CandyType  string `json:"candyType"`
	CandyCount int64  `json:"candyCount"`
}
type CandyBody struct {

	// change
	Change int64 `json:"change,omitempty"`

	// thanks
	Thanks string `json:"thanks,omitempty"`
}

func main() {

	caCert, err := ioutil.ReadFile("ca/minica.pem")
	if err != nil {
		fmt.Println(err)
		return
	}
	cert, err := x509.SystemCertPool()
	if err != nil {
		log.Println(err)
	}
	cert.AppendCertsFromPEM(caCert)

	config := &tls.Config{
		InsecureSkipVerify:    true,
		ClientAuth:            tls.RequireAndVerifyClientCert,
		RootCAs:               cert,
		GetCertificate:        utils.CertReqFunc("ca/client.candy.tld/cert.pem", "ca/client.candy.tld/key.pem"),
		VerifyPeerCertificate: utils.CertificateChains,
	}
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: config,
		},
	}

	k := flag.String("k", "", "candy type")
	c := flag.Int64("c", 0, "count of candy to buy")
	m := flag.Int64("m", 0, "amount of money")
	flag.Parse()
	order := OrderJSON{
		CandyType:  *k,
		Money:      *m,
		CandyCount: *c,
	}

	var data bytes.Buffer
	err = json.NewEncoder(&data).Encode(order)
	if err != nil {
		log.Fatal(err)
	}
	resp, err := client.Post("https://127.0.0.1:3333/v1/buy_candy", "application/json", bytes.NewBuffer(data.Bytes()))
	if err != nil {
		fmt.Println(err)
		return
	}
	var cb CandyBody
	// fmt.Println(resp.)
	if resp.StatusCode == 201 {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Println(err)
		}
		err = json.Unmarshal(body, &cb)
		if err != nil {
			fmt.Println(err)
			return
		}
		cb.Thanks = C.GoString(C.ask_cow(C.CString("Thank you!")))
		jsonData, err := json.Marshal(cb)
		fmt.Printf("%s", jsonData)
		defer resp.Body.Close()
	} else {
		_, err = io.Copy(os.Stdout, resp.Body)
		if err != nil {
			log.Fatal(err)
		}
	}

}
