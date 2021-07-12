package indirecteconomics

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"net/http"
	"time"
)

//fwlinkid =2
//https://www.hec.usace.army.mil/fwlink/?linkid=2&type=string

/*
	root += "?SIM=" + Simulation;
	root += "&ALT=" + Alternative;
	root += "&State=" + StateFIPS;
	root += "&CountyFIPS=" + CountyFIPS;
	root += "&LLR=" + (1 - LaborReduction).ToString();//format to 4 sig digits
	root += "&CLR=" + (1 - CapitalReduction).ToString();//format to 4 sig digits
	root += "&State_Name=" + StateAbbreviation;
	root += "&County_Name=" + CountyName + "_" + StateAbbreviation;//countyname_st
	root += "&Time=" + DateTime.Now.ToString();//format yyyyMMdd:HHmm:ss
*/
type EcamRequest struct {
	Simulation       string
	Alternative      string
	StateFIPS        string //2 digit fips
	CountyFIPS       string //5 digit fips? might be the additional 3 digits - not sure
	State            string
	County           string //countyname_stateabbreviation (2 character e.g OK)
	CapitalReduction float64
	LaborReduction   float64
}

func ComputeEcam(stateFips string, countyFips string, capitalLoss float64, laborloss float64) (EcamResult, error) {
	fwlinkurl := "https://www.hec.usace.army.mil/fwlink/?linkid=2&type=string"
	transCfg := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, // accept untrusted servers
	}
	client := &http.Client{Transport: transCfg}

	response, err := client.Get(fwlinkurl)
	if err != nil {
		return EcamResult{}, err
	}
	rc := response.Body
	sr := bufio.NewScanner(rc)
	sr.Scan()
	ecamRoot := sr.Text()
	ecamRoot += "?SIM=" + "Simulation"
	ecamRoot += "&ALT=" + "Alternative"
	ecamRoot += "&State=" + stateFips
	ecamRoot += "&CountyFIPS=" + countyFips
	ecamRoot += "&LLR=" + fmt.Sprintf("%f", (1-laborloss))
	ecamRoot += "&CLR=" + fmt.Sprintf("%f", (1-capitalLoss))
	ecamRoot += "&State_Name=" + "ST"
	ecamRoot += "&County_Name=" + "CN" + "_" + "ST" //countyname_st
	ecamRoot += "&Time=" + time.Now().Format("20121101:0304:05")
	response, err = client.Get(ecamRoot)
	if err != nil {
		return EcamResult{}, err
	}
	return ParseEcamResult(response)
}
