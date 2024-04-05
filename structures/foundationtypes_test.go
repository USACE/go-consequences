package structures

import (
	"fmt"
	"testing"

	"github.com/HydrologicEngineeringCenter/go-statistics/statistics"
)

func Test_FoundationTypes(t *testing.T) {

	m := make(map[string]FoundationHeightUncertainty, 0)
	m["Res1-1SNB_slab"] = FoundationHeightUncertainty{
		DefaultDistribution: statistics.TriangularDistribution{
			Min:        1,
			MostLikely: 1.5,
			Max:        2,
		},
		VzoneDistribution: statistics.TriangularDistribution{
			Min:        1,
			MostLikely: 1.5,
			Max:        2,
		},
	}
	m["Res1-1SNB_craw"] = FoundationHeightUncertainty{
		DefaultDistribution: statistics.EmpiricalDistribution{
			BinStarts: []float64{1, 2, 3, 4},
			BinWidth:  1,
			BinCounts: []int64{5, 3, 7, 2},
			MinValue:  1,
			MaxValue:  5,
		},
		VzoneDistribution: statistics.TriangularDistribution{
			Min:        1,
			MostLikely: 1.5,
			Max:        2,
		},
	}
	m["Res1-1SNB_Pier"] = FoundationHeightUncertainty{
		DefaultDistribution: statistics.EmpiricalDistribution{
			BinStarts: []float64{1, 2, 3, 4},
			BinWidth:  1,
			BinCounts: []int64{5, 3, 7, 2},
			MinValue:  1,
			MaxValue:  5,
		},
		VzoneDistribution: statistics.NormalDistribution{
			Mean:              1,
			StandardDeviation: 6,
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
