package crops

import (
	"github.com/USACE/go-consequences/hazards"
)

type Crop struct {
	ID                 int
	Name               string
	Yeild              float64
	PricePerUnit       float64
	ValuePerOutputUnit float64
	ProductionFunction string
	LossFunction       string
	CropSchedule       CropSchedule
}
type CropSchedule struct {
	StartPlantingDayOfYear int
	LastPlantingDayOfYear  int
	DaysToMaturity         int
}
type CropDamageCase byte

const (
	Unassigned              CropDamageCase = 0
	Impacted                CropDamageCase = 1
	NotImpactedDuringSeason CropDamageCase = 2
	PlantingDelayed         CropDamageCase = 4
	NotPlanted              CropDamageCase = 8
	SubstituteCrop          CropDamageCase = 16
)

func (cs CropSchedule) ComputeCropDamageCase(h hazards.Hazard_Event) CropDamageCase {
	//determine day of year of the hazard.
	hazard_start_doy := 0 //assign to zero to work through case logic.
	//determine duration of the hazard
	hazard_duration_days := 10 //assignment to start testing
	if hazard_start_doy <= cs.StartPlantingDayOfYear {
		//flood before start of planting.
		//determine if the crop start planting date is effected by the hazard
		if hazard_start_doy+hazard_duration_days < cs.StartPlantingDayOfYear {
			//what if it is a winter crop?
			if cs.StartPlantingDayOfYear+cs.DaysToMaturity > 365 {
				//winter crop
				harvestDoY := cs.StartPlantingDayOfYear + cs.DaysToMaturity - 365
				if harvestDoY > hazard_start_doy {
					return Impacted
				} else {
					return NotImpactedDuringSeason
				}
			}
		} else {
			//determine if hazard happened before planting and impacted planting season
			if hazard_start_doy+hazard_duration_days < cs.LastPlantingDayOfYear {
				return PlantingDelayed
			} else {
				return NotPlanted //what about substitutes?
			}
		}
	} else {
		//hazard after start of planting
		if cs.StartPlantingDayOfYear+cs.DaysToMaturity < hazard_start_doy {
			//hazard after harvest
			return NotImpactedDuringSeason
		} else {
			return Impacted
		}
	}
	return Unassigned
}
