package results

import (
	"encoding/json"
	"fmt"
)

func SampleResult() {
	header := []string{"structure fid", "structure damage", "content damage"}
	results := []interface{}{1, 5.0, 10.0}
	var ret = Consequences{Headers: header, Results: results}
	b, _ := json.Marshal(ret)
	fmt.Println(string(b))
}

func SampleResults() {
	header := []string{"structure fid", "structure damage", "content damage"}
	var rows []interface{}
	result := Consequences{Headers: header, Results: rows}
	for i := 0; i < 10; i++ {
		results := []interface{}{1 * i, 5.0 * float64(i), 10.0 * float64(i)}
		row := Consequences{Headers: header, Results: results}
		result.AddConsequence(row)
	}
	b, _ := json.Marshal(result)
	fmt.Println(string(b))
}
