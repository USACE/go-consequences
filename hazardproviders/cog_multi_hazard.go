package hazardproviders

import (
	"errors"
	"fmt"
	"time"

	"github.com/USACE/go-consequences/geography"
	"github.com/USACE/go-consequences/hazards"
)

type cogMultiHazardProvider struct {
	paramCogMap map[hazards.Parameter]cogReader
	startTime   time.Time
	//process     HazardFunction build this out later, it should be customizable, but not sure how to deal with the variety of info at the moment.
}

// InitMulti creates and produces an unexported cogMultiHazardProvider
func InitMulti(hpinfo HazardProviderInfo) (cogMultiHazardProvider, error) {
	tmpmap := make(map[hazards.Parameter]cogReader)
	for _, hp_param_and_path := range hpinfo.Hazards {
		cr, err := initCR(hp_param_and_path.FilePath)
		if err != nil {
			return cogMultiHazardProvider{}, err
		}
		tmpmap[hp_param_and_path.Hazard] = cr
	}

	return cogMultiHazardProvider{paramCogMap: tmpmap, startTime: hpinfo.StartTime}, nil
}
func (chp cogMultiHazardProvider) Close() {
	for _, v := range chp.paramCogMap {
		v.Close()
	}
}
func (chp cogMultiHazardProvider) Hazard(l geography.Location) (hazards.HazardEvent, error) {
	var h hazards.HazardEvent
	hd := hazards.HazardData{
		Depth:       -901,
		Velocity:    -901,
		ArrivalTime: time.Time{},
		Erosion:     -901,
		Duration:    -901,
		WaveHeight:  -901,
		Salinity:    false,
		Qualitative: "",
		DV:          -901,
	}
	for k, v := range chp.paramCogMap {
		hval, err := v.ProvideValue(l)
		if err != nil {
			return h, err
		}
		if k == hazards.ArrivalTime {
			//arrival time is more complicated than other parameters and it needs to be converted to be relative to a fixed date and time like a start time.
			sat := fmt.Sprintf("%fh", hval)
			duration, _ := time.ParseDuration(sat)
			t := chp.startTime.Add(duration)
			hd.SetParameter(k, t)
		} else {
			hd.SetParameter(k, hval)
		}

	}
	multi := hazards.HazardDataToMultiParameter(hd)
	return multi, nil //chp.process(hd, h)
}

func (chp cogMultiHazardProvider) HazardBoundary() (geography.BBox, error) {
	for _, v := range chp.paramCogMap {
		return v.GetBoundingBox() //what if the two grids are not the same extent?
	}
	return geography.BBox{}, errors.New("no values in the map of parameter and cogreaders")
}
