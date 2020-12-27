package nsi

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

type NsiProperties struct {
	Name      string  `json:"fd_id"`
	X         float64 `json:"x"`
	Y         float64 `json:"y"`
	Occtype   string  `json:"occtype"`
	FoundHt   float64 `json:"found_ht"`
	DamCat    string  `json:"st_damcat"`
	StructVal float64 `json:"val_struct"`
	ContVal   float64 `json:"val_cont"`
}
type NsiFeature struct {
	Properties NsiProperties `json:"properties"`
}
type NsiInventory struct {
	Features []NsiFeature
}

var apiUrl string = "https://nsi-dev.sec.usace.army.mil/nsiapi/structures" //this will only work behind the USACE firewall -
func GetByFips(fips string) NsiInventory {
	url := fmt.Sprintf("%s?fips=%s&fmt=fa", apiUrl, fips)
	return nsiApi(url)
}
func GetByBbox(bbox string) NsiInventory {
	url := fmt.Sprintf("%s?bbox=%s&fmt=fa", apiUrl, bbox)
	return nsiApi(url)
}
func nsiApi(url string) NsiInventory {
	inventory := NsiInventory{}
	transCfg := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, // accept untrusted servers
	}
	client := &http.Client{Transport: transCfg}

	response, err := client.Get(url)

	if err != nil {
		fmt.Println(err)
		return inventory
	}
	defer response.Body.Close()
	jsonData, err := ioutil.ReadAll(response.Body)
	features := make([]NsiFeature, 0)
	if err := json.Unmarshal(jsonData, &features); err != nil {
		fmt.Println("error unmarshalling NSI json " + err.Error() + " URL: " + url)
		s := string(jsonData)
		fmt.Println("first 1000 chars of jsonbody was: " + s[0:1000]) //s) //"last 100 chars of jsonbody was: " + s[len(s)-100:])
		return inventory
	}
	inventory.Features = features
	return inventory
}

type NsiStreamProcessor func(str NsiFeature)

/*
memory effecient structure compute methods
*/
func GetByFipsStream(fips string, nsp NsiStreamProcessor) error {
	url := fmt.Sprintf("%s?fips=%s&fmt=fs", apiUrl, fips)
	return nsiApiStream(url, nsp)
}
func GetByBboxStream(bbox string, nsp NsiStreamProcessor) error {
	url := fmt.Sprintf("%s?bbox=%s&fmt=fs", apiUrl, bbox)
	return nsiApiStream(url, nsp)
}
func nsiApiStream(url string, nsp NsiStreamProcessor) error {
	transCfg := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, // accept untrusted servers
	}
	client := &http.Client{Transport: transCfg}

	response, err := client.Get(url)

	if err != nil {
		fmt.Println(err)
		return err
	}
	defer response.Body.Close()
	dec := json.NewDecoder(response.Body)
	for {
		var n NsiFeature
		if err := dec.Decode(&n); err == io.EOF {
			break
		} else if err != nil {
			fmt.Printf("Error unmarshalling JSON record: %s.  Stopping Compute.\n", err)
			return err
		}
		nsp(n)
	}
	return nil
}
