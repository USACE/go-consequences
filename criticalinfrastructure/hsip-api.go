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
	Hospitals                Layer = iota //https://services1.arcgis.com/Hp6G80Pky0om7QvQ/arcgis/rest/services/Hospital/FeatureServer/0/query?outFields=*&where=1%3D1&f=geojson
	PowerPlants                           //https://services1.arcgis.com/Hp6G80Pky0om7QvQ/arcgis/rest/services/Plants_gdb/FeatureServer/0/query?outFields=*&where=1%3D1&f=geojson
	FireStations                          //Fire_Station
	WasteWater                            //Wastewater
	LawEnforcement                        //Local_Law_Enforcement_Locations
	EmergencyMedicalServices              //Emergency_Medical_Service_(EMS)_Stations_gdb
	AmtrakStations                        //https://geo.dot.gov/server/rest/services/Hosted/Amtrak_Stations_DS/FeatureServer/0/query?outFields=*&where=1%3D1&f=geojson

	// Beginning of HIFLD additions

	// services1.arcgis.com HIFLD additions
	BRSandEBSTransmitters                  //https://services1.arcgis.com/Hp6G80Pky0om7QvQ/arcgis/rest/services/Broadband_Radio_Service_(BRS)_and_Educational_Broadband_Service_(EBS)_Transmitters/FeatureServer/0/query?outFields=*&where=1%3D1&f=geojson Broadband Radio Service (BRS) and Educational Broadband Service (EBS) Transmitters
	CellularTowers                         //https://services1.arcgis.com/Hp6G80Pky0om7QvQ/arcgis/rest/services/Cellular_Towers_New/FeatureServer/0/query?outFields=*&where=1%3D1&f=geojson
	DialysisCenters                        //https://services1.arcgis.com/Hp6G80Pky0om7QvQ/arcgis/rest/services/Dialysis_Centers/FeatureServer/0/query?outFields=*&where=1%3D1&f=geojson
	EPAandFRSPowerPlants                   //https://services1.arcgis.com/Hp6G80Pky0om7QvQ/arcgis/rest/services/Environmental_Protection_Agency_EPA_Facility_Registry_Service_FRS_Power_Plants/FeatureServer/0/query?outFields=*&where=1%3D1&f=geojson
	FacilityInterests                      //https://services1.arcgis.com/Hp6G80Pky0om7QvQ/arcgis/rest/services/Facility_Interest/FeatureServer/0/query?outFields=*&where=1%3D1&f=geojson
	GeneratingUnits                        //https://services1.arcgis.com/Hp6G80Pky0om7QvQ/arcgis/rest/services/GeneratingUnits1/FeatureServer/0/query?outFields=*&where=1%3D1&f=geojson
	HurricaneEvacuationRoutes              //https://services1.arcgis.com/Hp6G80Pky0om7QvQ/arcgis/rest/services/Hurricane_Evacuation_Routes/FeatureServer/0/query?outFields=*&where=1%3D1&f=geojson
	LandMobileBroadcastTowers              //https://services1.arcgis.com/Hp6G80Pky0om7QvQ/arcgis/rest/services/Land_Mobile_Broadcast_Towers/FeatureServer/0/query?outFields=*&where=1%3D1&f=geojson
	LandMobileCommercialTransmissionTowers //https://services1.arcgis.com/Hp6G80Pky0om7QvQ/arcgis/rest/services/Land_Mobile_Commercial_Towers/FeatureServer/0/query?outFields=*&where=1%3D1&f=geojson
	LocalEmergencyOperationsCenterEOC      //https://services1.arcgis.com/Hp6G80Pky0om7QvQ/arcgis/rest/services/Land_Mobile_Commercial_Towers/FeatureServer/0/query?outFields=*&where=1%3D1&f=geojson
	LocalLawEnforcementLocations           //https://services1.arcgis.com/Hp6G80Pky0om7QvQ/arcgis/rest/services/Local_Law_Enforcement_Locations/FeatureServer/0/query?outFields=*&where=1%3D1&f=geojson
	MicrowaveServiceTowers                 //https://services1.arcgis.com/Hp6G80Pky0om7QvQ/arcgis/rest/services/Microwave_Service_Towers_New/FeatureServer/0/query?outFields=*&where=1%3D1&f=geojson
	NursingHomes                           //https://services1.arcgis.com/Hp6G80Pky0om7QvQ/arcgis/rest/services/NursingHomes/FeatureServer/0/query?outFields=*&where=1%3D1&f=geojson
	PagingTransmissionTowers               //https://services1.arcgis.com/Hp6G80Pky0om7QvQ/arcgis/rest/services/Paging_Transmission_Towers/FeatureServer/0/query?outFields=*&where=1%3D1&f=geojson
	Pharmacies                             //https://services1.arcgis.com/Hp6G80Pky0om7QvQ/arcgis/rest/services/Paging_Transmission_Towers/FeatureServer/0/query?outFields=*&where=1%3D1&f=geojson
	PSAP911ServiceAreaBoundaries           //https://services1.arcgis.com/Hp6G80Pky0om7QvQ/arcgis/rest/services/PSAP_911_Service_Area_Boundaries/FeatureServer/0/query?outFields=*&where=1%3D1&f=geojson
	PublicHealthDepartments                //https://services1.arcgis.com/Hp6G80Pky0om7QvQ/arcgis/rest/services/Public_Health_Departments/FeatureServer/0/query?outFields=*&where=1%3D1&f=geojson
	PublicRefrigeratedWarehouses           //https://services1.arcgis.com/Hp6G80Pky0om7QvQ/arcgis/rest/services/Public_Refrigerated_Warehouses/FeatureServer/0/query?outFields=*&where=1%3D1&f=geojson
	UrgentCareFacilities                   //https://services1.arcgis.com/Hp6G80Pky0om7QvQ/arcgis/rest/services/Urgent_Care_Facilities/FeatureServer/0/query?outFields=*&where=1%3D1&f=geojson
	VeteransHealthAdministrationFacilities //https://services1.arcgis.com/Hp6G80Pky0om7QvQ/arcgis/rest/services/Veterans_Health_Administration_Medical_Facilities/FeatureServer/0/query?outFields=*&where=1%3D1&f=geojson

	// The following root URLs are not the same as the main URL, kept for future reference
	// services.arcgis.com
	// ienccloud
	// geo.dot.gov
	// gis.fema.gov
	// carto.nationalmap.gov

	// // services.arcgis.com HIFLD additions
	// AviationFacilities //https://services.arcgis.com/xOi1kZaI0eWDREZv/ArcGIS/rest/services/Aviation_Facilities/FeatureServer/0/query?outFields=*&where=1%3D1&f=geojson

	// // ienccloud HIFLD additions
	// AiportsUsaceIenc //https://ienccloud.us/arcgis/rest/services/IENC_Feature_Classes/AIRPORT_AREA/MapServer/0/query?outFields=*&where=1%3D1&f=geojson
	// RoadsUSACEIEN    //https://ienccloud.us/arcgis/rest/services/IENC_Feature_Classes/ROADWAY_LINE/MapServer/0/query?outFields=*&where=1%3D1&f=geojson

	// // geo.dot.gov HIFLD additions
	// FerryTerminals                                 //https://geo.dot.gov/server/rest/services/Hosted/Ferry_Terminals_DS/FeatureServer/0/query?outFields=*&where=1%3D1&f=geojson
	// IntermodalFreightFacilitiesAirtoTruck          //https://geo.dot.gov/server/rest/services/Hosted/Intermodal_Freight_Facilities_Air_to_Truck_DS/FeatureServer/0/query?outFields=*&where=1%3D1&f=geojson
	// IntermodalFreightFacilitiesMarineRollonRolloff //https://geo.dot.gov/server/rest/services/Hosted/Intermodal_Freight_Facilities_Marine_Roll_on_Roll_off_DS/FeatureServer/0/query?outFields=*&where=1%3D1&f=geojson
	// IntermodalFreightFacilitiesRail_TOFC_COFC      //https://geo.dot.gov/server/rest/services/Hosted/Intermodal_Freight_Facilities_Rail_TOFC_COFC_DS/FeatureServer/0/query?outFields=*&where=1%3D1&f=geojson
	// RailroadBridges                                //https://geo.dot.gov/server/rest/services/Hosted/Railroad_Grade_Crossings_DS/FeatureServer/0/query?outFields=*&where=1%3D1&f=geojson
	// RailroadGradeCrossings                         //https://geo.dot.gov/server/rest/services/Hosted/Railroad_Grade_Crossings_DS/FeatureServer/0/query?outFields=*&where=1%3D1&f=geojson

	// // gis.fema.gov HIFLD additions
	// NationalShelterSystemFacilities //https://gis.fema.gov/arcgis/rest/services/NSS/FEMA_NSS/FeatureServer/5/query?outFields=*&where=1%3D1&f=geojson

	// // carto.nationalmap.gov
	// EmergencyMedicalServiceStations //https://carto.nationalmap.gov/arcgis/rest/services/structures/MapServer/51/query?outFields=*&where=1%3D1&f=geojson

	// // The below layers do not return a valud URL nor GeoJSON file when searching in the HIFLD database
	// // Not found or 404 Errors
	// AMTransmssionTowers                    //404 Error
	// ER_TSCA_Facilities                     //404 Error
	// EOA_FRS_WasterwaterTreatmentsPlants    //404 Error
	// FireTerminals                          //404 Error
	// FMTransmissionTowers                   //404 Error
	// RoadandRailroadTunnels                 //404 Error
	// RoadTunnels                            //404 Error
	// RoutesandStations                      //404 Error
	// SecondaryRoads578K                     //Cannot get GeoJSON
	// SecondaryRoadsinterstatesandUSHighways //Cannot get GeoJSON
	// StateEmergencyOperationsCenters_EOC    //Error in pulling data
	// StationsandTransfers                   //404 Error
	// TVDigitalStationTransmitters           //404 Error

)

