package nsi

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

//NSIStats is a struct describing the NSI Statistics endpoint json body return
type NSIStats struct {
	Count              int     `json:"num_structures"`
	MinYearBuilt       int     `json:"yrbuilt_min"`
	MaxYearBuilt       int     `json:"yrbuilt_max"`
	MedianYearBuilt    int     `json:"med_yr_blt_max"`
	MaxGroundElevation float64 `json:"ground_elv_max"`
	MinGroundElevation float64 `json:"ground_elv_min"`
	SumResUnits        int     `json:"resunits_sum"`
	SumEmp             int     `json:"empnum_sum"`
	SumTeachers        int     `json:"teachers_sum"`
	SumSqft            float64 `json:"sqft_sum"`
	SumPoP2AMO65       int     `json:"pop2amo65_sum"`
	SumPoP2AMU65       int     `json:"pop2amu65_sum"`
	SumPoP2PMO65       int     `json:"pop2pmo65_sum"`
	SumPoP2PMU65       int     `json:"pop2pmu65_sum"`
	SumStudents        int     `json:"students_sum"`
	SumStructVal       float64 `json:"val_struct_sum"`
	SumContVal         float64 `json:"val_cont_sum"`
	SumVehicVal        float64 `json:"val_vehic_sum"`
	MeanNumStory       float64 `json:"num_story_mean"`
	MeanSqft           float64 `json:"sqft_mean"`
}

var apiStatsURL string = "https://nsi-dev.sec.usace.army.mil/nsiapi/stats" //this will only work behind the USACE firewall -
//GetStatsByFips is intended to return NSIStats for a FIPS code however, it currently does not work.
func GetStatsByFips(fips string) NSIStats {
	url := fmt.Sprintf("%s?fips=%s", apiStatsURL, fips)
	return nsiStatsAPI(url)
}

//GetStatsByBbox returns an NSIStats for a bounding box.
func GetStatsByBbox(bbox string) NSIStats {
	url := fmt.Sprintf("%s?bbox=%s", apiStatsURL, bbox)
	return nsiStatsAPI(url)
}
func nsiStatsAPI(url string) NSIStats {
	transCfg := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, // accept untrusted servers
	}
	client := &http.Client{Transport: transCfg}

	response, err := client.Get(url)
	var stats NSIStats
	if err != nil {
		fmt.Println(err)
		return stats
	}
	defer response.Body.Close()
	jsonData, err := ioutil.ReadAll(response.Body)

	if err := json.Unmarshal(jsonData, &stats); err != nil {
		fmt.Println("error unmarshalling NSI json " + err.Error() + " URL: " + url)
		s := string(jsonData)
		fmt.Println("first 100 chars of jsonbody was: " + s[0:100]) //s) //"last 100 chars of jsonbody was: " + s[len(s)-100:])
		return stats
	}
	return stats
}
