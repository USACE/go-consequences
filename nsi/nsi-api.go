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

func (i NsiInventory) toStructures() []consequences.Structure {
	m := consequences.OccupancyTypeMap()
	defaultOcctype := m["RES1-1SNB"]
	structures := make([]consequences.Structure, len(i.Features))
	for idx, feature := range i.Features {
		structures[idx] = consequences.Structure{
			Name:      feature.Properties.Name,
			OccType:   defaultOcctype,
			DamCat:    feature.Properties.DamCat,
			StructVal: feature.Properties.StructVal,
			ContVal:   feature.Properties.ContVal,
			FoundHt:   feature.Properties.FoundHt,
		}
	}
	return structures
}

var apiUrl string = "https://nsi-dev.sec.usace.army.mil/nsiapi/structures" //this will only work behind the USACE firewall -

func GetByBbox(bbox string) []consequences.Structure {
	structures := make([]consequences.Structure, 0)
	transCfg := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, // accept untrusted servers
	}
	client := &http.Client{Transport: transCfg}
	url := fmt.Sprintf("%s?bbox=%s", apiUrl, bbox)
	response, err := client.Get(url)

	if err != nil {
		log.Fatal(err)
	}

	defer response.Body.Close()
	// UnmarshalJSON implements UnmarshalJSON interface
	jsonData, err := ioutil.ReadAll(response.Body)
	features := make([]NsiFeature, 0)

	if err := json.Unmarshal(jsonData, &features); err != nil {
		fmt.Println(err)
		return structures
	}
	inventory := NsiInventory{Features: features}
	structures = inventory.toStructures()
	return structures

}