func (l Layer) String() string {
	return [...]string{
		"Hospital",
		"Plants_gdb",
		"Fire_Station",
		"WasteWater",
		"Local_Law_Enforcement_Locations",
		"Emergency_Medical_Service_(EMS)_Stations_gdb",
		"Amtrak_Stations_DS",

		// Beginning of HIFLD additions

		// services1.arcgis.com HIFLD additions
		"Broadband_Radio_Service_(BRS)_and_Educational_Broadband_Service_(EBS)_Transmitters",
		"Cellular_Towers",
		"Dialysis_Centers",
		"Environmental_Protection_Agency_(EPA)_Facility_Registry_Service_(FRS)_Power_Plants",
		"Facility_Interests",
		"Generating_Units",
		// "Hospitals",
		"Hurricane_Evacuation_Routes",
		"Land_Mobile_Broadcast_Towers",
		"Land_Mobile_Commercial_Transmission_Towers",
		"Local_Emergency_Operations_Center_(EOC)",
		"Local_Law_Enforcement_Locations",
		"Microwave_Service_Towers",
		"Nursing_Homes",
		"Paging_Transmission_Towers",
		"Pharmacies",
		"PSAP_911_Service_Area_Boundaries",
		"Public_Health_Departments",
		"Public_Refrigerated_Warehouses",
		"Urgent_Care_Facilities",
		"Veterans_Health_Administration_Facilities",

		// The following root URLs are not the same as the main URL, kept for future reference
		// services.arcgis.com
		// ienccloud
		// geo.dot.gov
		// gis.fema.gov
		// carto.nationalmap.gov

		// // services.arcgis.com HIFLD additions
		// "Aviation_facilities",

		// // ienccloud HIFLD additions
		// "Airports_USACE_IENC",
		// "Roads_(USACE_IENC)",

		// // geo.dot.gov HIFLD additions
		// "Ferry_Terminals",
		// "Intermodal_Freight_Facilities_Air_to_Truck",
		// "Intermodal_Freight_Facilities_Marine_Roll_on_Roll_off",
		// "Intermodal_Freight_Facilities_Rail_TOFC_COFC",
		// "Railroad_Bridges",
		// "Railroad_Grade_Crossings",

		// // gis.fema.gov HIFLD additions
		// "National_Shelter_System_Facilities",

		// // carto.nationalmap.gov additions
		// "Emergency_Medical_Service_(EMS)_Stations",

		// // The below layers do not return a valud URL nor GeoJSON file when searching in the HIFLD database
		// // Not found or 404 Errors
		// "AM_Transmssion_Towers",
		// "EPA_Emergency_Response_(ER)_Toxic_Substances_Control_Act_(TSCA)_Facilities",
		// "EPA_Facility_Registry_Service_(FRS)_Wastewater_Treatment_Plants",
		// "Fire_Stations",
		// "FM_Transmission_Towers",
		// "Road_and_Railroad_Tunnels",
		// "Road_Tunnels",
		// "Routes_and_Stations",
		// "Secondary_Roads_578K",
		// "Secondary_Roads_interstates_and_US_Highways",
		// "State_Emergency_Operations_Centers_(EOC)",
		// "Stations_and_Transfers",
		// "TV_Digital_Station_Transmitters",
	}[l]
}

