package structureprovider

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/USACE/go-consequences/consequences"
	"github.com/USACE/go-consequences/geography"
	"github.com/USACE/go-consequences/structures"
)

//NsiProperties is a reflection of the JSON feature property attributes from the NSI-API
type NsiProperties struct {
	Name      string  `json:"fd_id"`
	X         float64 `json:"x"`
	Y         float64 `json:"y"`
	Occtype   string  `json:"occtype"`
	FoundHt   float64 `json:"found_ht"`
	FoundType string  `json:"found_type"`
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
type nsiStreamProvider struct {
	ApiURL string
}

func InitNSISP() nsiStreamProvider {
	return nsiStreamProvider{ApiURL: "https://nsi-dev.sec.usace.army.mil/nsiapi/structures"}
}
func (nsp nsiStreamProvider) ByFips(fipscode string, sp StreamProcessor) {
	url := fmt.Sprintf("%s?fips=%s&fmt=fs", nsp.ApiURL, fipscode)
	nsiStructureStream(url, sp)
}
func (nsp nsiStreamProvider) ByBbox(bbox geography.BBox, sp StreamProcessor) {
	url := fmt.Sprintf("%s?bbox=%s&fmt=fs", nsp.ApiURL, bbox.ToString())
	nsiStructureStream(url, sp)
}
func nsiStructureStream(url string, sp StreamProcessor) {
	m := structures.OccupancyTypeMap()
	//define a default occtype in case of emergancy
	defaultOcctype := m["RES1-1SNB"]
	transCfg := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, // accept untrusted servers
	}
	client := &http.Client{Transport: transCfg}

	response, err := client.Get(url)

	if err != nil {
		fmt.Println(err)
	}
	defer response.Body.Close()
	dec := json.NewDecoder(response.Body)
	for {
		var n NsiFeature
		if err := dec.Decode(&n); err == io.EOF {
			break
		} else if err != nil {
			fmt.Printf("Error unmarshalling JSON record: %s.  Stopping Compute.\n", err)
		}
		sp(NsiFeaturetoStructure(n, m, defaultOcctype))
	}
}

//NsiFeaturetoStructure converts an nsi.NsiFeature to a structures.Structure
func NsiFeaturetoStructure(f NsiFeature, m map[string]structures.OccupancyTypeStochastic, defaultOcctype structures.OccupancyTypeStochastic) structures.StructureStochastic {
	var occtype = defaultOcctype
	if ot, ok := m[f.Properties.Occtype]; ok {
		occtype = ot
	} else {
		occtype = defaultOcctype
		msg := "Using default " + f.Properties.Occtype + " not found"
		panic(msg)
	}
	return structures.StructureStochastic{
		OccType:   occtype,
		StructVal: consequences.ParameterValue{Value: f.Properties.StructVal},
		ContVal:   consequences.ParameterValue{Value: f.Properties.ContVal},
		FoundHt:   consequences.ParameterValue{Value: f.Properties.FoundHt},
		BaseStructure: structures.BaseStructure{
			Name:   f.Properties.Name,
			DamCat: f.Properties.DamCat,
			X:      f.Properties.X,
			Y:      f.Properties.Y,
		},
	}
}

//GetByFips returns an NsiInventory for a FIPS code
func GetByFips(fips string) NsiInventory {
	n := InitNSISP()
	url := fmt.Sprintf("%s?fips=%s&fmt=fa", n.ApiURL, fips)
	return nsiAPI(url)
}

//GetByBbox returns an NsiInventory for a Bounding Box
func GetByBbox(bbox string) NsiInventory {
	n := InitNSISP()
	url := fmt.Sprintf("%s?bbox=%s&fmt=fa", n.ApiURL, bbox)
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
	n := InitNSISP()
	url := fmt.Sprintf("%s?fips=%s&fmt=fs", n.ApiURL, fips)
	return nsiAPIStream(url, nsp)
}

//GetByBboxStream a streaming service for NsiFeature based on a bounding box
func GetByBboxStream(bbox string, nsp NsiStreamProcessor) error {
	n := InitNSISP()
	url := fmt.Sprintf("%s?bbox=%s&fmt=fs", n.ApiURL, bbox)
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
