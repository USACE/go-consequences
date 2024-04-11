package structures

import (
	"fmt"
	"testing"

	"github.com/HydrologicEngineeringCenter/go-statistics/statistics"
)

func Test_FoundationTypes(t *testing.T) {

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
	b, err := f.MarshalJSON()
	if err != nil {
		t.Fail()
	}
	fmt.Println(string(b))
	var fmap FoundationUncertainty
	err = fmap.Unmarshal(b)
	if err != nil {
		t.Fail()
	}
	b, err = fmap.MarshalJSON()
	if err != nil {
		t.Fail()
	}
	fmt.Println(string(b))
}
