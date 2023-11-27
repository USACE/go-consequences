package structures

import (
	"errors"
	"log"
	"math/rand"

	"github.com/HydrologicEngineeringCenter/go-statistics/paireddata"
	"github.com/USACE/go-consequences/consequences"
	"github.com/USACE/go-consequences/hazards"
)

func computeLifeLoss(e hazards.HazardEvent, s StructureDeterministic) (consequences.Result, error) {
	//reduce population somehow?
	if e.Has(hazards.Velocity) {
		sc, err := determineStability(s)
		if err != nil {
			return consequences.Result{}, err
		}
		//apply building stability criteria
		if evaluateStability(e, sc) {
			//select high fataility rate
			lethalityRate := HighLethality.Sample()
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
func submergenceCriteria(e hazards.HazardEvent, s StructureDeterministic) (consequences.Result, error) {
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
func evaluateMobility(s StructureDeterministic) Mobility {
	//determine based on age and disability
	return Mobile
}

var HighLethality LethalityCurve = LethalityCurve{
	data: paireddata.PairedData{},
}
var LowLethality LethalityCurve = LethalityCurve{
	data: paireddata.PairedData{},
}

type LethalityCurve struct {
	data paireddata.PairedData
}

func (lc LethalityCurve) Sample() float64 {
	return lc.data.SampleValue(rand.Float64())
}

// implement high and low lethality
func (lc LethalityCurve) SampleWithSeededRand(rand rand.Rand) float64 {
	return lc.data.SampleValue(rand.Float64())
}

type StabilityCriteria struct {
	curve paireddata.PairedData
	//curve paireddata.UncertaintyPairedData
}

func (sc StabilityCriteria) Evaluate(e hazards.HazardEvent) bool {
	depth := 0.0
	velocity := 0.0
	if e.Has(hazards.Depth) {
		//get depth from the hazard
	}
	if e.Has(hazards.Velocity) {
		//get velocity from the hazard
	}

	/*shortcut if velocity is less than the min value...
	if math.min(sc.curve.Yvals)>velocity{
		return false
	}*/
	comparisonDepth := sc.curve.SampleYValue(velocity)

	return depth >= comparisonDepth
}
func evaluateStability(e hazards.HazardEvent, sc StabilityCriteria) bool {
	return sc.Evaluate(e)
}

func determineStability(s StructureDeterministic) (StabilityCriteria, error) {
	//add construction type to optional parameters and provide default criteria
	return StabilityCriteria{}, errors.New("implement me")
}

type Mobility uint

const (
	Unknown   Mobility = 0 //0
	Mobile    Mobility = 1 //1
	NotMobile Mobility = 2 //2
)
