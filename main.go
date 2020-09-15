package main

import (
	"fmt"
  
	"github.com/USACE/go-consequences/consequences"
	"github.com/USACE/go-consequences/hazards"
	"github.com/USACE/go-consequences/paireddata"
)

func BaseStructure() consequences.Structure {
	//fake data to test
	xs := []float64{1.0, 2.0, 3.0, 4.0}
	ys := []float64{10.0, 20.0, 30.0, 40.0}
	cxs := []float64{1.0, 2.0, 3.0, 4.0}
	cys := []float64{5.0, 10.0, 15.0, 20.0}
	var dfun = paireddata.PairedData{Xvals: xs, Yvals: ys}
	var cdfun = paireddata.PairedData{Xvals: cxs, Yvals: cys}
	var o = consequences.OccupancyType{Name: "test", Structuredamfun: dfun, Contentdamfun: cdfun}
	var s = consequences.Structure{OccType: o, DamCat: "category", StructVal: 100.0, ContVal: 10.0, FoundHt: 0.0}
	return s
}
func ConvertBaseStructureToFire(s consequences.Structure) consequences.Structure {
	var fire = hazards.FireDamageFunction{}
	s.OccType.Structuredamfun = fire
	s.OccType.Contentdamfun = fire
	return s
}
func main() {

	var s = BaseStructure()
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
	s = ConvertBaseStructureToFire(s)
	ret = s.ComputeConsequences(f)
	fmt.Println("for a fire intensity of", f.Intensity, ret)

	f.Intensity = hazards.Medium
	ret = s.ComputeConsequences(f)
	fmt.Println("for a fire intensity of", f.Intensity, ret)

	f.Intensity = hazards.High
	ret = s.ComputeConsequences(f)
	fmt.Println("for a fire intensity of", f.Intensity, ret)
}
