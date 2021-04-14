package consequences

import (
	"encoding/json"
	"fmt"
	"strings"
)

type geoJsonVirtualResultsWriter struct {
	data                 strings.Builder
	S                    string
	HeaderHasBeenWritten bool
}

func InitVirtualGeoJsonResultsWriter() *geoJsonVirtualResultsWriter {
	var b strings.Builder
	return &geoJsonVirtualResultsWriter{data: b}
}
func (srw *geoJsonVirtualResultsWriter) Write(r Result) {
	if !srw.HeaderHasBeenWritten {
		fmt.Fprintf(&srw.data, "{\"type\": \"FeatureCollection\",\n\"features\":[")
		srw.HeaderHasBeenWritten = true
	}
	if srw.S != "" {
		fmt.Fprintf(&srw.data, srw.S)
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
	sp = strings.TrimRight(sp, ",\"")

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
func (srw *geoJsonVirtualResultsWriter) Close() {
	if srw.S != "" {
		srw.S = strings.TrimRight(srw.S, ",\n")
		fmt.Fprintf(&srw.data, srw.S)
		fmt.Fprintf(&srw.data, "\n]}")
		srw.S = ""
	}
}
func (srw *geoJsonVirtualResultsWriter) Bytes() []byte {
	srw.Close()
	return []byte(srw.data.String())
}
