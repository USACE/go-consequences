package lifeloss

import (
	"errors"

	"github.com/HydrologicEngineeringCenter/go-statistics/paireddata"
	"github.com/USACE/go-consequences/hazards"
	"github.com/USACE/go-consequences/structures"
)

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

func determineStability(s structures.StructureDeterministic) (StabilityCriteria, error) {
	//add construction type to optional parameters and provide default criteria
	return StabilityCriteria{}, errors.New("implement me")
}
