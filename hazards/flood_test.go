package hazards

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"
)

func TestDepth(t *testing.T) {
	d := DepthEvent{}
	d.SetDepth(2.5)
	if d.Depth() != 2.5 {
		t.Errorf("Expected %f, got %f", 2.5, d.Depth())
	}
}
func TestArrivalandDurationEvent(t *testing.T) {
	d := ArrivalandDurationEvent{}
	d.SetDuration(2.5)
	at := time.Date(1984, time.Month(1), 22, 0, 0, 0, 0, time.UTC)
	d.SetArrivalTime(at)
	if d.Duration() != 2.5 {
		t.Errorf("Expected %f, got %f", 2.5, d.Duration())
	}
	if d.ArrivalTime() != at {
		t.Errorf("Expected %s, got %s", at, d.ArrivalTime())
	}
}
func TestArrivalDepthandDurationEvent(t *testing.T) {
	d := ArrivalDepthandDurationEvent{}
	d.SetDuration(2.5)
	d.SetDepth(5)
	at := time.Date(1984, time.Month(1), 22, 0, 0, 0, 0, time.UTC)
	d.SetArrivalTime(at)
	if d.Depth() != 5 {
		t.Errorf("Expected %f, got %f", 5.0, d.Depth())
	}
	if d.Duration() != 2.5 {
		t.Errorf("Expected %f, got %f", 2.5, d.Duration())
	}
	if d.ArrivalTime() != at {
		t.Errorf("Expected %s, got %s", at, d.ArrivalTime())
	}
	s, _ := d.MarshalJSON()
	fmt.Printf(string(s) + "\n")
}
func TestDepthEventParameters(t *testing.T) {
	d := DepthEvent{depth: 2.5}
	if d.Has(Depth) {
		fmt.Println("Depth")
	}
	if d.Has(ArrivalTime) {
		fmt.Println("Arrival Time")
	}
	if d.Has(Erosion) {
		fmt.Println("Erosion")
	}
	if d.Has(Duration) {
		fmt.Println("Duration")
	}
	if d.Has(Velocity) {
		fmt.Println("Velocity")
	}
}
func TestArrivalandDurationEventParameters(t *testing.T) {
	at := time.Date(1984, time.Month(1), 22, 0, 0, 0, 0, time.UTC)
	d := ArrivalandDurationEvent{arrivalTime: at, durationInDays: 180}
	if d.Has(Depth) {
		fmt.Println("Depth")
	}
	if d.Has(ArrivalTime) {
		fmt.Println("Arrival Time")
	}
	if d.Has(Erosion) {
		fmt.Println("Erosion")
	}
	if d.Has(Duration) {
		fmt.Println("Duration")
	}
	if d.Has(Velocity) {
		fmt.Println("Velocity")
	}
}
func TestMarshalJSON(t *testing.T) {
	d := DepthEvent{}
	d.SetDepth(2.5)
	b, _ := json.Marshal(d)
	fmt.Println(string(b))
}
func TestMarshalParameterJSON(t *testing.T) {
	d := Default
	b, _ := json.Marshal(d)
	fmt.Println(string(b))
}
