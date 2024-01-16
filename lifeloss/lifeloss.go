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

type LifeLossEngine struct {
	LethalityCurves   map[LethalityZone]LethalityCurve
	StabilityCriteria map[string]StabilityCriteria
	WarningSystem     warning.WarningResponseSystem
	ResultsHeader     []string
	SeedGenerator     *rand.Rand
}

// Init in the lifeloss package creates a life loss engine with the default settings
func Init(seed int64, warningSystem warning.WarningResponseSystem) LifeLossEngine {
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
	rng := rand.New(rand.NewSource(seed))

	return LifeLossEngine{LethalityCurves: lethalityCurves, StabilityCriteria: stabilityCriteria, WarningSystem: warningSystem, SeedGenerator: rng}
}

func LifeLossHeader() []string {
	return []string{"ll_u65", "ll_o65", "ll_tot"}
}
func LifeLossDefaultResults() []interface{} {
	var ll_u65, ll_o65, ll_tot int32
	ll_u65 = 0
	ll_o65 = 0
	ll_tot = 0
	return []interface{}{ll_u65, ll_o65, ll_tot}
}
func (le LifeLossEngine) ComputeLifeLoss(e hazards.HazardEvent, s structures.StructureDeterministic) (consequences.Result, error) {
	//reduce population based off of the warning system's warning function
	rng := rand.New(rand.NewSource(le.SeedGenerator.Int63()))
	remainingPop, _ := le.WarningSystem.WarningFunction()(s, e)
	if e.Has(hazards.DV) && e.Has(hazards.Depth) || e.Has(hazards.Velocity) && e.Has(hazards.Depth) {
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
			llo65 := applylethalityRateToPopulation(lethalityRate, remainingPop.Pop2amo65, rng)
			//llo65 += applylethalityRateToPopulation(lethalityRate, remainingPop.Pop2pmo65, rng)
			llu65 := applylethalityRateToPopulation(lethalityRate, remainingPop.Pop2amu65, rng)
			//llu65 += applylethalityRateToPopulation(lethalityRate, remainingPop.Pop2pmu65, rng)
			result := consequences.Result{Headers: LifeLossHeader(), Result: []interface{}{llu65, llo65, llu65 + llo65}}
			//log.Println(result)
			return result, nil
		} else {
			return le.submergenceCriteria(e, s, remainingPop, rng)
		}
	} else {
		//apply submergence criteria
		return le.submergenceCriteria(e, s, remainingPop, rng)
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
func (lle LifeLossEngine) submergenceCriteria(e hazards.HazardEvent, s structures.StructureDeterministic, remainingPop structures.PopulationSet, rng *rand.Rand) (consequences.Result, error) {
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

		mobilitySet := evaluateMobility(remainingPop, rng)
		var llu65 int32 = 0
		var llo65 int32 = 0
		for k, v := range mobilitySet {
			//apply to the appropriate age/time of day
			//log.Println(v)
			if k == Mobile {
				if depth > float64(mobileDepthThreshold) {
					ret := lle.createLifeLossSet(v, lle.LethalityCurves[HighLethality], rng)
					llu65 += ret.Pop2amu65
					llo65 += ret.Pop2amo65
				} else {
					ret := lle.createLifeLossSet(v, lle.LethalityCurves[LowLethality], rng)
					llu65 += ret.Pop2amu65
					llo65 += ret.Pop2amo65
				}
			} else {
				if depth > float64(immobleDepthThreshold) {
					ret := lle.createLifeLossSet(v, lle.LethalityCurves[HighLethality], rng)
					llu65 += ret.Pop2amu65
					llo65 += ret.Pop2amo65
				} else {
					ret := lle.createLifeLossSet(v, lle.LethalityCurves[LowLethality], rng)
					llu65 += ret.Pop2amu65
					llo65 += ret.Pop2amo65
				}
			}
		}
		return consequences.Result{Headers: header, Result: []interface{}{llu65, llo65, llu65 + llo65}}, nil
	}
}
func (lle LifeLossEngine) createLifeLossSet(popset structures.PopulationSet, lc LethalityCurve, rng *rand.Rand) structures.PopulationSet {
	result := structures.PopulationSet{0, 0, 0, 0}
	result.Pop2amo65 = lle.evaluateLifeLoss(popset.Pop2amo65, lle.LethalityCurves[HighLethality], rng)
	result.Pop2pmo65 = lle.evaluateLifeLoss(popset.Pop2pmo65, lle.LethalityCurves[HighLethality], rng)
	result.Pop2amu65 = lle.evaluateLifeLoss(popset.Pop2amu65, lle.LethalityCurves[HighLethality], rng)
	result.Pop2pmu65 = lle.evaluateLifeLoss(popset.Pop2pmu65, lle.LethalityCurves[HighLethality], rng)
	//log.Println(result)
	return result
}
func (lle LifeLossEngine) evaluateLifeLoss(populationRemaining int32, lc LethalityCurve, rng *rand.Rand) int32 {
	var result int32 = 0
	var i int32 = 0
	for i = 0; i < populationRemaining; i++ {
		if lc.Sample() < rng.Float64() {
			result++
		}
	}
	return result
}
func evaluateMobility(s structures.PopulationSet, rng *rand.Rand) map[Mobility]structures.PopulationSet {
	//determine based on age and disability
	result := make(map[Mobility]structures.PopulationSet)
	mobileset := structures.PopulationSet{0, 0, 0, 0}
	notmobileset := structures.PopulationSet{0, 0, 0, 0}
	result[Mobile] = mobileset
	result[NotMobile] = notmobileset
	for i := 0; i < int(s.Pop2amo65); i++ {
		if rng.Float64() < .75 { //get this from nick
			popset := result[Mobile]
			popset.Pop2amo65 = popset.Pop2amo65 + 1
			result[Mobile] = popset
		} else {
			popset := result[NotMobile]
			popset.Pop2amo65 = popset.Pop2amo65 + 1
			result[NotMobile] = popset
		}
	}
	for i := 0; i < int(s.Pop2amu65); i++ {
		if rng.Float64() < .98 { //get this from nick
			popset := result[Mobile]
			popset.Pop2amu65 = popset.Pop2amu65 + 1
			result[Mobile] = popset
		} else {
			popset := result[NotMobile]
			popset.Pop2amu65 = popset.Pop2amu65 + 1
			result[NotMobile] = popset
		}
	}
	for i := 0; i < int(s.Pop2pmo65); i++ {
		if rng.Float64() < .75 {
			popset := result[Mobile]
			popset.Pop2pmo65 = popset.Pop2pmo65 + 1
			result[Mobile] = popset
		} else {
			popset := result[NotMobile]
			popset.Pop2pmo65 = popset.Pop2pmo65 + 1
			result[NotMobile] = popset
		}
	}
	for i := 0; i < int(s.Pop2pmu65); i++ {
		if rng.Float64() < .98 {
			popset := result[Mobile]
			popset.Pop2pmu65 = popset.Pop2pmu65 + 1
			result[Mobile] = popset
		} else {
			popset := result[NotMobile]
			popset.Pop2pmu65 = popset.Pop2pmu65 + 1
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
