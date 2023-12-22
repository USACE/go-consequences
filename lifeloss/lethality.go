package lifeloss

import (
	_ "embed"
	"math/rand"

	"github.com/HydrologicEngineeringCenter/go-statistics/paireddata"
)

//go:embed lowlethality.json
var DefaultLowLethalityBytes []byte

//go:embed highlethality.json
var DefaultHighLethalityBytes []byte

type LethalityZone int

const (
	LowLethality  LethalityZone = 1
	HighLethality LethalityZone = 2
)

type LethalityCurve struct {
	data paireddata.PairedData
}

func (lc LethalityCurve) Sample() float64 {
	return lc.data.SampleValue(rand.Float64())
}

// implement high and low lethality
func (lc LethalityCurve) SampleWithSeededRand(rand rand.Rand) float64 {
	return lc.data.SampleValue(rand.Float64())
}
