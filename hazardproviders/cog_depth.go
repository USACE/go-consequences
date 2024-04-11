package hazardproviders

import (
	"time"

	"github.com/USACE/go-consequences/geography"
	"github.com/USACE/go-consequences/hazards"
)

type cogHazardProvider struct {
	depthcr cogReader
	process HazardFunction
}

// Init creates and produces an unexported cogHazardProvider
func Init(fp string) (cogHazardProvider, error) {
	d, ed := initCR(fp)
	return cogHazardProvider{depthcr: d, process: DepthHazardFunction()}, ed
}
func Init_CustomFunction(fp string, function HazardFunction) (cogHazardProvider, error) {
	d, ed := initCR(fp)
	return cogHazardProvider{depthcr: d, process: function}, ed
}
func Init_Meters(fp string) (cogHazardProvider, error) {
	dm, edm := initCR_Meters(fp)
	return cogHazardProvider{depthcr: dm, process: DepthHazardFunction()}, edm
}
func Init_Meters_CustomFunction(fp string, function HazardFunction) (cogHazardProvider, error) {
	dm, edm := initCR_Meters(fp)
	return cogHazardProvider{depthcr: dm, process: function}, edm
}
func (chp cogHazardProvider) Close() {
	chp.depthcr.Close()
}
func (chp cogHazardProvider) Hazard(l geography.Location) (hazards.HazardEvent, error) {
	var h hazards.HazardEvent
	d, err := chp.depthcr.ProvideValue(l)
	if err != nil {
		return h, err
	}
	hd := hazards.HazardData{
		Depth:       d,
		Velocity:    0,
		ArrivalTime: time.Time{},
		Erosion:     0,
		Duration:    0,
		WaveHeight:  0,
		Salinity:    false,
		Qualitative: "",
	}
	return chp.process(hd, h)
}

func (chp cogHazardProvider) HazardBoundary() (geography.BBox, error) {
	return chp.depthcr.GetBoundingBox()
}
func (chp cogHazardProvider) SpatialReference() string {
	return chp.depthcr.SpatialReference()
}
func (chp cogHazardProvider) UpdateSpatialReference(sr_wkt string) {
	chp.depthcr.UpdateSpatialReference(sr_wkt)
}
