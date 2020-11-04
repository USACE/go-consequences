package hazards

import (
	"fmt"
	"testing"
)

func TestParameters(t *testing.T) {
	d := DepthEvent{Depth: 2.5, parameter: Depth}
	//p := Depth
	//d.parameter = p
	if d.Has(Depth) {
		fmt.Println("True")
	} else {
		fmt.Println("False!!")
	}
	if d.Has(ArrivalTime) {
		fmt.Println("False!!!")
	} else {
		fmt.Println("True!!")
	}
}
