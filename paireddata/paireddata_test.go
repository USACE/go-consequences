package paireddata

import "testing"

func TestSampleValue(t *testing.T) {
	x := []float64{1.0, 2.0, 3.0, 4.0}
	y := []float64{10.0, 20.0, 30.0, 40.0}
	pd := PairedData{Xvals: x, Yvals: y}

	got := pd.SampleValue(0)
	if got != 0.0 {
		t.Errorf("SampleValue(0) = %f; expected 0.0", got)
	}

	got = pd.SampleValue(0.99)
	if got != 0.0 {
		t.Errorf("SampleValue(0.99) = %f; expected 0.0", got)
	}

	got = pd.SampleValue(1.0)
	if got != 10.0 {
		t.Errorf("SampleValue(1.0) = %f; expected 10.0", got)
	}

	got = pd.SampleValue(2.5)
	if got != 25.0 {
		t.Errorf("SampleValue(2.5) = %f; expected 25.0", got)
	}

	got = pd.SampleValue(4.0)
	if got != 40.0 {
		t.Errorf("SampleValue(4.0) = %f; expected 40.0", got)
	}

	got = pd.SampleValue(4.1)
	if got != 40.0 {
		t.Errorf("SampleValue(4.1) = %f; expected 40.0", got)
	}
}
