package crops

import (
	"time"

	"github.com/USACE/go-consequences/hazards"
)

type productionFunction struct {
	harvestCost                           float64
	cumulativeMonthlyProductionCostsEarly []float64
	cumulativeMonthlyProductionCostsLate  []float64
	cumulativeMonthlyFixedCostsOnly       []float64
	productionCostLessHarvest             float64 //sum monthly or find max of cumulative...
	lossFromLatePlanting                  float64
}

func (p productionFunction) GetCumulativeMonthlyFixedCostsOnly() []float64 {
	return p.cumulativeMonthlyFixedCostsOnly
}

func (p productionFunction) GetExposedValue(e hazards.ArrivalandDurationEvent) float64 {
	hazardMonth := e.ArrivalTime.Month() //iota "enum"
	hazardMonthIndex := int(hazardMonth) - 1
	// We need to prorate monthly costs later
	// This assumes planting on time. Cannot deal with case where you are delayed THEN impacted.
	return p.cumulativeMonthlyProductionCostsEarly[hazardMonthIndex]
}

//NewProductionFunction is the constructor for the unexported productionFunction which represents the costs associated with producing a crop
func NewProductionFunction(mcfp []float64, mclp []float64, mfc []float64, cs CropSchedule, hc float64, latePlantingLoss float64) productionFunction {
	pf := productionFunction{
		harvestCost:          hc,
		lossFromLatePlanting: latePlantingLoss,
	}
	cmce, pclhe, fxe := cumulateMonthlyCosts(mcfp, mfc, cs.StartPlantingDate, cs.DaysToMaturity)
	pf.cumulativeMonthlyProductionCostsEarly = cmce
	cmcl, _, _ := cumulateMonthlyCosts(mclp, mfc, cs.LastPlantingDate, cs.DaysToMaturity)
	pf.cumulativeMonthlyProductionCostsLate = cmcl
	pf.productionCostLessHarvest = pclhe //is this appropriate should i store both to ensure proper accounting??
	pf.cumulativeMonthlyFixedCostsOnly = fxe
	return pf
}
func isLeapYear(year int) bool {
	leapFlag := false
	if year%4 == 0 {
		if year%100 == 0 {
			if year%400 == 0 {
				leapFlag = true
			} else {
				leapFlag = false
			}
		} else {
			leapFlag = true
		}
	} else {
		leapFlag = false
	}
	return leapFlag
}

//GetCumulativeMonthlyProductionCostsEarly provides access to the productionFunction's cumulative monthly production costs based on planting on the early date of the planting season
func (p productionFunction) GetCumulativeMonthlyProductionCostsEarly() []float64 {
	return p.cumulativeMonthlyProductionCostsEarly
}
func cumulateMonthlyCosts(mc []float64, fc []float64, start time.Time, daysToMaturity int) ([]float64, float64, []float64) {
	//this process assumes days to maturity is less than 1 year.
	totalCosts := 0.0
	totalFixed := 0.0
	cmc := make([]float64, 12)
	cfxc := make([]float64, 12)
	daysInYear := 365
	if isLeapYear(start.Year()) {
		daysInYear++
	}
	if daysToMaturity > daysInYear {
		panic("abort! abort! we hit an artery!")
	}
	startMonth := start.Month() //iota "enum"
	startMonthIndex := int(startMonth) - 1
	counter := 0
	updated := false
	year := start.Year()
	for ok := true; ok; ok = daysToMaturity > 0 {
		//compute days in the current month https://yourbasic.org/golang/last-day-month-date/
		t := time.Date(year, time.Month(startMonthIndex+counter+1), 0, 0, 0, 0, 0, time.UTC)
		daysInMonth := t.Day() //subtract the days in the current month from days to maturity.
		if counter == 0 {
			if !updated {
				daysInMonth -= start.Day() //remove the start of the month where crops weren't planted
			}
		}
		daysToMaturity -= daysInMonth
		if startMonthIndex+counter > 11 {
			if !updated {
				startMonthIndex = 0
				counter = 0
				year++
				updated = true
			}
		}
		totalFixed += fc[startMonthIndex+counter]
		cfxc[startMonthIndex+counter] = totalFixed
		totalCosts += mc[startMonthIndex+counter] + fc[startMonthIndex+counter]
		cmc[startMonthIndex+counter] = totalCosts

		counter++
	}
	return cmc, totalCosts, cfxc
}
