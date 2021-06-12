package crops

//Documentation: https://nassgeodata.gmu.edu/CropScape/devhelp/help.html

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

//StatisticsRow describes a row in the statistics result from the NASS stats endpoint
type StatisticsRow struct {
	Value    int     `json:"value"`
	Count    int     `json:"count"`
	Category string  `json:"category"`
	Color    string  `json:"color"`
	Acreage  float64 `json:"acreage"`
}

//StatisticsResult describes the structure of the result from the NASS stats endpoint
type StatisticsResult struct {
	Success      bool            `json:"success"`
	ErrorMessage string          `json:"errorMessage"`
	Rows         []StatisticsRow `json:"rows"`
}

//XMLStatsURLResponse is the xml return for a stats endpoint query
type XMLStatsURLResponse struct {
	XMLName   xml.Name `xml:"GetCDLStatResponse"`
	ReturnURL string   `xml:"returnURL"`
}

//XMLFileURLResponse is the xml return for the File NASS endpoint
type XMLFileURLResponse struct {
	XMLName   xml.Name `xml:"GetCDLFileResponse"`
	ReturnURL string   `xml:"returnURL"`
}

//XMLExtractResponse is the xml return for the Export NASS endpoint
type XMLExtractResponse struct {
	XMLName   xml.Name `xml:"ExtractCDLByValuesResponse"`
	ReturnURL string   `xml:"returnURL"`
}

//XMLCDLValueResponse is the xml return for a getCDLValue response from the NASS API
type XMLCDLValueResponse struct {
	XMLName xml.Name `xml:"GetCDLValueResponse"`
	Result  string   `xml:"Result"`
}

var apiStatsURL string = "http://nassgeodata.gmu.edu/axis2/services/CDLService/"

//GetStatsByBbox returns the statistics of crops in a bounding box in the projection of USA Contiguous Albers Equal Area Conic (USGS version).
func GetStatsByBbox(year string, minx string, miny string, maxx string, maxy string) StatisticsResult {
	url := fmt.Sprintf("%sGetCDLStat?year=%s&bbox=%s,%s,%s,%s&format=csv", apiStatsURL, year, minx, miny, maxx, maxy) //malformed json keys are not quoted.
	return nassStatsAPI(url)
}
func nassStatsAPI(url string) StatisticsResult {
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
	var resultURL XMLStatsURLResponse
	if err := xml.Unmarshal(xmlData, &resultURL); err != nil {
		fmt.Println("error unmarshalling NASS XML " + err.Error() + " URL: " + url)
		s := string(xmlData)
		fmt.Println(s)
		fmt.Println("first 100 chars of xmlbody was: " + s[0:100]) //s) //"last 100 chars of jsonbody was: " + s[len(s)-100:])
		return stats
	}
	response2, err2 := client.Get(resultURL.ReturnURL)
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

//GetCDLValue returns a crop type for a year and x,y coordinates in the projection of USA Contiguous Albers Equal Area Conic (USGS version).
func GetCDLValue(year string, x string, y string) Crop {
	url := fmt.Sprintf("%sGetCDLValue?year=%s&x=%s&y=%s", apiStatsURL, year, x, y) //malformed json keys are not quoted.
	return nassCDLValueAPI(url)
}
func nassCDLValueAPI(url string) Crop {
	transCfg := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, // accept untrusted servers
	}
	client := &http.Client{Transport: transCfg}
	response, err := client.Get(url)
	if err != nil {
		fmt.Println(err)
		return Crop{}
	}
	defer response.Body.Close()
	xmlData, err := ioutil.ReadAll(response.Body)
	var result XMLCDLValueResponse
	if err := xml.Unmarshal(xmlData, &result); err != nil {
		fmt.Println("error unmarshalling NASS XML " + err.Error() + " URL: " + url)
		s := string(xmlData)
		fmt.Println(s)
		fmt.Println("first 100 chars of xmlbody was: " + s[0:100]) //s) //"last 100 chars of jsonbody was: " + s[len(s)-100:])
		return Crop{}
	}
	nobrackets := strings.Trim(result.Result, "{}")
	kvs := strings.Split(nobrackets, ", ")
	value := strings.Split(kvs[2], ": ")[1]
	category := strings.Trim(strings.Split(kvs[3], ": ")[1], "\"")
	x, _ := strconv.ParseFloat(strings.Split(kvs[0], ": ")[1], 64)
	y, _ := strconv.ParseFloat(strings.Split(kvs[1], ": ")[1], 64)
	v, _ := strconv.Atoi(value)
	c := BuildCrop(byte(v), category)
	c = c.WithLocation(x, y)
	return c
}

