package consequences

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestResults(t *testing.T) {
	sampleResults()
}
func TestResult(t *testing.T) {
	sampleResult()
}

func sampleResult() {
	header := []string{"structure fid", "structure damage", "content damage"}
	results := []interface{}{1, 5.0, 10.0}
	var ret = Results{IsTable: false, Result: Result{Headers: header, Result: results}}
	b, _ := json.Marshal(ret)
	fmt.Println(string(b))
	fmt.Println(ret)
}

//sampleResults is similar to SampleResult but it stores and writes multiple consequence results in a single consequences struct
func sampleResults() {
	header := []string{"structure fid", "structure damage", "content damage"}
	var rows []interface{}
	result := Results{IsTable: true}
	result.Result.Headers = header
	result.Result.Result = rows
	for i := 0; i < 10; i++ {
		results := []interface{}{1 * i, 5.0 * float64(i), 10.0 * float64(i)}
		row := Result{Headers: header, Result: results}
		result.AddResult(row)
	}
	b, _ := json.Marshal(result)
	fmt.Println(string(b))
	fmt.Println(result)
}
