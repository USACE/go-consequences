package consequences

import (
	"github.com/USACE/go-consequences/hazards"
)

type occupancyType struct {
	name            string
	structuredamfun ValueSampler
	contentdamfun   ValueSampler
}

type Structure struct {
	occType                     occupancyType
	damCat                      string
	structVal, contVal, foundHt float64
}

func (s Structure) ComputeConsequences(d interface{}) ConsequenceDamageResult {
	header := []string{"structure damage", "content damage"}
	results := []interface{}{0.0, 0.0}
	var ret = ConsequenceDamageResult{Headers: header, Results: results}
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
