package nsi

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/USACE/go-consequences/consequences"
)

type NSIproperties struct {
	Name string  `json:"fd_id"`
	X    float64 `json:"x"`
	Y    float64 `json:"y"`
}
type NSIfeature struct {
	Properties NSIproperties `json:"properties"`
}

var apiUrl string = "https://nsi-dev.sec.usace.army.mil/nsiapi/structures" //this will only work behind the USACE firewall -

func GetByBbox(bbox string) []consequences.Structure {
	structures := make([]consequences.Structure, 0)
	transCfg := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, // accept untrusted servers
	}
	client := &http.Client{Transport: transCfg}
	url := fmt.Sprintf("%s?bbox=%s", apiUrl, bbox)
	fmt.Println(url)
	response, err := client.Get(url)

	if err != nil {
		log.Fatal(err)
	}

	defer response.Body.Close()
	// UnmarshalJSON implements UnmarshalJSON interface
	jsonData, err := ioutil.ReadAll(response.Body)
	c := make([]NSIfeature, 0)
	if err := json.Unmarshal(jsonData, &c); err != nil {
		return structures
	}
	m := consequences.OccupancyTypeMap()
	defaultOcctype := m["RES1-1SNB"]
	fmt.Print(defaultOcctype)

	fmt.Print(string(c[0].Properties.Name))
	return structures

}
