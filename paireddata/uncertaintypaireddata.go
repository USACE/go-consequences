package paireddata

import (
	"github.com/HydrologicEngineeringCenter/go-statistics/statistics"
)

// UncertaintyPairedData is paired data where Y is a distribution
type UncertaintyPairedData struct {
	Xvals []float64
	Yvals []statistics.ContinuousDistribution
}

//SampleValueSampler implements UncertaintyValueSamplerSampler interface on the UncertaintyPairedData struct to produce a deterministic paireddata value for a given random number between 0 and 1
func (up UncertaintyPairedData) SampleValueSampler(randomValue float64) ValueSampler {
	yVals := make([]float64, len(up.Yvals))
	for idx, dist := range up.Yvals {
		yVals[idx] = dist.InvCDF(randomValue)
	}

	return PairedData{Xvals: up.Xvals, Yvals: yVals}
}
func (p UncertaintyPairedData) CentralTendency() ValueSampler {
	yVals := make([]float64, len(p.Yvals))
	for idx, dist := range p.Yvals {
		yVals[idx] = dist.CentralTendency()
	}
	return PairedData{Xvals: p.Xvals, Yvals: yVals}
}
