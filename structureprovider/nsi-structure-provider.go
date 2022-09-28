package structureprovider

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/USACE/go-consequences/consequences"
	"github.com/USACE/go-consequences/geography"
	"github.com/USACE/go-consequences/structures"
)

//NsiProperties is a reflection of the JSON feature property attributes from the NSI-API
type NsiProperties struct {
	Name       int     `json:"fd_id"`
	X          float64 `json:"x"`
	Y          float64 `json:"y"`
	Occtype    string  `json:"occtype"`
	FoundHt    float64 `json:"found_ht"`
	FoundType  string  `json:"found_type"`
	DamCat     string  `json:"st_damcat"`
	StructVal  float64 `json:"val_struct"`
	ContVal    float64 `json:"val_cont"`
	CB         string  `json:"cbfips"`
	Pop2amu65  int32   `json:"pop2amu65"`
	Pop2amo65  int32   `json:"pop2amo65"`
	Pop2pmu65  int32   `json:"pop2pmu65"`
	Pop2pmo65  int32   `json:"pop2pmo65"`
	NumStories int32   `json:"num_story"`
}

//NsiFeature is a feature which contains the properties of a structure from the NSI API
type NsiFeature struct {
	Properties NsiProperties `json:"properties"`
}

type nsiStreamProvider struct {
	ApiURL          string
	OccTypeProvider structures.OccupancyTypeProvider
}

func InitNSISP() nsiStreamProvider {
	//this will only work with go1.16+
	otp := structures.JsonOccupancyTypeProvider{}
	otp.InitDefault()
	//url := urlFinder()
	url := "https://nsi.sec.usace.army.mil/nsiapi/structures"
	// TODO probably don't hard code a possibly changing url
	return nsiStreamProvider{ApiURL: url, OccTypeProvider: otp}
}
func InitNSISPwithOcctypeFilePath(occtypefp string) nsiStreamProvider {
	//this will only work with go1.16+
	otp := structures.JsonOccupancyTypeProvider{}
	otp.InitLocalPath(occtypefp)
	//url := urlFinder()
	url := "https://nsi.sec.usace.army.mil/nsiapi/structures"
	// TODO probably don't hard code a possibly changing url
	return nsiStreamProvider{ApiURL: url, OccTypeProvider: otp}
}
func urlFinder() string {
	url := "https://www.hec.usace.army.mil/fwlink/?linkid=1&type=string"
	transCfg := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, // accept untrusted servers
	}
	client := &http.Client{Transport: transCfg}

	response, err := client.Get(url)

	if err != nil {
		fmt.Println(err)
	}
	defer response.Body.Close()
	rootBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println(err)
	}
	return string(rootBytes)
}
func (nsp nsiStreamProvider) ByFips(fipscode string, sp consequences.StreamProcessor) {
	url := fmt.Sprintf("%s?fips=%s&fmt=fs", nsp.ApiURL, fipscode)
	nsp.nsiStructureStream(url, sp)
}
func (nsp nsiStreamProvider) ByBbox(bbox geography.BBox, sp consequences.StreamProcessor) {
	url := fmt.Sprintf("%s?bbox=%s&fmt=fs", nsp.ApiURL, bbox.ToString())
	nsp.nsiStructureStream(url, sp)
}
func (nsp nsiStreamProvider) ByJsonPost(jsonbody string, sp consequences.StreamProcessor) {
	body := []byte(jsonbody)
	url := fmt.Sprintf("%s?fmt=fs", nsp.ApiURL)
	nsp.nsiPostStructureStream(url, bytes.NewBuffer(body), sp)
}
func (nsp nsiStreamProvider) nsiStructureStream(url string, sp consequences.StreamProcessor) {
	m := nsp.OccTypeProvider.OccupancyTypeMap()
	fmt.Println(url)
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
	//b, err := ioutil.ReadAll(response.Body)
	//fmt.Println(string(b))
	for {
		var n NsiFeature
		if err := dec.Decode(&n); err == io.EOF {
			break
		} else if err != nil {
			fmt.Printf("Error unmarshalling JSON record: %s.  Stopping Compute.\n", err)
			if err == io.ErrUnexpectedEOF {
				break
			}
		}
		sp(NsiFeaturetoStructure(n, m, defaultOcctype))
	}
}
func (nsp nsiStreamProvider) nsiPostStructureStream(url string, body io.Reader, sp consequences.StreamProcessor) {
	m := nsp.OccTypeProvider.OccupancyTypeMap()
	//define a default occtype in case of emergancy
	defaultOcctype := m["RES1-1SNB"]
	transCfg := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, // accept untrusted servers
	}
	client := &http.Client{Transport: transCfg}

	response, err := client.Post(url, "application/json", body)

	if err != nil {
		fmt.Println(err)
	}
	defer response.Body.Close()
	dec := json.NewDecoder(response.Body)
	//b, err := ioutil.ReadAll(response.Body)
	//fmt.Println(string(b))
	for {
		var n NsiFeature
		if err := dec.Decode(&n); err == io.EOF {
			break
		} else if err != nil {
			fmt.Printf("Error unmarshalling JSON record: %s.  Stopping Compute.\n", err)
			if err == io.ErrUnexpectedEOF {
				break
			}
		}
		sp(NsiFeaturetoStructure(n, m, defaultOcctype))
	}
}

//NsiFeaturetoStructure converts an nsi.NsiFeature to a structures.Structure
func NsiFeaturetoStructure(f NsiFeature, m map[string]structures.OccupancyTypeStochastic, defaultOcctype structures.OccupancyTypeStochastic) structures.StructureStochastic {
	var occtype = defaultOcctype
	if otf, okf := m[f.Properties.Occtype+"-"+f.Properties.FoundType]; okf {
		occtype = otf
	} else {
		if ot, ok := m[f.Properties.Occtype]; ok {
			occtype = ot
		} else {
			occtype = defaultOcctype
			msg := "Using default " + f.Properties.Occtype + " not found"
			fmt.Print(msg) //panic(msg)
		}
	}
	return structures.StructureStochastic{
		OccType:    occtype,
		StructVal:  consequences.ParameterValue{Value: f.Properties.StructVal},
		ContVal:    consequences.ParameterValue{Value: f.Properties.ContVal},
		FoundHt:    consequences.ParameterValue{Value: f.Properties.FoundHt},
		FoundType:  f.Properties.FoundType,
		Pop2pmo65:  f.Properties.Pop2pmo65,
		Pop2pmu65:  f.Properties.Pop2pmu65,
		Pop2amo65:  f.Properties.Pop2amo65,
		Pop2amu65:  f.Properties.Pop2amu65,
		NumStories: f.Properties.NumStories,
		BaseStructure: structures.BaseStructure{
			Name:   strconv.Itoa(f.Properties.Name),
			CBFips: f.Properties.CB,
			DamCat: f.Properties.DamCat,
			X:      f.Properties.X,
			Y:      f.Properties.Y,
		},
	}
}

//NsiStreamProcessor is a function used to process an in memory NsiFeature through the NsiStreaming service endpoints
type NsiStreamProcessor func(str NsiFeature)
