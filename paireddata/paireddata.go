package paireddata

import (
	"sort"
)

// PairedData is paired data x and y values
type PairedData struct {
	Xvals []float64 `json:"xvalues"`
	Yvals []float64 `json:"yvalues"`
}

//SampleValue implements ValueSampler
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
func (p PairedData) IsMonotonicallyIncreasing() bool {
	monotonic := true
	prevYval := p.Yvals[0]
	for i := 1; i < len(p.Yvals); i++ {
		if prevYval > p.Yvals[i] {
			monotonic = false
		}
	}
	return monotonic
}
func (p *PairedData) ForceMonotonicInRange(min float64, max float64) {
	prevYval := min
	update := make([]float64, len(p.Yvals))
	for idx, y := range p.Yvals {
		if prevYval > y {
			update[idx] = prevYval
		} else {
			if y > max {
				update[idx] = max
				prevYval = max
			} else {
				update[idx] = y
				prevYval = y
			}
		}
	}
	p.Yvals = update
}
