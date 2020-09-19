package hazards

import "testing"

func TestSampleValue(t *testing.T) {
	fdf := FireDamageFunction{}
	var f = FireEvent{Intensity: Low}
	got := fdf.SampleValue(f.Intensity)
	if got != 33.3 {
		t.Errorf("SampleValue(Low) = %f; expected 33.3", got)
	}

	f.Intensity = Medium
	got = fdf.SampleValue(f.Intensity)
	if got != 50.0 {
		t.Errorf("SampleValue(Medium) = %f; expected 50.0", got)
	}

	f.Intensity = High
	got = fdf.SampleValue(f.Intensity)
	if got != 100.0 {
		t.Errorf("SampleValue(High) = %f; expected 100.0", got)
	}

	got = fdf.SampleValue(1.5)
	if got != 0.0 {
		t.Errorf("SampleValue(1.5) = %f; expected 0.0", got)
	}
}
