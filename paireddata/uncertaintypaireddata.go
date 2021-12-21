package paireddata

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/HydrologicEngineeringCenter/go-statistics/statistics"
)

// UncertaintyPairedData is paired data where Y is a distribution
type UncertaintyPairedData struct {
	Xvals []float64                           `json:"xvalues"`
	Yvals []statistics.ContinuousDistribution `json:"ydistributions"`
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
func (upd UncertaintyPairedData) MarshalJSON() ([]byte, error) {
	s := "{\"xvalues\":"
	ab, err := json.Marshal(upd.Xvals)
	if err != nil {
		return nil, errors.New("paireddata: could not marshal uncertain paired data xvalues")
	}
	s += fmt.Sprintf("%v", string(ab))
	s += ",\"ydistributions\":["
	for _, c := range upd.Yvals {
		dist, err := statistics.Marshal(c)
		if err != nil {
			return nil, errors.New("paireddata: could not marshal uncertain paired data yvalue")
		}
		s += fmt.Sprintf("%v,", dist)
	}
	s = strings.TrimRight(s, ",")
	s += "]}"
	return []byte(s), nil
}
func (upd *UncertaintyPairedData) UnmarshalJSON(b []byte) error {
	m := map[string]interface{}{}
	if err := json.Unmarshal(b, &m); err != nil {
		return err
	}
	valueBytes, err := json.Marshal(m["xvalues"])
	if err != nil {
		return err
	}
	var values []float64
	if err = json.Unmarshal(valueBytes, &values); err != nil {
		return err
	}
	upd.Xvals = values
	//now for the harder part of the distributions.
	distArrayBytes, err := json.Marshal(m["ydistributions"])
	if err != nil {
		return err
	}
	var distributions []statistics.ContinuousDistributionContainer
	if err = json.Unmarshal(distArrayBytes, &distributions); err != nil {
		return err
	}
	var ydist []statistics.ContinuousDistribution
	for _, c := range distributions {
		ydist = append(ydist, c.Value)
	}
	upd.Yvals = ydist
	return nil
}
