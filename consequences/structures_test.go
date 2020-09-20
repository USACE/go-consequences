package consequences

import (
	"math"
	"testing"

	"github.com/USACE/go-consequences/hazards"
	"github.com/USACE/go-consequences/paireddata"
)

func TestComputeConsequences(t *testing.T) {

	//build a basic structure with a defined depth damage relationship.
	x := []float64{1.0, 2.0, 3.0, 4.0}
	y := []float64{10.0, 20.0, 30.0, 40.0}
	pd := paireddata.PairedData{Xvals: x, Yvals: y}
	var o = OccupancyType{Name: "test", Structuredamfun: pd, Contentdamfun: pd}
	var s = Structure{OccType: o, DamCat: "category", StructVal: 100.0, ContVal: 100.0, FoundHt: 0.0}

	//test depth values
	var d = hazards.DepthEvent{Depth: 0.0}
	depths := []float64{0.0, 0.5, 1.0, 1.0001, 2.25, 2.5, 2.75, 3.99, 4, 5}
	expectedResults := []float64{0.0, 0.0, 10.0, 10.001, 22.5, 25.0, 27.5, 39.9, 40.0, 40.0}
	for idx := range depths {
		d.Depth = depths[idx]
		got := s.ComputeConsequences(d).Results[0].(float64)
		diff := expectedResults[idx] - got
		if math.Abs(diff) > .00000000000001 { //one more order of magnitude smaller causes 2.75 and 3.99 samples to fail.
			t.Errorf("ComputeConsequences(%f) = %f; expected %f", depths[idx], got, expectedResults[idx])
		}
	}

	s.FoundHt = 1.1 //test interpolation due to foundation height putting depth back in range
	got := s.ComputeConsequences(d).Results[0].(float64)
	if got != 39.0 {
		t.Errorf("ComputeConsequences(%f) = %f; expected %f", 39.0, got, 39.0)
	}

}
