package hazards

import (
	"fmt"
	"time"
)

//CoastalEvent describes a coastal event
type CoastalEvent struct {
	depth      float64 //still depth
	waveHeight float64 //continuous variable.
	salinity   bool    //default is false
}

func (d CoastalEvent) MarshalJSON() ([]byte, error) {
	s := fmt.Sprintf("{\"coastalevent\":{\"depth\":%f, \"waveheight\":%f,\"salinity\":%t}}", d.Depth(), d.WaveHeight(), d.Salinity())
	return []byte(s), nil
}
func (h CoastalEvent) Depth() float64 {
	return h.depth
}
func (h *CoastalEvent) SetDepth(d float64) {
	h.depth = d
}
func (h CoastalEvent) Velocity() float64 {
	return -901.0
}
func (h CoastalEvent) ArrivalTime() time.Time {
	return time.Time{}
}
func (h CoastalEvent) ArrivalTime2ft() time.Time {
	return time.Time{}
}
func (h CoastalEvent) Duration() float64 {
	return -901.0
}
func (h CoastalEvent) WaveHeight() float64 {
	return h.waveHeight
}
func (h *CoastalEvent) SetWaveHeight(d float64) {
	h.waveHeight = d
}
func (h CoastalEvent) Salinity() bool {
	return h.salinity
}
func (h *CoastalEvent) SetSalinity(d bool) {
	h.salinity = d
}
func (h CoastalEvent) Qualitative() string {
	return ""
}

//Parameters implements the HazardEvent interface
func (ad CoastalEvent) Parameters() Parameter {
	adp := Default
	adp = SetHasDepth(adp)
	if ad.WaveHeight() > 0.0 {
		adp = SetHasWaveHeight(adp)
	}
	if ad.Salinity() {
		adp = SetHasSalinity(adp)
	}
	return adp
}

//Has implements the HazardEvent Interface
func (ad CoastalEvent) Has(p Parameter) bool {
	adp := ad.Parameters()
	return adp&p != 0
}
