package structures

import (
	"errors"
	"math"
	"math/rand"

	"github.com/USACE/go-consequences/consequences"
	"github.com/USACE/go-consequences/geography"
	"github.com/USACE/go-consequences/hazards"
)

// BaseStructure represents a Structure name xy location and a damage category
type BaseStructure struct {
	Name                  string
	DamCat                string
	CBFips                string
	X, Y, GroundElevation float64
}

// StructureStochastic is a base structure with an occupancy type stochastic and parameter values for all parameters
type StructureStochastic struct {
	BaseStructure
	UseUncertainty                                         bool //defaults to false!
	OccType                                                OccupancyTypeStochastic
	FoundType, FirmZone, ConstructionType                  string
	StructVal, ContVal, FoundHt                            consequences.ParameterValue
	Pop2pmo65, Pop2pmu65, Pop2amo65, Pop2amu65, NumStories int32
}

// StructureDeterministic is a base strucure with a deterministic occupancy type and deterministic parameters
type StructureDeterministic struct {
	BaseStructure
	OccType                                                OccupancyTypeDeterministic
	FoundType, ConstructionType                            string
	StructVal, ContVal, FoundHt                            float64
	Pop2pmo65, Pop2pmu65, Pop2amo65, Pop2amu65, NumStories int32
}

// GetX implements consequences.Locatable
func (s BaseStructure) Location() geography.Location {
	return geography.Location{X: s.X, Y: s.Y}
}

// SampleStructure converts a structureStochastic into a structure deterministic based on an input seed
func (s StructureStochastic) SampleStructure(seed int64) StructureDeterministic {
	r := rand.New(rand.NewSource(seed))
	ot := OccupancyTypeDeterministic{} //Beware null errors!
	sv := 0.0
	cv := 0.0
	fh := 0.0
	if s.UseUncertainty {
		ot = s.OccType.SampleOccupancyType(r.Int63()) //this is super inefficient. At the time this is called we know the hazard.
		sv = s.StructVal.SampleValue(r.Float64())
		cv = s.ContVal.SampleValue(r.Float64())
		fh = s.FoundHt.SampleValue(r.Float64())
	} else {
		ot = s.OccType.CentralTendency()
		sv = s.StructVal.CentralTendency()
		cv = s.ContVal.CentralTendency()
		fh = s.FoundHt.CentralTendency()
	}

	return StructureDeterministic{
		OccType:          ot,
		StructVal:        sv,
		ContVal:          cv,
		FoundType:        s.FoundType,
		ConstructionType: s.ConstructionType,
		FoundHt:          fh,
		Pop2pmo65:        s.Pop2pmo65,
		Pop2pmu65:        s.Pop2pmu65,
		Pop2amo65:        s.Pop2amo65,
		Pop2amu65:        s.Pop2amu65,
		NumStories:       s.NumStories,
		BaseStructure:    BaseStructure{Name: s.Name, CBFips: s.CBFips, X: s.X, Y: s.Y, DamCat: s.DamCat, GroundElevation: s.GroundElevation}}
}

// Compute implements the consequences.Receptor interface on StrucutreStochastic
func (s StructureStochastic) Compute(d hazards.HazardEvent) (consequences.Result, error) {
	return s.SampleStructure(rand.Int63()).Compute(d) //this needs work so seeds can be controlled.
}

// Compute implements the consequences.Receptor interface on StrucutreDeterminstic
func (s StructureDeterministic) Compute(d hazards.HazardEvent) (consequences.Result, error) {
	/*add, addok := d.(hazards.ArrivalDepthandDurationEvent)
	if addok {
		return computeConsequencesWithReconstruction(add, s)
	}*/
	return computeConsequences(d, s)
}

