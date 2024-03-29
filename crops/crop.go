package crops

import (
	"errors"

	"github.com/USACE/go-consequences/consequences"
	"github.com/USACE/go-consequences/geography"
	"github.com/USACE/go-consequences/hazards"
)

//Crop describes a crop that can be used to compute agricultural consequences
type Crop struct {
	id                 byte
	name               string
	substituteName     string
	SubstituteCrop     Substitute
	x                  float64
	y                  float64
	totalMarketValue   float64 //Marketable value, yeild *pricePerUnit
	productionFunction productionFunction
	lossFunction       DamageFunction
	cropSchedule       CropSchedule
}
type Substitute struct {
	id                 byte
	name               string
	totalMarketValue   float64 //Marketable value, yeild *pricePerUnit
	productionFunction productionFunction
	lossFunction       DamageFunction
	cropSchedule       CropSchedule
}

func (s Substitute) toCrop() Crop {
	c := BuildCrop(s.id, s.name)
	c.WithProductionFunction(s.productionFunction)
	c.WithLossFunction(s.lossFunction)
	c.WithCropSchedule(s.cropSchedule)
	c.totalMarketValue = s.totalMarketValue
	return c
}
func (c Crop) toSubstitute() Substitute {
	s := Substitute{id: c.id, name: c.name, totalMarketValue: c.totalMarketValue, productionFunction: c.productionFunction, lossFunction: c.lossFunction, cropSchedule: c.cropSchedule}
	return s
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
func (c Crop) Location() geography.Location {
	return geography.Location{X: c.x, Y: c.y}
}

//GetTotalMarketValue returns crop.totalMarketValue
func (c Crop) GetTotalMarketValue() float64 {
	return c.totalMarketValue
}

//Compute implements concequence.Receptor on crop
func (c Crop) Compute(event hazards.HazardEvent) (consequences.Result, error) {
	//Check event to determine if it is an arrival time and duration event
	header := []string{"Crop", "x", "y", "Damage Outcome", "Damage", "Duration", "Arrival Time"}
	results := []interface{}{c.name, c.Location().X, c.Location().Y, Unassigned.String(), 0.0, 0.0, ""}
	var ret = consequences.Result{Headers: header, Result: results}
	var err error = nil
	da, ok := event.(hazards.ArrivalandDurationEvent)
	if ok {
		//determine cropdamageoutcome
		outcome := c.cropSchedule.ComputeCropDamageCase(da)

		results[5] = da.Duration()
		results[6] = da.ArrivalTime().Format("Mon Jan 2 15:04:05")
		//switch case on damageoutcome
		//compute damages
		damage := 0.0
		switch outcome {
		case Unassigned:
			//huh?
			damage = 0.0
			err = errors.New("Damage Outcome was Unassigned")
		case Impacted:
			damage = c.computeImpactedCase(da)
		case NotImpactedDuringSeason:
			damage = 0.0
			err = errors.New("Damage Outcome was Not Impacted During Season")
		case PlantingDelayed:
			damage = c.computeDelayedCase(da)
		case NotPlanted:
			damage, outcome = c.computeNotPlantedCase(da)
		case SubstituteCrop:
			// case for sbustitute crop not yet implemented
			//get the substitute, and compute damages on it... hope for no infinate loop.
			damage = c.computeSubstitueCase(da)
		default:
			damage = 0.0
			err = errors.New("Damage Outcome resulted in Default case")
		}
		results[3] = outcome.String()
		results[4] = damage
	}
	return ret, err
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
	//fmt.Println("loss = ", loss)
	return loss
}

func (c Crop) computeDelayedCase(e hazards.ArrivalandDurationEvent) float64 {
	// delayed loss is equivalent to total marketable value less harvest cost, times the percent loss due to late planting
	// Not using interpolated % loss for late plant
	plantingWindow := (c.cropSchedule.LastPlantingDate.Sub(c.cropSchedule.StartPlantingDate).Hours() / 24)
	//fmt.Println(plantingWindow)
	actualPlant := (e.ArrivalTime().AddDate(0, 0, int(e.Duration())))
	//fmt.Println(actualPlant)
	daysLate := (actualPlant.Sub(c.cropSchedule.StartPlantingDate.AddDate(actualPlant.Year(), 0, 0))).Hours() / 24
	//fmt.Println(daysLate)
	factor := (daysLate / plantingWindow) * c.productionFunction.lossFromLatePlanting / 100
	//fmt.Println("factor is : ", factor)
	return c.GetTotalMarketValue() * factor
}

func (c Crop) computeNotPlantedCase(e hazards.ArrivalandDurationEvent) (float64, CropDamageCase) {
	// Assume Loss is only fixed costs for entire year
	if c.substituteName == "" {
		return c.productionFunction.GetCumulativeMonthlyFixedCostsOnly()[11], NotPlanted
	} else {

		return c.computeSubstitueCase(e), SubstituteCrop
	}

}

func (c Crop) computeSubstitueCase(e hazards.ArrivalandDurationEvent) float64 {
	// TODO
	//fmt.Println("compute substitute, TODO implement me.")
	ocv := c.totalMarketValue - c.productionFunction.harvestCost
	sc := c.SubstituteCrop.toCrop()
	scv := sc.totalMarketValue - sc.productionFunction.harvestCost
	if scv <= ocv {
		return ocv - scv
	}
	return 0.0
}
