package hazards

import (
	"fmt"
	"testing"
	"time"
)

func TestDepthEventParameters(t *testing.T) {
	d := DepthEvent{depth: 2.5}
	if d.Has(Depth) {
		fmt.Println("Depth")
	}
	if d.Has(ArrivalTime) {
		fmt.Println("Arrival Time")
	}
	if d.Has(ArrivalTime2ft) {
		fmt.Println("Arrival Time 2ft")
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
	if d.Has(ArrivalTime2ft) {
		fmt.Println("Arrival Time 2ft")
	}
	if d.Has(Duration) {
		fmt.Println("Duration")
	}
	if d.Has(Velocity) {
		fmt.Println("Velocity")
	}
}
