package lifeloss

import (
	"encoding/json"
	"errors"
	"log"

	"github.com/HydrologicEngineeringCenter/go-statistics/paireddata"
	"github.com/USACE/go-consequences/consequences"
	"github.com/USACE/go-consequences/hazards"
	"github.com/USACE/go-consequences/structures"
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
	return LifeLossEngine{LethalityCurves: lethalityCurves, StabilityCriteria: stabilityCriteria}
}

func LifeLossHeader() []string {
	return []string{"ll_u65", "ll_o65", "ll_tot"}
}
func LifeLossDefaultResults() []interface{} {
	return []interface{}{0.0, 0.0, 0.0}
}
func (le LifeLossEngine) ComputeLifeLoss(e hazards.HazardEvent, s structures.StructureDeterministic) (consequences.Result, error) {
	//reduce population somehow?
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
			return consequences.Result{}, nil
		} else {
			return le.submergenceCriteria(e, s)
		}
	} else {
		//apply submergence criteria
		return le.submergenceCriteria(e, s)
	}
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
