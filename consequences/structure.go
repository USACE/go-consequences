package consequences

import (
	"github.com/USACE/go-consequences/hazards"
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
	//get the occupancy type map
	m := OccupancyTypeMap()
	// select a base structure type for testing
	var o = m["RES1-1SNB"]
	var s = Structure{OccType: o, DamCat: "category", StructVal: 100.0, ContVal: 10.0, FoundHt: 0.0}
	return s
}
func ConvertBaseStructureToFire(s Structure) Structure {
	var fire = hazards.FireDamageFunction{}
	s.OccType.Structuredamfun = fire
	s.OccType.Contentdamfun = fire
	return s
}