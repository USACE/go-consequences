package crops

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
//crops.CropType
func (c Crop) GetCropID() byte {
	return c.id
}
//crops.CropType
func (c Crop) GetCropName() string {
	return c.name
}
//consequences.Locatable
func (c Crop) GetX() float64 {
	return c.x
}
//consequences.Locatable
func (c Crop) GetY() float64 {
	return c.y
}
