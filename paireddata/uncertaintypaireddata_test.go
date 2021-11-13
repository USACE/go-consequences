package paireddata

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/HydrologicEngineeringCenter/go-statistics/statistics"
)

func Test_UncertiantyCentralTendency(t *testing.T) {

	contentsalinityxs := []float64{-1, -0.5, 0, 0.5, 1, 2, 3, 5, 7, 10}
	contentsalinityydists := make([]statistics.ContinuousDistribution, 10)
	contentsalinityydists[0], _ = statistics.InitDeterministic(0.0)
	contentsalinityydists[1], _ = statistics.InitDeterministic(0.0)
	contentsalinityydists[2], _ = statistics.Init([]float64{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}, []int64{250846, 523969, 523033, 481164, 496674, 219830, 146134, 71728, 25676, 8546})
	contentsalinityydists[3], _ = statistics.Init([]float64{3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32, 33, 34}, []int64{357, 1271, 2468, 4055, 5381, 6314, 6594, 6934, 7341, 7813, 7939, 8463, 8554, 6043, 3978, 3671, 3388, 2989, 2754, 2386, 2108, 1764, 1437, 1060, 804, 590, 482, 400, 315, 209, 151, 35})
	contentsalinityydists[4], _ = statistics.Init([]float64{7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32, 33, 34, 35, 36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 47, 48, 49, 50, 51, 52, 53}, []int64{187, 579, 919, 1430, 2080, 2293, 2355, 2305, 2454, 2549, 2565, 2762, 2789, 2992, 3052, 3287, 3253, 3443, 3548, 3700, 2270, 1608, 1666, 1725, 1784, 1809, 1822, 1799, 1838, 1707, 1591, 1447, 1339, 1271, 1143, 985, 917, 768, 656, 546, 407, 333, 203, 111, 102, 57, 14})
	contentsalinityydists[5], _ = statistics.Init([]float64{13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32, 33, 34, 35, 36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 47, 48, 49, 50, 51, 52, 53, 54, 55, 56, 57, 58, 59, 60, 61, 62, 63, 64}, []int64{120, 354, 607, 792, 1141, 1233, 1188, 1122, 1065, 1096, 1208, 1256, 1401, 1438, 1590, 1817, 1932, 2160, 2297, 2413, 2525, 2834, 2988, 2998, 1341, 1423, 1440, 1321, 1408, 1300, 1217, 1153, 1103, 1053, 964, 896, 816, 734, 685, 600, 505, 412, 379, 296, 215, 159, 124, 116, 77, 78, 39, 11})
	contentsalinityydists[6], _ = statistics.Init([]float64{20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32, 33, 34, 35, 36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 47, 48, 49, 50, 51, 52, 53, 54, 55, 56, 57, 58, 59, 60, 61, 62, 63, 64, 65, 66, 67, 68, 69, 70, 71, 72, 73, 74, 75, 76, 77, 78, 79, 80, 81, 82, 83}, []int64{0, 181, 317, 402, 579, 668, 597, 605, 641, 673, 759, 736, 830, 896, 1005, 953, 1019, 1073, 1247, 1411, 1509, 1751, 1891, 1440, 860, 816, 740, 703, 655, 612, 584, 600, 576, 509, 512, 493, 468, 409, 441, 410, 391, 376, 318, 308, 306, 267, 253, 215, 219, 207, 194, 164, 173, 141, 144, 119, 111, 97, 79, 66, 40, 38, 20, 7})
	contentsalinityydists[7], _ = statistics.Init([]float64{30, 31, 32, 33, 34, 35, 36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 47, 48, 49, 50, 51, 52, 53, 54, 55, 56, 57, 58, 59, 60, 61, 62, 63, 64, 65, 66, 67, 68, 69, 70, 71, 72, 73, 74, 75, 76, 77, 78, 79, 80, 81, 82, 83, 84, 85, 86, 87, 88, 89, 90, 91, 92, 93, 94}, []int64{41, 139, 203, 375, 506, 732, 933, 1104, 1315, 1456, 1511, 1570, 1625, 1763, 2162, 2478, 2912, 2795, 2278, 1742, 1337, 1378, 1235, 1139, 1116, 1036, 938, 863, 773, 746, 623, 533, 475, 388, 328, 279, 319, 341, 319, 333, 355, 314, 312, 347, 302, 304, 265, 263, 242, 237, 238, 182, 186, 181, 170, 135, 119, 115, 89, 98, 52, 52, 45, 21, 5})
	contentsalinityydists[8], _ = statistics.Init([]float64{36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 47, 48, 49, 50, 51, 52, 53, 54, 55, 56, 57, 58, 59, 60, 61, 62, 63, 64, 65, 66, 67, 68, 69, 70, 71, 72, 73, 74, 75, 76, 77, 78, 79, 80, 81, 82, 83, 84, 85, 86, 87, 88, 89, 90, 91, 92, 93, 94, 95, 96, 97, 98}, []int64{4, 4, 3, 6, 15, 18, 9, 26, 41, 59, 72, 71, 100, 110, 132, 145, 125, 137, 155, 125, 163, 150, 183, 192, 198, 63, 89, 68, 80, 72, 86, 95, 83, 79, 74, 69, 48, 78, 58, 51, 59, 61, 50, 69, 48, 40, 61, 37, 44, 31, 35, 20, 15, 30, 11, 11, 16, 10, 6, 2, 5, 2, 1})
	contentsalinityydists[9], _ = statistics.Init([]float64{42, 43, 44, 45, 46, 47, 48, 49, 50, 51, 52, 53, 54, 55, 56, 57, 58, 59, 60, 61, 62, 63, 64, 65, 66, 67, 68, 69, 70, 71, 72, 73, 74, 75, 76, 77, 78, 79, 80, 81, 82, 83, 84, 85, 86, 87, 88, 89, 90, 91, 92, 93, 94, 95, 96, 97, 98, 99}, []int64{5, 3, 3, 11, 14, 14, 14, 18, 37, 38, 58, 62, 79, 74, 108, 86, 106, 121, 126, 136, 126, 143, 149, 131, 124, 129, 157, 171, 172, 176, 128, 92, 88, 86, 100, 77, 78, 88, 73, 56, 61, 79, 54, 54, 51, 41, 35, 33, 25, 15, 36, 12, 18, 13, 6, 6, 2, 2})
	var contentsalinityStochastic = UncertaintyPairedData{Xvals: contentsalinityxs, Yvals: contentsalinityydists}

	contentsalinityStochastic.CentralTendency()

}

