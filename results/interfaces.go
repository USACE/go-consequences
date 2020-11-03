package results

import (
	"encoding/json"
	"fmt"
	"strings"
)

type Consequences struct {
	Headers []string      `json:"headers"`
	Rows    []interface{} `json:"rows"`
}

type Consequence struct {
	Headers []string
	Results []interface{}
}
type ConsequenceAddable interface {
	AddConsequence(c Consequence)
}

func (c *Consequences) AddConsequence(cr Consequence) {
	c.Headers = cr.Headers
	fmt.Println("Appending")
	fmt.Println(cr.Results)
	c.Rows = append(c.Rows, cr.Results)
	fmt.Println(c.Rows)
}

/*func (c Consequences) MarshalJSON() ([]byte, error) {
	return make([]byte, 0), nil
}*/
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

func (c Consequence) String() string {
	if len(c.Headers) != len(c.Results) {
		return "mismatched lengths"
	}
	var ret string = "the consequences were:"
	for i, h := range c.Headers {
		ret += " " + h + " = " + fmt.Sprintf("%f", c.Results[i].(float64)) + ","
	}
	return strings.Trim(ret, ",")
}
