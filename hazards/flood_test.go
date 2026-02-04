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
	fmt.Printf("%v\n", string(s))
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

func TestADDMulti(t *testing.T) {
	// create a series of hazardEvents
	var d1 = hazards.ArrivalDepthandDurationEvent{}
	d1.SetDuration(0)
	d1.SetDepth(1.0)
	t1 := time.Date(1984, time.Month(1), 1, 0, 0, 0, 0, time.UTC)
	d1.SetArrivalTime(t1)

	var d2 = hazards.ArrivalDepthandDurationEvent{}
	d2.SetDuration(5.0)
	d2.SetDepth(1.0)
	t2 := time.Date(1984, time.Month(1), 11, 0, 0, 0, 0, time.UTC)
	d2.SetArrivalTime(t2)

	var d3 = hazards.ArrivalDepthandDurationEvent{}
	d3.SetDuration(0.0)
	d3.SetDepth(1.0)
	t3 := time.Date(1984, time.Month(1), 21, 0, 0, 0, 0, time.UTC)
	d3.SetArrivalTime(t3)

	var d4 = hazards.ArrivalDepthandDurationEvent{}
	d4.SetDuration(0.0)
	d4.SetDepth(2.0)
	t4 := time.Date(1985, time.Month(1), 1, 0, 0, 0, 0, time.UTC)
	d4.SetArrivalTime(t4)

	var d5 = hazards.ArrivalDepthandDurationEvent{}
	d5.SetDuration(0.0)
	d5.SetDepth(2.0)
	t5 := time.Date(1985, time.Month(1), 11, 0, 0, 0, 0, time.UTC)
	d5.SetArrivalTime(t5)

	events := []hazards.ArrivalDepthandDurationEvent{d1, d2, d3, d4, d5}
	addm := hazards.ArrivalDepthandDurationEventMulti{Events: events}

	depths := []float64{}
	expected := []float64{1.0, 1.0, 1.0, 2.0, 2.0}

	for {
		depths = append(depths, addm.Depth())
		if addm.HasNext() {
			addm.Increment()
		} else {
			break
		}
	}

	// check that we properly read all depths and iterated through the events
	for i := range expected {
		if expected[i] != depths[i] {
			t.Errorf("Expected: %v. Got: %v\n", expected[i], depths[i])
		}
	}

	// check that we can reset the index
	addm.ResetIndex()
	if addm.Index() != 0 {
		t.Errorf("Index not reset")
	}
	// check that we can read the Parameters for the events
	for {
		if !addm.Has(hazards.Depth) {
			t.Errorf("Event didn't have Depth")
		}
		if !addm.Has(hazards.Duration) {
			t.Errorf("Event didn't have Duration")
		}
		if !addm.Has(hazards.ArrivalTime) {
			t.Errorf("Event didn't have ArrivalTime")
		}
		if addm.HasNext() {
			addm.Increment()
		} else {
			break
		}
	}
}
