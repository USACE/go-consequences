package consequences

import (
	"github.com/USACE/go-consequences/paireddata"
)

type OccupancyType struct {
	Name            string
	Structuredamfun paireddata.ValueSampler
	Contentdamfun   paireddata.ValueSampler
}

func OccupancyTypeMap() map[string]OccupancyType {
	m := make(map[string]OccupancyType)
	m["RES1-1SNB"] = res11snb()
	//m["RES1-1SWB"] = res11swb()
	return m
}
func res11snb() OccupancyType {
	structurexs := []float64{-2.0, -1.0, 0.0, 1.0, 2.0, 3.0, 4.0, 5.0, 6.0, 7.0, 8.0, 9.0, 10.0, 11.0, 12.0, 13.0, 14.0, 15.0, 16.0}
	structureys := []float64{0.0, 2.5, 13.399999618530273, 23.299999237060547, 32.099998474121094, 40.099998474121094, 47.099998474121094, 53.200000762939453, 58.599998474121094, 63.200000762939453, 67.199996948242188, 70.5, 73.199996948242188, 75.4000015258789, 77.199996948242188, 78.5, 79.5, 80.199996948242188, 80.699996948242188}
	contentxs := []float64{-2.0, -1.0, 0.0, 1.0, 2.0, 3.0, 4.0, 5.0, 6.0, 7.0, 8.0, 9.0, 10.0, 11.0, 12.0, 13.0}
	contentys := []float64{0.0, 2.4000000953674316, 8.1000003814697266, 13.300000190734863, 17.899999618530273, 22.0, 25.700000762939453, 28.799999237060547, 31.5, 33.799999237060547, 35.700000762939453, 37.200000762939453, 38.400001525878906, 39.200000762939453, 39.700000762939453, 40.0}
	var structuredamagefunction = paireddata.PairedData{Xvals: structurexs, Yvals: structureys}
	var contentdamagefunction = paireddata.PairedData{Xvals: contentxs, Yvals: contentys}
	return OccupancyType{Name: "RES1-1SNB", Structuredamfun: structuredamagefunction, Contentdamfun: contentdamagefunction}
}
func res11swb() OccupancyType {
	//edits are not complete on this occtype
	structurexs := []float64{-8.0, -7.0, -2.0, -1.0, 0.0, 1.0, 2.0, 3.0, 4.0, 5.0, 6.0, 7.0, 8.0, 9.0, 10.0, 11.0, 12.0, 13.0, 14.0, 15.0, 16.0}
	structureys := []float64{0.0, 2.5, 13.399999618530273, 23.299999237060547, 32.099998474121094, 40.099998474121094, 47.099998474121094, 53.200000762939453, 58.599998474121094, 63.200000762939453, 67.199996948242188, 70.5, 73.199996948242188, 75.4000015258789, 77.199996948242188, 78.5, 79.5, 80.199996948242188, 80.699996948242188}
	contentxs := []float64{-2.0, -1.0, 0.0, 1.0, 2.0, 3.0, 4.0, 5.0, 6.0, 7.0, 8.0, 9.0, 10.0, 11.0, 12.0, 13.0}
	contentys := []float64{0.0, 2.4000000953674316, 8.1000003814697266, 13.300000190734863, 17.899999618530273, 22.0, 25.700000762939453, 28.799999237060547, 31.5, 33.799999237060547, 35.700000762939453, 37.200000762939453, 38.400001525878906, 39.200000762939453, 39.700000762939453, 40.0}
	var structuredamagefunction = paireddata.PairedData{Xvals: structurexs, Yvals: structureys}
	var contentdamagefunction = paireddata.PairedData{Xvals: contentxs, Yvals: contentys}
	return OccupancyType{Name: "RES1-1SWB", Structuredamfun: structuredamagefunction, Contentdamfun: contentdamagefunction}
}
