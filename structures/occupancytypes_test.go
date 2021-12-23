package structures

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/HydrologicEngineeringCenter/go-statistics/statistics"
	"github.com/USACE/go-consequences/hazards"
	"github.com/USACE/go-consequences/paireddata"
)

func TestDamageFunctionFamily(t *testing.T) {
	//a map of hazard to damage function
	m := make(map[hazards.Parameter]DamageFunction)
	//a set of different hazard types.
	cep := hazards.Depth | hazards.Salinity | hazards.WaveHeight
	dep := hazards.Depth
	ce2p := hazards.Depth | hazards.WaveHeight

	//a fake deterministic damage function for coastal event with salinity
	cexs := []float64{1, 2, 3}
	ceys := []float64{1, 2, 3}
	var cedf = paireddata.PairedData{Xvals: cexs, Yvals: ceys}
	dfep := DamageFunction{}
	dfep.DamageFunction = cedf
	dfep.DamageDriver = hazards.Depth
	dfep.Source = "created for this test for coastal hazards"
	m[cep] = dfep

	//a fake deterministic damage function for a depth only event
	dexs := []float64{1, 2, 3}
	deys := []float64{2, 4, 6}
	var dedf = paireddata.PairedData{Xvals: dexs, Yvals: deys}
	dfdep := DamageFunction{}
	dfdep.DamageFunction = dedf
	dfdep.DamageDriver = hazards.Depth
	dfdep.Source = "created for this test for depth hazards"
	m[dep] = dfdep

	//a fake deterministic damage function for an event with depth and wave, but no salinity...
	ce2xs := []float64{1, 2, 3}
	ce2ys := []float64{3, 6, 9}
	var ce2df = paireddata.PairedData{Xvals: ce2xs, Yvals: ce2ys}
	dfcep := DamageFunction{}
	dfcep.DamageFunction = ce2df
	dfcep.DamageDriver = hazards.Depth
	dfcep.Source = "created for this test for coastal no salinity"
	m[ce2p] = dfcep

	//assign the fake damage function map as a family of damage functions.
	var df = DamageFunctionFamily{DamageFunctions: m}

	//fake instances of hazards (to match the hazard types.)
	ce := hazards.CoastalEvent{}
	ce.SetDepth(2)
	ce.SetSalinity(true)
	ce.SetWaveHeight(3.4)
	de := hazards.DepthEvent{}
	de.SetDepth(2)
	ce2 := hazards.CoastalEvent{}
	ce2.SetDepth(2)
	ce2.SetSalinity(false)
	ce2.SetWaveHeight(3.4)
	//confirm that for each hazard the correct damage function is pulled when requested and the proper damage value is computed.
	cv := df.DamageFunctions[ce.Parameters()].DamageFunction.SampleValue(ce.Depth())
	if cv != 2 {
		t.Errorf("Expected 2 got %v", cv)
	}
	dv := df.DamageFunctions[de.Parameters()].DamageFunction.SampleValue(de.Depth())
	if dv != 4 {
		t.Errorf("Expected 4 %v", dv)
	}
	c2v := df.DamageFunctions[ce2.Parameters()].DamageFunction.SampleValue(ce2.Depth())
	if c2v != 6 {
		t.Errorf("Expected 6 %v", c2v)
	}
}
func Test_occupancyCentralTendency(t *testing.T) {
	//a map of occupancy types
	m := OccupancyTypeMap()
	for name, ot := range m {
		fmt.Println("reading " + name)
		m2 := ot.CentralTendency()
		fmt.Println("computed " + m2.Name)
	}
}
func Test_occupancySample(t *testing.T) {
	//a map of occupancy types
	m := OccupancyTypeMap()
	for name, ot := range m {
		fmt.Println("reading " + name)
		m2 := ot.SampleOccupancyType(1234)
		fmt.Println("computed " + m2.Name)
	}
}
func Test_DamageFunctionStochastic_Marshal(t *testing.T) {

	//build a basic structure with a defined depth damage relationship.
	x := []float64{1.0, 2.0, 3.0, 4.0}
	y := arrayToDetermnisticDistributions([]float64{10.0, 20.0, 30.0, 40.0})
	pd := paireddata.UncertaintyPairedData{Xvals: x, Yvals: y}

	df := DamageFunctionStochastic{}
	df.Source = "fabricated"
	df.DamageFunction = pd
	df.DamageDriver = hazards.Depth

	b, err := json.Marshal(df)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(b))
	df2 := DamageFunctionStochastic{}
	err = json.Unmarshal(b, &df2)
	if err != nil {
		panic(err)
	}
	fmt.Println(df.Source)
}

