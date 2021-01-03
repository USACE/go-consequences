package crops

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

