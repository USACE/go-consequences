package consequences

import (
	"github.com/HenryGeorgist/go-statistics/statistics"
	"github.com/USACE/go-consequences/hazards"
)

type Structure struct {
	Name                              string
	OccType                           OccupancyType
	DamCat                            string
	StructVal, ContVal, FoundHt, X, Y float64
}
type ParameterValue struct {
	Value interface{}
}

//SampleValue on a ParameterValue is intended to help set structure values content values and foundaiton heights to uncertain parameters - this is a first draft of this interaction.
func (p ParameterValue) SampleValue(input interface{}) float64 {

	pval, okf := p.Value.(float64) //if the ParameterValue.Value is a float - pass it on back.
	if okf {
		return pval
	}
	pvaldist, okd := p.Value.(statistics.ContinuousDistribution)
	if okd {
		inval, ok := input.(float64)
		if ok {
			return pvaldist.InvCDF(inval)
		}
	}

	return 0
}

func (s Structure) ComputeConsequences(d interface{}) ConsequenceDamageResult { //what if we invert this general model to hazard.damage(consequence receptor)
	header := []string{"structure damage", "content damage"}
	results := []interface{}{0.0, 0.0}
	var ret = ConsequenceDamageResult{Headers: header, Results: results}
	de, ok := d.(hazards.DepthEvent)
	if ok {
		depth := de.Depth
		return computeFloodConsequences(depth, s)
	}
	def, okd := d.(float64)
	if okd {
		return computeFloodConsequences(def, s)
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
func computeFloodConsequences(d float64, s Structure) ConsequenceDamageResult {
	header := []string{"structure damage", "content damage"}
	results := []interface{}{0.0, 0.0}
	var ret = ConsequenceDamageResult{Headers: header, Results: results}
	depthAboveFFE := d - s.FoundHt
	damagePercent := s.OccType.Structuredamfun.SampleValue(depthAboveFFE) / 100 //assumes what type the damage array is in
	cdamagePercent := s.OccType.Contentdamfun.SampleValue(depthAboveFFE) / 100
	ret.Results[0] = damagePercent * s.StructVal
	ret.Results[1] = cdamagePercent * s.ContVal
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
