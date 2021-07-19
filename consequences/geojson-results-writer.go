package consequences

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"reflect"
	"strings"
)

type geoJsonResultsWriter struct {
	filepath             string
	w                    io.Writer
	S                    string
	HeaderHasBeenWritten bool
	HasClosed            bool
}

func InitGeoJsonResultsWriterFromFile(filepath string) (*geoJsonResultsWriter, error) {
	w, err := os.OpenFile(filepath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0600)
	if err != nil {
		return &geoJsonResultsWriter{}, err
	}
	return &geoJsonResultsWriter{filepath: filepath, w: w}, nil
}
func InitGeoJsonResultsWriter(w io.Writer) *geoJsonResultsWriter {
	return &geoJsonResultsWriter{filepath: "not applicapble", w: w}
}
func (srw *geoJsonResultsWriter) Write(r Result) {
	if !srw.HeaderHasBeenWritten {
		fmt.Fprintf(srw.w, "{\"type\": \"FeatureCollection\",\n\"features\":[")
		srw.HeaderHasBeenWritten = true
	}
	if srw.S != "" {
		fmt.Fprintf(srw.w, srw.S)
	}
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

	srw.S = s
}

func (srw *geoJsonResultsWriter) Close() {
	if srw.S != "" {
		srw.S = strings.TrimRight(srw.S, ",\n")
		fmt.Fprintf(srw.w, srw.S)
		fmt.Fprintf(srw.w, "\n]}")
		srw.S = ""
		srw.HasClosed = true
	}
	w2, ok := srw.w.(io.WriteCloser)
	if ok {
		w2.Close()
	}
}