func (l Layer) OccupancyType() string {
	return [...]string{
		"Hospital",
		"Power Plant",
		"Fire Station",
		"Waste Water Treatment Plant",
		"Local Law Enforcement",
		"Emergency Medical Service Station",
		"Amtrak Stations",

		// Beginning of HIFLD additions

		// services1.arcgis.com HIFLD additions
		"Broadband Radio Service (BRS) and Educational Broadband Service (EBS) Transmitters",
		"Cellular Towers",
		"Dialysis Centers",
		"Environmental Protection Agency (EPA) Facility Registry Service (FRS) Power Plants",
		"Facility Interests",
		"Generating Units",
		"Hurricane Evacuation Routes",
		"Land Mobile Broadcast Towers",
		"Land Mobile Commercial Transmission Towers",
		"Local Emergency Operations Center (EOC)",
		"Local Law Enforcement Locations",
		"Microwave Service Towers",
		"Nursing Homes",
		"Paging Transmission Towers",
		"Pharmacies",
		"PSAP 911 Service Area Boundaries",
		"Public Health Departments",
		"Public Refrigerated Warehouses",
		"Urgent Care Facilities",
		"Veterans Health Administration Facilities",

		// The following root URLs are not the same as the main URL, kept for future reference
		// services.arcgis.com
		// ienccloud
		// geo.dot.gov
		// gis.fema.gov
		// carto.nationalmap.gov

		// // services.arcgis.com HIFLD additions
		// "Aviation facilities",

		// // ienccloud HIFLD additions
		// "Airports USACE IENC",
		// "Roads (USACE IENC)",

		// // geo.dot.gov HIFLD additions
		// "Ferry Terminals",
		// "Intermodal Freight Facilities Air to Truck",
		// "Intermodal Freight Facilities Marine Roll on Roll off",
		// "Intermodal Freight Facilities Rail TOFC/ COFC",
		// "Railroad Bridges",
		// "Railroad Grade Crossings",

		// // gis.fema.gov HIFLD additions
		// "National Shelter System Facilities",

		// // carto.nationalmap.gov additions
		// "Emergency_Medical_Service_(EMS)_Stations",

		// // The below layers do not return a valud URL nor GeoJSON file when searching in the HIFLD database
		// // Not found or 404 Errors
		// "AM Transmssion Towers",
		// "EPA Emergency Response (ER) Toxic Substances Control Act (TSCA) Facilities",
		// "EPA Facility Registry Service (FRS) Wastewater Treatment Plants",
		// "Fire Stations",
		// "FM Transmission Towers",
		// "Road and Railroad Tunnels",
		// "Road Tunnels",
		// "Routes and Stations",
		// "Secondary Roads 578K",
		// "Secondary Roads interstates and US Highways",
		// "State Emergency Operations Centers (EOC)",
		// "Stations and Transfers",
		// "TV Digital Station Transmitters",
	}[l]
}

