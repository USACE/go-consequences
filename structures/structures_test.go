package structures

import (
	"math"
	"testing"
	"time"

	"github.com/HydrologicEngineeringCenter/go-statistics/statistics"
	"github.com/USACE/go-consequences/consequences"
	"github.com/USACE/go-consequences/hazards"
	"github.com/USACE/go-consequences/paireddata"
)

func TestComputeConsequences(t *testing.T) {

	//build a basic structure with a defined depth damage relationship.
	x := []float64{1.0, 2.0, 3.0, 4.0}
	y := []float64{10.0, 20.0, 30.0, 40.0}
	pd := paireddata.PairedData{Xvals: x, Yvals: y}
	sm := make(map[hazards.Parameter]paireddata.ValueSampler)
	var sdf = DamageFunctionFamily{DamageFunctions: sm}
	sdf.DamageFunctions[hazards.Default] = pd
	cm := make(map[hazards.Parameter]paireddata.ValueSampler)
	var cdf = DamageFunctionFamily{DamageFunctions: cm}
	cdf.DamageFunctions[hazards.Default] = pd
	var o = OccupancyTypeDeterministic{Name: "test", StructureDFF: sdf, ContentDFF: cdf}
	var s = StructureDeterministic{OccType: o, StructVal: 100.0, ContVal: 100.0, FoundHt: 0.0, BaseStructure: BaseStructure{DamCat: "category"}}

	//test depth values
	var d = hazards.DepthEvent{}
	depths := []float64{0.0, 0.5, 1.0, 1.0001, 2.25, 2.5, 2.75, 3.99, 4, 5}
	expectedResults := []float64{0.0, 0.0, 10.0, 10.001, 22.5, 25.0, 27.5, 39.9, 40.0, 40.0}
	for idx := range depths {
		d.SetDepth(depths[idx])
		r, err := s.Compute(d)
		if err != nil {
			panic(err)
		}
		dr, err := r.Fetch("structure damage")
		if err != nil {
			panic(err)
		}
		got := dr.(float64)
		diff := expectedResults[idx] - got
		if math.Abs(diff) > .00000000000001 { //one more order of magnitude smaller causes 2.75 and 3.99 samples to fail.
			t.Errorf("Compute(%f) = %f; expected %f", depths[idx], got, expectedResults[idx])
		}
	}
	//test interpolation due to foundation height putting depth back in range
	s.FoundHt = 1.1
	r, err := s.Compute(d)
	if err != nil {
		panic(err)
	}
	dr, err := r.Fetch("structure damage")
	if err != nil {
		panic(err)
	}
	got := dr.(float64)
	if got != 39.0 {
		t.Errorf("Compute(%f) = %f; expected %f", 39.0, got, 39.0)
	}

	//add a test for content value as well
	//add a test for different hazard types (float64 and fire?)
}
func TestComputeConsequencesUncertainty(t *testing.T) {

	//build a basic structure with a defined depth damage relationship.
	x := []float64{1.0, 2.0, 3.0, 4.0}
	y := []float64{10.0, 20.0, 30.0, 40.0}
	pd := paireddata.PairedData{Xvals: x, Yvals: y}
	sm := make(map[hazards.Parameter]interface{})
	var sdf = DamageFunctionFamilyStochastic{DamageFunctions: sm}
	sdf.DamageFunctions[hazards.Default] = pd
	cm := make(map[hazards.Parameter]interface{})
	var cdf = DamageFunctionFamilyStochastic{DamageFunctions: cm}
	cdf.DamageFunctions[hazards.Default] = pd

	var o = OccupancyTypeStochastic{Name: "test", StructureDFF: sdf, ContentDFF: cdf}

	sv := statistics.NormalDistribution{Mean: 0, StandardDeviation: 1}
	cv := statistics.NormalDistribution{Mean: 0, StandardDeviation: 1}
	spv := consequences.ParameterValue{Value: sv}
	cpv := consequences.ParameterValue{Value: cv}
	fhpv := consequences.ParameterValue{Value: 0}
	var s = StructureStochastic{OccType: o, StructVal: spv, ContVal: cpv, FoundHt: fhpv, BaseStructure: BaseStructure{DamCat: "category"}}
	s.UseUncertainty = true
	//test depth values
	var d = hazards.DepthEvent{}
	depths := []float64{0.0, 0.5, 1.0, 1.0001, 2.25, 2.5, 2.75, 3.99, 4, 5}
	expectedResults := []float64{0.0, 0.0, -.052138, -0.030335, -0.122390, -0.088922, -0.146414, 0.205319, 0.108698, -0.625010}
	for idx := range depths {
		d.SetDepth(depths[idx])
		r, err := s.Compute(d)
		if err != nil {
			panic(err)
		}
		dr, err := r.Fetch("structure damage")
		if err != nil {
			panic(err)
		}
		got := dr.(float64)
		diff := expectedResults[idx] - got
		if math.Abs(diff) > .000001 {
			t.Errorf("Compute(%f) = %f; expected %f", depths[idx], got, expectedResults[idx])
		}
	}
}
func TestComputeConsequencesWithReconstruction(t *testing.T) {

	//build a basic structure with a defined depth damage relationship.
	x := []float64{1.0, 2.0, 3.0, 4.0}
	y := []float64{10.0, 20.0, 30.0, 40.0}
	pd := paireddata.PairedData{Xvals: x, Yvals: y}
	sm := make(map[hazards.Parameter]paireddata.ValueSampler)
	var sdf = DamageFunctionFamily{DamageFunctions: sm}
	sdf.DamageFunctions[hazards.Default] = pd
	cm := make(map[hazards.Parameter]paireddata.ValueSampler)
	var cdf = DamageFunctionFamily{DamageFunctions: cm}
	cdf.DamageFunctions[hazards.Default] = pd
	var o = OccupancyTypeDeterministic{Name: "test", StructureDFF: sdf, ContentDFF: cdf}
	var s = StructureDeterministic{OccType: o, StructVal: 100.0, ContVal: 100.0, FoundHt: 0.0, BaseStructure: BaseStructure{DamCat: "category"}}

	//test depth values
	var d = hazards.ArrivalDepthandDurationEvent{}
	d.SetDuration(2.5)
	at := time.Date(1984, time.Month(1), 22, 0, 0, 0, 0, time.UTC)
	d.SetArrivalTime(at)
	depths := []float64{0.0, 0.5, 1.0, 1.0001, 2.25, 2.5, 2.75, 3.99, 4, 5}
	expectedResults := []float64{0.0, 0.0, 10.0, 10.001, 22.5, 25.0, 27.5, 39.9, 40.0, 40.0}
	for idx := range depths {
		d.SetDepth(depths[idx])
		r, err := s.Compute(d)
		if err != nil {
			panic(err)
		}
		out, err := r.Fetch("daystoreconstruction")
		if err != nil {
			panic(err)
		}
		got := out.(float64)

		diff := (expectedResults[idx]*1.8 + d.Duration()) - got //180.0/100=1.8
		if math.Abs(diff) > .0000000000001 {                    //one more order of magnitude smaller causes 2.75 and 3.99 samples to fail.
			t.Errorf("Compute(%f) = %f; expected %f", depths[idx], got, expectedResults[idx]*1.8+d.Duration())
		}
		out2, err := r.Fetch("rebuilddate")
		if err != nil {
			panic(err)
		}
		gotdate := out2.(time.Time)
		if gotdate.Equal(d.ArrivalTime().AddDate(0, 0, int(expectedResults[idx]*1.8+d.Duration()))) {
			t.Errorf("Compute(%f) = %s; expected %s", depths[idx], gotdate, d.ArrivalTime().AddDate(0, 0, int(expectedResults[idx]*1.8+d.Duration())))
		}
	}

}
