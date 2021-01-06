package crops

import (
	"testing"
	"time"

	"github.com/USACE/go-consequences/hazards"
)

func TestComputeCropDamageCase_FloodedBeforePlanting(t *testing.T) {
	at := time.Date(1984, time.Month(1), 22, 0, 0, 0, 0, time.UTC)
	h := hazards.ArrivalandDurationEvent{ArrivalTime: at, DurationInDays: 180}
	st := time.Date(1984, time.Month(7), 22, 0, 0, 0, 0, time.UTC)
	et := time.Date(1984, time.Month(7), 28, 0, 0, 0, 0, time.UTC)

	cs := CropSchedule{StartPlantingDate: st, LastPlantingDate: et, DaysToMaturity: 100}
	cdc := cs.ComputeCropDamageCase(h)
	expected := NotImpactedDuringSeason
	if cdc != expected {
		t.Errorf("ComputeCropDamageCase() = %v; expected %v", cdc, expected)
	}
}
func TestComputeCropDamageCase_FloodPostponedPlanting(t *testing.T) {
	at := time.Date(1984, time.Month(1), 22, 0, 0, 0, 0, time.UTC)
	h := hazards.ArrivalandDurationEvent{ArrivalTime: at, DurationInDays: 7}
	st := time.Date(1984, time.Month(1), 25, 0, 0, 0, 0, time.UTC)
	et := time.Date(1984, time.Month(1), 31, 0, 0, 0, 0, time.UTC)

	cs := CropSchedule{StartPlantingDate: st, LastPlantingDate: et, DaysToMaturity: 100}
	cdc := cs.ComputeCropDamageCase(h)
	expected := PlantingDelayed
	if cdc != expected {
		t.Errorf("ComputeCropDamageCase() = %v; expected %v", cdc, expected)
	}
}
func TestComputeCropDamageCase_FloodNoPlant(t *testing.T) {
	at := time.Date(1984, time.Month(1), 22, 0, 0, 0, 0, time.UTC)
	h := hazards.ArrivalandDurationEvent{ArrivalTime: at, DurationInDays: 12}
	st := time.Date(1984, time.Month(1), 25, 0, 0, 0, 0, time.UTC)
	et := time.Date(1984, time.Month(1), 31, 0, 0, 0, 0, time.UTC)

	cs := CropSchedule{StartPlantingDate: st, LastPlantingDate: et, DaysToMaturity: 100}
	cdc := cs.ComputeCropDamageCase(h)
	expected := NotPlanted
	if cdc != expected {
		t.Errorf("ComputeCropDamageCase() = %v; expected %v", cdc, expected)
	}
}
func TestComputeCropDamageCase_FloodAfterPlanting(t *testing.T) {
	at := time.Date(1984, time.Month(2), 1, 0, 0, 0, 0, time.UTC)
	h := hazards.ArrivalandDurationEvent{ArrivalTime: at, DurationInDays: 12}
	st := time.Date(1984, time.Month(1), 25, 0, 0, 0, 0, time.UTC)
	et := time.Date(1984, time.Month(1), 31, 0, 0, 0, 0, time.UTC)

	cs := CropSchedule{StartPlantingDate: st, LastPlantingDate: et, DaysToMaturity: 100}
	cdc := cs.ComputeCropDamageCase(h)
	expected := Impacted
	if cdc != expected {
		t.Errorf("ComputeCropDamageCase() = %v; expected %v", cdc, expected)
	}
}

//need a winter crop example.
func TestComputeCropDamageCase_FloodAfterPlanting_WinterCrop(t *testing.T) {
	at := time.Date(1984, time.Month(2), 1, 0, 0, 0, 0, time.UTC)
	h := hazards.ArrivalandDurationEvent{ArrivalTime: at, DurationInDays: 12}
	st := time.Date(1983, time.Month(12), 25, 0, 0, 0, 0, time.UTC)
	et := time.Date(1983, time.Month(12), 31, 0, 0, 0, 0, time.UTC)

	cs := CropSchedule{StartPlantingDate: st, LastPlantingDate: et, DaysToMaturity: 100}
	cdc := cs.ComputeCropDamageCase(h)
	expected := Impacted
	if cdc != expected {
		t.Errorf("ComputeCropDamageCase() = %v; expected %v", cdc, expected)
	}
}
