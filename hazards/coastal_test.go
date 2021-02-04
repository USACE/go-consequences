package hazards

import (
	"testing"
)

func TestCoastal_With_Salinity(t *testing.T) {
	d := CoastalEvent{Depth: 2.5, Salinity: true}
	if !d.Has(Salinity) {
		t.Error("Expected Salinity, but reported none.")
	}
}
func TestCoastal_With_Wave_NoSalt(t *testing.T) {
	d := CoastalEvent{Depth: 2.5, WaveHeight: 3.3}
	if !d.Has(WaveHeight) {
		t.Error("Expected Wave, but reported none.")
	} else {
		if d.WaveHeight != 3.3 {
			t.Error("Expected WaveHeight of 3.3, but got something else.")
		}
	}
}
func TestCoastal_With_Wave_With_Salt(t *testing.T) {
	d := CoastalEvent{Depth: 2.5, WaveHeight: 3.3, Salinity: true}
	if !d.Has(Salinity) {
		t.Error("Expected Salinity, but reported none.")
	}
	if !d.Has(WaveHeight) {
		t.Error("Expected Wave, but reported none.")
	} else {
		if d.WaveHeight != 3.3 {
			t.Error("Expected WaveHeight of 3.3, but got something else.")
		}
	}
}
