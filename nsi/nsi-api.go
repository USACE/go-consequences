package nsi

import (
	"crypto/tls"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
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

//SQLDataSet is a simple struct to store a sql dataset
type SQLDataSet struct {
	db *sql.DB
}

//OpenSQLNSIDataSet opens a sqldataset with the NSI data
func OpenSQLNSIDataSet(nsiLoc string) SQLDataSet {
	db, _ := sql.Open("sqlite3", nsiLoc)
	db.SetMaxOpenConns(1)
	return SQLDataSet{db: db}
}

var apiURL string = "https://nsi-dev.sec.usace.army.mil/nsiapi/structures" //this will only work behind the USACE firewall -
var nsiLoc string = "./nsiv2_29.gpkg?cache=shared&mode=rwc"                //this targets the location of the NSI - maybe get some way to prompt the user for this... cache=shared comes from https://github.com/mattn/go-sqlite3/issues/274 but unsure if it does anything

//GetByFips returns an NsiInventory for a FIPS code
func GetByFips(fips string) NsiInventory {
	url := fmt.Sprintf("%s?fips=%s&fmt=fa", apiURL, fips)
	return nsiAPI(url)
	// I haven't been able to test the commented out feature below
	// It should do the same thing as GetByFipsStream, but I am not sure if the error checking condition works

	// inv := nsiAPI(url)
	// if len(inv.Features) != 0 {
	// 	return inv
	// }

	// nsi := OpenSQLNSIDataSet(nsiLoc)
	// rows, err1 := nsi.db.Query("SELECT fd_id, x, y, cbfips, occtype, found_ht, found_type, st_damcat, val_struct, val_cont, pop2amu65, pop2amo65, pop2pmu65, pop2pmo65 FROM nsi WHERE cbfips LIKE '" + fips + "'%")
	// if err1 != nil {
	// 	log.Fatal(err1)
	// }
	// defer rows.Close()

	// var inventory NsiInventory
	// for rows.Next() { // Iterate and fetch the records from result cursor
	// 	feature := NsiFeature{}
	// 	err2 := rows.Scan(&feature.Properties.Name, &feature.Properties.X, &feature.Properties.Y, &feature.Properties.CB, &feature.Properties.Occtype, &feature.Properties.FoundHt, &feature.Properties.FoundType, &feature.Properties.DamCat, &feature.Properties.StructVal, &feature.Properties.ContVal, &feature.Properties.Pop2amu65, &feature.Properties.Pop2amo65, &feature.Properties.Pop2pmu65, &feature.Properties.Pop2pmo65)
	// 	if err2 != nil {
	// 		panic(err2)
	// 	}
	// 	inventory.Features = append(inventory.Features, feature)
	// }
	// return inventory
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
var nsiError error = errors.New("Not connected to USACE Firewall")

//GetByFipsStream a streaming service for NsiFeature based on a FIPs code
func GetByFipsStream(fips string, nsp NsiStreamProcessor) error {
	url := fmt.Sprintf("%s?fips=%s&fmt=fs", apiURL, fips)

	var curErr error = nsiAPIStream(url, nsp)
	// if we are behind the USACE Firewall, we go here
	if curErr.Error() != nsiError.Error() {
		return nsiAPIStream(url, nsp)
	}
	// if we are not behind the USACE Firewall, we access a local NSI database
	nsi := OpenSQLNSIDataSet(nsiLoc)
	rows, err1 := nsi.db.Query("SELECT fd_id, x, y, cbfips, occtype, found_ht, found_type, st_damcat, val_struct, val_cont, pop2amu65, pop2amo65, pop2pmu65, pop2pmo65 FROM nsi WHERE cbfips LIKE '" + fips + "%'")

	if err1 != nil {
		log.Fatal(err1)
	}
	defer rows.Close()

	for rows.Next() { // Iterate and fetch the records from result cursor
		feature := NsiFeature{}
		err2 := rows.Scan(&feature.Properties.Name, &feature.Properties.X, &feature.Properties.Y, &feature.Properties.CB, &feature.Properties.Occtype, &feature.Properties.FoundHt, &feature.Properties.FoundType, &feature.Properties.DamCat, &feature.Properties.StructVal, &feature.Properties.ContVal, &feature.Properties.Pop2amu65, &feature.Properties.Pop2amo65, &feature.Properties.Pop2pmu65, &feature.Properties.Pop2pmo65)
		if err2 != nil {
			panic(err2)
		}
		nsp(feature)
	}
	return nil
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
		// creating a new error format to return/match an error variable in the GetByFipsStream
		err1 := errors.New("Not connected to USACE Firewall")
		fmt.Println(err)
		return err1
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
