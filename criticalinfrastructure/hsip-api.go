package criticalinfrastructure

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/USACE/go-consequences/consequences"
	"github.com/USACE/go-consequences/geography"
	"github.com/USACE/go-consequences/hazards"
)

var hsip_root string = "https://services1.arcgis.com/Hp6G80Pky0om7QvQ/arcgis/rest/services/"
var esri_bbox string = "&geometryType=esriGeometryEnvelope&geometry="
var hsip_suffix string = "&outSR=4326&f=geojson"

// https://services1.arcgis.com/Hp6G80Pky0om7QvQ/arcgis/rest/services/Emergency_Medical_Service_(EMS)_Stations_gdb/FeatureServer/0/query?outFields=*&where=1%3D1&f=geojson
type Layer int

// Parameter types describe different parameters for hazards
const (
	Hospitals Layer = iota
	PowerPlants
	FireStations             //Fire_Station
	WasteWater               //Wastewater
	LawEnforcement           //Local_Law_Enforcement_Locations
	EmergencyMedicalServices //Emergency_Medical_Service_(EMS)_Stations_gdb
	AmtrakStations           //https://geo.dot.gov/server/rest/services/Hosted/Amtrak_Stations_DS/FeatureServer/0/query?outFields=*&where=1%3D1&f=geojson
)

func (l Layer) String() string {
	return [...]string{"Hospital", "Plants_gdb", "Fire_Station", "WasteWater", "Local_Law_Enforcement_Locations", "Emergency_Medical_Service_(EMS)_Stations_gdb", "Amtrak_Stations_DS"}[l]
}
func (l Layer) OccupancyType() string {
	return [...]string{"Hospital", "Power Plant", "Fire Station", "Waste Water Treatment Plant", "Local Law Enforcement", "Emergency Medical Service Station", "Amtrak Stations"}[l]
}
func (l Layer) DamageCategory() string {
	return [...]string{"Health and Medical", "Energy", "Safety and Security", "Water Systems", "Safety and Security", "Health and Medical", "Transportation"}[l]
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
	Attributes CriticalInfrastructureAttributes `json:"properties"`
	Geometry   geography.GeoJsonGeometry        `json:"geometry"`
}
type CriticalInfrastructureAttributes struct {
	Name string `json:"NAME"`
	//X              float64 `json:"LONGITUDE"`
	//Y              float64 `json:"LATITUDE"`
	DamageCategory string
	OccupancyType  string `json:"NAICS_DESC"`
}

var header = []string{"name", "x", "y", "Lifeline", "Dataset", "hazard"}

func (c CriticalInfrastructureFeature) Compute(h hazards.HazardEvent) (consequences.Result, error) {
	location := c.Geometry.ToLocation()
	results := []interface{}{c.Attributes.Name, location.X, location.Y, c.Attributes.DamageCategory, c.Attributes.OccupancyType, h}
	var ret = consequences.Result{Headers: header, Result: results}
	return ret, nil
}
func (c CriticalInfrastructureFeature) Location() geography.Location {
	location := c.Geometry.ToLocation()
	location.SRID = "4326"
	return location
}

func (h HsipProvider) ByBbox(bbox geography.BBox, sp consequences.StreamProcessor) {
	bbstring := fmt.Sprintf("%s%v,%v,%v,%v", esri_bbox, bbox.Bbox[0], bbox.Bbox[3], bbox.Bbox[2], bbox.Bbox[1])
	for _, l := range h.FilterList {
		queryString := fmt.Sprintf("%s%s/FeatureServer/0/query?outFields=*%s%s", hsip_root, l, bbstring, hsip_suffix)
		fmt.Println(queryString)
		processQuery(queryString, sp, l)
	}
}
func (h HsipProvider) ByFips(fipscode string, sp consequences.StreamProcessor) {
	log.Fatal("fips query not implemented for hsip provider")
}
func processQuery(url string, sp consequences.StreamProcessor, l Layer) {
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
	damcat := l.DamageCategory()
	occtype := l.OccupancyType()
	for _, cielement := range ci.Features {
		cielement.Attributes.DamageCategory = damcat
		cielement.Attributes.OccupancyType = occtype
		sp(cielement)
	}
}
