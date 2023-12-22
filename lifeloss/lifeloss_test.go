package lifeloss

import (
	"math"
	"testing"

	"github.com/HydrologicEngineeringCenter/go-statistics/paireddata"
	"github.com/USACE/go-consequences/hazards"
	"github.com/USACE/go-consequences/structures"
)

func TestComputeLifeloss(t *testing.T) {
	lle := Init()
	//build a basic structure with a defined depth damage relationship.
	x := []float64{1.0, 2.0, 3.0, 4.0}
	y := []float64{10.0, 20.0, 30.0, 40.0}
	pd := paireddata.PairedData{Xvals: x, Yvals: y}
	sm := make(map[hazards.Parameter]structures.DamageFunction)
	var sdf = structures.DamageFunctionFamily{DamageFunctions: sm}

	df := structures.DamageFunction{}
	df.Source = "fabricated"
	df.DamageFunction = pd
	df.DamageDriver = hazards.Depth

	sdf.DamageFunctions[hazards.Default] = df
	cm := make(map[hazards.Parameter]structures.DamageFunction)
	var cdf = structures.DamageFunctionFamily{DamageFunctions: cm}
	cdf.DamageFunctions[hazards.Default] = df
	componentmap := make(map[string]structures.DamageFunctionFamily)
	componentmap["contents"] = cdf
	componentmap["structure"] = sdf
	var o = structures.OccupancyTypeDeterministic{Name: "RES1`", ComponentDamageFunctions: componentmap}
	var s = structures.StructureDeterministic{
		BaseStructure:    structures.BaseStructure{DamCat: "category"},
		OccType:          o,
		FoundType:        "",
		ConstructionType: "M",
		StructVal:        100.0,
		ContVal:          100.0,
		FoundHt:          0.0,
		Pop2pmo65:        100,
		Pop2pmu65:        100,
		Pop2amo65:        100,
		Pop2amu65:        100,
		NumStories:       1,
	}
	//lle.WarningSystem.WarningFunction()(&s)
	//test depth values
	var d = hazards.DepthandDVEvent{}
	depths := []float64{0.0, 0.5, 1.0, 1.0001, 2.25, 2.5, 2.75, 3.99, 4, 5, 12}
	expectedResults := []float64{0.0, 0.0, 10.0, 10.001, 22.5, 25.0, 27.5, 39.9, 40.0, 40.0, 40.0}
	for idx := range depths {
		d.SetDepth(depths[idx])
		d.SetDV(75.3)
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
		lle.ComputeLifeLoss(d, s)
	}
	//test interpolation due to foundation height putting depth back in range
	s.FoundHt = 8.1
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

/*
func TestGenerateLowLethality(t *testing.T) {
	//low fatality rates
	frequencies := []float64{0.000000000000, 0.000003213863, 0.000006427726, 0.000009641590, 0.000012855450, 0.000016069320, 0.000019283180, 0.000022497040, 0.000025710910, 1.000000000000}
	proportionallifeloss := []float64{1.0, 1.0, 1.0, 0.5, 0.5, 0.2, 0.1428571, 0.04, 0.0, 0.0}
	lowFatalities := paireddata.PairedData{
		Xvals: frequencies,
		Yvals: proportionallifeloss,
	}
	bytes, err := json.Marshal(lowFatalities)
	if err != nil {
		t.Fail()
	}
	os.WriteFile("/workspaces/Go_Consequences/structures/lowlethality.json", bytes, fs.ModeAppend)
}
func TestGenerateHighLethality(t *testing.T) {
	//high fatality rates
	frequencies := []float64{0.0, 0.007246377, 0.01449275, 0.02173913, 0.02898551, 0.03623188, 0.04347826, 0.05072464, 0.05797102, 0.06521739, 0.07246377, 0.07971015, 0.08695652, 0.0942029, 0.1014493, 0.1086956, 0.115942, 0.1231884, 0.1304348, 0.1376812, 0.1449275, 0.1521739, 0.1594203, 0.1666667, 0.173913, 0.1811594, 0.1884058, 0.1956522, 0.2028985, 0.2101449, 0.2173913, 0.2246377, 0.2318841, 0.2391304, 0.2463768, 0.2536232, 0.2608696, 0.2681159, 0.2753623, 0.2826087, 0.2898551, 0.2971014, 0.3043478, 0.3115942, 0.3188406, 0.326087, 0.3333333, 0.3405797, 0.3478261, 0.3550725, 0.3623188, 0.3695652, 0.3768116, 0.384058, 0.3913043, 0.3985507, 0.4057971, 0.4130435, 0.4202898, 0.4275362, 0.4347826, 0.442029, 0.4492754, 0.4565217, 0.4637681, 0.4710145, 0.4782609, 0.4855072, 0.4927536, 0.5, 0.5072464, 0.5144928, 0.5217391, 0.5289855, 0.5362319, 0.5434783, 0.5507246, 0.557971, 0.5652174, 0.5724638, 0.5797101, 0.5869565, 0.5942029, 0.6014493, 0.6086956, 0.615942, 0.6231884, 0.6304348, 0.6376812, 0.6449276, 0.6521739, 0.6594203, 0.6666667, 0.6739131, 0.6811594, 0.6884058, 0.6956522, 0.7028986, 0.7101449, 0.7173913, 0.7246377, 0.7318841, 0.7391304, 0.7463768, 0.7536232, 0.7608696, 0.7681159, 0.7753623, 0.7826087, 0.7898551, 0.7971014, 0.8043478, 0.8115942, 0.8188406, 0.8260869, 0.8333333, 0.8405797, 0.8478261, 0.8550724, 0.8623188, 0.8695652, 0.8768116, 0.884058, 0.8913044, 0.8985507, 0.9057971, 0.9130435, 0.9202899, 0.9275362, 0.9347826, 0.942029, 0.9492754, 0.9565217, 0.9637681, 0.9710145, 0.9782609, 0.9855072, 1.0}
	proportionallifeloss := []float64{1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 0.975, 0.9708738, 0.96, 0.9565217, 0.9090909, 0.9, 0.8947368, 0.8571429, 0.8571429, 0.8333333, 0.8235294, 0.8, 0.8, 0.8, 0.75, 0.75, 0.75, 0.75, 0.75, 0.6666667, 0.6666667, 0.6666667, 0.6666667, 0.6666667, 0.5714286, 0.5714286, 0.55, 0.5, 0.5, 0.5, 0.5, 0.5, 0.5, 0.5, 0.5, 0.5, 0.5, 0.5, 0.5, 0.5, 0.4285714, 0.4285714, 0.3333333, 0.3333333, 0.3333333, 0.3333333, 0.3333333, 0.3333333, 0.3333333, 0.3333333, 0.2857143, 0.2, 0.2, 0.2, 0.1666667, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0}
	highFatalities := paireddata.PairedData{
		Xvals: frequencies,
		Yvals: proportionallifeloss,
	}
	bytes, err := json.Marshal(highFatalities)
	if err != nil {
		t.Fail()
	}
	os.WriteFile("/workspaces/Go_Consequences/structures/highlethality.json", bytes, fs.ModeAppend)
}
*/
