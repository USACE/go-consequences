package structures

import (
	"errors"
	"fmt"
	"math"
	"math/rand"
	"strings"
	"time"

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
type PopulationSet struct {
	Pop2pmo65, Pop2pmu65, Pop2amo65, Pop2amu65 int32
}

// StructureStochastic is a base structure with an occupancy type stochastic and parameter values for all parameters
type StructureStochastic struct {
	BaseStructure
	UseUncertainty                        bool //defaults to false!
	OccType                               OccupancyTypeStochastic
	FoundType, FirmZone, ConstructionType string
	StructVal, ContVal, FoundHt           consequences.ParameterValue
	NumStories                            int32
	PopulationSet
}

func (f *StructureStochastic) ApplyFoundationHeightUncertanty(fu *FoundationUncertainty) {
	queryString := "default_slab"
	default_FHU := fu.Values[queryString]
	if f.OccType.Name == "RES2" {
		queryString = "RES2"
	} else if f.OccType.Name == "RES3A" {
		queryString = "RES1_RES3A_RES3B"
	} else if f.OccType.Name == "RES3A" {
		queryString = "RES1_RES3A_RES3B"
	} else if strings.Contains(f.OccType.Name, "RES1") {
		queryString = "RES1_RES3A_RES3B"
	} else {
		queryString = "default"
	}
	if f.FoundType == "I" { //pile maps to peir
		queryString = fmt.Sprintf("%v_P", queryString)
	} else if f.FoundType == "W" { //wall maps to crawl space.
		queryString = fmt.Sprintf("%v_C", queryString)
	} else {
		queryString = fmt.Sprintf("%v_%v", queryString, f.FoundType)
	}
	FHU, ok := fu.Values[queryString]
	if !ok {
		FHU = default_FHU
	}
	if f.FirmZone == "V" {
		f.FoundHt = consequences.ParameterValue{
			Value: FHU.VzoneDistribution,
		}
	} else {
		f.FoundHt = consequences.ParameterValue{
			Value: FHU.DefaultDistribution,
		}
	}
}

// StructureDeterministic is a base strucure with a deterministic occupancy type and deterministic parameters
type StructureDeterministic struct {
	BaseStructure
	OccType                               OccupancyTypeDeterministic
	FoundType, FirmZone, ConstructionType string
	StructVal, ContVal, FoundHt           float64
	NumStories                            int32
	PopulationSet
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
		if fh < 0 {
			fh = 0.0
		}
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
		FirmZone:         s.FirmZone,
		FoundHt:          fh,
		PopulationSet:    PopulationSet{s.Pop2amo65, s.Pop2pmu65, s.Pop2amo65, s.Pop2amu65},
		NumStories:       s.NumStories,
		BaseStructure:    BaseStructure{Name: s.Name, CBFips: s.CBFips, X: s.X, Y: s.Y, DamCat: s.DamCat, GroundElevation: s.GroundElevation}}
}

// Compute implements the consequences.Receptor interface on StrucutreStochastic
func (s StructureStochastic) Compute(d hazards.HazardEvent) (consequences.Result, error) {
	return s.SampleStructure(rand.Int63()).Compute(d) //this needs work so seeds can be controlled.
}

// Compute implements the consequences.Receptor interface on StrucutreDeterminstic
func (s StructureDeterministic) Compute(d hazards.HazardEvent) (consequences.Result, error) {
	addMulti, ok := d.(hazards.MultiHazardEvent)
	if ok {
		return computeConsequencesMultiHazard(addMulti, s)
	}
	return computeConsequences(d, s)
}

