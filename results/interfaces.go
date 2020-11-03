package results

import (
	"fmt"
	"strings"
)

type Consequences struct {
	Headers []string      `json:"headers"`
	Results []interface{} `json:"results"`
	IsTable bool
}

//ConsequenceAddable gives me the ability to convert the results slice of interface into a table of slice of interface and store many results...
type ConsequenceAddable interface {
	AddConsequence(c Consequences) //is this too confusing? it works, but is it confusing?
}

func (c *Consequences) AddConsequence(cr Consequences) {
	c.IsTable = true
	c.Headers = cr.Headers
	c.Results = append(c.Results, cr.Results)
}

/* a better printed version of results - this is my preferred way to print, but it is more complex
func (c Consequence) MarshalJSON() ([]byte, error) {
	s := "{\"consequence\":{\""
	for i, val := range c.Headers {
		value, _ := json.Marshal(c.Results[i])
		s += val + "\":" + string(value) + ",\""
	}
	s = strings.TrimRight(s, ",\"")
	s += "}}"
	return []byte(s), nil
}
*/
func (c Consequences) String() string {
	if c.IsTable {
		return "Im a table!" //todo implement me
	}
	if len(c.Headers) != len(c.Results) {
		return "mismatched lengths"
	}
	var ret string = "the consequences were:"
	for i, h := range c.Headers {
		ret += " " + h + " = " + fmt.Sprintf("%f", c.Results[i].(float64)) + ","
	}
	return strings.Trim(ret, ",")
}
