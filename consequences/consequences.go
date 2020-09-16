package consequences

import (
	"github.com/USACE/go-consequences/hazards"
	"github.com/USACE/go-consequences/paireddata"
)

type Structure struct {
	OccType                     OccupancyType
	DamCat                      string
	StructVal, ContVal, FoundHt float64
}

func (s Structure) ComputeConsequences(d interface{}) ConsequenceDamageResult {
	header := []string{"structure damage", "content damage"}
	results := []interface{}{0.0, 0.0}
	var ret = ConsequenceDamageResult{Headers: header, Results: results}
	de, ok := d.(hazards.DepthEvent)
	if ok {
		depth := de.Depth
		depthAboveFFE := depth - s.FoundHt
		damagePercent := s.OccType.Structuredamfun.SampleValue(depthAboveFFE) / 100 //assumes what type the damage array is in
		cdamagePercent := s.OccType.Contentdamfun.SampleValue(depthAboveFFE) / 100
		ret.Results[0] = damagePercent * s.StructVal
		ret.Results[1] = cdamagePercent * s.ContVal
		return ret
	}
	def, okd := d.(float64)
	if okd {
		depthAboveFFE := def - s.FoundHt
		damagePercent := s.OccType.Structuredamfun.SampleValue(depthAboveFFE) / 100 //assumes what type the damage array is in
		cdamagePercent := s.OccType.Contentdamfun.SampleValue(depthAboveFFE) / 100
		ret.Results[0] = damagePercent * s.StructVal
		ret.Results[1] = cdamagePercent * s.ContVal
		return ret
	}
	fire, okf := d.(hazards.FireEvent)
	if okf {
		damagePercent := s.OccType.Structuredamfun.SampleValue(fire.Intensity) / 100 //assumes what type the damage array is in
		cdamagePercent := s.OccType.Contentdamfun.SampleValue(fire.Intensity) / 100
		ret.Results[0] = damagePercent * s.StructVal
		ret.Results[1] = cdamagePercent * s.ContVal
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
	var o = OccupancyType{Name: "test", Structuredamfun: dfun, Contentdamfun: cdfun}
	var s = Structure{OccType: o, DamCat: "category", StructVal: 100.0, ContVal: 10.0, FoundHt: 0.0}
	return s
}
func ConvertBaseStructureToFire(s Structure) Structure {
	var fire = hazards.FireDamageFunction{}
	s.OccType.Structuredamfun = fire
	s.OccType.Contentdamfun = fire
	return s
}
