package hazardproviders

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/USACE/go-consequences/geography"
	"github.com/USACE/go-consequences/hazards"
)

type jsonArrivalDepthDurationMultiHazardProvider struct {
	arrivals  []time.Time
	depthCRs  []cogReader
	durations []float64
	process   HazardFunction
}

type ADDProperties struct {
	Year      float64 `json:"year"`
	Month     float64 `json:"month"`
	Day       float64 `json:"day"`
	Depth     float64 `json:"depth"`
	Depthgrid string  `json:"depthgrid"`
	Duration  float64 `json:"duration"`
	Xmin      float64 `json:"xmin"`
	Xmax      float64 `json:"xmax"`
	Ymin      float64 `json:"ymin"`
	Ymax      float64 `json:"ymax"`
}

type ADDEvents struct {
	Events []ADDProperties `json:"events"`
}

func InitADDMHP(fp string) (jsonArrivalDepthDurationMultiHazardProvider, error) {
	fmt.Println("Connecting to: " + fp)
	c, err := os.ReadFile(fp)
	if err != nil {
		panic(err)
	}

	var events ADDEvents
	json.Unmarshal(c, &events)

	arrivalTimes := make([]time.Time, len(events.Events))
	durations := make([]float64, len(events.Events))
	depthCRs := make([]cogReader, len(events.Events))

	for i, e := range events.Events {
		at := time.Date(int(e.Year), time.Month(e.Month), int(e.Day), 0, 0, 0, 0, time.UTC)
		cr, err := initCR(e.Depthgrid)
		if err != nil {
			panic(err)
		}

		arrivalTimes[i] = at
		durations[i] = e.Duration
		depthCRs[i] = cr
	}

	return jsonArrivalDepthDurationMultiHazardProvider{
		arrivals:  arrivalTimes,
		depthCRs:  depthCRs,
		durations: durations,
		process:   ArrivalDepthAndDurationHazardFunction(),
	}, nil
}

func (j jsonArrivalDepthDurationMultiHazardProvider) Close() {
	for _, c := range j.depthCRs {
		c.Close()
	}
}

func (j jsonArrivalDepthDurationMultiHazardProvider) Hazard(l geography.Location) (hazards.HazardEvent, error) {
	var hm hazards.ArrivalDepthandDurationEventMulti
	for i, cr := range j.depthCRs {
		d, err := cr.ProvideValue(l)
		if err != nil {
			return hm, err
		}
		hd := hazards.HazardData{
			Depth:       d,
			Velocity:    0,
			ArrivalTime: j.arrivals[i],
			Erosion:     0,
			Duration:    j.durations[i],
			WaveHeight:  0,
			Salinity:    false,
			Qualitative: "",
		}
		var h hazards.HazardEvent
		h, err = j.process(hd, h)
		if err != nil {
			panic(err)
		}
		hm.Append(h.(hazards.ArrivalDepthandDurationEvent))
	}
	return &hm, nil
}

func (j jsonArrivalDepthDurationMultiHazardProvider) HazardBoundary() (geography.BBox, error) {
	// We'll probably want to do something different here.
	//   Probably allow user to define study bbox with a
	//   shapefile/geojson/directly entering bbox coords
	return j.depthCRs[0].GetBoundingBox()
}
