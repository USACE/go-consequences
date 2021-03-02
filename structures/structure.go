package structures

import (
	"math/rand"

	"github.com/USACE/go-consequences/consequences"
	"github.com/USACE/go-consequences/hazards"
)

//BaseStructure represents a Structure name xy location and a damage category
type BaseStructure struct {
	Name   string
	DamCat string
	X, Y   float64
}

//StructureStochastic is a base structure with an occupancy type stochastic and parameter values for all parameters
type StructureStochastic struct {
	BaseStructure
	UseUncertainty              bool //defaults to false!
	OccType                     OccupancyTypeStochastic
	StructVal, ContVal, FoundHt consequences.ParameterValue
}

//StructureDeterministic is a base strucure with a deterministic occupancy type and deterministic parameters
type StructureDeterministic struct {
	BaseStructure
	OccType                     OccupancyTypeDeterministic
	StructVal, ContVal, FoundHt float64
}

//GetX implements consequences.Locatable
func (s BaseStructure) GetX() float64 {
	return s.X
}

//GetY implements consequences.Locatable
func (s BaseStructure) GetY() float64 {
	return s.Y
}

//SampleStructure converts a structureStochastic into a structure deterministic based on an input seed
func (s StructureStochastic) SampleStructure(seed int64) StructureDeterministic {
	ot := OccupancyTypeDeterministic{} //Beware null errors!
	sv := 0.0
	cv := 0.0
	fh := 0.0
	if s.UseUncertainty {
		ot = s.OccType.SampleOccupancyType(seed)
		sv = s.StructVal.SampleValue(rand.Float64())
		cv = s.ContVal.SampleValue(rand.Float64())
		fh = s.FoundHt.SampleValue(rand.Float64())
	} else {
		ot = s.OccType.CentralTendency()
		sv = s.StructVal.CentralTendency()
		cv = s.ContVal.CentralTendency()
		fh = s.FoundHt.CentralTendency()
	}

	return StructureDeterministic{OccType: ot, StructVal: sv, ContVal: cv, FoundHt: fh, BaseStructure: BaseStructure{DamCat: s.DamCat}}
}

//Compute implements the consequences.Receptor interface on StrucutreStochastic
func (s StructureStochastic) Compute(d hazards.HazardEvent) consequences.Results {
	return s.SampleStructure(rand.Int63()).Compute(d) //this needs work so seeds can be controlled.
}

//Compute implements the consequences.Receptor interface on StrucutreDeterminstic
func (s StructureDeterministic) Compute(d hazards.HazardEvent) consequences.Results { //what if we invert this general model to hazard.damage(consequence receptor)
	return computeConsequences(d, s)
}

//the following two methods are legitimately the same - it seems i need an interface rather than a struct for a depthevent
//this area seems still in need of some refactoring for simplification.

func computeConsequences(e hazards.HazardEvent, s StructureDeterministic) consequences.Results {
	header := []string{"structure damage", "content damage"}
	results := []interface{}{0.0, 0.0}
	var ret = consequences.Result{Headers: header, Result: results}
	if e.Has(hazards.Depth) {
		depthAboveFFE := e.Depth() - s.FoundHt
		damagePercent := s.OccType.GetStructureDamageFunctionForHazard(e).SampleValue(depthAboveFFE) / 100 //assumes what type the damage array is in
		cdamagePercent := s.OccType.GetContentDamageFunctionForHazard(e).SampleValue(depthAboveFFE) / 100
		ret.Result[0] = damagePercent * s.StructVal
		ret.Result[1] = cdamagePercent * s.ContVal
	} else {
		ret.Result[0] = 0.0
		ret.Result[1] = 0.0
	}
	r := consequences.Results{IsTable: false, Result: ret}
	return r
}