func Test_DamageFunction_Marshal(t *testing.T) {
	//a fake deterministic damage function for a depth only event
	dexs := []float64{1, 2, 3}
	deys := []float64{2, 4, 6}
	var dedf = paireddata.PairedData{Xvals: dexs, Yvals: deys}
	df := DamageFunction{}
	df.DamageFunction = dedf
	df.DamageDriver = hazards.Depth
	df.Source = "created for testing marshaling"
	b, err := json.Marshal(df)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(b))
}
func Test_Erosion_DamageFunctionStochastic_Marshal(t *testing.T) {

	//build a basic structure with a defined depth damage relationship.
	x := []float64{10.0, 20.0, 30.0, 40.0, 50.0, 60.0, 70.0, 80.0, 90.0, 100.0}
	y := make([]statistics.ContinuousDistribution, 10)
	y[0] = statistics.TriangularDistribution{Min: 0.0, MostLikely: 0.0, Max: .5}
	y[1] = statistics.TriangularDistribution{Min: 0.5, MostLikely: 1.0, Max: 2.25}
	y[2] = statistics.TriangularDistribution{Min: 0.5, MostLikely: 1.75, Max: 4.5}
	y[3] = statistics.TriangularDistribution{Min: 0.5, MostLikely: 4.7, Max: 5.5}
	y[4] = statistics.TriangularDistribution{Min: 0.75, MostLikely: 4.8, Max: 6.5}
	y[5] = statistics.TriangularDistribution{Min: 0.75, MostLikely: 5.0, Max: 8.0}
	y[6] = statistics.TriangularDistribution{Min: 0.75, MostLikely: 7.25, Max: 9.0}
	y[7] = statistics.TriangularDistribution{Min: 1.0, MostLikely: 7.85, Max: 10.0}
	y[8] = statistics.TriangularDistribution{Min: 2.0, MostLikely: 8.0, Max: 11.0}
	y[9] = statistics.TriangularDistribution{Min: 3.5, MostLikely: 8.0, Max: 11.0}
	pd := paireddata.UncertaintyPairedData{Xvals: x, Yvals: y}

	df := DamageFunctionStochastic{}
	df.Source = "bhrercn"
	df.DamageFunction = pd
	df.DamageDriver = hazards.Erosion

	b, err := json.Marshal(df)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(b))

}

func Test_DamageFunctionFamilyStochastic_Marshal(t *testing.T) {
	m := make(map[hazards.Parameter]DamageFunctionStochastic)
	dffs := DamageFunctionFamilyStochastic{DamageFunctions: m}
	//build a basic structure with a defined depth damage relationship.
	x := []float64{1.0, 2.0, 3.0, 4.0}
	y := arrayToDetermnisticDistributions([]float64{10.0, 20.0, 30.0, 40.0})
	pd := paireddata.UncertaintyPairedData{Xvals: x, Yvals: y}

	df := DamageFunctionStochastic{}
	df.Source = "fabricated"
	df.DamageFunction = pd
	df.DamageDriver = hazards.Depth

	dffs.DamageFunctions[hazards.Default] = df
	dffs.DamageFunctions[hazards.Depth] = df
	dffs.DamageFunctions[hazards.Depth|hazards.ArrivalTime] = df

	b, err := json.Marshal(dffs)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(b))
	dffs2 := DamageFunctionFamilyStochastic{}
	err = json.Unmarshal(b, &dffs2)
	if err != nil {
		panic(err)
	}
	fmt.Println(dffs2.DamageFunctions[hazards.Default].Source)
}

func Test_OccupancyTypeStochastic_Marshal(t *testing.T) {
	ot := agr1()
	b, err := json.Marshal(ot)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(b))
	ot2 := OccupancyTypeStochastic{}
	err = json.Unmarshal(b, &ot2)
	if err != nil {
		panic(err)
	}
	fmt.Println(ot2.ComponentDamageFunctions["contents"].DamageFunctions[hazards.Default].Source)
}
func Test_JsonOcctypes_toFile(t *testing.T) {
	path := "/workspaces/Go_Consequences/data/occtypes.json"
	w, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		panic(err)
	}
	defer w.Close()
	//a map of occupancy types
	m := OccupancyTypeMap()
	c := OccupancyTypesContainer{OccupancyTypes: m}
	b, err := json.Marshal(c)
	if err != nil {
		panic(err)
	}
	_, err = w.Write(b)
	if err != nil {
		panic(err)
	}
}
