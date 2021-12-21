package paireddata

import (
	"encoding/json"
	"fmt"
	"testing"
)

func createTestData() PairedData {
	x := []float64{1.0, 2.0, 3.0, 4.0}
	y := []float64{10.0, 20.0, 30.0, 40.0}
	return PairedData{Xvals: x, Yvals: y}
}
func createUnhappyData() PairedData {
	x := []float64{1.0, 2.0, 3.0, 4.0}
	y := []float64{-10.0, 20.0, -30.0, 400.0}
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
func Test_Json(t *testing.T) {
	pd := createTestData()
	b, err := json.Marshal(pd)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(b))
	var pd2 PairedData
	err = json.Unmarshal(b, &pd2)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(pd2.SampleValue(3.0))
}
func Test_ForceNonNegativeMonotonic(t *testing.T) {
	pd := createUnhappyData()
	if pd.IsMonotonicallyIncreasing() {
		t.Error("say wha?")
	} else {
		pd.ForceMonotonicInRange(10.0, 500)
		if pd.IsMonotonicallyIncreasing() {
			for _, y := range pd.Yvals {
				fmt.Println(y)
			}
		} else {
			t.Error("say wha? take 2")
		}
	}

}
func Test_Compose(t *testing.T) {
	g := createTestData()
	y := []float64{8.0, 12.0, 16.0, 20.0}
	x := []float64{10.0, 20.0, 30.0, 40.0}
	f := PairedData{Yvals: y, Xvals: x}
	fog := f.Compose(g)
	fmt.Println(fog)
}
