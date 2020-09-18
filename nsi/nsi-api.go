package nsi

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/USACE/go-consequences/consequences"
)

var apiUrl string = "https://nsi-dev.sec.usace.army.mil/nsiapi/structures"

func GetByBbox(bbox string) []consequences.Structure {
	structures := make([]consequences.Structure, 0)
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
	m := consequences.OccupancyTypeMap()
	defaultOcctype := m["RES1-1SNB"]
	fmt.Print(defaultOcctype)
	fmt.Print(string(jsonData))
	return structures

}