func computeConsequences(e hazards.HazardEvent, s StructureDeterministic) (consequences.Result, error) {
	header := []string{"fd_id", "x", "y", "hazard", "damage category", "occupancy type", "structure damage", "content damage", "pop2amu65", "pop2amo65", "pop2pmu65", "pop2pmo65", "cbfips", "s_dam_per", "c_dam_per"}
	results := []interface{}{"updateme", 0.0, 0.0, e, "dc", "ot", 0.0, 0.0, 0, 0, 0, 0, "CENSUSBLOCKFIPS", 0, 0}
	var ret = consequences.Result{Headers: header, Result: results}
	var err error = nil
	sval := s.StructVal
	conval := s.ContVal
	sDamFun, sderr := s.OccType.GetComponentDamageFunctionForHazard("structure", e)
	if sderr != nil {
		return ret, sderr
	}
	cDamFun, cderr := s.OccType.GetComponentDamageFunctionForHazard("contents", e)
	if cderr != nil {
		return ret, cderr
	}
	if sDamFun.DamageDriver == hazards.Depth {
		damagefunctionMax := 24.0 //default in case it doesnt cast to paired data.
		damagefunctionMax = sDamFun.DamageFunction.Xvals[len(sDamFun.DamageFunction.Xvals)-1]
		representativeStories := math.Ceil(damagefunctionMax / 9.0)
		if s.NumStories > int32(representativeStories) {
			//there is great potential that the value of the structure is not representative of the damage function range.
			modifier := representativeStories / float64(s.NumStories)
			sval *= modifier
			conval *= modifier
		}
	} //else dont modify value because damage is not driven by depth
	if e.Has(sDamFun.DamageDriver) && e.Has(cDamFun.DamageDriver) {
		//they exist!
		sdampercent := 0.0
		cdampercent := 0.0
		switch sDamFun.DamageDriver {
		case hazards.Depth:
			depthAboveFFE := e.Depth() - s.FoundHt
			sdampercent = sDamFun.DamageFunction.SampleValue(depthAboveFFE) / 100 //assumes what type the damage array is in
			cdampercent = cDamFun.DamageFunction.SampleValue(depthAboveFFE) / 100
		case hazards.Erosion:
			sdampercent = sDamFun.DamageFunction.SampleValue(e.Erosion()) / 100 //assumes what type the damage array is in
			cdampercent = cDamFun.DamageFunction.SampleValue(e.Erosion()) / 100
		default:
			return consequences.Result{}, errors.New("structures: could not understand the damage driver")
		}

		ret.Result[0] = s.BaseStructure.Name
		ret.Result[1] = s.BaseStructure.X
		ret.Result[2] = s.BaseStructure.Y
		ret.Result[3] = e
		ret.Result[4] = s.BaseStructure.DamCat
		ret.Result[5] = s.OccType.Name
		ret.Result[6] = sdampercent * sval
		ret.Result[7] = cdampercent * conval
		ret.Result[8] = s.Pop2amu65
		ret.Result[9] = s.Pop2amo65
		ret.Result[10] = s.Pop2pmu65
		ret.Result[11] = s.Pop2pmo65
		ret.Result[12] = s.CBFips
		ret.Result[13] = sdampercent
		ret.Result[14] = cdampercent
	} else if e.Has(hazards.Qualitative) {
		//this was done primarily to support the NHC in categorizing structures in special zones in their classified surge grids.
		ret.Result[0] = s.BaseStructure.Name
		ret.Result[1] = s.BaseStructure.X
		ret.Result[2] = s.BaseStructure.Y
		ret.Result[3] = e
		ret.Result[4] = s.BaseStructure.DamCat
		ret.Result[5] = s.OccType.Name
		ret.Result[6] = 0.0
		ret.Result[7] = 0.0
		ret.Result[8] = s.Pop2amu65
		ret.Result[9] = s.Pop2amo65
		ret.Result[10] = s.Pop2pmu65
		ret.Result[11] = s.Pop2pmo65
		ret.Result[12] = s.CBFips
		ret.Result[13] = 0.0
		ret.Result[14] = 0.0
	} else {
		err = errors.New("structure: hazard did not contain valid parameters to impact a structure")
	}
	return ret, err
}

/*
func computeConsequencesWithReconstruction(e hazards.ArrivalDepthandDurationEvent, s StructureDeterministic) (consequences.Result, error) {
	header := []string{"fd_id", "x", "y", "hazard", "damage category", "occupancy type", "structure damage", "content damage", "pop2amu65", "pop2amo65", "pop2pmu65", "pop2pmo65", "cbfips", "daystoreconstruction", "rebuilddate"}
	results := []interface{}{"updateme", 0.0, 0.0, e, "dc", "ot", 0.0, 0.0, 0, 0, 0, 0, "CENSUSBLOCKFIPS", 0.0, time.Now()}
	var ret = consequences.Result{Headers: header, Result: results}
	var err error = nil

	if e.Has(hazards.Depth) { //currently the damage functions are depth based, so depth is required, the getstructuredamagefunctionforhazard method chooses approprate damage functions for a hazard.
		if e.Depth() < 0.0 {
			err = errors.New("depth above ground was less than zero")
		}
		if e.Depth() > 9999.0 {
			err = errors.New("depth above ground was greater than 9999")
		}
		depthAboveFFE := e.Depth() - s.FoundHt
		damagePercent := s.OccType.GetStructureDamageFunctionForHazard(e).SampleValue(depthAboveFFE) / 100 //assumes what type the damage array is in
		cdamagePercent := s.OccType.GetContentDamageFunctionForHazard(e).SampleValue(depthAboveFFE) / 100
		reconstructiondays := 0.0
		switch s.DamCat {
		case "RES":
			reconstructiondays = 30.0
		case "COM":
			reconstructiondays = 90.0
		case "IND":
			reconstructiondays = 270.0
		case "PUB":
			reconstructiondays = 360.0
		default:
			reconstructiondays = 180.0
		}
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
		ret.Result[12] = s.CBFips
		rebuilddays := (damagePercent * reconstructiondays) + e.Duration()
		ret.Result[13] = rebuilddays
		ret.Result[14] = e.ArrivalTime().AddDate(0, 0, int(rebuilddays)) //rounds to int
	} else {
		err = errors.New("Hazard did not contain depth")
	}
	return ret, err
}
*/
