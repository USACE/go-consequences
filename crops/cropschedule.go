package crops

import (
	"time"

	"github.com/USACE/go-consequences/hazards"
)

//CropSchedule stores the start and end of planting season, and time to maturity in days
type CropSchedule struct {
	StartPlantingDate time.Time
	LastPlantingDate  time.Time
	DaysToMaturity    int
}

//ComputeCropDamageCase evaluates a crop schedule against a hazard arrival and duration to determine impact on the crop season.
func (cs CropSchedule) ComputeCropDamageCase(h hazards.ArrivalandDurationEvent) CropDamageCase {
	//determine day of year of the hazard.
	hazardStartDoy := h.ArrivalTime().YearDay()
	//determine duration of the hazard
	hazardDurationDays := int(h.Duration())
	if hazardStartDoy <= cs.StartPlantingDate.YearDay() {
		//flood before start of planting.
		//determine if the crop start planting date is effected by the hazard
		if hazardStartDoy+hazardDurationDays < cs.StartPlantingDate.YearDay() {
			//what if it is a winter crop?
			if cs.StartPlantingDate.YearDay()+cs.DaysToMaturity > 365 {
				//winter crop
				harvestDoY := cs.StartPlantingDate.YearDay() + cs.DaysToMaturity - 365
				if harvestDoY > hazardStartDoy {
					return Impacted
				}
				return NotImpactedDuringSeason

			}
			return NotImpactedDuringSeason
		}
		//determine if hazard happened before planting and impacted planting season
		if hazardStartDoy+hazardDurationDays < cs.LastPlantingDate.YearDay() {
			return PlantingDelayed
		}
		return NotPlanted //what about substitutes?
	}
	//hazard after start of planting
	if cs.StartPlantingDate.YearDay()+cs.DaysToMaturity < hazardStartDoy {
		//hazard after harvest
		return NotImpactedDuringSeason
	}
	return Impacted
}
