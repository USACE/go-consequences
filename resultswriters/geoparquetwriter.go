package resultswriters

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/USACE/go-consequences/consequences"
	"github.com/apache/arrow/go/v14/parquet/file"
)

type geoParquetResultsWriter struct {
	filepath string
	w        file.Writer
}

func InitGeoParquetResultsWriterFromFile(filepath string) (*geoParquetResultsWriter, error) {

	return &geoParquetResultsWriter{filepath: filepath}, nil
}

func (srw *geoParquetResultsWriter) Write(r consequences.Result) {
	//properties
	sp := "\"properties\": {\""
	//get the properties from the result
	x := 0.0
	y := 0.0
	result := r.Result
	for i, val := range r.Headers {
		value, _ := json.Marshal(result[i])
		sp += val + "\":" + string(value) + ",\""
		if val == "x" {
			x = result[i].(float64)
		}
		if val == "y" {
			y = result[i].(float64)
		}
	}
	atype := reflect.TypeOf(result[len(result)-1])
	if atype.Kind() == reflect.String {
		sp = strings.TrimRight(sp, ",\"")
		sp += "\""
	} else {
		sp = strings.TrimRight(sp, ",\"")
	}

	//write out a feature
	s := "{\"type\": \"Feature\",\n\"geometry\": {\n\"type\": \"Point\",\n\"coordinates\": ["
	//get the x and y
	s += fmt.Sprintf("%g, ", x)  //x value
	s += fmt.Sprintf("%g]\n", y) //y value
	//close out geometry
	s += "},\n"
	s += sp + "}},\n" //this comma might be bad news...

	//srw.S = s
}

func (srw *geoParquetResultsWriter) Close() {
	srw.w.Close()
}
