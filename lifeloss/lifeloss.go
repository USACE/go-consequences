package lifeloss

import (
	"encoding/json"
	"log"
	"math/rand"

	"github.com/HydrologicEngineeringCenter/go-statistics/paireddata"
	"github.com/HydrologicEngineeringCenter/go-statistics/statistics"
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

type LifeLossProcess interface {
	RedistributePopulation(e hazards.HazardEvent, s structures.StructureDeterministic) (structures.PopulationSet, error)
	EvaluateStability(e hazards.HazardEvent, s structures.StructureDeterministic) (Stability, error)
	ComputeLifeLoss(e hazards.HazardEvent, s structures.StructureDeterministic, stability Stability) (consequences.Result, error)
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
	return []string{"ll_u65", "ll_o65", "ll_tot"} //@TODO: Consider adding structure stability state, fatality rates, and sampled hazard parameters
}
func LifeLossDefaultResults() []interface{} {
	var ll_u65, ll_o65, ll_tot int32
	ll_u65 = 0
	ll_o65 = 0
	ll_tot = 0
	return []interface{}{ll_u65, ll_o65, ll_tot}
}
func (le LifeLossEngine) RedistributePopulation(e hazards.HazardEvent, s structures.StructureDeterministic) (structures.StructureDeterministic, error) {
	remainingPop, _ := le.WarningSystem.WarningFunction()(s, e)
	sout := s.Clone()
	sout.PopulationSet = remainingPop
	return sout, nil
}
func (le LifeLossEngine) EvaluateStabilityCriteria(e hazards.HazardEvent, s structures.StructureDeterministic) (Stability, error) {
	if e.Has(hazards.DV) && e.Has(hazards.Depth) || e.Has(hazards.Velocity) && e.Has(hazards.Depth) {
		sc, err := le.determineStability(s)
		if err != nil {
			return Stable, err
		}
		result := sc.Evaluate(e)
		return result, nil
	} else {
		return Stable, nil
	}
}
func (le LifeLossEngine) ComputeLifeLoss(e hazards.HazardEvent, s structures.StructureDeterministic, stability Stability) (consequences.Result, error) {
	//reduce population based off of the warning system's warning function
	rng := rand.New(rand.NewSource(le.SeedGenerator.Int63()))
	remainingPop := s.PopulationSet
	//apply building stability criteria
	if stability == Collapsed {
		//log.Println("Stability Based Lifeloss")
		//select high fataility rate
		lethalityRate := le.LethalityCurves[HighLethality].Sample()
		log.Printf("high lethality rate: %v\n", lethalityRate)
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
	//log.Println("Submergence Based Lifeloss")
	header := LifeLossHeader()
	depth := e.Depth()
	if depth < 0.0 {
		//no lifeloss

		result := LifeLossDefaultResults()
		return consequences.Result{Headers: header, Result: result}, nil
	} else {
		//for all ages of individuals using different probabilities to assign mobility
		// for over and under 65 based on nsi attribute of "percent not mobile"
		immobleDepthThreshold := (float64(s.NumStories) - 1.0) * 9.0 //hard coded to 9 feet, probably should make it an attribute on the structure inventory?
		mobileDepthThreshold := (float64(s.NumStories) - 1.0) * 9.0
		hasAtticAccess := false //@TODO: probability of if attic access is determined by random number .95 (default)
		depthFromCeilingDistribution := statistics.TriangularDistribution{
			Min:        0.5,
			MostLikely: 1,
			Max:        1.5,
		}
		depthOnRoofDistribution := statistics.TriangularDistribution{
			Min:        3,
			MostLikely: 4,
			Max:        5,
		}
		depthFromFloorImmobileDistribution := statistics.TriangularDistribution{
			Min:        4,
			MostLikely: 5,
			Max:        6,
		}
		depthFromCeiling := depthFromCeilingDistribution.InvCDF(rng.Float64())
		immobleDepthThreshold += 5.0 + s.FoundHt
		mobileDepthThreshold += 9.0 + s.FoundHt - depthFromCeiling //9-1... height of ceiling minus 1 foot
		if hasAtticAccess {                                        //@TODO: "Roof access from attic is another random number .9(default)" - stored at the occupancy type level in lifesim.
			mobileDepthThreshold += depthFromCeiling + depthFromFloorImmobileDistribution.InvCDF(rng.Float64()) + depthOnRoofDistribution.InvCDF(rng.Float64()) //depth from ceiling + attic access + high hazard depth 4 feet above top of roof (should be a random number triangular distribution 4,5,6)
			immobleDepthThreshold += 9.0
		}

		mobilitySet := evaluateMobility(remainingPop, rng)
		var llu65 int32 = 0
		var llo65 int32 = 0
		lowLethality := lle.LethalityCurves[LowLethality].SampleWithSeededRand(rng)
		log.Printf("low lethality rate: %v\n", lowLethality)
		highLethality := lle.LethalityCurves[HighLethality].SampleWithSeededRand(rng)
		log.Printf("high lethality rate: %v\n", highLethality)
		for k, v := range mobilitySet {
			//apply to the appropriate age/time of day
			//log.Println(v)
			if k == Mobile {
				if depth > float64(mobileDepthThreshold) {
					ret := lle.createLifeLossSet(v, highLethality, rng)
					llu65 += ret.Pop2amu65
					llo65 += ret.Pop2amo65
				} else {
					ret := lle.createLifeLossSet(v, lowLethality, rng)
					llu65 += ret.Pop2amu65
					llo65 += ret.Pop2amo65
				}
			} else {
				if depth > float64(immobleDepthThreshold) {
					ret := lle.createLifeLossSet(v, highLethality, rng)
					llu65 += ret.Pop2amu65
					llo65 += ret.Pop2amo65
				} else {
					ret := lle.createLifeLossSet(v, lowLethality, rng)
					llu65 += ret.Pop2amu65
					llo65 += ret.Pop2amo65
				}
			}
		}
		return consequences.Result{Headers: header, Result: []interface{}{llu65, llo65, llu65 + llo65}}, nil
	}
}
func (lle LifeLossEngine) createLifeLossSet(popset structures.PopulationSet, lethalityRate float64, rng *rand.Rand) structures.PopulationSet {
	result := structures.PopulationSet{0, 0, 0, 0}
	result.Pop2amo65 = lle.evaluateLifeLoss(popset.Pop2amo65, lethalityRate, rng)
	result.Pop2pmo65 = lle.evaluateLifeLoss(popset.Pop2pmo65, lethalityRate, rng)
	result.Pop2amu65 = lle.evaluateLifeLoss(popset.Pop2amu65, lethalityRate, rng)
	result.Pop2pmu65 = lle.evaluateLifeLoss(popset.Pop2pmu65, lethalityRate, rng)
	//log.Println(result)
	return result
}
func (lle LifeLossEngine) evaluateLifeLoss(populationRemaining int32, lethalityRate float64, rng *rand.Rand) int32 {
	var result int32 = 0
	var i int32 = 0
	for i = 0; i < populationRemaining; i++ {
		if lethalityRate < rng.Float64() {
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
