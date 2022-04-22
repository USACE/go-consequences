package indirecteconomics

import (
	"bufio"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

type EcamResult struct {
	ProductionImpacts []EcamSectorResultOutput
	EmploymentImpacts []EcamSectorResultOutput
}

type EcamSectorResultOutput struct {
	Sector        string
	Benchmark     float64 //millions usd
	PercentChange float64
	Change        float64
}

func (er EcamResult) String() string {
	report := "Production\n"
	report += "Sector, Previous Output (millions USD), Post Shock Output (millions USD), Output Change (millions USD)\n"
	for _, r := range er.ProductionImpacts {
		report += fmt.Sprintf("%v, %f, %f, %f", r.Sector, r.Benchmark, r.Benchmark+r.Change, r.Change) + "\n" // i am not sure this is right
	}
	report += "\n"
	report += "Labor\n"
	report += "Sector, Previous Employment, Post Shock Employment, Employment Change\n"
	for _, r := range er.EmploymentImpacts {
		report += fmt.Sprintf("%v, %f, %f, %f", r.Sector, r.Benchmark, r.PercentChange, r.Change) + "\n" //i am not sure this is right either.
	}
	return report
}
func ParseEcamResult(webResponse *http.Response) (EcamResult, error) {
	//need significant error checking.  ref: hec2.ecam.ComputeECAM.java
	rc := webResponse.Body
	sr := bufio.NewScanner(rc)
	outputimpacts := make([]EcamSectorResultOutput, 0)
	laborimpacts := make([]EcamSectorResultOutput, 0)
	for sr.Scan() {
		if strings.Contains(sr.Text(), "Exit_Code") {
			break
		}
	}
	errorparts := strings.Split(sr.Text(), "|")
	if len(errorparts) >= 2 {
		ErrorCode := errorparts[1]
		if ErrorCode == "0" {
			sr.Scan() //Begin_Output_LF
			if sr.Text() == "BEGIN_OUTPUT_LF" {
				for sr.Scan() {
					if strings.Contains(sr.Text(), "END_OUTPUT_LF") {
						break
					}
					tmp, err := parseEcamSectorResult(sr.Text())
					if err != nil {
						return EcamResult{}, err
					}
					outputimpacts = append(outputimpacts, tmp)
				}
			} else {
				return EcamResult{}, errors.New("Expected BEGIN_OUTPUT_LF got " + sr.Text())
			}
			sr.Scan() //Begin_Employment_LF
			if sr.Text() == "BEGIN_EMPLOYMENT_LF" {
				for sr.Scan() {
					if strings.Contains(sr.Text(), "END_EMPLOYMENT_LF") {
						break
					}
					tmp, err := parseEcamSectorResult(sr.Text())
					if err != nil {
						return EcamResult{}, err
					}
					laborimpacts = append(laborimpacts, tmp)
				}
			} else {
				return EcamResult{}, errors.New("Expected BEGIN_EMPLOYMENT_LF got " + sr.Text())
			}
			return EcamResult{ProductionImpacts: outputimpacts, EmploymentImpacts: laborimpacts}, nil
		} else {
			return EcamResult{}, errors.New("ECAM server returned Error Code " + ErrorCode)
		}
	} else {
		b, err := ioutil.ReadAll(rc)
		if err != nil {
			return EcamResult{}, errors.New("ECAM server something we couldnt parse " + sr.Text() + " and we got this error when trying to read the response body " + err.Error())
		}
		return EcamResult{}, errors.New("ECAM server something we couldnt parse " + sr.Text() + " with body " + string(b))
	}

}

func parseEcamSectorResult(outputstring string) (EcamSectorResultOutput, error) {
	results := strings.Split(outputstring, "|")
	//ensure there are 4 outputs
	sector, err := abreviationToDetailedString(results[0])
	if err != nil {
		return EcamSectorResultOutput{}, err
	}
	if len(results) < 4 {
		return EcamSectorResultOutput{}, errors.New("did not find 4 values in the sector line for " + sector)
	} else {
		bench, err := strconv.ParseFloat(results[1], 64)
		if err != nil {
			return EcamSectorResultOutput{}, errors.New("could not parse benchmark " + results[1])
		}
		pctchange, err := strconv.ParseFloat(results[2], 64)
		if err != nil {
			return EcamSectorResultOutput{}, errors.New("could not parse percent change " + results[2])
		}
		change, err := strconv.ParseFloat(results[3], 64)
		if err != nil {
			return EcamSectorResultOutput{}, errors.New("could not parse change " + results[3])
		}
		return EcamSectorResultOutput{Sector: sector, Benchmark: bench, PercentChange: pctchange, Change: change}, nil
	}
}
func abreviationToDetailedString(abrv string) (string, error) {
	switch strings.ToUpper(abrv) {
	case "AGR":
		return "Agriculture", nil
	case "LVS":
		return "Livestock and ranching", nil
	case "FRS":
		return "Forestry", nil
	case "FSH":
		return "Fishing", nil
	case "CRU":
		return "Oil gas and coal Extraction", nil
	case "MIN":
		return "Minerals mining", nil
	case "PWR":
		return "Electric power generation and supply", nil
	case "GAS":
		return "Natural gas distribution", nil
	case "WTR":
		return "Water sewage and other systems", nil
	case "CON":
		return "Construction", nil
	case "FOD":
		return "Food processing", nil
	case "BEV":
		return "Beverages", nil
	case "TBC":
		return "Tobacco", nil
	case "TEX":
		return "Textiles and wearing apparel", nil
	case "WOD":
		return "Wood manufacturing", nil
	case "PPP":
		return "Paper printing and publishing", nil
	case "CHM":
		return "Chemical processing and refining", nil
	case "MAN":
		return "General manufacturing", nil
	case "ELE":
		return "Electronic instruments", nil
	case "CAR":
		return "Transportation equipment manufacturing", nil
	case "FRN":
		return "Furniture manufacturing", nil
	case "COM":
		return "Post and communications", nil
	case "TRN":
		return "Transportation services", nil
	case "TRD":
		return "Wholesale and retail distribution", nil
	case "INF":
		return "Information processing and publication", nil
	case "FIN":
		return "Financial services and insurance", nil
	case "REC":
		return "Recreation activities", nil
	case "SER":
		return "All other services", nil
	case "ORG":
		return "Non-government associations", nil
	case "GOV":
		return "State and federal government", nil
	case "RWJ":
		return "Rest of world adjustment", nil
	case "IVJ":
		return "Inventory valuation adjustment", nil
	case "DWE":
		return "Owner occupied dwellings", nil
	case "TOTAL":
		return "Total", nil
	default:
		return abrv, errors.New("Could not parse " + abrv)
	}
}