// Compute implements the consequences.Receptor interface on StrucutreDeterminstic
func (s StructureDeterministic) Clone() StructureDeterministic {
	return StructureDeterministic{
		OccType:          s.OccType,
		StructVal:        s.StructVal,
		ContVal:          s.ContVal,
		FoundType:        s.FoundType,
		ConstructionType: s.ConstructionType,
		FirmZone:         s.FirmZone,
		FoundHt:          s.FoundHt,
		PopulationSet:    PopulationSet{s.Pop2amo65, s.Pop2pmu65, s.Pop2amo65, s.Pop2amu65},
		NumStories:       s.NumStories,
		BaseStructure:    BaseStructure{Name: s.Name, CBFips: s.CBFips, X: s.X, Y: s.Y, DamCat: s.DamCat, GroundElevation: s.GroundElevation}}
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

func computeConsequencesWithReconstruction(e hazards.HazardEvent, s StructureDeterministic) (consequences.Result, error) {
	// NOTE: This version gets reconstruction as a damage function on the structure's occtype

	header := []string{"fd_id", "x", "y", "hazard", "damage category", "occupancy type", "structure damage", "content damage", "pop2amu65", "pop2amo65", "pop2pmu65", "pop2pmo65", "cbfips", "s_dam_per", "c_dam_per", "reconstruction_days"}
	results := []interface{}{"updateme", 0.0, 0.0, e, "dc", "ot", 0.0, 0.0, 0, 0, 0, 0, "CENSUSBLOCKFIPS", 0, 0, 0.0}
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

	rDamFun, rderr := s.OccType.GetComponentDamageFunctionForHazard("reconstruction", e)
	if rderr != nil {
		return ret, cderr
	}

	//TODO: Do we want to return the date that construction will be complete? Only useful if event has arrival time
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
	if e.Has(sDamFun.DamageDriver) && e.Has(cDamFun.DamageDriver) && e.Has(rDamFun.DamageDriver) {
		//they exist!
		sdampercent := 0.0
		cdampercent := 0.0
		reconstruction_days := 0.0

		switch sDamFun.DamageDriver {
		case hazards.Depth:
			depthAboveFFE := e.Depth() - s.FoundHt
			sdampercent = sDamFun.DamageFunction.SampleValue(depthAboveFFE) / 100 //assumes what type the damage array is in
			cdampercent = cDamFun.DamageFunction.SampleValue(depthAboveFFE) / 100
			duration := 0.0
			if e.Duration() > 0.0 { // nodata value for e.Duration == -901.0
				duration = e.Duration()
			}
			reconstruction_days = rDamFun.DamageFunction.SampleValue(sdampercent) + duration
		case hazards.Erosion:
			sdampercent = sDamFun.DamageFunction.SampleValue(e.Erosion()) / 100 //assumes what type the damage array is in
			cdampercent = cDamFun.DamageFunction.SampleValue(e.Erosion()) / 100
			reconstruction_days = rDamFun.DamageFunction.SampleValue(sdampercent)
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
		ret.Result[15] = math.Ceil(reconstruction_days)

	} else {
		err = errors.New("structure: hazard did not contain valid parameters to impact a structure")
	}
	return ret, err
}

func computeConsequencesMulti(events []hazards.HazardEvent, s StructureDeterministic) ([]consequences.Result, error) {

	var ret = make([]consequences.Result, len(events))
	var err error = nil

	// get damage functions for structure based on first hazard event (assumes same parameters) to prevent repeated lookups
	sDamFun, sderr := s.OccType.GetComponentDamageFunctionForHazard("structure", events[0])
	if sderr != nil {
		return ret, sderr
	}
	cDamFun, cderr := s.OccType.GetComponentDamageFunctionForHazard("contents", events[0])
	if cderr != nil {
		return ret, cderr
	}
	rDamFun, rderr := s.OccType.GetComponentDamageFunctionForHazard("reconstruction", events[0])
	if rderr != nil {
		return ret, cderr
	}

	sval := s.StructVal
	svalcurr := sval
	sDamageFactor := 0.0 // this is the current pct_damage to the structure
	conval := s.ContVal
	convalcurr := conval
	cDamageFactor := 0.0 // this is the current pct_damage to the contents

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

	for i, e := range events {
		if !e.Has(hazards.ArrivalTime) {
			return ret, errors.New("structures: hazard event does not have ArrivalTime")
		}

		if i > 0 {
			// update structure value and damagefactor to reflect construction progress from previous event
			tc0, err := ret[i-1].Fetch("completion_date") //
			if err != nil {
				return ret, errors.New("structures: unable to get completion date for previous hazard event")
			}
			last_completion_time := tc0.(time.Time)

			// calculate reconstruction progress assuming linear rebuild
			t0 := events[i-1].ArrivalTime().AddDate(0, 0, int(events[i-1].Duration()))                 // this is the time reconstruction began
			pct_complete := (float64(e.ArrivalTime().Sub(t0)) / float64(last_completion_time.Sub(t0))) // this is the percentage of the reconstruction that is complete

			if pct_complete > 1.0 {
				pct_complete = 1.0
			}

			sPctloss_rebuilt := sDamageFactor * pct_complete
			cPctloss_rebuilt := cDamageFactor * pct_complete

			// update structure damage factor to reflect completed construction
			sDamageFactor = sDamageFactor - sPctloss_rebuilt
			if sDamageFactor < 0.0 {
				sDamageFactor = 0
			}

			cDamageFactor = cDamageFactor - cPctloss_rebuilt
			if cDamageFactor < 0.0 {
				cDamageFactor = 0
			}

			// update structure value to reflect completed construction
			svalcurr = sval * (1 - sDamageFactor)
			convalcurr = conval * (1 - cDamageFactor)
		}

		header := []string{"fd_id", "structure damage", "content damage", "s_dam_per", "c_dam_per", "reconstruction_days", "completion_date", "structure_value", "content_value"}
		values := []interface{}{"updateme", 0.0, 0.0, 0.0, 0.0, 0.0, time.Time{}, 0.0, 0.0}
		result := consequences.Result{Headers: header, Result: values}

		if e.Has(sDamFun.DamageDriver) && e.Has(cDamFun.DamageDriver) && e.Has(rDamFun.DamageDriver) {
			//they exist!
			sdampercent := 0.0
			sdamage := 0.0
			cdampercent := 0.0
			cdamage := 0.0

			reconstruction_days := 0.0
			completion_date := time.Time{}

			switch sDamFun.DamageDriver {
			case hazards.Depth:
				depthAboveFFE := e.Depth() - s.FoundHt
				sdampercent = sDamFun.DamageFunction.SampleValue(depthAboveFFE) / 100 //assumes what type the damage array is in
				cdampercent = cDamFun.DamageFunction.SampleValue(depthAboveFFE) / 100
				sdamage = svalcurr * sdampercent
				cdamage = convalcurr * cdampercent

				sDamageFactor = 1 - (1-sDamageFactor)*(1-sdampercent)
				cDamageFactor = 1 - (1-cDamageFactor)*(1-cdampercent)
				// total time to complete reconstruction consists of three parts
				//	1. Time between the start and end of the event. This is e.Duration()
				//	2. Time between the end of the event and the beginning of reconstruction.
				//		- In reality, this would depend on a lot but simplest assumption is that reconstruction can begin as soon as event ends.
				//	3. Time between reconstruction start and reconstruction end. This is the value returned from the damage function

				arrival := e.ArrivalTime() // do we need a check that the ArrivalTime is not just the default time.Time{}?

				duration := 0.0
				if e.Duration() > 0.0 { // nodata value for e.Duration == -901.0
					duration = e.Duration()
				}

				// calculate reconstruction_days based on damageFactor to account for potential remaining damage from previous events
				reconstruction_days = math.Ceil(rDamFun.DamageFunction.SampleValue(sDamageFactor) + duration)
				completion_date = arrival.AddDate(0, 0, int(reconstruction_days))

			case hazards.Erosion:
				sdampercent = sDamFun.DamageFunction.SampleValue(e.Erosion()) / 100 //assumes what type the damage array is in
				cdampercent = cDamFun.DamageFunction.SampleValue(e.Erosion()) / 100
				sdamage = svalcurr * sdampercent
				cdamage = convalcurr * cdampercent

				sDamageFactor = 1 - (1-sDamageFactor)*(1-sdampercent)
				cDamageFactor = 1 - (1-cDamageFactor)*(1-cdampercent)
				arrival := e.ArrivalTime()
				// calculate reconstruction_days based on damageFactor to account for potential remaining damage from previous events
				reconstruction_days = math.Ceil(rDamFun.DamageFunction.SampleValue(sDamageFactor))
				completion_date = arrival.AddDate(0, 0, int(reconstruction_days))

			default:
				return ret, errors.New("structures: could not understand the damage driver")
			}

			svalcurr = svalcurr * (1 - sDamageFactor)
			convalcurr = convalcurr * (1 - cDamageFactor)

			result.Result[0] = s.BaseStructure.Name
			result.Result[1] = sdamage
			result.Result[2] = cdamage
			result.Result[3] = sdampercent
			result.Result[3] = cdampercent
			result.Result[5] = math.Ceil(reconstruction_days)
			result.Result[6] = completion_date
			result.Result[7] = svalcurr
			result.Result[8] = convalcurr

		} else {
			err = errors.New("structure: hazard did not contain valid parameters to impact a structure")
		}
		ret[i] = result
	}
	return ret, err
}

func computeConsequencesMultiHazard(event hazards.MultiHazardEvent, s StructureDeterministic) (consequences.Result, error) {
	// this function needs to return a single result. Not a slice of results
	// Make a nested Result where each column is itself a Result

	// Rethinking the results format. The current version repeats the structure info for each result.
	// The main result body can include the structure info, and we can simply store the details of the iteration results in a column as json.
	mainHeader := []string{
		"fd_id", "x", "y", "damage category", "occupancy type",
		"pop2amu65", "pop2amo65", "pop2pmu65", "pop2pmo65", "cbfips",
		"original structure value", "original content value", "final structure value", "final content value", //
		"original ffe", "final ffe", "times rebuilt", "times raised", "StructureTotalLoss", "ContentsTotalLoss",
		"hazard results",
	}
	subResultsHeader := make([]string, 0)
	subResultsResult := make([]interface{}, 0)
	subResult := consequences.Result{Headers: subResultsHeader, Result: subResultsResult}
	mainResults := []interface{}{
		"updateme", 0.0, 0.0, "damcat", "occtype",
		0.0, 0.0, 0.0, 0.0, "cbfips",
		0.0, 0.0, 0.0, 0.0,
		0.0, 0.0, 0.0, 0.0, 0.0, 0.0,
		subResult,
	}
	var ret = consequences.Result{Headers: mainHeader, Result: mainResults}
	var err error = nil

	// damage functions for structure
	sDamFun, sderr := s.OccType.GetComponentDamageFunctionForHazard("structure", event)
	if sderr != nil {
		return ret, sderr
	}
	cDamFun, cderr := s.OccType.GetComponentDamageFunctionForHazard("contents", event)
	if cderr != nil {
		return ret, cderr
	}
	rDamFun, rderr := s.OccType.GetComponentDamageFunctionForHazard("reconstruction", event)
	if rderr != nil {
		return ret, cderr
	}

	// variables for tracking value and damage across hazards
	sval := s.StructVal
	svalcurr := sval
	sDamageFactor := 0.0 // this is the current pct_damage to the structure
	conval := s.ContVal
	convalcurr := conval
	cDamageFactor := 0.0 // this is the current pct_damage to the contents

	// adjust value for tall structures
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

	for {
		// Calculate reconstruction from previous hazard if we aren't on the first one
		if event.HasPrevious() {

			pr, err := subResult.Fetch(fmt.Sprintf("%d", event.Index()-1))
			// pr, err := ret.Fetch(fmt.Sprintf("%d", event.Index()-1))
			if err != nil {
				return ret, fmt.Errorf("structures: Unable to fetch previous result at index = %v", event.Index())
			}
			previous_result := pr.(consequences.Result)
			tc0, err := previous_result.Fetch("completion_date")
			if err != nil {
				return ret, errors.New("structures: unable to get completion date for previous hazard event")
			}
			last_completion_time := tc0.(time.Time)

			previous_event, err := event.Previous()
			if err != nil {
				return ret, err
			}

			// calculate reconstruction progress assuming linear rebuild
			t0 := previous_event.ArrivalTime().AddDate(0, 0, int(previous_event.Duration()))               // this is the time reconstruction began
			pct_complete := (float64(event.ArrivalTime().Sub(t0)) / float64(last_completion_time.Sub(t0))) // this is the percentage of the reconstruction that is complete

			if pct_complete > 1.0 {
				pct_complete = 1.0
			}

			sPctloss_rebuilt := sDamageFactor * pct_complete
			cPctloss_rebuilt := cDamageFactor * pct_complete

			// update structure damage factor to reflect completed construction
			sDamageFactor = sDamageFactor - sPctloss_rebuilt
			if sDamageFactor < 0.0 {
				sDamageFactor = 0
			}

			cDamageFactor = cDamageFactor - cPctloss_rebuilt
			if cDamageFactor < 0.0 {
				cDamageFactor = 0
			}

			// update structure value to reflect completed construction
			svalcurr = sval * (1 - sDamageFactor)
			convalcurr = conval * (1 - cDamageFactor)
		}

		header := []string{"hazard", "structure damage", "content damage", "s_dam_per", "c_dam_per", "reconstruction_days", "completion_date", "structure_value", "content_value"}
		values := []interface{}{event.This(), 0.0, 0.0, 0.0, 0.0, 0.0, time.Time{}, 0.0, 0.0}
		result := consequences.Result{Headers: header, Result: values}

		if event.Has(sDamFun.DamageDriver) && event.Has(cDamFun.DamageDriver) && event.Has(rDamFun.DamageDriver) {
			//they exist!
			sdampercent := 0.0
			sdamage := 0.0
			cdampercent := 0.0
			cdamage := 0.0

			reconstruction_days := 0.0
			completion_date := time.Time{}

			switch sDamFun.DamageDriver {
			case hazards.Depth:
				depthAboveFFE := event.Depth() - s.FoundHt
				sdampercent = sDamFun.DamageFunction.SampleValue(depthAboveFFE) / 100 //assumes what type the damage array is in
				cdampercent = cDamFun.DamageFunction.SampleValue(depthAboveFFE) / 100
				sdamage = svalcurr * sdampercent
				cdamage = convalcurr * cdampercent

				sDamageFactor = 1 - (1-sDamageFactor)*(1-sdampercent)
				cDamageFactor = 1 - (1-cDamageFactor)*(1-cdampercent)
				// total time to complete reconstruction consists of three parts
				//	1. Time between the start and end of the event. This is e.Duration()
				//	2. Time between the end of the event and the beginning of reconstruction.
				//		- In reality, this would depend on a lot but simplest assumption is that reconstruction can begin as soon as event ends.
				//	3. Time between reconstruction start and reconstruction end. This is the value returned from the damage function

				arrival := event.ArrivalTime() // do we need a check that the ArrivalTime is not just the default time.Time{}?

				duration := 0.0
				if event.Duration() > 0.0 { // nodata value for e.Duration == -901.0
					duration = event.Duration()
				}

				// calculate reconstruction_days based on damageFactor to account for potential remaining damage from previous events
				reconstruction_days = math.Ceil(rDamFun.DamageFunction.SampleValue(sDamageFactor) + duration)
				completion_date = arrival.AddDate(0, 0, int(reconstruction_days))

			case hazards.Erosion:
				sdampercent = sDamFun.DamageFunction.SampleValue(event.Erosion()) / 100 //assumes what type the damage array is in
				cdampercent = cDamFun.DamageFunction.SampleValue(event.Erosion()) / 100
				sdamage = svalcurr * sdampercent
				cdamage = convalcurr * cdampercent

				sDamageFactor = 1 - (1-sDamageFactor)*(1-sdampercent)
				cDamageFactor = 1 - (1-cDamageFactor)*(1-cdampercent)
				arrival := event.ArrivalTime()
				// calculate reconstruction_days based on damageFactor to account for potential remaining damage from previous events
				reconstruction_days = math.Ceil(rDamFun.DamageFunction.SampleValue(sDamageFactor))
				completion_date = arrival.AddDate(0, 0, int(reconstruction_days))

			default:
				return ret, errors.New("structures: could not understand the damage driver")
			}

			svalcurr = svalcurr * (1 - sDamageFactor)
			convalcurr = convalcurr * (1 - cDamageFactor)

			result.Result[1] = sdamage
			result.Result[2] = cdamage
			result.Result[3] = sdampercent
			result.Result[3] = cdampercent
			result.Result[5] = math.Ceil(reconstruction_days)
			result.Result[6] = completion_date
			result.Result[7] = svalcurr
			result.Result[8] = convalcurr

		} else {
			err = errors.New("structure: hazard did not contain valid parameters to impact a structure")
		}
		subResultsHeader = append(subResultsHeader, fmt.Sprintf("%d", event.Index()))
		subResultsResult = append(subResultsResult, result)
		subResult = consequences.Result{Headers: subResultsHeader, Result: subResultsResult}
		ret.Result[0] = s.BaseStructure.Name
		ret.Result[1] = s.BaseStructure.X
		ret.Result[2] = s.BaseStructure.Y
		ret.Result[3] = s.BaseStructure.DamCat
		ret.Result[4] = s.OccType.Name
		ret.Result[5] = s.Pop2amu65
		ret.Result[6] = s.Pop2amo65
		ret.Result[7] = s.Pop2pmu65
		ret.Result[8] = s.Pop2pmo65
		ret.Result[9] = s.CBFips
		ret.Result[10] = sval
		ret.Result[11] = conval
		ret.Result[12] = svalcurr
		ret.Result[13] = convalcurr
		ret.Result[20] = subResult

		if event.HasNext() {
			event.Increment() // go to the next event and restart loop
		} else {
			break
		}
	}
	return ret, err
}
