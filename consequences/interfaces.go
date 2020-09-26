package consequences

import (
	"fmt"
	"strings"
)

type ConsequenceReceptor interface {
	ComputeConsequences(event interface{}) ConsequenceDamageResult
}
type ConsequenceDamageResult struct {
	Headers []string
	Results []interface{}
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
