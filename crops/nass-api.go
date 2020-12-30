package crops

import (
	"crypto/tls"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

type StatisticsRow struct {
	Value    int     `json:"value"`
	Count    int     `json:"count"`
	Category string  `json:"category"`
	Color    string  `json:"color"`
	Acreage  float64 `json:"acreage"`
}
type StatisticsResult struct {
	Success      bool            `json:"success"`
	ErrorMessage string          `json:"errorMessage"`
	Rows         []StatisticsRow `json:"rows"`
}
type XmlResponse struct {
	XMLName   xml.Name `xml:"GetCDLStatResponse"`
	ReturnUrl string   `xml:"returnURL"`
}

var apiStatsUrl string = "http://nassgeodata.gmu.edu/axis2/services/CDLService/GetCDLStat"

func GetStatsByBbox(year string, minx string, miny string, maxx string, maxy string) StatisticsResult {
	url := fmt.Sprintf("%s?year=%s&bbox=%s,%s,%s,%s&format=csv", apiStatsUrl, year, minx, miny, maxx, maxy)
	return nassStatsApi(url)
}
func nassStatsApi(url string) StatisticsResult {
	transCfg := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, // accept untrusted servers
	}
	client := &http.Client{Transport: transCfg}
	var stats StatisticsResult
	response, err := client.Get(url)
	if err != nil {
		fmt.Println(err)
		return stats
	}
	defer response.Body.Close()
	xmlData, err := ioutil.ReadAll(response.Body)
	var resultURL XmlResponse
	if err := xml.Unmarshal(xmlData, &resultURL); err != nil {
		fmt.Println("error unmarshalling NASS XML " + err.Error() + " URL: " + url)
		s := string(xmlData)
		fmt.Println(s)
		fmt.Println("first 100 chars of xmlbody was: " + s[0:100]) //s) //"last 100 chars of jsonbody was: " + s[len(s)-100:])
		return stats
	}
	response2, err2 := client.Get(resultURL.ReturnUrl)
	if err2 != nil {
		fmt.Println(err2)
		return stats
	}
	defer response2.Body.Close()
	jsonData, err2 := ioutil.ReadAll(response2.Body)
	s := string(jsonData)
	if s == "" {
		stats.Success = false
	} else {
		stats.Success = true
		rows := strings.Split(s, "\r\n")
		for i := 1; i < len(rows); i++ {
			values := strings.Split(rows[i], ", ")
			v, _ := strconv.Atoi(values[0])
			c, _ := strconv.Atoi(values[2])
			a, _ := strconv.ParseFloat(values[3], 64)
			stats.Rows = append(stats.Rows, StatisticsRow{Value: v, Category: values[1], Count: c, Acreage: a})
		}
	}

	return stats
}
