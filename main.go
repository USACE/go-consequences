package main

import (
	"fmt"

	"github.com/USACE/go-consequences/damage"
)

func main() {

	var s = damage.BaseStructure()
	var d = damage.DepthEvent{Depth: 3.0}

	depths := []float64{0.0, 0.5, 1.0, 1.0001, 2.25}
	for idx := range depths {
		d.Depth = depths[idx]
		fmt.Println("for a depth of", d.Depth, s.ComputeConsequences(d))
	}

	// Update Foundation Height
	// s.foundHt = 1.1 //test interpolation due to foundation height putting depth back in range
	// ret = s.ComputeConsequences(d)
	// fmt.Println("for a depth of", d.depth, ret)

	// // Fire Tests
	// var f = fireEvent{intensity: low}
	// s = ConvertBaseStructureToFire(s)
	// ret = s.ComputeConsequences(f)
	// fmt.Println("for a fire intensity of", f.intensity, ret)

	// f = fireEvent{intensity: medium}
	// s = ConvertBaseStructureToFire(s)
	// ret = s.ComputeConsequences(f)
	// fmt.Println("for a fire intensity of", f.intensity, ret)

	// f = fireEvent{intensity: high}
	// s = ConvertBaseStructureToFire(s)
	// ret = s.ComputeConsequences(f)
	// fmt.Println("for a fire intensity of", f.intensity, ret)
}
