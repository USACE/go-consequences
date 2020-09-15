package main

import (
	"fmt"

	"github.com/USACE/go-consequences/consequences"
	"github.com/USACE/go-consequences/hazards"
	"github.com/USACE/go-consequences/paireddata"
)

type occupancyType struct {
	name            string
	structuredamfun consequences.ValueSampler
	contentdamfun   consequences.ValueSampler
}
type fireDamageFunction struct {
}

type Structure struct {
	occType                     occupancyType
	damCat                      string
	structVal, contVal, foundHt float64
}

func (s Structure) ComputeConsequences(d interface{}) consequences.ConsequenceDamageResult {
	header := []string{"structure damage", "content damage"}
	results := []interface{}{0.0, 0.0}
	var ret = consequences.ConsequenceDamageResult{Headers: header, Results: results}
	de, ok := d.(hazards.DepthEvent)
	if ok {
		depth := de.Depth
		depthAboveFFE := depth - s.foundHt
		damagePercent := s.occType.structuredamfun.SampleValue(depthAboveFFE) / 100 //assumes what type the damage array is in
		cdamagePercent := s.occType.contentdamfun.SampleValue(depthAboveFFE) / 100
		ret.Results[0] = damagePercent * s.structVal
		ret.Results[1] = cdamagePercent * s.contVal
		return ret
	}
	def, okd := d.(float64)
	if okd {
		depthAboveFFE := def - s.foundHt
		damagePercent := s.occType.structuredamfun.SampleValue(depthAboveFFE) / 100 //assumes what type the damage array is in
		cdamagePercent := s.occType.contentdamfun.SampleValue(depthAboveFFE) / 100
		ret.Results[0] = damagePercent * s.structVal
		ret.Results[1] = cdamagePercent * s.contVal
		return ret
	}
	fire, okf := d.(hazards.FireEvent)
	if okf {
		damagePercent := s.occType.structuredamfun.SampleValue(fire.Intensity) / 100 //assumes what type the damage array is in
		cdamagePercent := s.occType.contentdamfun.SampleValue(fire.Intensity) / 100
		ret.Results[0] = damagePercent * s.structVal
		ret.Results[1] = cdamagePercent * s.contVal
		return ret
	}
	return ret
}
func BaseStructure() Structure {
	//fake data to test
	xs := []float64{1.0, 2.0, 3.0, 4.0}
	ys := []float64{10.0, 20.0, 30.0, 40.0}
	cxs := []float64{1.0, 2.0, 3.0, 4.0}
	cys := []float64{5.0, 10.0, 15.0, 20.0}
	var dfun = paireddata.PairedData{Xvals: xs, Yvals: ys}
	var cdfun = paireddata.PairedData{Xvals: cxs, Yvals: cys}
	var o = occupancyType{name: "test", structuredamfun: dfun, contentdamfun: cdfun}
	var s = Structure{occType: o, damCat: "category", structVal: 100.0, contVal: 10.0, foundHt: 0.0}
	return s
}
func ConvertBaseStructureToFire(s Structure) Structure {
	var fire = hazards.FireDamageFunction{}
	s.occType.structuredamfun = fire
	s.occType.contentdamfun = fire
	return s
}
func main() {

	var s = BaseStructure()
	var d = hazards.DepthEvent{Depth: 3.0}

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

	s.foundHt = 1.1 //test interpolation due to foundation height putting depth back in range
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
