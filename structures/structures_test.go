package structures

import (
	"fmt"
	"math"
	"math/rand"
	"testing"
	"time"

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
	expectedResults := []float64{0.0, 0.0, 0.138100, -0.117163, -0.198414, -0.234834, -0.022169, -0.721810, -0.178571, 0.362431}
	var seed int64 = 1234
	r := rand.New(rand.NewSource(seed))

	for idx := range depths {
		d.SetDepth(depths[idx])
		sd := s.SampleStructure(r.Int63())
		r, err := sd.Compute(d)
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

	xr := []float64{0.0, 1.0}
	yr := []float64{0, 100.0}
	pdr := paireddata.PairedData{Xvals: xr, Yvals: yr}
	pdrdf := DamageFunction{}
	pdrdf.DamageFunction = pdr
	pdrdf.DamageDriver = hazards.Depth
	pdrdf.Source = "created for this test"
	dm := make(map[hazards.Parameter]DamageFunction)
	var rdf = DamageFunctionFamily{DamageFunctions: dm}
	rdf.DamageFunctions[hazards.Default] = pdrdf

	components := make(map[string]DamageFunctionFamily)
	components["structure"] = sdf
	components["contents"] = cdf
	components["reconstruction"] = rdf

	var o = OccupancyTypeDeterministic{Name: "test", ComponentDamageFunctions: components}
	var s = StructureDeterministic{OccType: o, StructVal: 100.0, ContVal: 100.0, FoundHt: 0.0, BaseStructure: BaseStructure{DamCat: "category"}}

	//test depth values
	var d = hazards.ArrivalDepthandDurationEvent{}
	d.SetDuration(3)
	at := time.Date(1984, time.Month(1), 22, 0, 0, 0, 0, time.UTC)
	d.SetArrivalTime(at)
	depths := []float64{0.0, 0.5, 1.0, 1.0001, 2.25, 2.5, 2.75, 3.99, 4, 5}
	expectedResults := []float64{0.0, 0.0, 10.0, 10.001, 22.5, 25.0, 27.5, 39.9, 40.0, 40.0}
	for idx := range depths {
		d.SetDepth(depths[idx])
		r, err := computeConsequencesWithReconstruction(d, s)
		if err != nil {
			panic(err)
		}

		out, err := r.Fetch("reconstruction_days")
		if err != nil {
			panic(err)
		}
		got := out.(float64)
		expect := math.Ceil(expectedResults[idx] + d.Duration())
		fmt.Printf("Reconstruction time was %3.2f days.\n", got)
		diff := expect - got //180.0/100=1.8
		if math.Abs(diff) > .0000000000001 {
			t.Errorf("Compute(%f) = %f; expected %f", depths[idx], got, expect)
		}

	}

}

func TestComputeConsequencesMulti(t *testing.T) {

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

	xr := []float64{0.0, 1.0}
	yr := []float64{0, 100.0}
	pdr := paireddata.PairedData{Xvals: xr, Yvals: yr}
	pdrdf := DamageFunction{}
	pdrdf.DamageFunction = pdr
	pdrdf.DamageDriver = hazards.Depth
	pdrdf.Source = "created for this test"
	dm := make(map[hazards.Parameter]DamageFunction)
	var rdf = DamageFunctionFamily{DamageFunctions: dm}
	rdf.DamageFunctions[hazards.Default] = pdrdf

	components := make(map[string]DamageFunctionFamily)
	components["structure"] = sdf
	components["contents"] = cdf
	components["reconstruction"] = rdf

	var o = OccupancyTypeDeterministic{Name: "test", ComponentDamageFunctions: components}
	var s = StructureDeterministic{OccType: o, StructVal: 100.0, ContVal: 100.0, FoundHt: 0.0, BaseStructure: BaseStructure{DamCat: "category"}}

	// create a series of hazardEvents
	var d1 = hazards.ArrivalDepthandDurationEvent{}
	d1.SetDuration(0)
	d1.SetDepth(1.0)
	t1 := time.Date(1984, time.Month(1), 1, 0, 0, 0, 0, time.UTC)
	d1.SetArrivalTime(t1)

	var d2 = hazards.ArrivalDepthandDurationEvent{}
	d2.SetDuration(5.0)
	d2.SetDepth(1.0)
	t2 := time.Date(1984, time.Month(1), 11, 0, 0, 0, 0, time.UTC)
	d2.SetArrivalTime(t2)

	var d3 = hazards.ArrivalDepthandDurationEvent{}
	d3.SetDuration(0.0)
	d3.SetDepth(1.0)
	t3 := time.Date(1984, time.Month(1), 21, 0, 0, 0, 0, time.UTC)
	d3.SetArrivalTime(t3)

	var d4 = hazards.ArrivalDepthandDurationEvent{}
	d4.SetDuration(0.0)
	d4.SetDepth(2.0)
	t4 := time.Date(1985, time.Month(1), 1, 0, 0, 0, 0, time.UTC)
	d4.SetArrivalTime(t4)

	var d5 = hazards.ArrivalDepthandDurationEvent{}
	d5.SetDuration(0.0)
	d5.SetDepth(2.0)
	t5 := time.Date(1985, time.Month(1), 11, 0, 0, 0, 0, time.UTC)
	d5.SetArrivalTime(t5)

	events := []hazards.HazardEvent{d1, d2, d3, d4, d5}

	et1 := time.Date(1984, time.Month(1), 11, 0, 0, 0, 0, time.UTC)
	et2 := time.Date(1984, time.Month(1), 26, 0, 0, 0, 0, time.UTC)
	et3 := time.Date(1984, time.Month(2), 5, 0, 0, 0, 0, time.UTC)
	et4 := time.Date(1985, time.Month(1), 21, 0, 0, 0, 0, time.UTC)
	et5 := time.Date(1985, time.Month(2), 8, 0, 0, 0, 0, time.UTC)

	expectedResults := []time.Time{et1, et2, et3, et4, et5}
	expectedDmgs := []float64{10.0, 10.0, 9.5, 20.0, 18.0}
	// event 1 (10% dmg) arrives first (obviously). No adjustments required. dmg = 100*0.1 = 10.0
	// event 2 (10% dmg) arrives after event 1 completes reconstruction. No adjustments required. dmg = 100*0.1 = 10.0
	// event 3 interrupts event 2 reconstruction @ 50% complete
	// 		event 2 had dmg=10, so structure was repaired by 5 when event 3 hit.
	// 		structure value when event 3 hits is 100-(10-5)=95.
	// 		event 3 does 10% damage ==> expected damage = 0.1*95 = 9.5
	//		damageFactor to calculate reconstruction = 1 - (1-sDamageFactor)*(1-sdampercent) = 1 - (1-0.05)*(1-.1) = 0.145 ==> 14.5 reconstruction days (rounds to 15)
	//		completion date should be 1984/2/5
	// event 4 occurs after event 3 reconstruction complete. No adjustments required. dmg = 100*0.2 = 20
	// event 5 interrupts event 4 reconstruction @ 50% complete
	// 		structure value when event 5 hits is 100 - (20-10) = 90.
	// 		event 5 does 20% damage ==> expected damage = 0.2*90 = 18.0
	//		damageFactor to calculate reconstruction = 1 - (1-sDamageFactor)*(1-sdampercent) = 1 - (1-0.1)*(1-.2) = 0.28 ==> 28 reconstruction days

	results, err := computeConsequencesMulti(events, s)
	if err != nil {
		panic(err)
	}

	for idx := range results {
		out, err := results[idx].Fetch("completion_date")
		if err != nil {
			panic(err)
		}

		dif := expectedResults[idx].Sub(out.(time.Time))
		fmt.Printf("Completion date was %v. Expected: %v. Diff: %v\n", out, expectedResults[idx], dif)
		if math.Abs(float64(dif)) > float64(time.Minute) { // if the error is greater than about 1 minute
			t.Errorf("Completion date was %v. Expected: %v. Diff: %v\n", out, expectedResults[idx], dif)

		}

		dmgout, err := results[idx].Fetch("structure damage")
		if err != nil {
			panic(err)
		}
		fmt.Printf("Damage was %3.2f. Expected: %3.2f\n", dmgout, expectedDmgs[idx])
		if math.Abs(dmgout.(float64)-float64(expectedDmgs[idx])) > 0.000000001 {
			t.Errorf("Damage was %3.2f. Expected: %3.2f\n", dmgout, expectedDmgs[idx])
		}

	}

}

func TestComputeConsequencesMultiHazard(t *testing.T) {
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

	xr := []float64{0.0, 1.0}
	yr := []float64{0, 100.0}
	pdr := paireddata.PairedData{Xvals: xr, Yvals: yr}
	pdrdf := DamageFunction{}
	pdrdf.DamageFunction = pdr
	pdrdf.DamageDriver = hazards.Depth
	pdrdf.Source = "created for this test"
	dm := make(map[hazards.Parameter]DamageFunction)
	var rdf = DamageFunctionFamily{DamageFunctions: dm}
	rdf.DamageFunctions[hazards.Default] = pdrdf

	components := make(map[string]DamageFunctionFamily)
	components["structure"] = sdf
	components["contents"] = cdf
	components["reconstruction"] = rdf

	var o = OccupancyTypeDeterministic{Name: "test", ComponentDamageFunctions: components}
	var s = StructureDeterministic{OccType: o, StructVal: 100.0, ContVal: 100.0, FoundHt: 0.0, BaseStructure: BaseStructure{DamCat: "category"}}

	// create a series of hazardEvents
	var d1 = hazards.ArrivalDepthandDurationEvent{}
	d1.SetDuration(0)
	d1.SetDepth(1.0)
	t1 := time.Date(1984, time.Month(1), 1, 0, 0, 0, 0, time.UTC)
	d1.SetArrivalTime(t1)

	var d2 = hazards.ArrivalDepthandDurationEvent{}
	d2.SetDuration(5.0)
	d2.SetDepth(1.0)
	t2 := time.Date(1984, time.Month(1), 11, 0, 0, 0, 0, time.UTC)
	d2.SetArrivalTime(t2)

	var d3 = hazards.ArrivalDepthandDurationEvent{}
	d3.SetDuration(0.0)
	d3.SetDepth(1.0)
	t3 := time.Date(1984, time.Month(1), 21, 0, 0, 0, 0, time.UTC)
	d3.SetArrivalTime(t3)

	var d4 = hazards.ArrivalDepthandDurationEvent{}
	d4.SetDuration(0.0)
	d4.SetDepth(2.0)
	t4 := time.Date(1985, time.Month(1), 1, 0, 0, 0, 0, time.UTC)
	d4.SetArrivalTime(t4)

	var d5 = hazards.ArrivalDepthandDurationEvent{}
	d5.SetDuration(0.0)
	d5.SetDepth(2.0)
	t5 := time.Date(1985, time.Month(1), 11, 0, 0, 0, 0, time.UTC)
	d5.SetArrivalTime(t5)

	events := []hazards.ArrivalDepthandDurationEvent{d5, d1, d2, d3, d4}

	addMulti := &hazards.ArrivalDepthandDurationEventMulti{Events: events} //need to use the pointer reference because methods on MultiHazardEvent require pointers

	// test that we can sort events by ArrivalTime
	if !addMulti.IsSorted() {
		fmt.Println("Sorting...")
		addMulti.Sort()
	}

	et1 := time.Date(1984, time.Month(1), 11, 0, 0, 0, 0, time.UTC)
	et2 := time.Date(1984, time.Month(1), 26, 0, 0, 0, 0, time.UTC)
	et3 := time.Date(1984, time.Month(2), 5, 0, 0, 0, 0, time.UTC)
	et4 := time.Date(1985, time.Month(1), 21, 0, 0, 0, 0, time.UTC)
	et5 := time.Date(1985, time.Month(2), 8, 0, 0, 0, 0, time.UTC)

	expectedResults := []time.Time{et1, et2, et3, et4, et5}
	expectedDmgs := []float64{10.0, 10.0, 9.5, 20.0, 18.0}

	results, err := s.Compute(addMulti)
	if err != nil {
		panic(err)
	}

	hr, err := results.Fetch("hazard results")
	if err != nil {
		panic(err)
	}
	hazardResults := hr.(consequences.Result)

	for i := range expectedResults {
		r, err := hazardResults.Fetch(fmt.Sprintf("%d", i))
		if err != nil {
			panic(err)
		}
		result := r.(consequences.Result)
		out, err := result.Fetch("completion_date")
		if err != nil {
			panic(err)
		}

		dif := expectedResults[i].Sub(out.(time.Time))
		fmt.Printf("Completion date was %v. Expected: %v. Diff: %v\n", out, expectedResults[i], dif)
		if math.Abs(float64(dif)) > float64(time.Minute) { // if the error is greater than about 1 minute
			t.Errorf("Completion date was %v. Expected: %v. Diff: %v\n", out, expectedResults[i], dif)

		}

		dmgout, err := result.Fetch("structure damage")
		if err != nil {
			panic(err)
		}
		fmt.Printf("Damage was %3.2f. Expected: %3.2f\n", dmgout, expectedDmgs[i])
		if math.Abs(dmgout.(float64)-float64(expectedDmgs[i])) > 0.000000001 {
			t.Errorf("Damage was %3.2f. Expected: %3.2f\n", dmgout, expectedDmgs[i])
		}

	}
}
