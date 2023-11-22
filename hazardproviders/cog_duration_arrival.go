package hazardproviders

import (
	"errors"
	"fmt"
	"time"

	"github.com/USACE/go-consequences/geography"
	"github.com/USACE/go-consequences/hazards"
)

type cogDurationAndArrivalHazardProvider struct {
	durationCR cogReader
	arrivalCR  cogReader //decimal days
	startTime  time.Time
	process    HazardFunction
}

// Init creates and produces an unexported cogHazardProvider
func InitDaAHP(durationfp string, arrivalfp string, startTime time.Time) (cogDurationAndArrivalHazardProvider, error) {
	d, ed := initCR(durationfp)
	a, ad := initCR(arrivalfp)
	var et error
	et = nil
	if ed != nil {
		if ad != nil {
			et = errors.New(ed.Error() + ad.Error())
		}
		et = ed
	}
	if ad != nil {
		et = ad
	}
	return cogDurationAndArrivalHazardProvider{durationCR: d, arrivalCR: a, startTime: startTime, process: ArrivalAndDurationHazardFunction()}, et
}
func (chp cogDurationAndArrivalHazardProvider) Close() {
	chp.durationCR.Close()
	chp.arrivalCR.Close()
}
func (chp cogDurationAndArrivalHazardProvider) Hazard(l geography.Location) (hazards.HazardEvent, error) {
	var h hazards.HazardEvent
	d, err := chp.durationCR.ProvideValue(l)
	if err != nil {
		return h, err
	}
	a, err := chp.arrivalCR.ProvideValue(l)
	if err != nil {
		return h, err
	}
	sat := fmt.Sprintf("%fh", a)
	duration, _ := time.ParseDuration(sat)
	t := chp.startTime.Add(duration)
	hd := hazards.HazardData{
		Depth:       0,
		Velocity:    0,
		ArrivalTime: t,
		Erosion:     0,
		Duration:    d,
		WaveHeight:  0,
		Salinity:    false,
		Qualitative: "",
	}
	return chp.process(hd, h)
}

func (chp cogDurationAndArrivalHazardProvider) HazardBoundary() (geography.BBox, error) {
	return chp.durationCR.GetBoundingBox() //what if the two grids are not the same extent?
}
