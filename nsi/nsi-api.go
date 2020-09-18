package nsi

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

var apiUrl string = "https://nsi-dev.sec.usace.army.mil/nsiapi/structures"

func GetByBbox(bbox string) {
	transCfg := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, // accept untrusted servers
	}
	client := &http.Client{Transport: transCfg}
	url := fmt.Sprintf("%s?bbox=%s", apiUrl, bbox)
	//fmt.Println(url)
	response, err := client.Get(url)

	if err != nil {
		log.Fatal(err)
	}

	defer response.Body.Close()

	jsonData, err := ioutil.ReadAll(response.Body)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Print(string(jsonData))

}
