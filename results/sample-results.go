package results

import (
	"encoding/json"
	"fmt"
)

//SampleResult is a method to create a consequence result marshal it to json and print the original struct to string
func SampleResult() {
	header := []string{"structure fid", "structure damage", "content damage"}
	results := []interface{}{1, 5.0, 10.0}
	var ret = Consequence{Headers: header, Results: results}
	b, _ := json.Marshal(ret)
	fmt.Println(string(b))
	fmt.Println(ret)
}

//SampleResults is similar to SampleResult but it stores and writes multiple consequence results in a single consequences struct
func SampleResults() {
	header := []string{"structure fid", "structure damage", "content damage"}
	var rows []interface{}
	result := Consequences{IsTable: true}
	result.Headers = header
	result.Results = rows
	for i := 0; i < 10; i++ {
		results := []interface{}{1 * i, 5.0 * float64(i), 10.0 * float64(i)}
		row := Consequence{Headers: header, Results: results}
		result.AddConsequence(row)
	}
	b, _ := json.Marshal(result)
	fmt.Println(string(b))
	fmt.Println(result)
}
