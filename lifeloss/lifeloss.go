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
	if e.Has(hazards.Velocity) {
		sc, err := le.determineStability(s)
		if err != nil {
			return consequences.Result{}, err
		}
		//apply building stability criteria
		if sc.Evaluate(e) {
			//select high fataility rate
			lethalityRate := 987654321.0 //HighLethality.Sample()
			//apply same fatality rate to everyone
			log.Println(lethalityRate)
			return consequences.Result{}, errors.New("under construction")
		} else {
			return submergenceCriteria(e, s)
		}
	} else {
		//apply submergence criteria
		return submergenceCriteria(e, s)
	}
}
func submergenceCriteria(e hazards.HazardEvent, s structures.StructureDeterministic) (consequences.Result, error) {
	//apply submergence criteria
	depth := 0.0
	//set depth
	if depth > 0.0 {
		//no lifeloss
		return consequences.Result{}, errors.New("under construction")
	} else {
		//for all individuals using different probabilities for over and under 65 based on nsi attributes
		mobility := evaluateMobility(s)
		if mobility == Mobile {

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

	return le.StabilityCriteria["woodanchored"], nil
}
