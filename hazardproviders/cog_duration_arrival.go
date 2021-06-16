package hazardproviders

import (
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
func InitDaAHP(durationfp string, arrivalfp string, startTime time.Time) cogDurationAndArrivalHazardProvider {
	return cogDurationAndArrivalHazardProvider{durationCR: initCR(durationfp), arrivalCR: initCR(arrivalfp), startTime: startTime}
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
