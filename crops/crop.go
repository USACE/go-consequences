package crops

import (
	"github.com/USACE/go-consequences/consequences"
	"github.com/USACE/go-consequences/hazards"
)

//Crop describes a crop that can be used to compute agricultural consequences
type Crop struct {
	id                 byte
	name               string
	x                  float64
	y                  float64
	yeild              float64
	pricePerUnit       float64
	valuePerOutputUnit float64
	productionFunction productionFunction
	lossFunction       DamageFunction
	cropSchedule       CropSchedule
}

//GetCropID fulfils the crops.CropType interface
func (c Crop) GetCropID() byte {
	return c.id
}

//GetCropName fulfils the crops.CropType interface
func (c Crop) GetCropName() string {
	return c.name
}

//GetX fulfils the consequences.Locatable interface
func (c Crop) GetX() float64 {
	return c.x
}

//GetY fulfils the consequences.Locatable interface
func (c Crop) GetY() float64 {
	return c.y
}

//ComputeConsequences implements concequence receptor on crop
func (c Crop) ComputeConsequences(event interface{}) consequences.Results {
	//Check event to determine if it is an arrival time and duration event
	header := []string{"Crop", "Damage Outcome", "Damage"}
	results := []interface{}{c.name, Unassigned, 0.0}
	var ret = consequences.Result{Headers: header, Result: results}
	da, ok := event.(hazards.ArrivalandDurationEvent)
	if ok {
		//determine cropdamageoutcome
		outcome := c.cropSchedule.ComputeCropDamageCase(da)
		results[1] = outcome
		//switch case on damageoutcome
		//compute damages
		damage := 0.0
		switch outcome {
		case Unassigned:
			//huh?
			damage = 0.0
		case Impacted:
			damage = 10
		case NotImpactedDuringSeason:
			damage = 0.0
		case PlantingDelayed:
			damage = 1.0
		case NotPlanted:
			damage = 0.0 //fixed costs?
		case SubstituteCrop:
			//get the substitute, and compute damages on it... hope for no infinate loop.
			damage = 0.0
		default:
			damage = 0.0
		}
		results[2] = damage
	}

	r := consequences.Results{IsTable: false, Result: ret}
	return r
}
