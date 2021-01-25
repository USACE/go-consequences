package consequences

import (
	"github.com/HenryGeorgist/go-statistics/statistics"
)

//Receptor is an interface for all things that can have consequences from a hazard event
type Receptor interface {
	Compute(event interface{}) Results
}

//
// work in progress
//

//Inventory provides a struct to allow for a slice of ConcequenceReceptor
type Inventory struct {
	Inventory []Receptor
}

//Provider defines an interface to provide a consequences Inventory
type Provider interface {
	GetInventoryBoundingBox(bb BoundingBox) (Inventory, error)
	GetInventoryFIPS(fc FIPS) (Inventory, error)
	GetInventoryFile(filePath string) (Inventory, error)
	//ProvideStructure(location Locatable) ConsequencesReceptor
	//ProvideStructure(fdId string) ConsequencesReceptor
}

//
// End work in progress
//

//Locatable is an interface that defines that a thing can have an x and y location
type Locatable interface {
	GetX() float64
	GetY() float64
}

//BoundingBox represents a rectangular area by extents.
type BoundingBox struct {
	//need to support multiple needs here everyone treats this different...
}

//FIPS is the Federal Information Processing Standard
type FIPS struct {
	Code string
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
