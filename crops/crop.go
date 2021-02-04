package crops

import (
	"fmt"

	"github.com/USACE/go-consequences/consequences"
	"github.com/USACE/go-consequences/hazards"
)

//Crop describes a crop that can be used to compute agricultural consequences
type Crop struct {
	id                 byte
	name               string
	substituteName     string
	x                  float64
	y                  float64
	totalMarketValue   float64 //Marketable value, yeild *pricePerUnit
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

//WithOutput allows the setting of the yeild per acre and price per unit of output and resulting value per output
func (c *Crop) WithOutput(cropYeild float64, price float64) Crop {
	c.totalMarketValue = cropYeild * price
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

//GetTotalMarketValue returns crop.totalMarketValue
func (c Crop) GetTotalMarketValue() float64 {
	return c.totalMarketValue
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
			damage = c.computeImpactedCase(da)
		case NotImpactedDuringSeason:
			damage = 0.0
		case PlantingDelayed:
			damage = c.computeDelayedCase(da)
		case NotPlanted:
			damage = c.computeNotPlantedCase(da)
		case SubstituteCrop:
			// case for sbustitute crop not yet implemented
			//get the substitute, and compute damages on it... hope for no infinate loop.
			damage = c.computeSubstitueCase(da)
		default:
			damage = 0.0
		}
		results[2] = damage
	}

	r := consequences.Results{IsTable: false, Result: ret}
	return r
}
func (c Crop) computeImpactedCase(e hazards.ArrivalandDurationEvent) float64 {
	// Determine crop damage percent based on damage dur curve and event dur
	dmgfactor := c.lossFunction.ComputeDamagePercent(e) / 100
	exposedProductionValue := c.productionFunction.GetExposedValue(e)
	totalProductionCost := c.productionFunction.productionCostLessHarvest
	percentProductionValue := exposedProductionValue / totalProductionCost
	totalMarketValue := c.GetTotalMarketValue()
	totalMarketValueLessHarvestCost := totalMarketValue - c.productionFunction.harvestCost
	loss := dmgfactor * percentProductionValue * totalMarketValueLessHarvestCost
	fmt.Println("loss = ", loss)
	return loss
}

func (c Crop) computeDelayedCase(e hazards.ArrivalandDurationEvent) float64 {
	// delayed loss is equivalent to total marketable value less harvest cost, times the percent loss due to late planting
	// Not using interpolated % loss for late plant
	plantingWindow := (c.cropSchedule.LastPlantingDate.Sub(c.cropSchedule.StartPlantingDate).Hours() / 24)
	fmt.Println(plantingWindow)
	actualPlant := (e.ArrivalTime.AddDate(0, 0, int(e.DurationInDays)))
	fmt.Println(actualPlant)
	daysLate := (actualPlant.Sub(c.cropSchedule.StartPlantingDate.AddDate(actualPlant.Year(), 0, 0))).Hours() / 24
	fmt.Println(daysLate)
	factor := (daysLate / plantingWindow) * c.productionFunction.lossFromLatePlanting / 100
	fmt.Println("factor is : ", factor)
	return c.GetTotalMarketValue() * factor
}

func (c Crop) computeNotPlantedCase(e hazards.ArrivalandDurationEvent) float64 {
	// Assume Loss is only fixed costs for entire year
	return c.productionFunction.GetCumulativeMonthlyFixedCostsOnly()[12]
}

func (c Crop) computeSubstitueCase(e hazards.ArrivalandDurationEvent) float64 {
	// TODO
	return 0.0
}
