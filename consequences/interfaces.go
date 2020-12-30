package consequences

import (
	"fmt"
	"strings"

	"github.com/HenryGeorgist/go-statistics/statistics"
)

type ConsequenceReceptor interface {
	ComputeConsequences(event interface{}) ConsequenceDamageResult
}
type ConsequenceDamageResult struct {
	Headers []string
	Results []interface{}
}
type Locatable interface {
	GetX() float64
	GetY() float64
}
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
func (c ConsequenceDamageResult) MarshalJSON() ([]byte, error) {
	return make([]byte, 0), nil
}
func (c ConsequenceDamageResult) String() string {
	if len(c.Headers) != len(c.Results) {
		return "mismatched lengths"
	}
	var ret string = "the consequences were:"
	for i, h := range c.Headers {
		ret += " " + h + " = " + fmt.Sprintf("%f", c.Results[i].(float64)) + ","
	}
	return strings.Trim(ret, ",")
}
