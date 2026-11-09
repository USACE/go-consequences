package structures

import (
	"encoding/json"
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
	BaseFloodElevation                    float64 // Would BFE ever be stochastic and require type to be ParameterValue?
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
	OccType                                         OccupancyTypeDeterministic
	FoundType, FirmZone, ConstructionType           string
	StructVal, ContVal, FoundHt, BaseFloodElevation float64
	NumStories                                      int32
	PopulationSet
}

// GetX implements consequences.Locatable
func (s BaseStructure) Location() geography.Location {
	return geography.Location{X: s.X, Y: s.Y, Z: s.GroundElevation}
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
		OccType:            ot,
		StructVal:          sv,
		ContVal:            cv,
		FoundType:          s.FoundType,
		ConstructionType:   s.ConstructionType,
		FirmZone:           s.FirmZone,
		FoundHt:            fh,
		BaseFloodElevation: s.BaseFloodElevation,
		PopulationSet:      PopulationSet{s.Pop2amo65, s.Pop2pmu65, s.Pop2amo65, s.Pop2amu65},
		NumStories:         s.NumStories,
		BaseStructure:      BaseStructure{Name: s.Name, CBFips: s.CBFips, X: s.X, Y: s.Y, DamCat: s.DamCat, GroundElevation: s.GroundElevation}}
}

// Compute implements the consequences.Receptor interface on StrucutreStochastic
func (s StructureStochastic) Compute(d hazards.HazardEvent) (consequences.Result, error) {
	return s.SampleStructure(rand.Int63()).Compute(d) //this needs work so seeds can be controlled.
}

// Compute implements the consequences.Receptor interface on StrucutreDeterminstic
func (s StructureDeterministic) Compute(d hazards.HazardEvent) (consequences.Result, error) {
	multi, ok := d.(hazards.MultiHazardEvent)
	if ok {
		multiFrequency, ok2 := multi.(hazards.MultiFrequencyHazardEvent)
		if ok2 {
			return computeConsequencesMultiFrequency(multiFrequency, s)
		}
		return computeConsequencesMultiHazard(multi, s)
	}
	return computeConsequences(d, s)
}