//GetCDLFileByFIPS stores a NASS CDL Geotif for a given year and county FIPS
func GetCDLFileByFIPS(year string, fips string) (nassTiffReader, error) {
	url := fmt.Sprintf("%sGetCDLFile?year=%s&fips=%s", apiStatsURL, year, fips)
	return nassFileAPI(url)
}

//GetCDLFileByBbox stores a NASS CDL Geotif for a given year and bounding box
func GetCDLFileByBbox(year string, minx string, miny string, maxx string, maxy string) (nassTiffReader, error) {
	url := fmt.Sprintf("%sGetCDLStat?year=%s&bbox=%s,%s,%s,%s&format=csv", apiStatsURL, year, minx, miny, maxx, maxy)
	return nassFileAPI(url)
}
func nassFileAPI(url string) (nassTiffReader, error) {
	transCfg := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, // accept untrusted servers
	}
	client := &http.Client{Transport: transCfg}
	response, err := client.Get(url)
	if err != nil {
		fmt.Println(err)
		return nassTiffReader{}, err
	}
	defer response.Body.Close()
	xmlData, err := ioutil.ReadAll(response.Body)
	var resultURL XMLFileURLResponse
	if err := xml.Unmarshal(xmlData, &resultURL); err != nil {
		fmt.Println("error unmarshalling NASS XML " + err.Error() + " URL: " + url)
		s := string(xmlData)
		fmt.Println(s)
		fmt.Println("first 100 chars of xmlbody was: " + s[0:100]) //s) //"last 100 chars of jsonbody was: " + s[len(s)-100:])
		return nassTiffReader{}, err
	}
	response2, err2 := client.Get(resultURL.ReturnURL)
	if err2 != nil {
		fmt.Println(err2)
		return nassTiffReader{}, err2
	}
	defer response2.Body.Close()
	//vsimem
	fileparts := strings.Split(resultURL.ReturnURL, "/")
	outfile := "/workspaces/Go_Consequences/data/" + fileparts[len(fileparts)-1]

	out, err := os.Create(outfile)
	if err != nil {
		fmt.Println(err)
		return nassTiffReader{}, err
	}
	defer out.Close()
	io.Copy(out, response2.Body)

	ret := Init(outfile)
	return ret, nil
}

//GetCDLFileByFIPSFiltered provides a filtered geotif for a fips code. croptype is the list of values to keep and can be provided as a comma separated array of values to include
func GetCDLFileByFIPSFiltered(year string, fips string, cropType string) bool {
	url := fmt.Sprintf("%sGetCDLFile?year=%s&fips=%s", apiStatsURL, year, fips)
	return nassFilteredFileAPI(url, cropType)
}
func nassFilteredFileAPI(url string, cropType string) bool {
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
	var resultURL XMLFileURLResponse
	if err := xml.Unmarshal(xmlData, &resultURL); err != nil {
		fmt.Println("error unmarshalling NASS XML " + err.Error() + " URL: " + url)
		s := string(xmlData)
		fmt.Println(s)
		fmt.Println("first 100 chars of xmlbody was: " + s[0:100]) //s) //"last 100 chars of jsonbody was: " + s[len(s)-100:])
		return false
	}
	// now filter the file
	url2 := fmt.Sprintf("%sExtractCDLByValues?file=%s&values=%s", apiStatsURL, resultURL.ReturnURL, cropType)
	response2, err2 := client.Get(url2)
	if err2 != nil {
		fmt.Println(err2)
		return false
	}
	defer response2.Body.Close()
	filteredxmlData, err := ioutil.ReadAll(response2.Body)
	var filteredresultURL XMLExtractResponse
	if err := xml.Unmarshal(filteredxmlData, &filteredresultURL); err != nil {
		fmt.Println("error unmarshalling NASS XML " + err.Error() + " URL: " + url)
		s := string(filteredxmlData)
		fmt.Println(s)
		fmt.Println("first 100 chars of xmlbody was: " + s[0:100]) //s) //"last 100 chars of jsonbody was: " + s[len(s)-100:])
		return false
	}
	response3, err3 := client.Get(filteredresultURL.ReturnURL)
	if err3 != nil {
		fmt.Println(err3)
		return false
	}
	defer response3.Body.Close()
	fileparts := strings.Split(filteredresultURL.ReturnURL, "/")
	outfile := "C:\\Temp\\agtesting\\" + fileparts[len(fileparts)-1]
	out, err := os.Create(outfile)
	if err != nil {
		fmt.Println(err)
		return false
	}
	defer out.Close()
	io.Copy(out, response3.Body)

	return true
}
