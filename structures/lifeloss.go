package structures

import (
	_ "embed"
	"errors"
	"log"
	"math/rand"

	"github.com/HydrologicEngineeringCenter/go-statistics/paireddata"
	"github.com/USACE/go-consequences/consequences"
	"github.com/USACE/go-consequences/hazards"
)

//go:embed lowlethality.json
var DefaultLowLethalityBytes []byte

//go:embed highlethality.json
var DefaultHighLethalityBytes []byte

type LethalityZone int

const (
	LowLethality  LethalityZone = 1
	HighLethality LethalityZone = 2
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

func LifeLossHeader() []string {
	return []string{"ll_u65", "ll_o65", "ll_tot"}
}
func LifeLossDefaultResults() []interface{} {
	return []interface{}{0.0, 0.0, 0.0}
}
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
		depth = e.Depth()
	}
	if e.Has(hazards.Velocity) {
		//get velocity from the hazard
		velocity = e.Velocity()
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
