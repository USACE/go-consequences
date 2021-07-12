package indirecteconomics

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"net/http"
	"time"
)

type CapitalAndLabor struct {
	Capital float64
	Labor   int64
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
