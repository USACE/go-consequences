package hazards

import (
	"fmt"
	"testing"
)

func TestParameters(t *testing.T) {
	d := DepthEvent{Depth: 2.5, parameter: Default}
	//p := Depth
	d.parameter = SetHasDepth(d.parameter)
	d.parameter = SetHasVelocity(d.parameter)
	//d.parameter = p
	fmt.Println(d.parameter.String())
	if d.Has(Depth) {
		fmt.Println("Has Depth")
	} else {
		fmt.Println("Hasnt depth")
	}
	if d.Has(ArrivalTime) {
		fmt.Println("Has Arrival Time")
	} else {
		fmt.Println("Hasnt arrival time")
	}
	if d.Has(ArrivalTime2ft) {
		fmt.Println("Has Arrival Time 2ft")
	} else {
		fmt.Println("Hasnt arrival time 2ft")
	}
	if d.Has(Velocity) {
		fmt.Println("Has Velocity")
	} else {
		fmt.Println("Hasnt Velocity")
	}
}
