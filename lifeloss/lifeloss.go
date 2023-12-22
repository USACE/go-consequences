package lifeloss

import (
	"errors"
	"log"

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

func LifeLossHeader() []string {
	return []string{"ll_u65", "ll_o65", "ll_tot"}
}
func LifeLossDefaultResults() []interface{} {
	return []interface{}{0.0, 0.0, 0.0}
}
func computeLifeLoss(e hazards.HazardEvent, s structures.StructureDeterministic) (consequences.Result, error) {
	//reduce population somehow?
	if e.Has(hazards.Velocity) {
		sc, err := determineStability(s)
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
