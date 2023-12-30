package lifeloss

import (
	"encoding/json"
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

type PopulationSet struct {
	o652am int
	o652pm int
	u652am int
	u652pm int
}
type LifeLossEngine struct {
	LethalityCurves   map[LethalityZone]LethalityCurve
	StabilityCriteria map[string]StabilityCriteria
	WarningSystem     warning.WarningResponseSystem
	ResultsHeader     []string
}

// Init in the lifeloss package creates a life loss engine with the default settings
func Init() LifeLossEngine {
	//initialize the default high lethality rate relationship
	var HighPD paireddata.PairedData
	json.Unmarshal(DefaultHighLethalityBytes, &HighPD)
	//initialize the default low lethality rate relationships
	var LowPD paireddata.PairedData
	json.Unmarshal(DefaultLowLethalityBytes, &LowPD)
	//create a lethality curve instance for High Lethality
	High := LethalityCurve{data: HighPD}
	//create a lethality curve instance for Low Lethality
	Low := LethalityCurve{data: LowPD}
	//create a map of lethality zone to lethality curve
	lethalityCurves := make(map[LethalityZone]LethalityCurve)
	lethalityCurves[HighLethality] = High
	lethalityCurves[LowLethality] = Low
	//initalize the stability criteria
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
			log.Println("Stability Based Lifeloss")
			//select high fataility rate
			lethalityRate := le.LethalityCurves[HighLethality].Sample()
			//apply same fatality rate to everyone
			//log.Println(lethalityRate)
			llo65 := applylethalityRateToPopulation(lethalityRate, s.Pop2amo65, rng)
			//llo65 += applylethalityRateToPopulation(lethalityRate, s.Pop2pmo65, rng)
			llu65 := applylethalityRateToPopulation(lethalityRate, s.Pop2amu65, rng)
			//llu65 += applylethalityRateToPopulation(lethalityRate, s.Pop2pmu65, rng)
			result := consequences.Result{Headers: LifeLossHeader(), Result: []interface{}{llu65, llo65, llu65 + llo65}}
			//log.Println(result)
			return result, nil
		} else {
			return le.submergenceCriteria(e, s, rng)
		}
	} else {
		//apply submergence criteria
		return le.submergenceCriteria(e, s, rng)
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
func (lle LifeLossEngine) submergenceCriteria(e hazards.HazardEvent, s structures.StructureDeterministic, rng *rand.Rand) (consequences.Result, error) {
	//apply submergence criteria
	log.Println("Submergence Based Lifeloss")
	header := LifeLossHeader()
	depth := e.Depth()
	if depth < 0.0 {
		//no lifeloss

		result := LifeLossDefaultResults()
		return consequences.Result{Headers: header, Result: result}, nil
	} else {
		//for all ages of individuals using different probabilities to assign mobility
		// for over and under 65 based on nsi attribute of "percent not mobile"
		immobleDepthThreshold := (float64(s.NumStories) - 1.0) * 9.0
		mobileDepthThreshold := (float64(s.NumStories) - 1.0) * 9.0
		hasAtticAccess := false
		immobleDepthThreshold += 5.0 + s.FoundHt
		mobileDepthThreshold += 8.0 + s.FoundHt //9-1...
		if hasAtticAccess {
			mobileDepthThreshold += 1.0 + 6.0 + 4.0
			immobleDepthThreshold += 9.0
		}

		mobilitySet := evaluateMobility(s, rng)
		llu65 := 0
		llo65 := 0
		for k, v := range mobilitySet {
			//apply to the appropriate age/time of day
			//log.Println(v)
			if k == Mobile {
				if depth > float64(mobileDepthThreshold) {
					ret := lle.createLifeLossSet(v, lle.LethalityCurves[HighLethality], rng)
					llu65 += ret.u652am
					llo65 += ret.o652am
				} else {
					ret := lle.createLifeLossSet(v, lle.LethalityCurves[LowLethality], rng)
					llu65 += ret.u652am
					llo65 += ret.o652am
				}
			} else {
				if depth > float64(immobleDepthThreshold) {
					ret := lle.createLifeLossSet(v, lle.LethalityCurves[HighLethality], rng)
					llu65 += ret.u652am
					llo65 += ret.o652am
				} else {
					ret := lle.createLifeLossSet(v, lle.LethalityCurves[LowLethality], rng)
					llu65 += ret.u652am
					llo65 += ret.o652am
				}
			}
		}
		return consequences.Result{Headers: header, Result: []interface{}{llu65, llo65, llu65 + llo65}}, nil
	}
}
func (lle LifeLossEngine) createLifeLossSet(popset PopulationSet, lc LethalityCurve, rng *rand.Rand) PopulationSet {
	result := PopulationSet{0, 0, 0, 0}
	result.o652am = lle.evaluateLifeLoss(popset.o652am, lle.LethalityCurves[HighLethality], rng)
	result.o652pm = lle.evaluateLifeLoss(popset.o652pm, lle.LethalityCurves[HighLethality], rng)
	result.u652am = lle.evaluateLifeLoss(popset.u652am, lle.LethalityCurves[HighLethality], rng)
	result.u652pm = lle.evaluateLifeLoss(popset.u652pm, lle.LethalityCurves[HighLethality], rng)
	//log.Println(result)
	return result
}
func (lle LifeLossEngine) evaluateLifeLoss(populationRemaining int, lc LethalityCurve, rng *rand.Rand) int {
	result := 0
	for i := 0; i < populationRemaining; i++ {
		if lc.Sample() < rng.Float64() {
			result++
		}
	}
	return result
}
func evaluateMobility(s structures.StructureDeterministic, rng *rand.Rand) map[Mobility]PopulationSet {
	//determine based on age and disability
	result := make(map[Mobility]PopulationSet)
	mobileset := PopulationSet{0, 0, 0, 0}
	notmobileset := PopulationSet{0, 0, 0, 0}
	result[Mobile] = mobileset
	result[NotMobile] = notmobileset
	for i := 0; i < int(s.Pop2amo65); i++ {
		if rng.Float64() < .75 { //get this from nick
			popset := result[Mobile]
			popset.o652am = popset.o652am + 1
			result[Mobile] = popset
		} else {
			popset := result[NotMobile]
			popset.o652am = popset.o652am + 1
			result[NotMobile] = popset
		}
	}
	for i := 0; i < int(s.Pop2amu65); i++ {
		if rng.Float64() < .98 { //get this from nick
			popset := result[Mobile]
			popset.u652am = popset.u652am + 1
			result[Mobile] = popset
		} else {
			popset := result[NotMobile]
			popset.u652am = popset.u652am + 1
			result[NotMobile] = popset
		}
	}
	for i := 0; i < int(s.Pop2pmo65); i++ {
		if rng.Float64() < .75 {
			popset := result[Mobile]
			popset.o652pm = popset.o652pm + 1
			result[Mobile] = popset
		} else {
			popset := result[NotMobile]
			popset.o652pm = popset.o652pm + 1
			result[NotMobile] = popset
		}
	}
	for i := 0; i < int(s.Pop2pmu65); i++ {
		if rng.Float64() < .98 {
			popset := result[Mobile]
			popset.u652pm = popset.u652pm + 1
			result[Mobile] = popset
		} else {
			popset := result[NotMobile]
			popset.u652pm = popset.u652pm + 1
			result[NotMobile] = popset
		}
	}
	return result
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
