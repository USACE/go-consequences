package hazardproviders

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/USACE/go-consequences/geography"
	"github.com/USACE/go-consequences/hazards"
)

type csvArrivalDepthDurationMultiHazardProvider struct {
	f         *os.File
	arrivals  []time.Time
	depths    []float64
	durations []float64
	process   HazardFunction
	bbox      geography.BBox
}

func InitCSV(fp string, b geography.BBox) (csvArrivalDepthDurationMultiHazardProvider, error) {
	fmt.Println("Connecting to: " + fp)
	file, err := os.Open(fp)
	if err != nil {
		panic(err)
	}

	reader := csv.NewReader(file)

	rows, err := reader.ReadAll()
	if err != nil {
		panic(err)
	}
	arrivals := make([]time.Time, len(rows)-1)
	depths := make([]float64, len(rows)-1)
	durations := make([]float64, len(rows)-1)

	for i, row := range rows[1:] {
		year, err := strconv.Atoi(row[2])
		if err != nil {
			panic(err)
		}
		month, err := strconv.Atoi(strings.TrimSpace(row[3]))
		if err != nil {
			panic(err)
		}
		day, err := strconv.Atoi(strings.TrimSpace(row[4]))
		if err != nil {
			panic(err)
		}
		at := time.Date(year, time.Month(month), int(day), 0, 0, 0, 0, time.UTC)
		depth, err := strconv.ParseFloat(row[5], 64)
		if err != nil {
			panic(err)
		}
		dur, err := strconv.ParseFloat(row[6], 64)
		if err != nil {
			panic(err)
		}

		arrivals[i] = at
		depths[i] = depth
		durations[i] = dur
	}

	return csvArrivalDepthDurationMultiHazardProvider{
		f:         file,
		arrivals:  arrivals,
		depths:    depths,
		durations: durations,
		process:   ArrivalDepthAndDurationHazardFunction(),
		bbox:      b,
	}, nil

}

func (c csvArrivalDepthDurationMultiHazardProvider) Close() {
	c.f.Close()
}

func (c csvArrivalDepthDurationMultiHazardProvider) Hazard(l geography.Location) (hazards.HazardEvent, error) {
	var hm hazards.ArrivalDepthandDurationEventMulti
	if c.bbox.Contains(l) {
		for i, d := range c.depths {
			hd := hazards.HazardData{
				Depth:       d,
				Velocity:    0,
				ArrivalTime: c.arrivals[i],
				Erosion:     0,
				Duration:    c.durations[i],
				WaveHeight:  0,
				Salinity:    false,
				Qualitative: "",
			}
			var h hazards.HazardEvent
			h, err := c.process(hd, h)
			if err != nil {
				panic(err)
			}
			hm.Append(h.(hazards.ArrivalDepthandDurationEvent))
		}
	}
	return &hm, nil
}

func (c csvArrivalDepthDurationMultiHazardProvider) HazardBoundary() (geography.BBox, error) {
	return c.bbox, nil
}
