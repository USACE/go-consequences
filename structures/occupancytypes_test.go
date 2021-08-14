package structures

import (
	"fmt"
	"testing"

	"github.com/USACE/go-consequences/hazards"
	"github.com/USACE/go-consequences/paireddata"
)

func TestDamageFunctionFamily(t *testing.T) {
	//a map of hazard to damage function
	m := make(map[hazards.Parameter]paireddata.ValueSampler)
	//a set of different hazard types.
	cep := hazards.Depth | hazards.Salinity | hazards.WaveHeight
	dep := hazards.Depth
	ce2p := hazards.Depth | hazards.WaveHeight

	//a fake deterministic damage function for coastal event with salinity
	cexs := []float64{1, 2, 3}
	ceys := []float64{1, 2, 3}
	var cedf = paireddata.PairedData{Xvals: cexs, Yvals: ceys}
	m[cep] = cedf

	//a fake deterministic damage function for a depth only event
	dexs := []float64{1, 2, 3}
	deys := []float64{2, 4, 6}
	var dedf = paireddata.PairedData{Xvals: dexs, Yvals: deys}
	m[dep] = dedf

	//a fake deterministic damage function for an event with depth and wave, but no salinity...
	ce2xs := []float64{1, 2, 3}
	ce2ys := []float64{3, 6, 9}
	var ce2df = paireddata.PairedData{Xvals: ce2xs, Yvals: ce2ys}
	m[ce2p] = ce2df

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
	cv := df.DamageFunctions[ce.Parameters()].SampleValue(ce.Depth)
	if cv != 2 {
		t.Errorf("Expected 2")
	}
	dv := df.DamageFunctions[de.Parameters()].SampleValue(de.Depth)
	if dv != 4 {
		t.Errorf("Expected 4")
	}
	c2v := df.DamageFunctions[ce2.Parameters()].SampleValue(ce2.Depth)
	if c2v != 6 {
		t.Errorf("Expected 6")
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
