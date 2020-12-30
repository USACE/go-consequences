package crops

import (
	"crypto/tls"
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
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
type XmlStatsURLResponse struct {
	XMLName   xml.Name `xml:"GetCDLStatResponse"`
	ReturnUrl string   `xml:"returnURL"`
}
type XmlFileURLResponse struct {
	XMLName   xml.Name `xml:"GetCDLFileResponse"`
	ReturnUrl string   `xml:"returnURL"`
}
type XmlCDLValueResponse struct {
	XMLName xml.Name `xml:"GetCDLValueResponse"`
	Result  string   `xml:"Result"`
}

var apiStatsUrl string = "http://nassgeodata.gmu.edu/axis2/services/CDLService/"

func GetStatsByBbox(year string, minx string, miny string, maxx string, maxy string) StatisticsResult {
	url := fmt.Sprintf("%sGetCDLStat?year=%s&bbox=%s,%s,%s,%s&format=csv", apiStatsUrl, year, minx, miny, maxx, maxy) //malformed json keys are not quoted.
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
	var resultURL XmlStatsURLResponse
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

func GetCDLValue(year string, x string, y string) string {
	url := fmt.Sprintf("%sGetCDLValue?year=%s&x=%s&y=%s", apiStatsUrl, year, x, y) //malformed json keys are not quoted.
	return nassCDLValueApi(url)
}
func nassCDLValueApi(url string) string {
	transCfg := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, // accept untrusted servers
	}
	client := &http.Client{Transport: transCfg}
	response, err := client.Get(url)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	defer response.Body.Close()
	xmlData, err := ioutil.ReadAll(response.Body)
	var result XmlCDLValueResponse
	if err := xml.Unmarshal(xmlData, &result); err != nil {
		fmt.Println("error unmarshalling NASS XML " + err.Error() + " URL: " + url)
		s := string(xmlData)
		fmt.Println(s)
		fmt.Println("first 100 chars of xmlbody was: " + s[0:100]) //s) //"last 100 chars of jsonbody was: " + s[len(s)-100:])
		return ""
	}
	return result.Result
}
func GetCDLFileByFIPS(year string, fips string) bool {
	url := fmt.Sprintf("%sGetCDLFile?year=%s&fips=%s", apiStatsUrl, year, fips)
	return nassFileApi(url)
}
func nassFileApi(url string) bool {
	transCfg := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, // accept untrusted servers
	}
	client := &http.Client{Transport: transCfg}
	response, err := client.Get(url)
	if err != nil {
		fmt.Println(err)
		return false
	}
	defer response.Body.Close()
	xmlData, err := ioutil.ReadAll(response.Body)
	var resultURL XmlFileURLResponse
	if err := xml.Unmarshal(xmlData, &resultURL); err != nil {
		fmt.Println("error unmarshalling NASS XML " + err.Error() + " URL: " + url)
		s := string(xmlData)
		fmt.Println(s)
		fmt.Println("first 100 chars of xmlbody was: " + s[0:100]) //s) //"last 100 chars of jsonbody was: " + s[len(s)-100:])
		return false
	}
	response2, err2 := client.Get(resultURL.ReturnUrl)
	if err2 != nil {
		fmt.Println(err2)
		return false
	}
	defer response2.Body.Close()
	fileparts := strings.Split(resultURL.ReturnUrl, "/")
	outfile := "C:\\Temp\\agtesting\\" + fileparts[len(fileparts)-1]
	out, err := os.Create(outfile)
	if err != nil {
		fmt.Println(err)
		return false
	}
	defer out.Close()
	io.Copy(out, response2.Body)

	return true
}
