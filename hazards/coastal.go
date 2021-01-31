package hazards

//CoastalEvent describes a coastal event
type CoastalEvent struct {
	Depth      float64 //still depth
	WaveHeight float64 //continuous variable.
	Salinity   bool    //default is false
}

//Parameters implements the HazardEvent interface
func (ad CoastalEvent) Parameters() Parameter {
	adp := Default
	adp = SetHasDepth(adp)
	if ad.WaveHeight > 0.0 {
		adp = SetHasWaveHeight(adp)
	}
	if ad.Salinity {
		adp = SetHasSalinity(adp)
	}
	return adp
}

//Has implements the HazardEvent Interface
func (ad CoastalEvent) Has(p Parameter) bool {
	adp := ad.Parameters()
	return adp&p != 0
}
