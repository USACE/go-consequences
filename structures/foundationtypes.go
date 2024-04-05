package structures

import (
	"encoding/json"
	"fmt"
	"os"
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

func InitFoundationUncertaintyFromFile(file string) (*FoundationUncertainty, error) {
	fbytes, err := os.ReadFile(file)
	fu := FoundationUncertainty{
		Values: map[string]FoundationHeightUncertainty{},
	}
	if err != nil {
		return &fu, err
	}
	err = fu.Unmarshal(fbytes)
	return &fu, err
}
func InitFoundationUncertainty() (*FoundationUncertainty, error) {
	//todo update to file in resources
	m := make(map[string]FoundationHeightUncertainty, 0)
	m["default_slab"] = FoundationHeightUncertainty{
		DefaultDistribution: statistics.NormalDistribution{
			Mean:              0.53,
			StandardDeviation: 1.01,
		},
		VzoneDistribution: statistics.NormalDistribution{
			Mean:              0.44,
			StandardDeviation: 1.08,
		},
	}
	m["default_craw"] = FoundationHeightUncertainty{
		DefaultDistribution: statistics.NormalDistribution{
			Mean:              0.53,
			StandardDeviation: 1.01,
		},
		VzoneDistribution: statistics.NormalDistribution{
			Mean:              1.74,
			StandardDeviation: 0.78,
		},
	}
	m["default_base"] = FoundationHeightUncertainty{
		DefaultDistribution: statistics.NormalDistribution{
			Mean:              0.53,
			StandardDeviation: 1.01,
		},
		VzoneDistribution: statistics.NormalDistribution{
			Mean:              2.19,
			StandardDeviation: 1.56,
		},
	}
	m["default_pier"] = FoundationHeightUncertainty{
		DefaultDistribution: statistics.NormalDistribution{
			Mean:              0.53,
			StandardDeviation: 1.01,
		},
		VzoneDistribution: statistics.NormalDistribution{
			Mean:              7.6,
			StandardDeviation: 3.46,
		},
	}
	//RES2
	m["RES2_slab"] = FoundationHeightUncertainty{
		DefaultDistribution: statistics.NormalDistribution{
			Mean:              2.0,
			StandardDeviation: 0.93,
		},
		VzoneDistribution: statistics.NormalDistribution{
			Mean:              2.0,
			StandardDeviation: 0.93,
		},
	}
	m["RES2_craw"] = FoundationHeightUncertainty{
		DefaultDistribution: statistics.NormalDistribution{
			Mean:              2.0,
			StandardDeviation: 0.93,
		},
		VzoneDistribution: statistics.NormalDistribution{
			Mean:              2.0,
			StandardDeviation: 0.93,
		},
	}
	m["RES2_base"] = FoundationHeightUncertainty{
		DefaultDistribution: statistics.NormalDistribution{
			Mean:              2.0,
			StandardDeviation: 0.93,
		},
		VzoneDistribution: statistics.NormalDistribution{
			Mean:              2.0,
			StandardDeviation: 0.93,
		},
	}
	m["RES2_pier"] = FoundationHeightUncertainty{
		DefaultDistribution: statistics.NormalDistribution{
			Mean:              2.0,
			StandardDeviation: 0.93,
		},
		VzoneDistribution: statistics.NormalDistribution{
			Mean:              2.0,
			StandardDeviation: 0.93,
		},
	}
	//RES1,RES3A RES3B
	m["RES1_RES3A_RES3B_slab"] = FoundationHeightUncertainty{
		DefaultDistribution: statistics.NormalDistribution{
			Mean:              0.77,
			StandardDeviation: 0.91,
		},
		VzoneDistribution: statistics.NormalDistribution{
			Mean:              0.77,
			StandardDeviation: 0.91,
		},
	}
	m["RES1_RES3A_RES3B_craw"] = FoundationHeightUncertainty{
		DefaultDistribution: statistics.NormalDistribution{
			Mean:              1.74,
			StandardDeviation: 0.78,
		},
		VzoneDistribution: statistics.NormalDistribution{
			Mean:              1.74,
			StandardDeviation: 0.78,
		},
	}
	m["RES1_RES3A_RES3B_base"] = FoundationHeightUncertainty{
		DefaultDistribution: statistics.NormalDistribution{
			Mean:              2.91,
			StandardDeviation: 1.56,
		},
		VzoneDistribution: statistics.NormalDistribution{
			Mean:              2.91,
			StandardDeviation: 1.56,
		},
	}
	m["RES1_RES3A_RES3B_pier"] = FoundationHeightUncertainty{
		DefaultDistribution: statistics.NormalDistribution{
			Mean:              7.6,
			StandardDeviation: 3.46,
		},
		VzoneDistribution: statistics.NormalDistribution{
			Mean:              7.6,
			StandardDeviation: 3.46,
		},
	}
	f := FoundationUncertainty{
		Values: m,
	}
	return &f, nil
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
