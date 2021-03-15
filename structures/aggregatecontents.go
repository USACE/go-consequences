package structures

import (
	"math/rand"

	"github.com/HenryGeorgist/go-statistics/data"
	"github.com/HenryGeorgist/go-statistics/statistics"
	"github.com/USACE/go-consequences/hazards"
	"github.com/USACE/go-consequences/paireddata"
	"github.com/USACE/go-consequences/structures"
)

func aggregateTriangular(min []float64, mostLikely []float64, max []float64) statistics.TriangularDistribution {

	const seed = 54321
	src := rand.NewSource(seed)
	rnd := rand.New(src)

	const N = 1000
	// aggregate non-zero distributions
	if max[1] != 0 {
		triDist1 := statistics.TriangularDistribution{Min: min[0], MostLikely: mostLikely[0], Max: max[0]}
		triDist2 := statistics.TriangularDistribution{Min: min[1], MostLikely: mostLikely[1], Max: max[1]}
		// initialize histogram in which to store aggregated distributions
		histogram := data.Init(1, min[0], max[1])
		// randomly sample each distribution, store in histogram
		for k := 0; k < N; k++ {
			probability := rnd.Float64()
			val1 := triDist1.InvCDF(probability)
			val2 := triDist2.InvCDF(probability)
			histogram.AddObservation(val1)
			histogram.AddObservation(val2)
		}
		// pull summary stats of aggregated sample
		return statistics.TriangularDistribution{Min: histogram.InvCDF(0), MostLikely: histogram.InvCDF(.5), Max: histogram.InvCDF(1)}
	} else {
		// zero-valued distributions
		return statistics.TriangularDistribution{Min: 0, MostLikely: 0, Max: 0}
	}

}

func comEng() structures.OccupancyTypeStochastic {

	// 2D array of most likely content damage engineeered structures
	contentDamageFunctionArrayMostLikely := [][]float64{
		{0, 0, 5, 18, 35, 39, 43, 47, 70, 75}, //non-perishable
		{0, 0, 5, 18, 35, 39, 43, 47, 70, 75}} //perishable
	mostLikely := transpose(contentDamageFunctionArrayMostLikely)
	// 2D array of minimum content damage engineered structures
	contentDamageFunctionArrayMin := [][]float64{
		{0, 0, 0, 4, 10, 22, 27, 33, 44, 48}, //non-perishable
		{0, 0, 0, 5, 17, 28, 37, 43, 50, 50}} //perishable
	min := transpose(contentDamageFunctionArrayMin)
	// 2D array of maximum content damage engineered structures
	contentDamageFunctionArrayMax := [][]float64{
		{0, 0, 5, 15, 22, 35, 44, 50, 55, 70}, //non-perishable
		{0, 0, 8, 28, 50, 58, 65, 65, 90, 90}} //perishable
	max := transpose(contentDamageFunctionArrayMax)

	structurexs := []float64{-1, -0.5, 0, 0.5, 1, 2, 3, 5, 7, 10}
	structureydists := make([]statistics.ContinuousDistribution, 10)
	structureydists[0] = statistics.TriangularDistribution{Min: 0, MostLikely: 0, Max: 0}
	structureydists[1] = statistics.TriangularDistribution{Min: 0, MostLikely: 0, Max: 0}
	structureydists[2] = statistics.TriangularDistribution{Min: 0, MostLikely: 5, Max: 9}
	structureydists[3] = statistics.TriangularDistribution{Min: 5, MostLikely: 10, Max: 17}
	structureydists[4] = statistics.TriangularDistribution{Min: 12, MostLikely: 20, Max: 27}
	structureydists[5] = statistics.TriangularDistribution{Min: 18, MostLikely: 30, Max: 36}
	structureydists[6] = statistics.TriangularDistribution{Min: 28, MostLikely: 35, Max: 43}
	structureydists[7] = statistics.TriangularDistribution{Min: 33, MostLikely: 40, Max: 48}
	structureydists[8] = statistics.TriangularDistribution{Min: 43, MostLikely: 53, Max: 60}
	structureydists[9] = statistics.TriangularDistribution{Min: 48, MostLikely: 58, Max: 69}
	contentxs := []float64{-1, -0.5, 0, 0.5, 1, 2, 3, 5, 7, 10}
	contentydists := make([]statistics.ContinuousDistribution, 10)
	for i := 0; i < len(contentxs); i++ {
		contentydists[i] = aggregateTriangular(min[i], mostLikely[i], max[i])
	}
	var structuredamagefunctionStochastic = paireddata.UncertaintyPairedData{Xvals: structurexs, Yvals: structureydists}
	var contentdamagefunctionStochastic = paireddata.UncertaintyPairedData{Xvals: contentxs, Yvals: contentydists}

	sm := make(map[hazards.Parameter]interface{})
	var sdf = structures.DamageFunctionFamilyStochastic{DamageFunctions: sm}

	cm := make(map[hazards.Parameter]interface{})
	var cdf = structures.DamageFunctionFamilyStochastic{DamageFunctions: cm}
	//Default hazard.
	sdf.DamageFunctions[hazards.Default] = structuredamagefunctionStochastic
	cdf.DamageFunctions[hazards.Default] = contentdamagefunctionStochastic
	//Depth Hazard
	sdf.DamageFunctions[hazards.Depth] = structuredamagefunctionStochastic
	cdf.DamageFunctions[hazards.Depth] = contentdamagefunctionStochastic

	return structures.OccupancyTypeStochastic{Name: "comEng", StructureDFF: sdf, ContentDFF: cdf}
}

