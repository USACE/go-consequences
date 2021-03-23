package structures

import (
	"math/rand"

	"github.com/USACE/go-consequences/consequences"
	"github.com/USACE/go-consequences/geography"
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
	UseUncertainty                             bool //defaults to false!
	OccType                                    OccupancyTypeStochastic
	StructVal, ContVal, FoundHt                consequences.ParameterValue
	Pop2pmo65, Pop2pmu65, Pop2amo65, Pop2amu65 int32
}

//StructureDeterministic is a base strucure with a deterministic occupancy type and deterministic parameters
type StructureDeterministic struct {
	BaseStructure
	OccType                                    OccupancyTypeDeterministic
	StructVal, ContVal, FoundHt                float64
	Pop2pmo65, Pop2pmu65, Pop2amo65, Pop2amu65 int32
}

//GetX implements consequences.Locatable
func (s BaseStructure) Location() geography.Location {
	return geography.Location{X: s.X, Y: s.Y}
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

	return StructureDeterministic{
		OccType:       ot,
		StructVal:     sv,
		ContVal:       cv,
		FoundHt:       fh,
		Pop2pmo65:     s.Pop2pmo65,
		Pop2pmu65:     s.Pop2pmu65,
		Pop2amo65:     s.Pop2amo65,
		Pop2amu65:     s.Pop2amu65,
		BaseStructure: BaseStructure{Name: s.Name, X: s.X, Y: s.Y, DamCat: s.DamCat}}
}

//Compute implements the consequences.Receptor interface on StrucutreStochastic
func (s StructureStochastic) Compute(d hazards.HazardEvent) consequences.Result {
	return s.SampleStructure(rand.Int63()).Compute(d) //this needs work so seeds can be controlled.
}

//Compute implements the consequences.Receptor interface on StrucutreDeterminstic
func (s StructureDeterministic) Compute(d hazards.HazardEvent) consequences.Result { //what if we invert this general model to hazard.damage(consequence receptor)
	return computeConsequences(d, s)
}

//the following two methods are legitimately the same - it seems i need an interface rather than a struct for a depthevent
//this area seems still in need of some refactoring for simplification.

func computeConsequences(e hazards.HazardEvent, s StructureDeterministic) consequences.Result {
	//header := []string{"structure damage", "content damage"}
	header := []string{"fd_id", "x", "y", "hazard", "damage category", "occupancy type", "structure damage", "content damage", "pop2amu65", "pop2amo65", "pop2pmu65", "pop2pmo65"}
	results := []interface{}{"updateme", 0.0, 0.0, e, "dc", "ot", 0.0, 0.0, 0, 0, 0, 0}
	var ret = consequences.Result{Headers: header, Result: results}
	if e.Has(hazards.Depth) {
		depthAboveFFE := e.Depth() - s.FoundHt
		damagePercent := s.OccType.GetStructureDamageFunctionForHazard(e).SampleValue(depthAboveFFE) / 100 //assumes what type the damage array is in
		cdamagePercent := s.OccType.GetContentDamageFunctionForHazard(e).SampleValue(depthAboveFFE) / 100
		ret.Result[0] = s.BaseStructure.Name
		ret.Result[1] = s.BaseStructure.X
		ret.Result[2] = s.BaseStructure.Y
		ret.Result[3] = e
		ret.Result[4] = s.BaseStructure.DamCat
		ret.Result[5] = s.OccType.Name
		ret.Result[6] = damagePercent * s.StructVal
		ret.Result[7] = cdamagePercent * s.ContVal
		ret.Result[8] = s.Pop2amu65
		ret.Result[9] = s.Pop2amo65
		ret.Result[10] = s.Pop2pmu65
		ret.Result[11] = s.Pop2pmo65
	}
	return ret
}
