package hazardproviders

import (
	"github.com/USACE/go-consequences/geography"
	"github.com/USACE/go-consequences/hazards"
)

type cogHazardProvider struct {
	depthcr cogReader
}

//Init creates and produces an unexported cogHazardProvider
func Init(fp string) (cogHazardProvider, error) {
	d, ed := initCR(fp)
	return cogHazardProvider{depthcr: d}, ed
}
func Init_Meters(fp string) (cogHazardProvider, error) {
	dm, edm := initCR_Meters(fp)
	return cogHazardProvider{depthcr: dm}, edm
}
func (chp cogHazardProvider) Close() {
	chp.depthcr.Close()
}
func (chp cogHazardProvider) ProvideHazard(l geography.Location) (hazards.HazardEvent, error) {
	h := hazards.DepthEvent{}
	d, err := chp.depthcr.ProvideValue(l)
	h.SetDepth(d)
	return h, err
}
func (chp cogHazardProvider) ProvideHazardBoundary() (geography.BBox, error) {
	return chp.depthcr.GetBoundingBox()
}
