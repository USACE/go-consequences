package crops

type Crop struct {
	ID                 int
	Name               string
	Yeild              float64
	PricePerUnit       float64
	ValuePerOutputUnit float64
	ProductionFunction productionFunction
	LossFunction       DamageFunction
	CropSchedule       CropSchedule
}
func (c Crop) GetX() float64{
	
}