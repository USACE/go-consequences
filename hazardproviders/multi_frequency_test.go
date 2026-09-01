package hazardproviders

import (
	"testing"

	"github.com/USACE/go-consequences/geography"
	"github.com/USACE/go-consequences/hazards"
)

func TestInitMultiFrequency(t *testing.T) {
	file := "/workspaces/go-consequences/data/lifecycle/test_multi-frequency.csv"
	expectedDepths := []float64{1, 10, 30, 45, 59, 78, 89, 102, 140, 180, 240, 330, 350, 370}
	expectedFreqs := []float64{.99, .95, .9, .8, .7, .6, .5, .4, .3, .2, .1, .01, .002, .001}

	mfHP, err := InitMultiFrequencyDepthProvider(file, "CSV", "")
	if err != nil {
		panic(err)
	}
	defer mfHP.Close()

	loc := geography.Location{
		X: -71.45,
		Y: 42.95,
	}

	haz, err := mfHP.Hazard(loc)
	if err != nil {
		panic(err)
	}
	h := haz.(*hazards.DepthEventMultiFrequency)

	for {
		edepth := expectedDepths[h.Index()]
		efreq := expectedFreqs[h.Index()]

		if h.Depth() != edepth {
			t.Errorf("Event at index %d had Depth = %v. Expected: %3.2f", h.Index(), h.Depth(), edepth)
		}
		if h.Frequency() != efreq {
			t.Errorf("Event at index %d had Frequency = %v. Expected: %3.2f", h.Index(), h.Duration(), efreq)
		}

		if h.HasNext() {
			h.Increment()
		} else {
			break
		}
	}
}
