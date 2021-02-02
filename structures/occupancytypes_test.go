package structures

import (
	"testing"

	"github.com/USACE/go-consequences/hazards"
	"github.com/USACE/go-consequences/paireddata"
)

func TestDamageFunctionFamily(t *testing.T) {
	m := make(map[hazards.Parameter]paireddata.ValueSampler)

	cep := hazards.Depth | hazards.Salinity | hazards.WaveHeight
	dep := hazards.Depth
	ce2p := hazards.Depth | hazards.WaveHeight

	cexs := []float64{1, 2, 3}
	ceys := []float64{1, 2, 3}
	var cedf = paireddata.PairedData{Xvals: cexs, Yvals: ceys}
	m[cep] = cedf

	dexs := []float64{1, 2, 3}
	deys := []float64{2, 4, 6}
	var dedf = paireddata.PairedData{Xvals: dexs, Yvals: deys}
	m[dep] = dedf

	ce2xs := []float64{1, 2, 3}
	ce2ys := []float64{3, 6, 9}
	var ce2df = paireddata.PairedData{Xvals: ce2xs, Yvals: ce2ys}
	m[ce2p] = ce2df
	var df = DamageFunctionFamily{DamageFunctions: m}
	ce := hazards.CoastalEvent{Depth: 2, Salinity: true, WaveHeight: 3.4}
	de := hazards.DepthEvent{Depth: 2}
	ce2 := hazards.CoastalEvent{Depth: 2, Salinity: false, WaveHeight: 3.4}
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
