package main

import (
	"fmt"

	"github.com/USACE/go-consequences/consequences"
	"github.com/USACE/go-consequences/hazards"
)

func main() {

	var s = consequences.BaseStructure()
	var d = hazards.DepthEvent{Depth: 3.0}

	depths := []float64{0.0, 0.5, 1.0, 1.0001, 2.25}
	for idx := range depths {
		d.Depth = depths[idx]
		fmt.Println("for a depth of", d.Depth, s.ComputeConsequences(d))
	}
	//simplified compute
	ret := s.ComputeConsequences(d)
	fmt.Println("for a depth of", d.Depth, ret)

	d.Depth = 0.0 // test lower case
	ret = s.ComputeConsequences(d)
	fmt.Println("for a depth of", d.Depth, ret)

	d.Depth = .5 // should return 0
	ret = s.ComputeConsequences(d)
	fmt.Println("for a depth of", d.Depth, ret)

	d.Depth = 1.0 // test lowest valid case
	ret = s.ComputeConsequences(d)
	fmt.Println("for a depth of", d.Depth, ret)

	d.Depth = 1.0001 // test lowest interp case
	ret = s.ComputeConsequences(d)
	fmt.Println("for a depth of", d.Depth, ret)

	d.Depth = 2.25 //test interpolation case
	ret = s.ComputeConsequences(d)
	fmt.Println("for a depth of", d.Depth, ret)

	d.Depth = 2.5 //test interpolation case
	ret = s.ComputeConsequences(d)
	fmt.Println("for a depth of", d.Depth, ret)

	d.Depth = 2.75 //test interpolation case
	ret = s.ComputeConsequences(d)
	fmt.Println("for a depth of", d.Depth, ret)

	d.Depth = 3.99 // test highest interp case
	ret = s.ComputeConsequences(d)
	fmt.Println("for a depth of", d.Depth, ret)

	d.Depth = 4.0 // test highest valid case
	ret = s.ComputeConsequences(d)
	fmt.Println("for a depth of", d.Depth, ret)

	d.Depth = 5.0 //test upper case
	ret = s.ComputeConsequences(d)
	fmt.Println("for a depth of", d.Depth, ret)

	s.FoundHt = 1.1 //test interpolation due to foundation height putting depth back in range
	ret = s.ComputeConsequences(d)
	fmt.Println("for a depth of", d.Depth, ret)

	var f = hazards.FireEvent{Intensity: hazards.Low}
	s = consequences.ConvertBaseStructureToFire(s)
	ret = s.ComputeConsequences(f)
	fmt.Println("for a fire intensity of", f.Intensity, ret)

	f.Intensity = hazards.Medium
	ret = s.ComputeConsequences(f)
	fmt.Println("for a fire intensity of", f.Intensity, ret)

	f.Intensity = hazards.High
	ret = s.ComputeConsequences(f)
	fmt.Println("for a fire intensity of", f.Intensity, ret)
}
