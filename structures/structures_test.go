package structures

import (
	"math"
	"testing"

	"github.com/HydrologicEngineeringCenter/go-statistics/paireddata"
	"github.com/HydrologicEngineeringCenter/go-statistics/statistics"
	"github.com/USACE/go-consequences/consequences"
	"github.com/USACE/go-consequences/hazards"
)

func TestComputeConsequences(t *testing.T) {

	//build a basic structure with a defined depth damage relationship.
	x := []float64{1.0, 2.0, 3.0, 4.0}
	y := []float64{10.0, 20.0, 30.0, 40.0}
	pd := paireddata.PairedData{Xvals: x, Yvals: y}
	sm := make(map[hazards.Parameter]DamageFunction)
	var sdf = DamageFunctionFamily{DamageFunctions: sm}

	df := DamageFunction{}
	df.Source = "fabricated"
	df.DamageFunction = pd
	df.DamageDriver = hazards.Depth

	sdf.DamageFunctions[hazards.Default] = df
	cm := make(map[hazards.Parameter]DamageFunction)
	var cdf = DamageFunctionFamily{DamageFunctions: cm}
	cdf.DamageFunctions[hazards.Default] = df
	componentmap := make(map[string]DamageFunctionFamily)
	componentmap["contents"] = cdf
	componentmap["structure"] = sdf
	var o = OccupancyTypeDeterministic{Name: "test", ComponentDamageFunctions: componentmap}
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
	y := arrayToDetermnisticDistributions([]float64{10.0, 20.0, 30.0, 40.0})
	pd := paireddata.UncertaintyPairedData{Xvals: x, Yvals: y}
	sm := make(map[hazards.Parameter]DamageFunctionStochastic)
	var sdf = DamageFunctionFamilyStochastic{DamageFunctions: sm}

	df := DamageFunctionStochastic{}
	df.Source = "fabricated"
	df.DamageFunction = pd
	df.DamageDriver = hazards.Depth

	sdf.DamageFunctions[hazards.Default] = df
	cm := make(map[hazards.Parameter]DamageFunctionStochastic)
	var cdf = DamageFunctionFamilyStochastic{DamageFunctions: cm}
	cdf.DamageFunctions[hazards.Default] = df
	componentmap := make(map[string]DamageFunctionFamilyStochastic)
	componentmap["contents"] = cdf
	componentmap["structure"] = sdf
	var o = OccupancyTypeStochastic{Name: "test", ComponentDamageFunctions: componentmap}

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

func TestComputeConsequences_erosion(t *testing.T) {
	df := erosionDamageFunction()
	//build a basic structure with a defined depth damage relationship.
	ot := OccupancyTypeStochastic{}
	ot.Name = "fake"
	cm := make(map[string]DamageFunctionFamilyStochastic)
	dfm := make(map[hazards.Parameter]DamageFunctionStochastic)
	dfm[hazards.Erosion] = df
	dffs := DamageFunctionFamilyStochastic{DamageFunctions: dfm}
	cm["contents"] = dffs
	cm["structure"] = dffs
	ot.ComponentDamageFunctions = cm
	var s = StructureStochastic{OccType: ot, StructVal: consequences.ParameterValue{Value: 100.0}, ContVal: consequences.ParameterValue{Value: 100.0}, FoundHt: consequences.ParameterValue{Value: 0.0}, BaseStructure: BaseStructure{DamCat: "category"}}

	//test depth values
	e := hazards.NewCoastalEvent(hazards.CoastalEvent{})
	erosions := []float64{0.0, 5, 10, 10.5, 20, 25, 50, 75, 85, 95}
	expectedResults := []float64{0.0, 0.0, 0.00, 0.05, 1, 1.375, 4.8, 7.55, 7.925, 8} //these need to be updated
	for idx := range erosions {
		e.SetErosion(erosions[idx])
		r, err := s.Compute(e)
		if err != nil {
			if e.Erosion() != 0.0 {
				t.Errorf("Compute(%f) = %v; expected %v", erosions[idx], err.Error(), "structure: hazard did not contain valid parameters to impact a structure")
			}
		} else {
			dr, err := r.Fetch("structure damage")
			if err != nil {
				panic(err)
			}
			got := dr.(float64)
			diff := expectedResults[idx] - got
			if math.Abs(diff) > .00000000000001 { //one more order of magnitude smaller causes 2.75 and 3.99 samples to fail.
				t.Errorf("Compute(%f) = %f; expected %f", erosions[idx], got, expectedResults[idx])
			}
		}

	}
	//test confirm foundation height does not impact results.
	s.FoundHt = consequences.ParameterValue{Value: statistics.DeterministicDistribution{Value: 1.1}}
	r, err := s.Compute(e)
	if err != nil {
		panic(err)
	}
	dr, err := r.Fetch("structure damage")
	if err != nil {
		panic(err)
	}
	got := dr.(float64)
	if got != expectedResults[len(erosions)-1] {
		t.Errorf("Compute(%f) = %f; expected %f", e.Erosion(), got, expectedResults[len(erosions)-1])
	}
}

/*
func TestComputeConsequencesWithReconstruction(t *testing.T) {

	//build a basic structure with a defined depth damage relationship.
	x := []float64{1.0, 2.0, 3.0, 4.0}
	y := []float64{10.0, 20.0, 30.0, 40.0}
	pd := paireddata.PairedData{Xvals: x, Yvals: y}
	pddf := DamageFunction{}
	pddf.DamageFunction = pd
	pddf.DamageDriver = hazards.Depth
	pddf.Source = "created for this test"
	sm := make(map[hazards.Parameter]DamageFunction)
	var sdf = DamageFunctionFamily{DamageFunctions: sm}
	sdf.DamageFunctions[hazards.Default] = pddf
	cm := make(map[hazards.Parameter]DamageFunction)
	var cdf = DamageFunctionFamily{DamageFunctions: cm}
	cdf.DamageFunctions[hazards.Default] = pddf
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
		if math.Abs(diff) > .0000000000001 {
			t.Errorf("Compute(%f) = %f; expected %f", depths[idx], got, expectedResults[idx]*1.8+d.Duration())
		}
		out2, err := r.Fetch("rebuilddate")
		if err != nil {
			panic(err)
		}
		gotdate := out2.(time.Time)
		if !gotdate.Equal(d.ArrivalTime().AddDate(0, 0, int(expectedResults[idx]*1.8+d.Duration()))) {
			t.Errorf("Compute(%f) = %s; expected %s", depths[idx], gotdate, d.ArrivalTime().AddDate(0, 0, int(expectedResults[idx]*1.8+d.Duration())))
		}
	}

}
*/
