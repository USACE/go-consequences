package consequences

import (
	"github.com/HydrologicEngineeringCenter/go-statistics/statistics"
	"github.com/USACE/go-consequences/geography"
	"github.com/USACE/go-consequences/hazards"
)

//Receptor is an interface for all things that can have consequences from a hazard event
type Receptor interface {
	Compute(event hazards.HazardEvent) (Result, error)
	Location() geography.Location
}

//Inventory provides a struct to allow for a slice of ConcequenceReceptor
type Inventory struct {
	Inventory []Receptor
}
type StreamProvider interface {
	ByFips(fipscode string, sp StreamProcessor)
	ByBbox(bbox geography.BBox, sp StreamProcessor)
}
type StreamProcessor func(str Receptor)
type ResultsWriter interface {
	Write(Result)
	Close()
}

//ParameterValue is a way to allow parameters to be either a scalar or a distribution.
type ParameterValue struct {
	Value interface{}
}

//CentralTendency on a ParameterValue is intended to help set structure values content values and foundaiton heights to central tendencies.
func (p ParameterValue) CentralTendency() float64 {
	pval, okf := p.Value.(float64) //if the ParameterValue.Value is a float - pass it on back.
	if okf {
		return pval
	}
	pvaldist, okd := p.Value.(statistics.ContinuousDistribution)
	if okd {
		return pvaldist.CentralTendency()
	}
	return 0
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
