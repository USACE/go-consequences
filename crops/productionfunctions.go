package crops

import "time"

type productionFunction struct {
	harvestCost                      float64
	cumulativeMonthlyProductionCosts []float64
	productionCostLessHarvest        float64 //sum monthly or find max of cumulative...
	lossFromLatePlanting             float64
}

func NewProductionFunction(mc []float64, cs CropSchedule, hc float64, latePlantingLoss float64) productionFunction {
	pf := productionFunction{
		harvestCost:          hc,
		lossFromLatePlanting: latePlantingLoss,
	}
	cmc, pclh := cumulateMonthlyCosts(mc, cs)
	pf.cumulativeMonthlyProductionCosts = cmc
	pf.productionCostLessHarvest = pclh
	return pf
}
func isLeapYear(year int) bool{  
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
func cumulateMonthlyCosts(mc []float64, cs CropSchedule) ([]float64, float64) {
	//check for winter crops.
	totalCosts := 0.0
	cmc := make([]float64, 12)
	if cs.StartPlantingDate.YearDay()+cs.DaysToMaturity > 365 {
		//winter crop.
		daysInYear := 365
		if isLeapYear(cs.StartPlantingDate.Year()){
			daysInYear +=1
		}

	} else {
		//contained between 0 and 365 days
		startMonth := cs.StartPlantingDate.Month() //iota "enum"
		startMonthIndex := int(startMonth)
		counter := 0
		daysToMaturity := cs.DaysToMaturity
		for ok := true; ok; ok = daysToMaturity > 0 {
			//compute days in the current month https://yourbasic.org/golang/last-day-month-date/
			t := time.Date(cs.StartPlantingDate.Year(), time.Month(startMonthIndex+counter+1), 0, 0, 0, 0, 0, time.UTC)
			daysInMonth := t.Day() //subtract the days in the current month from days to maturity.
			daysToMaturity -= daysInMonth
			totalCosts += mc[startMonthIndex+counter]
			cmc[startMonthIndex+counter] = totalCosts
			counter++
		}

	}

	return cmc, totalCosts
}
