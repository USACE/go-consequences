package hazards_test

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/USACE/go-consequences/hazards"
)

func ExampleDepthEvent() {
	d := hazards.DepthEvent{}
	d.SetDepth(2.5)
	fmt.Println(d.Has(hazards.Depth))
	fmt.Println(d.Has(hazards.Velocity))
	fmt.Println(d.Depth())
	// Output:
	// true
	// false
	// 2.5
}
func TestDepth(t *testing.T) {
	d := hazards.DepthEvent{}
	d.SetDepth(2.5)
	if d.Depth() != 2.5 {
		t.Errorf("Expected %f, got %f", 2.5, d.Depth())
	}
}
func TestArrivalandDurationEvent(t *testing.T) {
	d := hazards.ArrivalandDurationEvent{}
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
	d := hazards.ArrivalDepthandDurationEvent{}
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
	d := hazards.DepthEvent{}
	d.SetDepth(2.5)
	if d.Has(hazards.Depth) {
		fmt.Println("Depth")
	}
	if d.Has(hazards.ArrivalTime) {
		fmt.Println("Arrival Time")
	}
	if d.Has(hazards.Erosion) {
		fmt.Println("Erosion")
	}
	if d.Has(hazards.Duration) {
		fmt.Println("Duration")
	}
	if d.Has(hazards.Velocity) {
		fmt.Println("Velocity")
	}
}
func TestArrivalandDurationEventParameters(t *testing.T) {
	at := time.Date(1984, time.Month(1), 22, 0, 0, 0, 0, time.UTC)
	d := hazards.ArrivalandDurationEvent{}
	d.SetArrivalTime(at)
	d.SetDuration(180)
	if d.Has(hazards.Depth) {
		fmt.Println("Depth")
	}
	if d.Has(hazards.ArrivalTime) {
		fmt.Println("Arrival Time")
	}
	if d.Has(hazards.Erosion) {
		fmt.Println("Erosion")
	}
	if d.Has(hazards.Duration) {
		fmt.Println("Duration")
	}
	if d.Has(hazards.Velocity) {
		fmt.Println("Velocity")
	}
}
func TestMarshalJSON(t *testing.T) {
	d := hazards.DepthEvent{}
	d.SetDepth(2.5)
	b, _ := json.Marshal(d)
	fmt.Println(string(b))
}
func TestMarshalParameterJSON(t *testing.T) {
	d := hazards.Default
	b, _ := json.Marshal(d)
	fmt.Println(string(b))
}
