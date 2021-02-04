package nsi

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

//NsiProperties is a reflection of the JSON feature property attributes from the NSI-API
type NsiProperties struct {
	Name      string  `json:"fd_id"`
	X         float64 `json:"x"`
	Y         float64 `json:"y"`
	Occtype   string  `json:"occtype"`
	FoundHt   float64 `json:"found_ht"`
	DamCat    string  `json:"st_damcat"`
	StructVal float64 `json:"val_struct"`
	ContVal   float64 `json:"val_cont"`
	CB        string  `json:"cbfips"`
	Pop2amu65 int32   `json:"pop2amu65"`
	Pop2amo65 int32   `json:"pop2amo65"`
	Pop2pmu65 int32   `json:"pop2pmu65"`
	Pop2pmo65 int32   `json:"pop2pmo65"`
}

//NsiFeature is a feature which contains the properties of a structure from the NSI API
type NsiFeature struct {
	Properties NsiProperties `json:"properties"`
}

//NsiInventory is a slice of NsiFeature that describes a complete json feature array return or feature collection return
type NsiInventory struct {
	Features []NsiFeature
}

var apiURL string = "https://nsi-dev.sec.usace.army.mil/nsiapi/structures" //this will only work behind the USACE firewall -
//GetByFips returns an NsiInventory for a FIPS code
func GetByFips(fips string) NsiInventory {
	url := fmt.Sprintf("%s?fips=%s&fmt=fa", apiURL, fips)
	return nsiAPI(url)
}

//GetByBbox returns an NsiInventory for a Bounding Box
func GetByBbox(bbox string) NsiInventory {
	url := fmt.Sprintf("%s?bbox=%s&fmt=fa", apiURL, bbox)
	return nsiAPI(url)
}
func nsiAPI(url string) NsiInventory {
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

//NsiStreamProcessor is a function used to process an in memory NsiFeature through the NsiStreaming service endpoints
type NsiStreamProcessor func(str NsiFeature)

/*
memory effecient structure compute methods
*/

//GetByFipsStream a streaming service for NsiFeature based on a FIPs code
func GetByFipsStream(fips string, nsp NsiStreamProcessor) error {
	url := fmt.Sprintf("%s?fips=%s&fmt=fs", apiURL, fips)
	return nsiAPIStream(url, nsp)
}

//GetByBboxStream a streaming service for NsiFeature based on a bounding box
func GetByBboxStream(bbox string, nsp NsiStreamProcessor) error {
	url := fmt.Sprintf("%s?bbox=%s&fmt=fs", apiURL, bbox)
	return nsiAPIStream(url, nsp)
}
func nsiAPIStream(url string, nsp NsiStreamProcessor) error {
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
	//resp, _ := ioutil.ReadAll(response.Body)
	//s := string(resp)
	//fmt.Println(s)
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
