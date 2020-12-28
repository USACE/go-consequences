package crops

import (
	"time"

	"github.com/USACE/go-consequences/hazards"
)

type CropSchedule struct {
	StartPlantingDate time.Time
	LastPlantingDate  time.Time
	DaysToMaturity    int
}

func (cs CropSchedule) ComputeCropDamageCase(h hazards.ArrivalandDurationEvent) CropDamageCase {
	//determine day of year of the hazard.
	hazard_start_doy := h.ArrivalTime.YearDay()
	//determine duration of the hazard
	hazard_duration_days := h.DurationInDays
	if hazard_start_doy <= cs.StartPlantingDate.YearDay() {
		//flood before start of planting.
		//determine if the crop start planting date is effected by the hazard
		if hazard_start_doy+hazard_duration_days < cs.StartPlantingDate.YearDay() {
			//what if it is a winter crop?
			if cs.StartPlantingDate.YearDay()+cs.DaysToMaturity > 365 {
				//winter crop
				harvestDoY := cs.StartPlantingDate.YearDay() + cs.DaysToMaturity - 365
				if harvestDoY > hazard_start_doy {
					return Impacted
				} else {
					return NotImpactedDuringSeason
				}
			} else {
				return NotImpactedDuringSeason
			}
		} else {
			//determine if hazard happened before planting and impacted planting season
			if hazard_start_doy+hazard_duration_days < cs.LastPlantingDate.YearDay() {
				return PlantingDelayed
			} else {
				return NotPlanted //what about substitutes?
			}
		}
	} else {
		//hazard after start of planting
		if cs.StartPlantingDate.YearDay()+cs.DaysToMaturity < hazard_start_doy {
			//hazard after harvest
			return NotImpactedDuringSeason
		} else {
			return Impacted
		}
	}
}
