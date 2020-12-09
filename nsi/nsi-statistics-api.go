package nsi

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type NSIStats struct {
	Count                 int     `json:"num_structures"`
	Min_Year_Built        int     `json:"yrbuilt_min"`
	Max_Year_Built        int     `json:"yrbuilt_max"`
	Max_Median_Year_Built int     `json:"med_yr_blt_max"`
	Max_Ground_Elevation  float64 `json:"ground_elv_max"`
	Min_Ground_Elevation  float64 `json:"ground_elv_min"`
	Sum_ResUnits          int     `json:"resunits_sum"`
	Sum_Emp               int     `json:"empnum_sum"`
	Sum_Teachers          int     `json:"teachers_sum"`
	Sum_Sqft              float64 `json:"sqft_sum"`
	Sum_PoP_2AMO65        int     `json:"pop2amo65_sum"`
	Sum_PoP_2AMU65        int     `json:"pop2amu65_sum"`
	Sum_PoP_2PMO65        int     `json:"pop2pmo65_sum"`
	Sum_PoP_2PMU65        int     `json:"pop2pmu65_sum"`
	Sum_Students          int     `json:"students_sum"`
	Sum_Struct_Val        float64 `json:"val_struct_sum"`
	Sum_Cont_Val          float64 `json:"val_cont_sum"`
	Sum_Vehic_Val         float64 `json:"val_vehic_sum"`
	Mean_NumStory         float64 `json:"num_story_mean"`
	Mean_Sqft             float64 `json:"sqft_mean"`
}

var apiStatsUrl string = "https://nsi-dev.sec.usace.army.mil/nsiapi/stats" //this will only work behind the USACE firewall -
func GetStatsByFips(fips string) NSIStats {
	url := fmt.Sprintf("%s?fips=%s", apiStatsUrl, fips)
	return nsiStatsApi(url)
}
func GetStatsByBbox(bbox string) NSIStats {
	url := fmt.Sprintf("%s?bbox=%s", apiStatsUrl, bbox)
	return nsiStatsApi(url)
}
func nsiStatsApi(url string) NSIStats {
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
