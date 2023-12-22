package lifeloss

import (
	"encoding/json"
	"errors"
	"log"
	"math/rand"

	"github.com/HydrologicEngineeringCenter/go-statistics/paireddata"
	"github.com/USACE/go-consequences/consequences"
	"github.com/USACE/go-consequences/hazards"
	"github.com/USACE/go-consequences/structures"
	"github.com/USACE/go-consequences/warning"
)

type Mobility uint

const (
	Unknown   Mobility = 0 //0
	Mobile    Mobility = 1 //1
	NotMobile Mobility = 2 //2
)

type LifeLossEngine struct {
	LethalityCurves   map[LethalityZone]LethalityCurve
	StabilityCriteria map[string]StabilityCriteria
	WarningSystem     warning.WarningResponseSystem
	ResultsHeader     []string
}

func Init() LifeLossEngine {
	var HighPD paireddata.PairedData
	json.Unmarshal(DefaultHighLethalityBytes, &HighPD)
	var LowPD paireddata.PairedData
	json.Unmarshal(DefaultLowLethalityBytes, &LowPD)
	High := LethalityCurve{data: HighPD}
	Low := LethalityCurve{data: LowPD}
	lethalityCurves := make(map[LethalityZone]LethalityCurve)
	lethalityCurves[HighLethality] = High
	lethalityCurves[LowLethality] = Low
	stabilityCriteria := make(map[string]StabilityCriteria)
	stabilityCriteria["woodunanchored"] = RescDamWoodUnanchored
	stabilityCriteria["woodanchored"] = RescDamWoodAnchored
	stabilityCriteria["masonryconcretebrick"] = RescDamMasonryConcreteBrick
	return LifeLossEngine{LethalityCurves: lethalityCurves, StabilityCriteria: stabilityCriteria, WarningSystem: warning.ComplianceBasedWarningSystem{ComplianceRate: .75}}
}

func LifeLossHeader() []string {
	return []string{"ll_u65", "ll_o65", "ll_tot"}
}
func LifeLossDefaultResults() []interface{} {
	return []interface{}{0.0, 0.0, 0.0}
}
func (le LifeLossEngine) ComputeLifeLoss(e hazards.HazardEvent, s structures.StructureDeterministic) (consequences.Result, error) {
	//reduce population based off of the warning system's warning function
	rng := rand.New(rand.NewSource(123454))
	le.WarningSystem.WarningFunction()(&s)
	if e.Has(hazards.DV) {
		sc, err := le.determineStability(s)
		if err != nil {
			return consequences.Result{}, err
		}
		//apply building stability criteria
		if sc.Evaluate(e) == Collapsed {
			//select high fataility rate
			lethalityRate := le.LethalityCurves[HighLethality].Sample()
			//apply same fatality rate to everyone
			log.Println(lethalityRate)
			llo65 := applylethalityRateToPopulation(lethalityRate, s.Pop2amo65, rng)
			//llo65 += applylethalityRateToPopulation(lethalityRate, s.Pop2pmo65, rng)
			llu65 := applylethalityRateToPopulation(lethalityRate, s.Pop2amu65, rng)
			//llu65 += applylethalityRateToPopulation(lethalityRate, s.Pop2pmu65, rng)
			return consequences.Result{Headers: LifeLossHeader(), Result: []interface{}{llu65, llo65, llu65 + llo65}}, nil
		} else {
			return le.submergenceCriteria(e, s)
		}
	} else {
		//apply submergence criteria
		return le.submergenceCriteria(e, s)
	}
}
func applylethalityRateToPopulation(lethalityrate float64, population int32, rng *rand.Rand) int32 {
	result := 0
	for i := 0; i < int(population); i++ {
		if rng.Float64() < lethalityrate {
			result++
		}
	}
	return int32(result)
}
func (lle LifeLossEngine) submergenceCriteria(e hazards.HazardEvent, s structures.StructureDeterministic) (consequences.Result, error) {
	//apply submergence criteria
	header := LifeLossHeader()
	depth := e.Depth()
	if depth < 0.0 {
		//no lifeloss

		result := LifeLossDefaultResults()
		return consequences.Result{Headers: header, Result: result}, nil
	} else {
		//for all individuals using different probabilities for over and under 65 based on nsi attributes
		immobleDepthThreshold := (float64(s.NumStories) - 1.0) * 9.0
		mobileDepthThreshold := (float64(s.NumStories) - 1.0) * 9.0
		hasAtticAccess := false
		immobleDepthThreshold += 5.0 + s.FoundHt
		mobileDepthThreshold += 8.0 + s.FoundHt //9-1...
		if hasAtticAccess {
			mobileDepthThreshold += 1.0 + 6.0 + 4.0
			immobleDepthThreshold += 9.0
		}
		mobility := evaluateMobility(s)
		if mobility == Mobile {
			if depth > float64(mobileDepthThreshold) {
				log.Println(lle.LethalityCurves[HighLethality].Sample())
			} else {
				log.Println(lle.LethalityCurves[LowLethality].Sample())
			}
		} else {
			if depth > float64(immobleDepthThreshold) {
				log.Println(lle.LethalityCurves[HighLethality].Sample())
			} else {
				log.Println(lle.LethalityCurves[LowLethality].Sample())
			}
		}
	}
	return consequences.Result{}, errors.New("under construction")
}
func evaluateMobility(s structures.StructureDeterministic) Mobility {
	//determine based on age and disability
	return Mobile
}

func (le LifeLossEngine) determineStability(s structures.StructureDeterministic) (StabilityCriteria, error) {
	//add construction type to optional parameters and provide default criteria
	if s.OccType.Name == "RES2" {
		return le.StabilityCriteria["woodunanchored"], nil
	}
	//get construction type.
	if s.ConstructionType == "M" {
		return le.StabilityCriteria["masonryconcretebrick"], nil
	}
	if s.ConstructionType == "S" {
		return le.StabilityCriteria["masonryconcretebrick"], nil
	}
	return le.StabilityCriteria["woodanchored"], nil
}
