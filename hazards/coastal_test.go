package hazards_test

import (
	"testing"

	"github.com/USACE/go-consequences/hazards"
)

func TestCoastal_With_Salinity(t *testing.T) {
	d := hazards.CoastalEvent{}
	d.SetDepth(2.5)
	d.SetSalinity(true)

	if !d.Has(hazards.Salinity) {
		t.Error("Expected Salinity, but reported none.")
	}
}

func TestCoastal_With_Wave_NoSalt(t *testing.T) {
	d := hazards.CoastalEvent{}
	d.SetDepth(2.5)
	d.SetWaveHeight(3.3)
	if !d.Has(hazards.WaveHeight) {
		t.Error("Expected Wave, but reported none.")
	} else {
		if d.WaveHeight() != 3.3 {
			t.Error("Expected WaveHeight of 3.3, but got something else.")
		}
	}
}

func TestCoastal_With_Wave_With_Salt(t *testing.T) {
	d := hazards.CoastalEvent{}
	d.SetDepth(2.5)
	d.SetWaveHeight(3.3)
	d.SetSalinity(true)
	if !d.Has(hazards.Salinity) {
		t.Error("Expected Salinity, but reported none.")
	}
	if !d.Has(hazards.WaveHeight) {
		t.Error("Expected Wave, but reported none.")
	} else {
		if d.WaveHeight() != 3.3 {
			t.Error("Expected WaveHeight of 3.3, but got something else.")
		}
	}
}

func Test_CoastalWithErosion(t *testing.T) {
	c := hazards.CoastalEvent{}
	c.SetErosion(20)
	d := hazards.NewCoastalEvent(c)

	if d.Has(hazards.WaveHeight) {
		t.Error("Did not expected Wave.")
	}

	if d.Has(hazards.Salinity) {
		t.Error("Did not expected Salinity.")
	}

	if d.Has(hazards.Depth) {
		t.Error("Did not expected Depth.")
	}

	if !d.Has(hazards.Erosion) {
		t.Error("Expected Erosion, but reported none.")
	}

	if d.Erosion() != 20 {
		t.Error("Expected PercentEroded of 20, but got something else.")
	}
}

func TestCoastalMultiFrequency(t *testing.T) {
	depths := []float64{1, 10, 30, 45, 59, 78, 89, 102, 140, 180, 240, 330, 350, 370}
	waveheights := []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14}
	freqs := []float64{.99, .95, .9, .8, .7, .6, .5, .4, .3, .2, .1, .01, .002, .001}

	Events := make([]hazards.CoastalFrequencyEvent, len(depths))
	for i, d := range depths {
		ei := hazards.CoastalFrequencyEvent{}
		ei.SetDepth(d)
		ei.SetWaveHeight(waveheights[i])
		ei.SetFrequency(freqs[i])
		Events[i] = ei
	}

	dmf := hazards.MultiFrequencyCoastalEvent{Events: Events}

	testAllFreqs := dmf.HazardFrequencies()
	testFreqs := make([]float64, len(freqs))
	testDepths := make([]float64, len(depths))
	for {
		testFreqs[dmf.Index()] = dmf.Frequency()
		testDepths[dmf.Index()] = dmf.Depth()
		if dmf.HasNext() {
			dmf.Increment()
		} else {
			if !(dmf.Index() == len(Events)-1) {
				t.Errorf("Length mismatch between input hazards and output hazards")
			}
			break
		}
	}

	for i := range testFreqs {
		if testFreqs[i] != freqs[i] {
			t.Errorf("Wrong Frequency from Frequency(). Expected: %v. Got: %v\n", testFreqs[i], freqs[i])
		}
		if testAllFreqs[i] != freqs[i] {
			t.Errorf("Wrong Frequency from HazardFrequencies(). Expected: %v. Got: %v\n", testFreqs[i], freqs[i])
		}
		if testDepths[i] != depths[i] {
			t.Errorf("Wrong Depth from Depth(). Expected: %v. Got: %v\n", testDepths[i], depths[i])
		}
	}
}
