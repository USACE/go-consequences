package nsi

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/USACE/go-consequences/consequences"
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

func (i NsiInventory) toStructures() []consequences.StructureStochastic {
	m := consequences.OccupancyTypeMap()
	defaultOcctype := m["RES1-1SNB"]
	var occtype = defaultOcctype
	structures := make([]consequences.StructureStochastic, len(i.Features))
	for idx, feature := range i.Features {

		if ot, ok := m[feature.Properties.Occtype]; ok {
			occtype = ot
		} else {
			occtype = defaultOcctype
			msg := "Using default " + feature.Properties.Occtype + " not found"
			panic(msg)
		}
		structures[idx] = consequences.StructureStochastic{
			Name:      feature.Properties.Name,
			OccType:   occtype,
			DamCat:    feature.Properties.DamCat,
			StructVal: consequences.ParameterValue{Value: feature.Properties.StructVal},
			ContVal:   consequences.ParameterValue{Value: feature.Properties.ContVal},
			FoundHt:   consequences.ParameterValue{Value: feature.Properties.FoundHt},
			X:         feature.Properties.X,
			Y:         feature.Properties.Y,
		}
	}
	return structures
}

var apiUrl string = "https://nsi-dev.sec.usace.army.mil/nsiapi/structures" //this will only work behind the USACE firewall -
func GetByFips(fips string) []consequences.StructureStochastic {
	url := fmt.Sprintf("%s?fips=%s", apiUrl, fips)
	return nsiApi(url)
}
func GetByBbox(bbox string) []consequences.StructureStochastic {
	url := fmt.Sprintf("%s?bbox=%s", apiUrl, bbox)
	return nsiApi(url)
}
func nsiApi(url string) []consequences.StructureStochastic {
	structures := make([]consequences.StructureStochastic, 0)
	transCfg := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, // accept untrusted servers
	}
	client := &http.Client{Transport: transCfg}

	response, err := client.Get(url)

	if err != nil {
		fmt.Println(err)
		return structures
	}
	defer response.Body.Close()
	jsonData, err := ioutil.ReadAll(response.Body)
	features := make([]NsiFeature, 0)

	if err := json.Unmarshal(jsonData, &features); err != nil {
		fmt.Println("error unmarshalling NSI json " + err.Error() + " URL: " + url)
		s := string(jsonData)
		fmt.Println("last 100 chars of jsonbody was: " + s[len(s)-10:])
		return structures
	}
	inventory := NsiInventory{Features: features}
	structures = inventory.toStructures()
	return structures
}

type NsiStreamProcessor func(str consequences.StructureStochastic)

/*
memory effecient structure compute methods
*/
func GetByFipsStream(fips string, nsp NsiStreamProcessor) error {
	url := fmt.Sprintf("%s?fips=%s&stream=true", apiUrl, fips)
	return nsiApiStream(url, nsp)
}
func GetByBboxStream(bbox string, nsp NsiStreamProcessor) error {
	url := fmt.Sprintf("%s?bbox=%s&stream=true", apiUrl, bbox)
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

	m := consequences.OccupancyTypeMap()
	defaultOcctype := m["RES1-1SNB"]
	var occtype = defaultOcctype
	for {
		var n NsiFeature
		if err := dec.Decode(&n); err == io.EOF {
			break
		} else if err != nil {
			fmt.Printf("Error unmarshalling JSON record: %s.  Stopping Compute.\n", err)
			return err
		}
		if ot, ok := m[n.Properties.Occtype]; ok {
			occtype = ot
		} else {
			occtype = defaultOcctype
			msg := "Using default " + n.Properties.Occtype + " not found"
			return errors.New(msg)
		}
		nsp(consequences.StructureStochastic{
			Name:      n.Properties.Name,
			OccType:   occtype,
			DamCat:    n.Properties.DamCat,
			StructVal: consequences.ParameterValue{Value: n.Properties.StructVal},
			ContVal:   consequences.ParameterValue{Value: n.Properties.ContVal},
			FoundHt:   consequences.ParameterValue{Value: n.Properties.FoundHt},
			X:         n.Properties.X,
			Y:         n.Properties.Y,
		})
	}
	return nil
}
