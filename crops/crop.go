package crops

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