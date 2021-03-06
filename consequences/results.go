package consequences

import (
	"encoding/json"
	"strings"
)

//Result is a container to store a list of headers and a list of results
type Result struct {
	Headers []string      `json:"headers"`
	Result  []interface{} `json:"result"`
}

//Results stores a consequence struct and a boolean flagging the result as a table or not
type Results struct {
	IsTable bool
	Result  `json:"results"`
}

//ResultAddable gives me the ability to convert the results slice of interface into a table of slice of interface and store many results...
type ResultAddable interface {
	AddResult(c Result) //is this too confusing? it works, but is it confusing?
}

//AddResult fulfils the ConsequenceAddable interface on the Consequences struct
func (c *Results) AddResult(cr Result) {
	c.IsTable = true
	//todo check headers for equivalency...
	c.Headers = cr.Headers
	c.Result.Result = append(c.Result.Result, cr.Result)
}

//MarshalJSON a better printed version of results - this is my preferred way to print, but it is more complex
func (c Results) MarshalJSON() ([]byte, error) {
	s := "{\"consequences\":["
	for _, result := range c.Result.Result {
		s += "{\"consequence\":{\""
		vals, ok := result.([]interface{})
		if ok {
			for i, val := range c.Headers {
				value, _ := json.Marshal(vals[i])
				s += val + "\":" + string(value) + ",\""
			}
			s = strings.TrimRight(s, ",\"")
			s += "}},"
		}
	}
	s = strings.TrimRight(s, ",")
	s += "]}"
	return []byte(s), nil
}
func (c Result) MarshalJSON() ([]byte, error) {
	s := "{\"consequence\":{\""
	result := c.Result
	for i, val := range c.Headers {
		value, _ := json.Marshal(result[i])
		s += val + "\":" + string(value) + ",\""
	}
	s = strings.TrimRight(s, ",\"")
	s += "}}"
	return []byte(s), nil
}
