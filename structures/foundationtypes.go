package structures

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/HydrologicEngineeringCenter/go-statistics/statistics"
)

type FoundationHeightUncertainty struct {
	DefaultDistribution statistics.ContinuousDistribution `json:"default"`
	VzoneDistribution   statistics.ContinuousDistribution `json:"vzone"`
}
type FoundationHeightUncertaintyContainer struct {
	DefaultDistribution statistics.ContinuousDistributionContainer `json:"default"`
	VzoneDistribution   statistics.ContinuousDistributionContainer `json:"vzone"`
}
type FoundationUncertainty struct {
	Values map[string]FoundationHeightUncertainty `json:"values"`
}
type FoundationUncertaintyContainer struct {
	Values map[string]FoundationHeightUncertaintyContainer `json:"values"`
}

func (fu *FoundationUncertainty) MarshalJSON() ([]byte, error) {
	var sb strings.Builder
	sb.WriteString("{\"values\":{")
	for k, v := range fu.Values {
		dist, err := v.MarshalJSON()
		if err != nil {
			return nil, err
		}
		sb.WriteString(fmt.Sprintf("\"%v\":%v,", k, string(dist)))
	}
	s := sb.String()
	s = strings.Trim(s, ",")
	s = fmt.Sprintf("%v}}", s)
	return []byte(s), nil
}
func (fu *FoundationUncertainty) Unmarshal(b []byte) error {
	var tmp FoundationUncertaintyContainer
	err := json.Unmarshal(b, &tmp)
	if err != nil {
		return err
	}
	fu.Values = make(map[string]FoundationHeightUncertainty)
	for k, v := range tmp.Values {
		fhu := FoundationHeightUncertainty{
			DefaultDistribution: v.DefaultDistribution.Value,
			VzoneDistribution:   v.VzoneDistribution.Value,
		}
		fu.Values[k] = fhu
	}
	return nil
}
func (fhu FoundationHeightUncertainty) Sample(firmZone string) statistics.ContinuousDistribution {
	if fhu.VzoneDistribution != nil {
		if strings.Compare("v", firmZone) == 0 {
			return fhu.VzoneDistribution
		}
	}
	return fhu.DefaultDistribution
}
func (fhu *FoundationHeightUncertainty) MarshalJSON() ([]byte, error) {
	distvs, err := statistics.Marshal(fhu.VzoneDistribution)
	if err != nil {
		return []byte{}, err
	}
	dists, err := statistics.Marshal(fhu.DefaultDistribution)
	if err != nil {
		return []byte{}, err
	}
	s := fmt.Sprintf("{\"default\":%v,\"vzone\":%v}", dists, distvs)
	return []byte(s), nil
}
func (fhu *FoundationHeightUncertainty) UnmarshalJSON(data []byte) error {
	var tmp FoundationHeightUncertaintyContainer
	err := json.Unmarshal(data, &tmp)
	if err != nil {
		return err
	}
	fhu.VzoneDistribution = tmp.VzoneDistribution.Value
	fhu.DefaultDistribution = tmp.DefaultDistribution.Value
	return nil
}
