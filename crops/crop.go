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

//BuildCrop builds a crop since the properties of crop are not exported
func BuildCrop(cropid byte, cropname string) Crop {
	return Crop{id: cropid, name: cropname}
}

//WithLocation allows the construction of a location on a crop
func (c *Crop) WithLocation(xloc float64, yloc float64) Crop {
	c.x = xloc
	c.y = yloc
	return *c
}

//WithOutput allows the setting of the yeild per acre and price per unit of output
func (c *Crop) WithOutput(cropYeild float64, price float64) Crop {
	c.yeild = cropYeild
	c.pricePerUnit = price
	return *c
}

//WithProductionFunction allows the setting of the production function
func (c *Crop) WithProductionFunction(pf productionFunction) Crop {
	c.productionFunction = pf
	return *c
}

//WithLossFunction allows the setting of the loss function
func (c *Crop) WithLossFunction(lf DamageFunction) Crop {
	c.lossFunction = lf
	return *c
}

//WithCropSchedule allows the setting of the cropschedule
func (c *Crop) WithCropSchedule(cs CropSchedule) Crop {
	c.cropSchedule = cs
	return *c
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
