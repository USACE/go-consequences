package results

type Consequence struct {
	Headers []string      `json:"headers"`
	Results []interface{} `json:"results"`
}
type Consequences struct {
	IsTable bool
	Consequence
}

//ConsequenceAddable gives me the ability to convert the results slice of interface into a table of slice of interface and store many results...
type ConsequenceAddable interface {
	AddConsequence(c Consequence) //is this too confusing? it works, but is it confusing?
}

func (c *Consequences) AddConsequence(cr Consequence) {
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
	return "Im not a table?"
}
