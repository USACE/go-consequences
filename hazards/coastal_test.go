package hazards

import (
	"testing"
)

func TestCoastal_With_Salinity(t *testing.T) {
	d := CoastalEvent{depth: 2.5, salinity: true}
	if !d.Has(Salinity) {
		t.Error("Expected Salinity, but reported none.")
	}
}

func TestCoastal_With_Wave_NoSalt(t *testing.T) {
	d := CoastalEvent{depth: 2.5, waveHeight: 3.3}
	if !d.Has(WaveHeight) {
		t.Error("Expected Wave, but reported none.")
	} else {
		if d.WaveHeight() != 3.3 {
			t.Error("Expected WaveHeight of 3.3, but got something else.")
		}
	}
}

func TestCoastal_With_Wave_With_Salt(t *testing.T) {
	d := CoastalEvent{depth: 2.5, waveHeight: 3.3, salinity: true}
	if !d.Has(Salinity) {
		t.Error("Expected Salinity, but reported none.")
	}
	if !d.Has(WaveHeight) {
		t.Error("Expected Wave, but reported none.")
	} else {
		if d.WaveHeight() != 3.3 {
			t.Error("Expected WaveHeight of 3.3, but got something else.")
		}
	}
}

func Test_CoastalWithErosion(t *testing.T) {
	d := NewCoastalEvent(CoastalEvent{percentEroded: 20})

	if d.Has(WaveHeight) {
		t.Error("Did not expected Wave.")
	}

	if d.Has(Salinity) {
		t.Error("Did not expected Salinity.")
	}

	if d.Has(Depth) {
		t.Error("Did not expected Depth.")
	}

	if !d.Has(Erosion) {
		t.Error("Expected Erosion, but reported none.")
	}

	if d.Erosion() != 20 {
		t.Error("Expected PercentEroded of 20, but got something else.")
	}
}
