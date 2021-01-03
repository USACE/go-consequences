package consequences

import (
	"github.com/HenryGeorgist/go-statistics/statistics"
)

//ConsequenceReceptor is an interface for all things that can have consequences from a hazard event
type ConsequenceReceptor interface {
	ComputeConsequences(event interface{}) Results
}

//Locatable is an interface that defines that a thing can have an x and y location
type Locatable interface {
	GetX() float64
	GetY() float64
}

//ParameterValue is a way to allow parameters to be either a scalar or a distribution.
type ParameterValue struct {
	Value interface{}
}

//SampleValue on a ParameterValue is intended to help set structure values content values and foundaiton heights to uncertain parameters - this is a first draft of this interaction.
func (p ParameterValue) SampleValue(input interface{}) float64 {

	pval, okf := p.Value.(float64) //if the ParameterValue.Value is a float - pass it on back.
	if okf {
		return pval
	}
	pvaldist, okd := p.Value.(statistics.ContinuousDistribution)
	if okd {
		inval, ok := input.(float64)
		if ok {
			return pvaldist.InvCDF(inval)
		}
	}

	return 0
}
