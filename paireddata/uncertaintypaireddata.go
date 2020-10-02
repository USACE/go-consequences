package paireddata

import (
	"github.com/HenryGeorgist/go-statistics/statistics"
)

// UncertaintyPairedData is paired data where Y is a distribution
type UncertaintyPairedData struct {
	Xvals []float64
	Yvals []statistics.ContinuousDistribution
}

/* Needs work - this confuses things elsewhere..
// SampleValue implements UncertianValueSampler
func (p UncertaintyPairedData) SampleValue(inputValue interface{}, randomValue float64) float64 {
	xval, ok := inputValue.(float64)
	if !ok {
		return 0.0
	}
	if xval < p.Xvals[0] {
		return 0.0 //xval is less than lowest x value
	}
	size := len(p.Xvals)
	if xval >= p.Xvals[size-1] {
		return p.Yvals[size-1].InvCDF(randomValue) //xval yeilds largest y value
	}
	if xval == p.Xvals[0] {
		return p.Yvals[0].InvCDF(randomValue)
	}
	upper := sort.SearchFloat64s(p.Xvals, xval)
	//interpolate
	lower := upper - 1 // safe because we trapped the 0 case earlier
	lowerY := p.Yvals[lower].InvCDF(randomValue)
	upperY := p.Yvals[upper].InvCDF(randomValue)
	slope := (upperY - lowerY) / (p.Xvals[upper] - p.Xvals[lower])
	return lowerY + slope*(xval-p.Xvals[lower])
}
*/
func (p UncertaintyPairedData) SampleValueSampler(randomValue float64) ValueSampler {
	yVals := make([]float64, len(p.Yvals))
	for idx, dist := range p.Yvals {
		yVals[idx] = dist.InvCDF(randomValue)
	}
	return PairedData{Xvals: p.Xvals, Yvals: yVals}
}
