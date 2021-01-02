package paireddata

import (
	"sort"
)


// PairedData is paired data x and y values 
type PairedData struct {
	Xvals []float64
	Yvals []float64
}

// SampleValue implements ValueSampler
func (p PairedData) SampleValue(inputValue interface{}) float64 {
	xval, ok := inputValue.(float64)
	if !ok {
		return 0.0
	}
	if xval < p.Xvals[0] {
		return 0.0 //xval is less than lowest x value
	}
	size := len(p.Xvals)
	if xval >= p.Xvals[size-1] {
		return p.Yvals[size-1] //xval yeilds largest y value
	}
	if xval == p.Xvals[0] {
		return p.Yvals[0]
	}
	upper := sort.SearchFloat64s(p.Xvals, xval)
	//interpolate
	lower := upper - 1 // safe because we trapped the 0 case earlier
	slope := (p.Yvals[upper] - p.Yvals[lower]) / (p.Xvals[upper] - p.Xvals[lower])
	a := p.Yvals[lower]
	return a + slope*(xval-p.Xvals[lower])
}