func comNonEng() structures.OccupancyTypeStochastic {

	// 2D array of most likely content damage engineeered structures
	contentDamageFunctionArrayMostLikely := [][]float64{
		{0, 0, 1, 8, 12, 18, 25, 39, 50, 60},  //non-perishable
		{0, 0, 2, 15, 30, 42, 64, 71, 80, 87}} //perishable
	mostLikely := transpose(contentDamageFunctionArrayMostLikely)
	// 2D array of minimum content damage engineered structures
	contentDamageFunctionArrayMin := [][]float64{
		{0, 0, 0, 3, 7, 13, 20, 30, 40, 45}, //non-perishable
		{0, 0, 0, 5, 9, 15, 23, 30, 35, 41}} //perishable
	min := transpose(contentDamageFunctionArrayMin)
	// 2D array of maximum content damage engineered structures
	contentDamageFunctionArrayMax := [][]float64{
		{0, 0, 4, 18, 28, 38, 49, 64, 72, 90},   //non-perishable
		{0, 0, 10, 35, 54, 65, 84, 95, 99, 100}} //perishable
	max := transpose(contentDamageFunctionArrayMax)

	structurexs := []float64{-1, -0.5, 0, 0.5, 1, 2, 3, 5, 7, 10}
	structureydists := make([]statistics.ContinuousDistribution, 10)
	structureydists[0] = statistics.TriangularDistribution{Min: 0, MostLikely: 0, Max: 0}
	structureydists[1] = statistics.TriangularDistribution{Min: 0, MostLikely: 0, Max: 10}
	structureydists[2] = statistics.TriangularDistribution{Min: 0, MostLikely: 5, Max: 15}
	structureydists[3] = statistics.TriangularDistribution{Min: 5, MostLikely: 12, Max: 20}
	structureydists[4] = statistics.TriangularDistribution{Min: 10, MostLikely: 20, Max: 30}
	structureydists[5] = statistics.TriangularDistribution{Min: 15, MostLikely: 28, Max: 42}
	structureydists[6] = statistics.TriangularDistribution{Min: 20, MostLikely: 35, Max: 55}
	structureydists[7] = statistics.TriangularDistribution{Min: 28, MostLikely: 45, Max: 65}
	structureydists[8] = statistics.TriangularDistribution{Min: 35, MostLikely: 55, Max: 75}
	structureydists[9] = statistics.TriangularDistribution{Min: 40, MostLikely: 60, Max: 78}
	contentxs := []float64{-1, -0.5, 0, 0.5, 1, 2, 3, 5, 7, 10}
	contentydists := make([]statistics.ContinuousDistribution, 10)
	for i := 0; i < len(contentxs); i++ {
		contentydists[i] = aggregateTriangular(min[i], mostLikely[i], max[i])
	}
	var structuredamagefunctionStochastic = paireddata.UncertaintyPairedData{Xvals: structurexs, Yvals: structureydists}
	var contentdamagefunctionStochastic = paireddata.UncertaintyPairedData{Xvals: contentxs, Yvals: contentydists}

	sm := make(map[hazards.Parameter]interface{})
	var sdf = structures.DamageFunctionFamilyStochastic{DamageFunctions: sm}

	cm := make(map[hazards.Parameter]interface{})
	var cdf = structures.DamageFunctionFamilyStochastic{DamageFunctions: cm}
	//Default hazard.
	sdf.DamageFunctions[hazards.Default] = structuredamagefunctionStochastic
	cdf.DamageFunctions[hazards.Default] = contentdamagefunctionStochastic
	//Depth Hazard
	sdf.DamageFunctions[hazards.Depth] = structuredamagefunctionStochastic
	cdf.DamageFunctions[hazards.Depth] = contentdamagefunctionStochastic

	return structures.OccupancyTypeStochastic{Name: "comNonEng", StructureDFF: sdf, ContentDFF: cdf}

}

func transpose(slice [][]float64) [][]float64 {
	xl := len(slice[0])
	yl := len(slice)
	result := make([][]float64, xl)
	for i := range result {
		result[i] = make([]float64, yl)
	}
	for i := 0; i < xl; i++ {
		for j := 0; j < yl; j++ {
			result[i][j] = slice[j][i]
		}
	}
	return result
}
