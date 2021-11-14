package paireddata

import (
	"testing"

	"github.com/HydrologicEngineeringCenter/go-statistics/statistics"
)

func Test_Uncertainty_centralTendency(t *testing.T) {
	upd := createSampleData()
	vs := upd.CentralTendency()
	pd, ok := vs.(PairedData)
	if ok {
		for idx, x := range pd.Xvals {
			if pd.Yvals[idx] != x {
				t.Error("values dont match")
			}
		}
	} else {
		t.Error("did not yeild paireddata")
	}

}
func createSampleData() UncertaintyPairedData {
	structurexs := []float64{0.0, 1.0, 2.0, 3.0, 4.0, 5.0, 6.0, 7.0, 8.0, 9.0, 10.0, 11.0, 12.0, 13.0, 14.0, 15.0, 16.0}
	structureydists := make([]statistics.ContinuousDistribution, 17)
	structureydists[0] = statistics.NormalDistribution{Mean: 0, StandardDeviation: 1.0}
	structureydists[1] = statistics.NormalDistribution{Mean: 1.0, StandardDeviation: 1.0}
	structureydists[2] = statistics.NormalDistribution{Mean: 2.0, StandardDeviation: 1.0}
	structureydists[3] = statistics.NormalDistribution{Mean: 3.0, StandardDeviation: 1.0}
	structureydists[4] = statistics.NormalDistribution{Mean: 4.0, StandardDeviation: 1.0}
	structureydists[5] = statistics.NormalDistribution{Mean: 5.0, StandardDeviation: 1.0}
	structureydists[6] = statistics.NormalDistribution{Mean: 6.0, StandardDeviation: 1.0}
	structureydists[7] = statistics.NormalDistribution{Mean: 7.0, StandardDeviation: 1.0}
	structureydists[8] = statistics.NormalDistribution{Mean: 8.0, StandardDeviation: 1.0}
	structureydists[9] = statistics.NormalDistribution{Mean: 9.0, StandardDeviation: 1.0}
	structureydists[10] = statistics.NormalDistribution{Mean: 10.0, StandardDeviation: 1.0}
	structureydists[11] = statistics.NormalDistribution{Mean: 11.0, StandardDeviation: 1.0}
	structureydists[12] = statistics.NormalDistribution{Mean: 12.0, StandardDeviation: 1.0}
	structureydists[13] = statistics.NormalDistribution{Mean: 13.0, StandardDeviation: 1.0}
	structureydists[14] = statistics.NormalDistribution{Mean: 14.0, StandardDeviation: 1.0}
	structureydists[15] = statistics.NormalDistribution{Mean: 15.0, StandardDeviation: 1.0}
	structureydists[16] = statistics.NormalDistribution{Mean: 16.0, StandardDeviation: 1.0}

	return UncertaintyPairedData{Xvals: structurexs, Yvals: structureydists}
}
