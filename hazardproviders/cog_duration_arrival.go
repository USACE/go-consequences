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
}

//Init creates and produces an unexported cogHazardProvider
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
	return cogDurationAndArrivalHazardProvider{durationCR: d, arrivalCR: a, startTime: startTime}, et
}
func (chp cogDurationAndArrivalHazardProvider) Close() {
	chp.durationCR.Close()
	chp.arrivalCR.Close()
}
func (chp cogDurationAndArrivalHazardProvider) ProvideHazard(l geography.Location) (hazards.HazardEvent, error) {
	h := hazards.ArrivalandDurationEvent{}
	d, err := chp.durationCR.ProvideValue(l)
	if err != nil {
		return h, err
	}
	h.SetDuration(d)
	a, err := chp.arrivalCR.ProvideValue(l)
	if err != nil {
		return h, err
	}
	sat := fmt.Sprintf("%fh", a)
	duration, _ := time.ParseDuration(sat)
	t := chp.startTime.Add(duration)
	h.SetArrivalTime(t)
	return h, nil
}
func (chp cogDurationAndArrivalHazardProvider) ProvideHazardBoundary() (geography.BBox, error) {
	return chp.durationCR.GetBoundingBox() //what if the two grids are not the same extent?
}