// Compute implements the consequences.Receptor interface on StrucutreDeterminstic
func (s StructureDeterministic) Clone() StructureDeterministic {
	return StructureDeterministic{
		OccType:            s.OccType,
		StructVal:          s.StructVal,
		ContVal:            s.ContVal,
		FoundType:          s.FoundType,
		ConstructionType:   s.ConstructionType,
		FirmZone:           s.FirmZone,
		FoundHt:            s.FoundHt,
		BaseFloodElevation: s.BaseFloodElevation,
		PopulationSet:      PopulationSet{s.Pop2amo65, s.Pop2pmu65, s.Pop2amo65, s.Pop2amu65},
		NumStories:         s.NumStories,
		BaseStructure:      BaseStructure{Name: s.Name, CBFips: s.CBFips, X: s.X, Y: s.Y, DamCat: s.DamCat, GroundElevation: s.GroundElevation}}
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

// computeConsequencesMultiHazard computes lifecycle consequences for a structure exposed to multiple potential damage drivers (flood depth, erosion, wave attack).
func computeConsequencesMultiHazard(event hazards.MultiHazardEvent, s StructureDeterministic) (consequences.Result, error) {

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

	sDepthDamFun, sderr := s.OccType.GetComponentDamageFunctionForDamageDriver("structure", event.Parameters()^hazards.Erosion) // can't simple use hazards.Depth as hazard may have salinity or waveheight params (^ is bitwise XOR operator)
	if sderr != nil {
		return ret, sderr
	}
	sErosionDamFun, seerr := s.OccType.GetComponentDamageFunctionForDamageDriver("structure", hazards.Erosion)
	if seerr != nil {
		// GetComponentDamageFunctionForDamageDriver does not return a default damage function if the Parameter argument does not include Depth.
		// Should we define a default erosion damage function for the function to return when an Erosion df is requested?
		return ret, seerr
	}

	cDepthDamFun, cderr := s.OccType.GetComponentDamageFunctionForDamageDriver("contents", event.Parameters()^hazards.Erosion) // can't simple use hazards.Depth as hazard may have salinity or waveheight params (^ is bitwise XOR operator)
	if cderr != nil {
		return ret, cderr
	}
	cErosionDamFun, ceerr := s.OccType.GetComponentDamageFunctionForDamageDriver("contents", hazards.Erosion)
	if ceerr != nil {
		return ret, ceerr
	}

	rDamFun, rderr := s.OccType.GetComponentDamageFunctionForDamageDriver("reconstruction", hazards.Depth) // Would it make more sense for reconstruction to be treated as a damage driver rather than a component?
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
	cumulativeStructureLoss := 0.0
	cumulativeContentLoss := 0.0
	timesRebuilt := 0

	// instead of adjusting the structure value of tall structures for depth-damage calcs, we will modify the damage factor
	// returned by the depth damage function. This will leave structure value unmodified for erosion calculations
	depthDamageModifier := 1.0
	damagefunctionMax := 24.0 //default in case it doesnt cast to paired data.
	damagefunctionMax = sDepthDamFun.DamageFunction.Xvals[len(sDepthDamFun.DamageFunction.Xvals)-1]
	representativeStories := math.Ceil(damagefunctionMax / 9.0)
	if s.NumStories > int32(representativeStories) {
		//there is great potential that the value of the structure is not representative of the damage function range.
		depthDamageModifier = representativeStories / float64(s.NumStories)
	}

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

			if sPctloss_rebuilt > 0 || cPctloss_rebuilt > 0 {
				timesRebuilt++
			}

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

		if event.Has(sDepthDamFun.DamageDriver) && event.Has(sErosionDamFun.DamageDriver) &&
			event.Has(cDepthDamFun.DamageDriver) && event.Has(cErosionDamFun.DamageDriver) &&
			event.Has(rDamFun.DamageDriver) {

			s_depth_dam_pct := 0.0
			c_depth_dam_pct := 0.0
			s_erosion_dam_pct := 0.0
			c_erosion_dam_pct := 0.0
			s_wave_dam_pct := 0.0 // will we have specific wave attack damage functions, or are we getting adjusted depth damage functions based on WaveHeight parameter?
			c_wave_dam_pct := 0.0

			sdampercent := 0.0
			cdampercent := 0.0
			sdamage := 0.0
			cdamage := 0.0

			reconstruction_days := 0.0
			completion_date := time.Time{}

			depthAboveFFE := event.Depth() - s.FoundHt
			erosion_pct := event.Erosion()
			//1. get damage pct for each damage driver
			//2. multiply depth damage pct by depthDamageModifier
			//3. combine damage pcts according to equation to get event damage pct (sdampercent = depthpct + erosionpct + wavepct - (d*w) - (d*e) - (w*e) + (d*w*e))
			//4. update structure damage factor sDamageFactor = 1 - (1-sDamageFactor)*(1-sdampercent)
			//5. calculate reconstruction days
			s_depth_dam_pct = sDepthDamFun.DamageFunction.SampleValue(depthAboveFFE)
			c_depth_dam_pct = cDepthDamFun.DamageFunction.SampleValue(depthAboveFFE)
			s_depth_dam_pct = s_depth_dam_pct * depthDamageModifier
			c_depth_dam_pct = c_depth_dam_pct * depthDamageModifier

			s_erosion_dam_pct = sErosionDamFun.DamageFunction.SampleValue(erosion_pct)
			c_erosion_dam_pct = cErosionDamFun.DamageFunction.SampleValue(erosion_pct)

			sdampercent = s_depth_dam_pct + s_erosion_dam_pct + s_wave_dam_pct -
				(s_depth_dam_pct * s_erosion_dam_pct) - (s_depth_dam_pct * s_wave_dam_pct) - (s_erosion_dam_pct * s_wave_dam_pct) +
				(s_depth_dam_pct * s_erosion_dam_pct * s_wave_dam_pct)

			cdampercent = c_depth_dam_pct + c_erosion_dam_pct + c_wave_dam_pct -
				(c_depth_dam_pct * c_erosion_dam_pct) - (c_depth_dam_pct * c_wave_dam_pct) - (c_erosion_dam_pct * c_wave_dam_pct) +
				(c_depth_dam_pct * c_erosion_dam_pct * c_wave_dam_pct)

			sDamageFactor = 1 - (1-sDamageFactor)*(1-sdampercent)
			cDamageFactor = 1 - (1-cDamageFactor)*(1-cdampercent)

			sdamage = svalcurr - (svalcurr * sDamageFactor)
			cdamage = convalcurr - (convalcurr * cDamageFactor)

			arrival := event.ArrivalTime() // do we need a check that the ArrivalTime is not just the default time.Time{}?

			duration := 0.0
			if event.Duration() > 0.0 { // nodata value for e.Duration == -901.0
				duration = event.Duration()
			}

			// calculate reconstruction_days based on damageFactor to account for potential remaining damage from previous events
			reconstruction_days = math.Ceil(rDamFun.DamageFunction.SampleValue(sDamageFactor) + duration)
			completion_date = arrival.AddDate(0, 0, int(reconstruction_days))

			svalcurr = sval * (1 - sDamageFactor)
			convalcurr = conval * (1 - cDamageFactor)
			cumulativeContentLoss += cdamage
			cumulativeStructureLoss += sdamage
			result.Result[1] = sdamage
			result.Result[2] = cdamage
			result.Result[3] = sdampercent
			result.Result[4] = cdampercent
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
		ret.Result[16] = int32(timesRebuilt)
		ret.Result[18] = cumulativeStructureLoss
		ret.Result[19] = cumulativeContentLoss
		ret.Result[20] = subResult

		if event.HasNext() {
			event.Increment() // go to the next event and restart loop
		} else {
			break
		}
	}
	return ret, err
}

// ComputeEAD takes an array of damages and frequencies and integrates the curve. we should probably refactor this into paired data as a function.
// This func is copied from the compute module to prevent cyclical import
func ComputeEAD(damages []float64, freq []float64) float64 {
	triangle := 0.0
	square := 0.0
	x1 := 1.0 // create a triangle to the first probability space - linear interpolation is probably a problem, maybe use log linear interpolation for the triangle
	y1 := 0.0
	eadT := 0.0
	for i := 0; i < len(freq); i++ {
		xdelta := x1 - freq[i]
		square = xdelta * y1
		triangle = ((xdelta) * (damages[i] - y1)) / 2.0
		eadT += square + triangle
		x1 = freq[i]
		y1 = damages[i]
	}
	if x1 != 0.0 {
		xdelta := x1 - 0.0
		eadT += xdelta * y1 //no extrapolation, just continue damages out as if it were truth for all remaining probability.

	}
	return eadT
}

func computeConsequencesMultiFrequency(event hazards.MultiFrequencyHazardEvent, s StructureDeterministic) (consequences.Result, error) {

	header := []string{"fd_id", "x", "y", "damage category", "occupancy type", "val_structure", "val_contents", "struct_ead", "cont_ead", "pop2amu65", "pop2amo65", "pop2pmu65", "pop2pmo65", "cbfips"}
	results := []interface{}{"updateme", 0.0, 0.0, "dc", "ot", 0.0, 0.0, 0.0, 0.0, 0, 0, 0, 0, "CENSUSBLOCKFIPS"}
	var ret = consequences.Result{Headers: header, Result: results}
	sval := s.StructVal
	conval := s.ContVal
	sDamFun, sderr := s.OccType.GetComponentDamageFunctionForHazard("structure", event)
	if sderr != nil {
		return ret, sderr
	}
	cDamFun, cderr := s.OccType.GetComponentDamageFunctionForHazard("contents", event)
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
	sdams := make([]float64, len(event.Frequencies()))
	cdams := make([]float64, len(event.Frequencies()))
	successes := 0
	for {

		freq := event.Frequency()
		ret.Headers = append(ret.Headers, fmt.Sprintf("%2.6fS", freq))
		ret.Headers = append(ret.Headers, fmt.Sprintf("%2.6fC", freq))
		ret.Headers = append(ret.Headers, fmt.Sprintf("%2.6fH", freq))

		if event.Has(sDamFun.DamageDriver) && event.Has(cDamFun.DamageDriver) {
			switch sDamFun.DamageDriver {
			case hazards.Depth:
				depthAboveFFE := event.Depth() - s.FoundHt
				spct := sDamFun.DamageFunction.SampleValue(depthAboveFFE) / 100 //assumes what type the damage array is in
				cpct := cDamFun.DamageFunction.SampleValue(depthAboveFFE) / 100
				sdam := spct * sval
				cdam := cpct * conval
				sdams[event.Index()] = sdam
				cdams[event.Index()] = cdam

				ret.Result = append(ret.Result, sdam)
				ret.Result = append(ret.Result, cdam)
				hazardstring := ""
				bytes, err := json.Marshal(event.This())
				if err == nil {
					hazardstring = string(bytes)
				}
				ret.Result = append(ret.Result, hazardstring)
				successes++
			case hazards.Erosion:
				spct := sDamFun.DamageFunction.SampleValue(event.Erosion()) / 100 //assumes what type the damage array is in
				cpct := cDamFun.DamageFunction.SampleValue(event.Erosion()) / 100
				sdam := spct * sval
				cdam := cpct * conval
				sdams[event.Index()] = sdam
				cdams[event.Index()] = cdam

				ret.Result = append(ret.Result, sdam)
				ret.Result = append(ret.Result, cdam)
				hazardstring := ""
				bytes, err := json.Marshal(event.This())
				if err == nil {
					hazardstring = string(bytes)
				}
				ret.Result = append(ret.Result, hazardstring)
				successes++
			default:
				sdams[event.Index()] = 0.0
				cdams[event.Index()] = 0.0
				ret.Result = append(ret.Result, 0.0)
				ret.Result = append(ret.Result, 0.0)
				hazardstring := ""
				bytes, err := json.Marshal(event.This())
				if err == nil {
					hazardstring = string(bytes)
				}
				ret.Result = append(ret.Result, hazardstring)
				// Do we want to error out on this? Is it possible that one but not all of the hazards in the series could have an invalid damage driver?
				return consequences.Result{}, fmt.Errorf("structures: could not understand the damage driver for event in MultiFrequencyEvent at Index %v", event.Index())
			}
		} else {
			// We don't raise an error here because there might be no hazard for a high frequency event (e.g. 2yr flood), but there may still be one for the
			// lower frequency events
			sdams[event.Index()] = 0.0
			cdams[event.Index()] = 0.0
			ret.Result = append(ret.Result, -9999)
			ret.Result = append(ret.Result, -9999)
			hazardstring := ""
			bytes, err := json.Marshal(event.This())
			if err == nil {
				hazardstring = string(bytes)
			}
			ret.Result = append(ret.Result, hazardstring)
		}

		if event.HasNext() {
			event.Increment()
		} else {
			break
		}
	}
	sEAD := ComputeEAD(sdams, event.Frequencies())
	cEAD := ComputeEAD(cdams, event.Frequencies())
	ret.Result[0] = s.BaseStructure.Name
	ret.Result[1] = s.BaseStructure.X
	ret.Result[2] = s.BaseStructure.Y
	ret.Result[3] = s.BaseStructure.DamCat
	ret.Result[4] = s.OccType.Name
	ret.Result[5] = sval
	ret.Result[6] = conval
	ret.Result[7] = sEAD
	ret.Result[8] = cEAD
	ret.Result[9] = s.Pop2amu65
	ret.Result[10] = s.Pop2amo65
	ret.Result[11] = s.Pop2pmu65
	ret.Result[12] = s.Pop2pmo65
	ret.Result[13] = s.CBFips
	if successes == 0 {
		return ret, errors.New("structure: hazard did not contain valid parameters to impact a structure")
	}
	return ret, nil
}
