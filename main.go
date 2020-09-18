package main

import (
	"fmt"

	"github.com/USACE/go-consequences/consequences"
	"github.com/USACE/go-consequences/hazards"
	"github.com/USACE/go-consequences/nsi"
)

func main() {

	var s = consequences.BaseStructure()
	var d = hazards.DepthEvent{Depth: 3.0}
	depths := []float64{3.0, 0.0, 0.5, 1.0, 1.0001, 2.25, 2.5, 2.75, 3.99, 4, 5}
	for idx := range depths {
		d.Depth = depths[idx]
		fmt.Println("for a depth of", d.Depth, s.ComputeConsequences(d))
	}

	s.FoundHt = 1.1 //test interpolation due to foundation height putting depth back in range
	ret := s.ComputeConsequences(d)
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

	var bbox string = "-81.58418,30.25165,-81.58161,30.26939,-81.55898,30.26939,-81.55281,30.24998,-81.58418,30.25165"
	structures := nsi.GetByBbox(bbox)

	for i, str := range structures {
		fmt.Println(i, "at structure", str.Name, "for a depth of", d.Depth, str.ComputeConsequences(d))
	}

}
