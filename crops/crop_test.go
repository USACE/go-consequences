package crops

import (
	"testing"
	"time"

	"github.com/USACE/go-consequences/hazards"
)

func TestComputeCropDamage_FloodedBeforePlanting(t *testing.T) {
	//setup
	//hazard definition
	at := time.Date(1984, time.Month(1), 22, 0, 0, 0, 0, time.UTC)
	h := hazards.ArrivalandDurationEvent{ArrivalTime: at, DurationInDays: 10}
	//construct a TestCrop
	c := createTestCrop()

	//compute
	cd := c.ComputeConsequences(h)
	//expected results
	expectedcase := NotImpactedDuringSeason
	expecteddamage := 0.0

	//test
	if cd.Result.Result[1] != expectedcase {
		t.Errorf("ComputeConsequence() = %v; expected %v", cd.Result.Result[1], expectedcase)
	}
	if cd.Result.Result[2] != expecteddamage {
		t.Errorf("ComputeConsequence() = %v; expected %v", cd.Result.Result[2], expecteddamage)
	}
}

func TestComputeCropDamage_FloodedAfterPlanting(t *testing.T) {
	//setup
	//hazard definition
	at := time.Date(1984, time.Month(7), 29, 0, 0, 0, 0, time.UTC)
	h := hazards.ArrivalandDurationEvent{ArrivalTime: at, DurationInDays: 10}
	//construct a TestCrop
	c := createTestCrop()

	//compute
	cd := c.ComputeConsequences(h)
	//expected results
	expectedcase := Impacted
	expecteddamage := 10.0 //temporary value for testing

	//test
	if cd.Result.Result[1] != expectedcase {
		t.Errorf("ComputeConsequence() = %v; expected %v", cd.Result.Result[1], expectedcase)
	}
	if cd.Result.Result[2] != expecteddamage {
		t.Errorf("ComputeConsequence() = %v; expected %v", cd.Result.Result[2], expecteddamage)
	}
}
func TestReadFromXML(t *testing.T) {
	//"C:\\Temp\\agtesting\\Corn.crop"
	path := "C:\\Temp\\agtesting\\Corn.crop"
	c := ReadFromXML(path)
	if c.GetCropName() != "Corn" {
		t.Error("Did not parse corn")
	}
}
func createTestCrop() Crop {
	//Crop Schedule
	st := time.Date(1984, time.Month(7), 22, 0, 0, 0, 0, time.UTC)
	et := time.Date(1984, time.Month(7), 28, 0, 0, 0, 0, time.UTC)
	cs := CropSchedule{StartPlantingDate: st, LastPlantingDate: et, DaysToMaturity: 100}

	//Production Function
	mc := []float64{1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0}
	pf := NewProductionFunction(mc, cs, 1.0, .1)

	//Damage Function
	one := []float64{1.1, 2.1, 3.1, 4.1, 5.1, 6.1, 7.1, 8.1, 9.1, 10.1, 11.1, 12.1}
	two := []float64{1.2, 2.2, 3.2, 4.2, 5.2, 6.2, 7.2, 8.2, 9.2, 10.2, 11.2, 12.2}
	three := []float64{1.3, 2.3, 3.3, 4.3, 5.3, 6.3, 7.3, 8.3, 9.3, 10.3, 11.3, 12.3}
	four := []float64{1.4, 2.4, 3.4, 4.4, 5.4, 6.4, 7.4, 8.4, 9.4, 10.4, 11.4, 12.4}

	m := make(map[float64][]float64)
	m[1.0] = one
	m[2.0] = two
	m[3.0] = three
	m[4.0] = four

	df := DamageFunction{DurationDamageCurves: m}

	//Crop
	c := BuildCrop(1, "Corn")
	//omit location
	c = c.WithOutput(100.00, 5.24) //yeild and price
	c = c.WithProductionFunction(pf)
	c = c.WithLossFunction(df)
	c = c.WithCropSchedule(cs)
	return c
}