// TODO: What does it sit under
func (l Layer) DamageCategory() string {
	return [...]string{
		"Health and Medical",
		"Energy",
		"Safety and Security",
		"Water Systems",
		"Safety and Security",
		"Health and Medical",
		"Transportation",

		// Beginning of HIFLD additions
		// Based on the FEMA Lifeline Designation and the naming of past categories

		// services1.arcgis.com HIFLD additions
		"Communications",
		"Communications",
		"Health & Medical",
		"Energy",
		"Energy",
		"Energy",
		"Transportation",
		"Communications",
		"Communications",
		"Safety & Security",
		"Safety & Security",
		"Communications",
		"Health & Medical",
		"Communications",
		"Health & Medical",
		"Safety & Security",
		"Health & Medical",
		"Food, Hydration, Shelter",
		"Health & Medical",
		"Health & Medical",

		// The following root URLs are not the same as the main URL, kept for future reference
		// services.arcgis.com
		// ienccloud
		// geo.dot.gov
		// gis.fema.gov
		// carto.nationalmap.gov

		// // services.arcgis.com HIFLD additions
		// "Transportation",

		// // ienccloud HIFLD additions
		// "Transportation",
		// "Transportation",

		// // geo.dot.gov HIFLD additions
		// "Transportation",
		// "Transportation",
		// "Transportation",
		// "Transportation",
		// "Transportation",
		// "Transportation",

		// // gis.fema.gov HIFLD additions
		// "Food, Hydration, Shelter",

		// // carto.nationalmap.gov additions
		// "Health & Medical",

		// // The below layers do not return a valud URL nor GeoJSON file when searching in the HIFLD database
		// // Not found or 404 Errors
		// "Communications",
		// "Hazardous Materials",
		// "Water Systems",
		// "Safety & Security",
		// "Communications",
		// "Transportation",
		// "Transportation",
		// "Transportation",
		// "Transportation",
		// "Transportation",
		// "Safety & Security",
		// "Transportation",
		// "Communications",
	}[l]
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
