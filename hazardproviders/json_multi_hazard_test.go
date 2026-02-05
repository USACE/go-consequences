package hazardproviders

import (
	"fmt"
	"testing"

	"github.com/USACE/go-consequences/geography"
)

func TestInitADDMHP(t *testing.T) {
	file := "/workspaces/go-consequences/data/lifecycle/test_arrival-depth-duration_hazards.json"

	ADDMHP, err := InitADDMHP(file)
	if err != nil {
		panic(err)
	}

	loc := geography.Location{
		X:    -71.481,
		Y:    43.001,
		SRID: "",
	}

	h, err := ADDMHP.Hazard(loc)
	if err != nil {
		panic(err)
	}

	for {
		fmt.Printf(
			"%d: Depth: %3.2f, Duration: %3.2f, Arrival: %v\n",
			h.Index(), h.Depth(), h.Duration(), h.ArrivalTime(),
		)
		if h.HasNext() {
			h.Increment()
		} else {
			break
		}
	}
}
