package crops

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
func cumulateMonthlyCosts(mc []float64, cs CropSchedule) ([]float64, float64) {
	//check for winter crops.
	totalCosts := 0.0
	cmc := make([]float64, 12)
	if cs.StartPlantingDayOfYear+cs.DaysToMaturity > 365 {
		//winter crop.
	} else {
		//within a year

	}

	return cmc, totalCosts
}
