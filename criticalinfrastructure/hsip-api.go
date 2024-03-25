package criticalinfrastructure

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/USACE/go-consequences/consequences"
	"github.com/USACE/go-consequences/geography"
	"github.com/USACE/go-consequences/hazards"
)

var hsip_root string = "https://services1.arcgis.com/Hp6G80Pky0om7QvQ/arcgis/rest/services/"
var esri_bbox string = "&geometryType=esriGeometryEnvelope&geometry="
var hsip_suffix string = "&outSR=4326&f=json"

type Layer int

// Parameter types describe different parameters for hazards
const (
	Hospitals Layer = iota
)

func (l Layer) String() string {
	return [...]string{"Hospital"}[l]
}

type HsipProvider struct {
	FilterList []Layer
}

func InitHsipProvider(list []Layer) HsipProvider {
	return HsipProvider{
		FilterList: list,
	}
}

type CriticalInfrastructureReturn struct {
	Features []CriticalInfrastructureFeature `json:"features"`
}
type CriticalInfrastructureFeature struct {
	Attributes CriticalInfrastructureAttributes `json:"attributes"`
}
type CriticalInfrastructureAttributes struct {
	Name           string  `json:"NAME"`
	X              float64 `json:"LATITUDE"`
	Y              float64 `json:"LONGITUDE"`
	DamageCategory string
	OccupancyType  string `json:"NAICS_DESC"`
}

func (c CriticalInfrastructureFeature) Compute(h hazards.HazardEvent) (consequences.Result, error) {
	fmt.Println(c)
	return consequences.Result{}, nil
}
func (c CriticalInfrastructureFeature) Location() geography.Location {
	return geography.Location{
		X:    c.Attributes.X,
		Y:    c.Attributes.Y,
		SRID: "4326",
	}
}

func (h HsipProvider) ByBbox(bbox geography.BBox, sp consequences.StreamProcessor) {
	bbstring := fmt.Sprintf("%s%v,%v,%v,%v", esri_bbox, bbox.Bbox[0], bbox.Bbox[1], bbox.Bbox[2], bbox.Bbox[3])
	for _, l := range h.FilterList {
		queryString := fmt.Sprintf("%s%s/FeatureServer/0/query?outFields=*%s%s", hsip_root, l, bbstring, hsip_suffix)
		fmt.Println(queryString)
		processQuery(queryString, sp)
	}
}
func processQuery(url string, sp consequences.StreamProcessor) {
	transCfg := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, // accept untrusted servers
	}
	client := &http.Client{Transport: transCfg}

	response, err := client.Get(url)

	if err != nil {
		fmt.Println(err)
	}
	defer response.Body.Close()
	//dec := json.NewDecoder(response.Body)
	b, err := io.ReadAll(response.Body)
	//fmt.Println(string(b))
	var ci CriticalInfrastructureReturn
	json.Unmarshal(b, &ci)

	for _, cielement := range ci.Features {
		sp(cielement)
	}
}
