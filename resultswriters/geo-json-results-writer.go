package resultswriters

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/USACE/go-consequences/consequences"
)

type geoJsonResultsWriter struct {
	w                io.Writer
	hasWrittenHeader bool
	firstFeature     bool
}

func InitGeoJsonResultsWriter(w io.Writer) *geoJsonResultsWriter {
	return &geoJsonResultsWriter{
		w:            w,
		firstFeature: true,
	}
}

func (gjw *geoJsonResultsWriter) Write(r consequences.Result) {
	if !gjw.hasWrittenHeader {
		fmt.Fprint(gjw.w, `{"type":"FeatureCollection","features":[`)
		gjw.hasWrittenHeader = true
	}

	if !gjw.firstFeature {
		fmt.Fprint(gjw.w, ",")
	}
	gjw.firstFeature = false

	// 1. Unmarshal the existing result format
	// Based on your output, r.MarshalJSON() returns {"consequence": { "x":..., "y":..., ... }}
	var wrapper struct {
		Consequence map[string]interface{} `json:"consequence"`
	}
	b, _ := r.MarshalJSON()
	json.Unmarshal(b, &wrapper)

	data := wrapper.Consequence
	x, _ := data["x"].(float64)
	y, _ := data["y"].(float64)

	// 2. Map to GeoJSON Feature structure
	feature := map[string]interface{}{
		"type": "Feature",
		"geometry": map[string]interface{}{
			"type":        "Point",
			"coordinates": []float64{x, y},
		},
		"properties": data,
	}

	// 3. Remove coordinates from properties to avoid redundancy (optional)
	delete(data, "x")
	delete(data, "y")

	// 4. Stream to writer
	resultByte, _ := json.Marshal(feature)
	gjw.w.Write(resultByte)
}

func (gjw *geoJsonResultsWriter) Close() {
	if gjw.hasWrittenHeader {
		fmt.Fprint(gjw.w, "]}")
	}
	if w2, ok := gjw.w.(io.WriteCloser); ok {
		w2.Close()
	}
}