// func Test_UncertaintyJson(t *testing.T) {
// 	pd := createSampleData()
// 	b, err := json.Marshal(pd)
// 	if err != nil {
// 		panic(err)
// 	}
// 	fmt.Println(string(b))
// 	pd2 := json.Unmarshal(b, &UncertaintyPairedData{})
// 	fmt.Println(pd2)
// }
func createSampleData() UncertaintyPairedData {
	structurexs := []float64{-2.0, -1.0, 0.0, 1.0, 2.0, 3.0, 4.0, 5.0, 6.0, 7.0, 8.0, 9.0, 10.0, 11.0, 12.0, 13.0, 14.0, 15.0, 16.0}
	structureydists := make([]statistics.ContinuousDistribution, 19)
	structureydists[0] = statistics.NormalDistribution{Mean: 0, StandardDeviation: 0}
	structureydists[1] = statistics.NormalDistribution{Mean: 2.5, StandardDeviation: 0.30000001192092896}
	structureydists[2] = statistics.NormalDistribution{Mean: 13.399999618530273, StandardDeviation: 1.2000000476837158}
	structureydists[3] = statistics.NormalDistribution{Mean: 23.299999237060547, StandardDeviation: 1.6000000238418579}
	structureydists[4] = statistics.NormalDistribution{Mean: 32.099998474121094, StandardDeviation: 1.6000000238418579}
	structureydists[5] = statistics.NormalDistribution{Mean: 40.099998474121094, StandardDeviation: 1.7999999523162842}
	structureydists[6] = statistics.NormalDistribution{Mean: 47.099998474121094, StandardDeviation: 1.8999999761581421}
	structureydists[7] = statistics.NormalDistribution{Mean: 53.200000762939453, StandardDeviation: 2}
	structureydists[8] = statistics.NormalDistribution{Mean: 58.599998474121094, StandardDeviation: 2.0999999046325684}
	structureydists[9] = statistics.NormalDistribution{Mean: 63.200000762939453, StandardDeviation: 2.2000000476837158}
	structureydists[10] = statistics.NormalDistribution{Mean: 67.199996948242188, StandardDeviation: 2.2999999523162842}
	structureydists[11] = statistics.NormalDistribution{Mean: 70.5, StandardDeviation: 2.2999999523162842}
	structureydists[12] = statistics.NormalDistribution{Mean: 73.199996948242188, StandardDeviation: 2.3499999046325684}
	structureydists[13] = statistics.NormalDistribution{Mean: 75.4000015258789, StandardDeviation: 2.3900001049041748}
	structureydists[14] = statistics.NormalDistribution{Mean: 77.199996948242188, StandardDeviation: 2.4000000953674316}
	structureydists[15] = statistics.NormalDistribution{Mean: 78.5, StandardDeviation: 2.4100000858306885}
	structureydists[16] = statistics.NormalDistribution{Mean: 79.5, StandardDeviation: 2.4200000762939453}
	structureydists[17] = statistics.NormalDistribution{Mean: 80.199996948242188, StandardDeviation: 2.4300000667572021}
	structureydists[18] = statistics.NormalDistribution{Mean: 80.699996948242188, StandardDeviation: 2.4300000667572021}

	return UncertaintyPairedData{Xvals: structurexs, Yvals: structureydists}
}
