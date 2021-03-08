package crops

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"
)

//xmlCrop is used for reading xml files not for any other real purpose
type xmlCrop struct {
	ID                    byte         `xml:"id"`
	Name                  string       `xml:"name"`
	Yeild                 float64      `xml:"Yield"`
	Unit                  string       `xml:"Unit"`
	PricePerUnit          float64      `xml:"UnitPrice"`
	HarvestCost           float64      `xml:"HarvestCost"`
	FirstPlantDate        string       `xml:"FirstPlantDate"`
	LastPlantDate         string       `xml:"LastPlantDate"`
	HarvestDate           string       `xml:"HarvestDate"`
	MonthlyFixedCost      string       `xml:"MonthlyFixedCost"`
	MonthlyFirstPlantCost string       `xml:"MonthlyFirstPlantCost"`
	MonthlyLastPlantCost  string       `xml:"MonthlyLastPlantCost"`
	PercentLossLastPlant  float64      `xml:"PctLossLastPlant"`
	DryoutPeriod          int32        `xml:"DryoutPeriod"`
	SubstituteCropID      byte         `xml:"SubstituteCrop"`
	Durations             xmlDurations `xml:"Durations"`
}
type xmlDurations struct {
	XMLName  xml.Name `xml:"Durations"`
	Duration []string `xml:"Duration"`
}

//ReadFromXML reads crop schedules loss functions and production functions from xml (HEC-FIA format)
func ReadFromXML(filePath string) Crop {

	xmlFile, err := os.Open(filePath)
	if err != nil {
		fmt.Println(err)
	}
	defer xmlFile.Close()
	bytes, _ := ioutil.ReadAll(xmlFile)

	var c xmlCrop
	if errxml := xml.Unmarshal(bytes, &c); err != nil {
		fmt.Println(errxml)
	}
	ret := BuildCrop(c.ID, c.Name)
	ret = ret.WithOutput(c.Yeild, c.PricePerUnit)
	//parse the cropschedule
	st := xmltoTime(c.FirstPlantDate)
	et := xmltoTime(c.LastPlantDate)
	ht := xmltoTime(c.HarvestDate)
	dtm := st.YearDay() - ht.YearDay()
	if dtm < 0 {
		dtm += 365
	}
	cs := CropSchedule{StartPlantingDate: st, LastPlantingDate: et, DaysToMaturity: dtm}
	ret.WithCropSchedule(cs)
	//parse the loss function
	lf := xmltoLossFunction(c.Durations)
	ret.WithLossFunction(lf)
	//parse the production function
	pf := xmltoProductionFunction(c.MonthlyFirstPlantCost, c.MonthlyLastPlantCost, c.MonthlyFixedCost, cs, c.HarvestCost, c.PercentLossLastPlant, c.Yeild, c.PricePerUnit)
	ret.WithProductionFunction(pf)
	return ret
}
func xmltoProductionFunction(mcfps string, mclps string, mfcs string, cs CropSchedule, hc float64, lpl float64, yeild float64, price float64) productionFunction {
	totalValue := yeild * price
	mcfpvals := strings.Split(mcfps, ",")
	mclpvals := strings.Split(mclps, ",")
	mfcvals := strings.Split(mfcs, ",")
	//convert to floats (this code is not yet correct. Monthly costs are listed as fractions for variable and dollars for fixed..)
	mcfp := make([]float64, len(mcfpvals))
	mclp := make([]float64, len(mclpvals))
	mfc := make([]float64, len(mfcvals))
	totalFixedCosts := 0.0
	totalVariableCostsFP := 0.0
	totalVariableCostsLP := 0.0
	for i := 0; i < len(mcfpvals); i++ {
		f, _ := strconv.ParseFloat(mcfpvals[i], 64)
		mcfp[i] = f
		totalVariableCostsFP += f
		f2, _ := strconv.ParseFloat(mclpvals[i], 64)
		mclp[i] = f2
		totalVariableCostsLP += f2
		f3, _ := strconv.ParseFloat(mfcvals[i], 64)
		totalFixedCosts += f3
		mfc[i] = f3
	}
	if (totalFixedCosts + totalVariableCostsFP) > totalValue {
		panic("Costs are higher than product value! I DECLARE BANKRUPTCY")
	}
	if (totalFixedCosts + totalVariableCostsLP) > totalValue {
		panic("Costs are higher than product value! I DECLARE BANKRUPTCY")
	}

	pf := NewProductionFunction(mcfp, mclp, mfc, cs, hc, lpl)
	return pf
}
func xmltoTime(ddMMM string) time.Time {
	//not sure how this handles leap years - it always assigns to year 0000 currently.
	const layout = "02Jan"
	t, err := time.Parse(layout, ddMMM)
	if err != nil {
		fmt.Println(err)
	}
	return t
}
func xmltoLossFunction(input xmlDurations) DamageFunction {
	m := make(map[float64][]float64)
	for _, s := range input.Duration {
		vals := strings.Split(s, ",")
		//convert to floats
		damages := make([]float64, len(vals)-1)
		for i := 1; i < len(vals); i++ {
			f, _ := strconv.ParseFloat(vals[i], 64)
			damages[i-1] = f
		}
		//add to map
		f, _ := strconv.ParseFloat(vals[0], 64)
		m[f] = damages
	}
	//construct damagefunction
	df := DamageFunction{DurationDamageCurves: m}
	return df
}
