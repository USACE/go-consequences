package paireddata

import "testing"

func createTestData() PairedData {
	x := []float64{1.0, 2.0, 3.0, 4.0}
	y := []float64{10.0, 20.0, 30.0, 40.0}
	return PairedData{Xvals: x, Yvals: y}
}
func TestSampleValue_belowRange(t *testing.T) {
	//setup
	pd := createTestData()
	//perform action
	got := pd.SampleValue(0)
	//test result
	if got != 0.0 {
		//report error
		t.Errorf("SampleValue(0) = %f; expected 0.0", got)
	}
}
func TestSampleValue_belowRange_precision(t *testing.T) {
	//setup
	pd := createTestData()
	//perform action
	got := pd.SampleValue(0.99)
	//test result
	if got != 0.0 {
		//report error
		t.Errorf("SampleValue(0.99) = %f; expected 0.0", got)
	}
}
func TestSampleValue_lowestValue(t *testing.T) {
	//setup
	pd := createTestData()
	//perform action
	got := pd.SampleValue(1.0)
	//test result
	if got != 10.0 {
		//report error
		t.Errorf("SampleValue(1.0) = %f; expected 10.0", got)
	}
}
func TestSampleValue_betweenOrdinates(t *testing.T) {
	//setup
	pd := createTestData()
	//perform action
	got := pd.SampleValue(2.5)
	//test result
	if got != 25.0 {
		//report error
		t.Errorf("SampleValue(2.5) = %f; expected 25.0", got)
	}
}
func TestSampleValue_highestValue(t *testing.T) {
	//setup
	pd := createTestData()
	//perform action
	got := pd.SampleValue(4.0)
	//test result
	if got != 40.0 {
		//report error
		t.Errorf("SampleValue(4.0) = %f; expected 40.0", got)
	}
}
func TestSampleValue_aboveRange(t *testing.T) {
	//setup
	pd := createTestData()
	//perform action
	got := pd.SampleValue(4.1)
	//test result
	if got != 40.0 {
		//report error
		t.Errorf("SampleValue(4.1) = %f; expected 40.0", got)
	}
}
