package consequences

import (
	"math/rand"

	"github.com/USACE/go-consequences/paireddata"
)

type OccupancyType struct { //this is mutable
	Name            string
	Structuredamfun interface{} //if i make this an empty interface, it could be a value sampler, or an uncertainty valuesampler sampler...
	Contentdamfun   interface{} //if i make this an empty interface, it could be a value sampler, or an uncertainty valuesampler sampler...
}
type OccupancyTypeM struct { //need to swap - this is immutable not mutable
	Name            string
	Structuredamfun paireddata.ValueSampler
	Contentdamfun   paireddata.ValueSampler
}

func (o OccupancyType) SampleOccupancyType(seed int64) OccupancyTypeM {
	sd, oks := o.Structuredamfun.(paireddata.ValueSampler)
	cd, okc := o.Contentdamfun.(paireddata.ValueSampler)
	if oks && okc {
		return OccupancyTypeM{Name: o.Name, Structuredamfun: sd, Contentdamfun: cd}
	}
	rand.Seed(seed)
	if oks {
		cd2, okc1 := o.Contentdamfun.(paireddata.UncertaintyValueSamplerSampler)
		if okc1 {
			cd = cd2.SampleValueSampler(rand.Float64())
		} else {
			cd = nil
		}
	} else {
		sd2, oks1 := o.Structuredamfun.(paireddata.UncertaintyValueSamplerSampler)
		if oks1 {
			sd = sd2.SampleValueSampler(rand.Float64())
		} else {
			sd = nil
		}
	}
	if okc {
		sd3, oks2 := o.Structuredamfun.(paireddata.UncertaintyValueSamplerSampler)
		if oks2 {
			sd = sd3.SampleValueSampler(rand.Float64())
		} else {
			sd = nil
		}
	} else {
		cd3, okc2 := o.Contentdamfun.(paireddata.UncertaintyValueSamplerSampler)
		if okc2 {
			cd = cd3.SampleValueSampler(rand.Float64())
		} else {
			cd = nil
		}
	}

	return OccupancyTypeM{Name: o.Name, Structuredamfun: sd, Contentdamfun: cd}
}
func OccupancyTypeMap() map[string]OccupancyType {
	m := make(map[string]OccupancyType)
	m["AGR1"] = agr1()
	m["COM1"] = com1()
	m["COM2"] = com2()
	m["COM3"] = com3()
	m["COM4"] = com4()
	m["COM5"] = com5()
	m["COM6"] = com6()
	m["COM7"] = com7()
	m["COM8"] = com8()
	m["COM9"] = com9()
	m["COM10"] = com10()
	m["EDU1"] = edu1()
	m["EDU2"] = edu2()
	m["GOV1"] = gov1()
	m["GOV2"] = gov2()
	m["IND1"] = ind1()
	m["IND2"] = ind2()
	m["IND3"] = ind3()
	m["IND4"] = ind4()
	m["IND5"] = ind5()
	m["IND6"] = ind6()
	m["REL1"] = rel1()
	m["RES1-1SNB"] = res11snb()
	m["RES1-1SWB"] = res11swb()
	m["RES1-2SNB"] = res12snb()
	m["RES1-2SWB"] = res12swb()
	m["RES1-3SNB"] = res13snb()
	m["RES1-3SWB"] = res13swb()
	m["RES1-SLNB"] = res1slnb()
	m["RES1-SLWB"] = res1slwb()
	m["RES2"] = res2()
	m["RES3A"] = res3a()
	m["RES3B"] = res3b()
	m["RES3C"] = res3c()
	m["RES3D"] = res3d()
	m["RES3E"] = res3e()
	m["RES3F"] = res3f()
	m["RES4"] = res4()
	m["RES5"] = res5()
	m["RES6"] = res6()

	return m
}
func res11snb() OccupancyType {
	structurexs := []float64{-2.0, -1.0, 0.0, 1.0, 2.0, 3.0, 4.0, 5.0, 6.0, 7.0, 8.0, 9.0, 10.0, 11.0, 12.0, 13.0, 14.0, 15.0, 16.0}
	structureys := []float64{0.0, 2.5, 13.399999618530273, 23.299999237060547, 32.099998474121094, 40.099998474121094, 47.099998474121094, 53.200000762939453, 58.599998474121094, 63.200000762939453, 67.199996948242188, 70.5, 73.199996948242188, 75.4000015258789, 77.199996948242188, 78.5, 79.5, 80.199996948242188, 80.699996948242188}
	contentxs := []float64{-2.0, -1.0, 0.0, 1.0, 2.0, 3.0, 4.0, 5.0, 6.0, 7.0, 8.0, 9.0, 10.0, 11.0, 12.0, 13.0}
	contentys := []float64{0.0, 2.4000000953674316, 8.1000003814697266, 13.300000190734863, 17.899999618530273, 22.0, 25.700000762939453, 28.799999237060547, 31.5, 33.799999237060547, 35.700000762939453, 37.200000762939453, 38.400001525878906, 39.200000762939453, 39.700000762939453, 40.0}
	var structuredamagefunction = paireddata.PairedData{Xvals: structurexs, Yvals: structureys}
	var contentdamagefunction = paireddata.PairedData{Xvals: contentxs, Yvals: contentys}
	return OccupancyType{Name: "RES1-1SNB", Structuredamfun: structuredamagefunction, Contentdamfun: contentdamagefunction}
}
func res11swb() OccupancyType {
	structurexs := []float64{-8, -7, -6, -5, -4, -3, -2, -1, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16}
	structureys := []float64{0, 0.699999988079071, 0.800000011920929, 2.4000000953674316, 5.1999998092651367, 9, 13.800000190734863, 19.399999618530273, 25.5, 32, 38.700000762939453, 45.5, 52.200000762939453, 58.599998474121094, 64.5, 69.800003051757812, 74.199996948242188, 77.699996948242188, 80.0999984741211, 81.0999984741211, 81.0999984741211, 81.0999984741211, 81.0999984741211, 81.0999984741211, 81.0999984741211}
	contentxs := []float64{-8, -7, -6, -5, -4, -3, -2, -1, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16}
	contentys := []float64{0.10000000149011612, 0.800000011920929, 2.0999999046325684, 3.7000000476837158, 5.6999998092651367, 8, 10.5, 13.199999809265137, 16, 18.899999618530273, 21.799999237060547, 24.700000762939453, 27.399999618530273, 30, 32.400001525878906, 34.5, 36.299999237060547, 37.700000762939453, 38.599998474121094, 39.099998474121094, 39.099998474121094, 39.099998474121094, 39.099998474121094, 39.099998474121094, 39.099998474121094}
	var structuredamagefunction = paireddata.PairedData{Xvals: structurexs, Yvals: structureys}
	var contentdamagefunction = paireddata.PairedData{Xvals: contentxs, Yvals: contentys}
	return OccupancyType{Name: "RES1-1SWB", Structuredamfun: structuredamagefunction, Contentdamfun: contentdamagefunction}
}
func agr1() OccupancyType {
	structurexs := []float64{-4, -3, -2, -1, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24}
	structureys := []float64{0, 0, 0, 0, 0, 6, 11, 15, 19, 25, 30, 35, 41, 46, 51, 57, 63, 70, 75, 79, 82, 84, 87, 89, 90, 92, 93, 95, 96}
	contentxs := []float64{-4, -3, -2, -1, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24}
	contentys := []float64{0, 0, 0, 0, 6, 20, 43, 58, 65, 66, 66, 67, 70, 75, 76, 76, 76, 77, 77, 77, 78, 78, 78, 79, 79, 79, 79, 80, 80}
	var structuredamagefunction = paireddata.PairedData{Xvals: structurexs, Yvals: structureys}
	var contentdamagefunction = paireddata.PairedData{Xvals: contentxs, Yvals: contentys}
	return OccupancyType{Name: "AGR1", Structuredamfun: structuredamagefunction, Contentdamfun: contentdamagefunction}
}
func com1() OccupancyType {
	structurexs := []float64{-4, -3, -2, -1, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24}
	structureys := []float64{0, 0, 0, 0, 1, 9, 14, 16, 18, 20, 23, 26, 30, 34, 38, 42, 47, 51, 55, 58, 61, 64, 67, 69, 71, 74, 76, 78, 80}
	contentxs := []float64{-4, -3, -2, -1, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24}
	contentys := []float64{0, 0, 0, 0, 2, 26, 42, 56, 68, 78, 83, 85, 87, 88, 89, 90, 91, 92, 92, 92, 93, 93, 94, 94, 94, 94, 94, 94, 94}
	var structuredamagefunction = paireddata.PairedData{Xvals: structurexs, Yvals: structureys}
	var contentdamagefunction = paireddata.PairedData{Xvals: contentxs, Yvals: contentys}
	return OccupancyType{Name: "COM1", Structuredamfun: structuredamagefunction, Contentdamfun: contentdamagefunction}
}
func com2() OccupancyType {
	structurexs := []float64{-4, -3, -2, -1, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24}
	structureys := []float64{0, 0, 0, 0, 0, 5, 8, 11, 13, 16, 19, 22, 25, 29, 32, 37, 41, 45, 49, 52, 55, 58, 61, 63, 66, 68, 70, 71, 73}
	contentxs := []float64{-4, -3, -2, -1, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24}
	contentys := []float64{0, 0, 0, 0, 3, 16, 27, 36, 49, 57, 63, 69, 72, 76, 80, 82, 84, 86, 87, 87, 88, 89, 89, 89, 89, 89, 89, 89, 89}
	var structuredamagefunction = paireddata.PairedData{Xvals: structurexs, Yvals: structureys}
	var contentdamagefunction = paireddata.PairedData{Xvals: contentxs, Yvals: contentys}
	return OccupancyType{Name: "COM2", Structuredamfun: structuredamagefunction, Contentdamfun: contentdamagefunction}
}
func com3() OccupancyType {
	structurexs := []float64{-4, -3, -2, -1, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24}
	structureys := []float64{0, 0, 0, 0, 0, 9, 12, 13, 16, 19, 22, 25, 28, 32, 35, 39, 43, 47, 50, 54, 57, 61, 64, 68, 71, 75, 78, 82, 85}
	contentxs := []float64{-4, -3, -2, -1, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24}
	contentys := []float64{0, 0, 0, 0, 4, 29, 46, 67, 79, 85, 91, 92, 92, 93, 94, 96, 96, 97, 97, 98, 98, 98, 98, 98, 98, 98, 98, 98, 98}
	var structuredamagefunction = paireddata.PairedData{Xvals: structurexs, Yvals: structureys}
	var contentdamagefunction = paireddata.PairedData{Xvals: contentxs, Yvals: contentys}
	return OccupancyType{Name: "COM3", Structuredamfun: structuredamagefunction, Contentdamfun: contentdamagefunction}
}
func com4() OccupancyType {
	structurexs := []float64{-4, -3, -2, -1, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24}
	structureys := []float64{0, 0, 0, 0, 2, 11, 16, 22, 28, 35, 38, 41, 44, 47, 50, 54, 57, 59, 62, 66, 68, 70, 72, 74, 76, 77, 78, 79, 80}
	contentxs := []float64{-4, -3, -2, -1, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24}
	contentys := []float64{0, 0, 0, 0, 2, 18, 25, 35, 43, 49, 52, 55, 57, 58, 60, 65, 67, 68, 69, 70, 71, 71, 72, 72, 72, 72, 72, 72, 72}
	var structuredamagefunction = paireddata.PairedData{Xvals: structurexs, Yvals: structureys}
	var contentdamagefunction = paireddata.PairedData{Xvals: contentxs, Yvals: contentys}
	return OccupancyType{Name: "COM4", Structuredamfun: structuredamagefunction, Contentdamfun: contentdamagefunction}
}
func com5() OccupancyType {
	structurexs := []float64{-4, -3, -2, -1, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24}
	structureys := []float64{0, 0, 0, 0, 0, 11, 11, 12, 13, 15, 17, 19, 22, 24, 28, 31, 34, 37, 40, 44, 48, 51, 55, 59, 63, 67, 71, 75, 79}
	contentxs := []float64{-4, -3, -2, -1, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24}
	contentys := []float64{0, 0, 0, 0, 0, 50, 74, 83, 100, 100, 100, 100, 100, 100, 100, 100, 100, 100, 100, 100, 100, 100, 100, 100, 100, 100, 100, 100, 100}
	var structuredamagefunction = paireddata.PairedData{Xvals: structurexs, Yvals: structureys}
	var contentdamagefunction = paireddata.PairedData{Xvals: contentxs, Yvals: contentys}
	return OccupancyType{Name: "COM5", Structuredamfun: structuredamagefunction, Contentdamfun: contentdamagefunction}
}
func com6() OccupancyType {
	structurexs := []float64{-4, -3, -2, -1, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24}
	structureys := []float64{0, 0, 0, 0, 0, 0, 0, 20, 25, 30, 35, 40, 43, 47, 50, 53, 55, 57, 60, 60, 60, 60, 60, 60, 60, 60, 60, 60, 60}
	contentxs := []float64{-4, -3, -2, -1, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24}
	contentys := []float64{0, 0, 0, 0, 0, 0, 0, 10, 20, 30, 65, 72, 78, 85, 95, 95, 95, 95, 96, 96, 96, 96, 96, 96, 96, 96, 96, 96, 96}
	var structuredamagefunction = paireddata.PairedData{Xvals: structurexs, Yvals: structureys}
	var contentdamagefunction = paireddata.PairedData{Xvals: contentxs, Yvals: contentys}
	return OccupancyType{Name: "COM6", Structuredamfun: structuredamagefunction, Contentdamfun: contentdamagefunction}
}
func com7() OccupancyType {
	structurexs := []float64{-4, -3, -2, -1, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24}
	structureys := []float64{0, 0, 0, 0, 2, 11, 12, 13, 14, 16, 17, 18, 20, 22, 24, 27, 30, 34, 37, 41, 44, 48, 51, 54, 56, 59, 61, 64, 66}
	contentxs := []float64{-4, -3, -2, -1, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24}
	contentys := []float64{0, 0, 0, 0, 0, 28, 51, 60, 63, 67, 71, 72, 74, 77, 81, 86, 92, 94, 97, 99, 100, 100, 100, 100, 100, 100, 100, 100, 100}
	var structuredamagefunction = paireddata.PairedData{Xvals: structurexs, Yvals: structureys}
	var contentdamagefunction = paireddata.PairedData{Xvals: contentxs, Yvals: contentys}
	return OccupancyType{Name: "COM7", Structuredamfun: structuredamagefunction, Contentdamfun: contentdamagefunction}
}
func com8() OccupancyType {
	structurexs := []float64{-4, -3, -2, -1, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24}
	structureys := []float64{0, 0, 0, 0, 1, 9, 11, 12, 14, 16, 18, 20, 22, 26, 29, 33, 37, 41, 45, 50, 53, 57, 60, 63, 66, 69, 73, 76, 78}
	contentxs := []float64{-4, -3, -2, -1, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24}
	contentys := []float64{0, 0, 0, 0, 13, 45, 55, 64, 73, 77, 80, 82, 83, 85, 87, 89, 90, 91, 92, 93, 94, 95, 96, 96, 96, 96, 96, 96, 96}
	var structuredamagefunction = paireddata.PairedData{Xvals: structurexs, Yvals: structureys}
	var contentdamagefunction = paireddata.PairedData{Xvals: contentxs, Yvals: contentys}
	return OccupancyType{Name: "COM8", Structuredamfun: structuredamagefunction, Contentdamfun: contentdamagefunction}
}
func com9() OccupancyType {
	structurexs := []float64{-4, -3, -2, -1, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24}
	structureys := []float64{0, 0, 0, 0, 0, 2, 4, 5, 5, 5, 6, 8, 10, 12, 15, 20, 24, 29, 35, 42, 49, 56, 62, 68, 74, 80, 86, 92, 98}
	contentxs := []float64{-4, -3, -2, -1, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24}
	contentys := []float64{0, 0, 0, 0, 0, 4, 6, 8, 9, 10, 12, 17, 22, 30, 41, 57, 66, 73, 79, 84, 90, 97, 98, 100, 100, 100, 100, 100, 100}
	var structuredamagefunction = paireddata.PairedData{Xvals: structurexs, Yvals: structureys}
	var contentdamagefunction = paireddata.PairedData{Xvals: contentxs, Yvals: contentys}
	return OccupancyType{Name: "COM9", Structuredamfun: structuredamagefunction, Contentdamfun: contentdamagefunction}
}
func com10() OccupancyType {
	structurexs := []float64{-4, -3, -2, -1, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24}
	structureys := []float64{0, 0, 0, 0, 0, 3, 5, 6, 7, 8, 10, 13, 17, 21, 25, 30, 35, 41, 47, 52, 58, 65, 71, 76, 81, 86, 91, 96, 100}
	contentxs := []float64{-4, -3, -2, -1, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24}
	contentys := []float64{0, 0, 0, 0, 0, 11, 17, 20, 23, 25, 29, 35, 42, 51, 63, 77, 93, 100, 100, 100, 100, 100, 100, 100, 100, 100, 100, 100, 100}
	var structuredamagefunction = paireddata.PairedData{Xvals: structurexs, Yvals: structureys}
	var contentdamagefunction = paireddata.PairedData{Xvals: contentxs, Yvals: contentys}
	return OccupancyType{Name: "COM10", Structuredamfun: structuredamagefunction, Contentdamfun: contentdamagefunction}
}
func edu1() OccupancyType {
	structurexs := []float64{-4, -3, -2, -1, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24}
	structureys := []float64{0, 0, 0, 0, 0, 5, 7, 9, 9, 10, 11, 13, 15, 17, 20, 24, 28, 33, 39, 45, 52, 59, 64, 69, 74, 79, 84, 89, 94}
	contentxs := []float64{-4, -3, -2, -1, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24}
	contentys := []float64{0, 0, 0, 0, 0, 27, 38, 53, 64, 68, 70, 72, 75, 79, 83, 88, 94, 100, 100, 100, 100, 100, 100, 100, 100, 100, 100, 100, 100}
	var structuredamagefunction = paireddata.PairedData{Xvals: structurexs, Yvals: structureys}
	var contentdamagefunction = paireddata.PairedData{Xvals: contentxs, Yvals: contentys}
	return OccupancyType{Name: "EDU1", Structuredamfun: structuredamagefunction, Contentdamfun: contentdamagefunction}
}
func edu2() OccupancyType {
	structurexs := []float64{-4, -3, -2, -1, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24}
	structureys := []float64{0, 0, 0, 0, 0, 5, 7, 9, 9, 10, 11, 13, 15, 17, 20, 24, 28, 33, 39, 45, 52, 59, 64, 69, 74, 79, 84, 89, 94}
	contentxs := []float64{-4, -3, -2, -1, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24}
	contentys := []float64{0, 0, 0, 0, 0, 27, 38, 53, 64, 68, 70, 72, 75, 79, 83, 88, 94, 100, 100, 100, 100, 100, 100, 100, 100, 100, 100, 100, 100}
	var structuredamagefunction = paireddata.PairedData{Xvals: structurexs, Yvals: structureys}
	var contentdamagefunction = paireddata.PairedData{Xvals: contentxs, Yvals: contentys}
	return OccupancyType{Name: "EDU2", Structuredamfun: structuredamagefunction, Contentdamfun: contentdamagefunction}
}
func gov1() OccupancyType {
	structurexs := []float64{-4, -3, -2, -1, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24}
	structureys := []float64{0, 0, 0, 0, 0, 5, 8, 13, 14, 14, 15, 17, 19, 22, 26, 31, 37, 44, 51, 59, 65, 70, 74, 79, 83, 87, 91, 95, 98}
	contentxs := []float64{-4, -3, -2, -1, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24}
	contentys := []float64{0, 0, 0, 0, 0, 30, 59, 74, 83, 90, 100, 100, 100, 100, 100, 100, 100, 100, 100, 100, 100, 100, 100, 100, 100, 100, 100, 100, 100}
	var structuredamagefunction = paireddata.PairedData{Xvals: structurexs, Yvals: structureys}
	var contentdamagefunction = paireddata.PairedData{Xvals: contentxs, Yvals: contentys}
	return OccupancyType{Name: "GOV1", Structuredamfun: structuredamagefunction, Contentdamfun: contentdamagefunction}
}
func gov2() OccupancyType {
	structurexs := []float64{-4, -3, -2, -1, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24}
	structureys := []float64{0, 0, 0, 0, 0, 7, 10, 11, 12, 15, 17, 20, 23, 27, 31, 35, 40, 44, 48, 52, 56, 60, 64, 68, 72, 76, 80, 84, 88}
	contentxs := []float64{-4, -3, -2, -1, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24}
	contentys := []float64{0, 0, 0, 0, 0, 8, 20, 38, 55, 70, 81, 89, 98, 100, 100, 100, 100, 100, 100, 100, 100, 100, 100, 100, 100, 100, 100, 100, 100}
	var structuredamagefunction = paireddata.PairedData{Xvals: structurexs, Yvals: structureys}
	var contentdamagefunction = paireddata.PairedData{Xvals: contentxs, Yvals: contentys}
	return OccupancyType{Name: "GOV2", Structuredamfun: structuredamagefunction, Contentdamfun: contentdamagefunction}
}
func ind1() OccupancyType {
	structurexs := []float64{-4, -3, -2, -1, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24}
	structureys := []float64{0, 0, 0, 0, 1, 10, 12, 15, 19, 22, 26, 30, 35, 39, 42, 48, 50, 51, 53, 54, 55, 55, 56, 56, 57, 57, 57, 58, 58}
	contentxs := []float64{-4, -3, -2, -1, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24}
	contentys := []float64{0, 0, 0, 0, 0, 15, 24, 34, 41, 47, 52, 57, 60, 63, 64, 66, 68, 69, 72, 73, 73, 73, 74, 74, 74, 74, 75, 75, 75}
	var structuredamagefunction = paireddata.PairedData{Xvals: structurexs, Yvals: structureys}
	var contentdamagefunction = paireddata.PairedData{Xvals: contentxs, Yvals: contentys}
	return OccupancyType{Name: "IND1", Structuredamfun: structuredamagefunction, Contentdamfun: contentdamagefunction}
}
func ind2() OccupancyType {
	structurexs := []float64{-4, -3, -2, -1, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24}
	structureys := []float64{0, 0, 0, 0, 1, 9, 14, 17, 22, 26, 30, 32, 35, 37, 39, 43, 46, 48, 50, 51, 54, 55, 57, 59, 60, 62, 63, 65, 66}
	contentxs := []float64{-4, -3, -2, -1, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24}
	contentys := []float64{0, 0, 0, 0, 0, 9, 23, 35, 44, 52, 58, 62, 65, 68, 70, 73, 74, 77, 78, 78, 79, 80, 80, 80, 80, 81, 81, 81, 81}
	var structuredamagefunction = paireddata.PairedData{Xvals: structurexs, Yvals: structureys}
	var contentdamagefunction = paireddata.PairedData{Xvals: contentxs, Yvals: contentys}
	return OccupancyType{Name: "IND2", Structuredamfun: structuredamagefunction, Contentdamfun: contentdamagefunction}
}
func ind3() OccupancyType {
	structurexs := []float64{-4, -3, -2, -1, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24}
	structureys := []float64{0, 0, 0, 0, 0, 13, 14, 19, 22, 25, 28, 30, 33, 34, 36, 39, 40, 42, 42, 43, 43, 44, 44, 44, 44, 44, 45, 45, 45}
	contentxs := []float64{-4, -3, -2, -1, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24}
	contentys := []float64{0, 0, 0, 0, 2, 20, 41, 51, 62, 67, 71, 73, 76, 78, 79, 82, 83, 84, 86, 87, 87, 88, 88, 88, 88, 88, 88, 88, 88}
	var structuredamagefunction = paireddata.PairedData{Xvals: structurexs, Yvals: structureys}
	var contentdamagefunction = paireddata.PairedData{Xvals: contentxs, Yvals: contentys}
	return OccupancyType{Name: "IND3", Structuredamfun: structuredamagefunction, Contentdamfun: contentdamagefunction}
}
func ind4() OccupancyType {
	structurexs := []float64{-4, -3, -2, -1, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24}
	structureys := []float64{0, 0, 0, 0, 0, 10, 14, 18, 22, 26, 34, 41, 42, 42, 45, 47, 49, 50, 50, 50, 50, 50, 50, 50, 50, 50, 50, 50, 50}
	contentxs := []float64{-4, -3, -2, -1, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24}
	contentys := []float64{0, 0, 0, 0, 15, 20, 26, 31, 37, 40, 44, 48, 53, 56, 57, 60, 62, 63, 63, 63, 64, 65, 65, 65, 65, 65, 65, 65, 65}
	var structuredamagefunction = paireddata.PairedData{Xvals: structurexs, Yvals: structureys}
	var contentdamagefunction = paireddata.PairedData{Xvals: contentxs, Yvals: contentys}
	return OccupancyType{Name: "IND4", Structuredamfun: structuredamagefunction, Contentdamfun: contentdamagefunction}
}
func ind5() OccupancyType {
	structurexs := []float64{-4, -3, -2, -1, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24}
	structureys := []float64{0, 0, 0, 0, 0, 13, 14, 19, 22, 25, 28, 30, 33, 34, 36, 39, 40, 42, 42, 43, 43, 44, 44, 44, 44, 44, 45, 45, 45}
	contentxs := []float64{-4, -3, -2, -1, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24}
	contentys := []float64{0, 0, 0, 0, 2, 20, 41, 51, 62, 67, 71, 73, 76, 78, 79, 82, 83, 84, 86, 87, 87, 88, 88, 88, 88, 88, 88, 88, 88}
	var structuredamagefunction = paireddata.PairedData{Xvals: structurexs, Yvals: structureys}
	var contentdamagefunction = paireddata.PairedData{Xvals: contentxs, Yvals: contentys}
	return OccupancyType{Name: "IND5", Structuredamfun: structuredamagefunction, Contentdamfun: contentdamagefunction}
}
func ind6() OccupancyType {
	structurexs := []float64{-4, -3, -2, -1, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24}
	structureys := []float64{0, 0, 0, 0, 0, 22, 31, 37, 43, 47, 50, 54, 57, 61, 63, 64, 65, 67, 68, 69, 70, 71, 72, 73, 74, 75, 76, 76, 77}
	contentxs := []float64{-4, -3, -2, -1, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24}
	contentys := []float64{0, 0, 0, 0, 0, 20, 35, 47, 56, 59, 66, 69, 71, 72, 78, 79, 80, 80, 81, 81, 81, 82, 82, 82, 83, 83, 83, 83, 83}
	var structuredamagefunction = paireddata.PairedData{Xvals: structurexs, Yvals: structureys}
	var contentdamagefunction = paireddata.PairedData{Xvals: contentxs, Yvals: contentys}
	return OccupancyType{Name: "IND6", Structuredamfun: structuredamagefunction, Contentdamfun: contentdamagefunction}
}
func rel1() OccupancyType {
	structurexs := []float64{-4, -3, -2, -1, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24}
	structureys := []float64{0, 0, 0, 0, 0, 10, 11, 11, 12, 12, 13, 14, 14, 15, 17, 19, 24, 30, 38, 45, 52, 58, 64, 69, 74, 78, 82, 85, 88}
	contentxs := []float64{-4, -3, -2, -1, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24}
	contentys := []float64{0, 0, 0, 0, 10, 52, 72, 85, 92, 95, 98, 99, 100, 100, 100, 100, 100, 100, 100, 100, 100, 100, 100, 100, 100, 100, 100, 100, 100}
	var structuredamagefunction = paireddata.PairedData{Xvals: structurexs, Yvals: structureys}
	var contentdamagefunction = paireddata.PairedData{Xvals: contentxs, Yvals: contentys}
	return OccupancyType{Name: "REL1", Structuredamfun: structuredamagefunction, Contentdamfun: contentdamagefunction}
}
func res12snb() OccupancyType {
	structurexs := []float64{-2, -1, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16}
	structureys := []float64{0, 3, 9.3000001907348633, 15.199999809265137, 20.899999618530273, 26.299999237060547, 31.399999618530273, 36.200000762939453, 40.700000762939453, 44.900001525878906, 48.799999237060547, 52.400001525878906, 55.700000762939453, 58.700000762939453, 61.400001525878906, 63.799999237060547, 65.9000015258789, 67.699996948242188, 69.199996948242188}
	contentxs := []float64{-2, -1, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16}
	contentys := []float64{0, 1, 5, 8.6999998092651367, 12.199999809265137, 15.5, 18.5, 21.299999237060547, 23.899999618530273, 26.299999237060547, 28.399999618530273, 30.299999237060547, 32, 33.400001525878906, 34.700000762939453, 35.599998474121094, 36.400001525878906, 36.900001525878906, 37.200000762939453}
	var structuredamagefunction = paireddata.PairedData{Xvals: structurexs, Yvals: structureys}
	var contentdamagefunction = paireddata.PairedData{Xvals: contentxs, Yvals: contentys}
	return OccupancyType{Name: "RES1-2SNB", Structuredamfun: structuredamagefunction, Contentdamfun: contentdamagefunction}
}
func res12swb() OccupancyType {
	structurexs := []float64{-8, -7, -6, -5, -4, -3, -2, -1, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16}
	structureys := []float64{1.7000000476837158, 1.7000000476837158, 1.8999999761581421, 2.9000000953674316, 4.6999998092651367, 7.1999998092651367, 10.199999809265137, 13.899999618530273, 17.899999618530273, 22.299999237060547, 27, 31.899999618530273, 36.900001525878906, 41.900001525878906, 46.900001525878906, 51.799999237060547, 56.400001525878906, 60.799999237060547, 64.800003051757812, 68.4000015258789, 71.4000015258789, 73.699996948242188, 75.4000015258789, 76.4000015258789, 76.4000015258789}
	contentxs := []float64{-8, -7, -6, -5, -4, -3, -2, -1, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16}
	contentys := []float64{0, 1, 2.2999999523162842, 3.7000000476837158, 5.1999998092651367, 6.8000001907348633, 8.3999996185302734, 10.100000381469727, 11.899999618530273, 13.800000190734863, 15.699999809265137, 17.700000762939453, 19.799999237060547, 22, 24.299999237060547, 26.700000762939453, 29.100000381469727, 31.700000762939453, 34.400001525878906, 37.200000762939453, 40, 43, 46.099998474121094, 49.299999237060547, 52.599998474121094}
	var structuredamagefunction = paireddata.PairedData{Xvals: structurexs, Yvals: structureys}
	var contentdamagefunction = paireddata.PairedData{Xvals: contentxs, Yvals: contentys}
	return OccupancyType{Name: "RES1-2SWB", Structuredamfun: structuredamagefunction, Contentdamfun: contentdamagefunction}
}
func res13snb() OccupancyType {
	structurexs := []float64{-2, -1, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16}
	structureys := []float64{0, 3, 9.3000001907348633, 15.199999809265137, 20.899999618530273, 26.299999237060547, 31.399999618530273, 36.200000762939453, 40.700000762939453, 44.900001525878906, 48.799999237060547, 52.400001525878906, 55.700000762939453, 58.700000762939453, 61.400001525878906, 63.799999237060547, 65.9000015258789, 67.699996948242188, 69.199996948242188}
	contentxs := []float64{-2, -1, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16}
	contentys := []float64{0, 1, 5, 8.6999998092651367, 12.199999809265137, 15.5, 18.5, 21.299999237060547, 23.899999618530273, 26.299999237060547, 28.399999618530273, 30.299999237060547, 32, 33.400001525878906, 34.700000762939453, 35.599998474121094, 36.400001525878906, 36.900001525878906, 37.200000762939453}
	var structuredamagefunction = paireddata.PairedData{Xvals: structurexs, Yvals: structureys}
	var contentdamagefunction = paireddata.PairedData{Xvals: contentxs, Yvals: contentys}
	return OccupancyType{Name: "RES1-3SNB", Structuredamfun: structuredamagefunction, Contentdamfun: contentdamagefunction}
}
func res13swb() OccupancyType {
	structurexs := []float64{-8, -7, -6, -5, -4, -3, -2, -1, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16}
	structureys := []float64{1.7000000476837158, 1.7000000476837158, 1.8999999761581421, 2.9000000953674316, 4.6999998092651367, 7.1999998092651367, 10.199999809265137, 13.899999618530273, 17.899999618530273, 22.299999237060547, 27, 31.899999618530273, 36.900001525878906, 41.900001525878906, 46.900001525878906, 51.799999237060547, 56.400001525878906, 60.799999237060547, 64.800003051757812, 68.4000015258789, 71.4000015258789, 73.699996948242188, 75.4000015258789, 76.4000015258789, 76.4000015258789}
	contentxs := []float64{-8, -7, -6, -5, -4, -3, -2, -1, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16}
	contentys := []float64{0, 1, 2.2999999523162842, 3.7000000476837158, 5.1999998092651367, 6.8000001907348633, 8.3999996185302734, 10.100000381469727, 11.899999618530273, 13.800000190734863, 15.699999809265137, 17.700000762939453, 19.799999237060547, 22, 24.299999237060547, 26.700000762939453, 29.100000381469727, 31.700000762939453, 34.400001525878906, 37.200000762939453, 40, 43, 46.099998474121094, 49.299999237060547, 52.599998474121094}
	var structuredamagefunction = paireddata.PairedData{Xvals: structurexs, Yvals: structureys}
	var contentdamagefunction = paireddata.PairedData{Xvals: contentxs, Yvals: contentys}
	return OccupancyType{Name: "RES1-3SWB", Structuredamfun: structuredamagefunction, Contentdamfun: contentdamagefunction}
}
func res1slnb() OccupancyType {
	structurexs := []float64{-2, -1, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16}
	structureys := []float64{0, 6.4000000953674316, 7.1999998092651367, 9.3999996185302734, 12.899999618530273, 17.399999618530273, 22.799999237060547, 28.899999618530273, 35.5, 42.299999237060547, 49.200000762939453, 56.099998474121094, 62.599998474121094, 68.5999984741211, 73.9000015258789, 78.4000015258789, 81.699996948242188, 83.800003051757812, 84.4000015258789}
	contentxs := []float64{-2, -1, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16}
	contentys := []float64{0, 2.2000000476837158, 2.9000000953674316, 4.6999998092651367, 7.5, 11.100000381469727, 15.300000190734863, 20.100000381469727, 25.200000762939453, 30.5, 35.700000762939453, 40.900001525878906, 45.799999237060547, 50.200000762939453, 54.099998474121094, 57.200000762939453, 59.400001525878906, 60.5, 60.5}
	var structuredamagefunction = paireddata.PairedData{Xvals: structurexs, Yvals: structureys}
	var contentdamagefunction = paireddata.PairedData{Xvals: contentxs, Yvals: contentys}
	return OccupancyType{Name: "RES1-SLNB", Structuredamfun: structuredamagefunction, Contentdamfun: contentdamagefunction}
}
func res1slwb() OccupancyType {
	structurexs := []float64{-8, -7, -6, -5, -4, -3, -2, -1, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16}
	structureys := []float64{0, 0, 2.5, 3.0999999046325684, 4.6999998092651367, 7.1999998092651367, 10.399999618530273, 14.199999809265137, 18.5, 23.200000762939453, 28.200000762939453, 33.400001525878906, 38.599998474121094, 43.799999237060547, 48.799999237060547, 53.5, 57.799999237060547, 61.599998474121094, 64.800003051757812, 67.199996948242188, 68.800003051757812, 69.300003051757812, 69.300003051757812, 69.300003051757812, 69.300003051757812}
	contentxs := []float64{-8, -7, -6, -5, -4, -3, -2, -1, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16}
	contentys := []float64{0.60000002384185791, 0.699999988079071, 1.3999999761581421, 2.4000000953674316, 3.7999999523162842, 5.4000000953674316, 7.3000001907348633, 9.3999996185302734, 11.600000381469727, 13.800000190734863, 16.100000381469727, 18.200000762939453, 20.200000762939453, 22.100000381469727, 23.600000381469727, 24.899999618530273, 25.799999237060547, 26.299999237060547, 26.299999237060547, 26.299999237060547, 26.299999237060547, 26.299999237060547, 26.299999237060547, 26.299999237060547, 26.299999237060547}
	var structuredamagefunction = paireddata.PairedData{Xvals: structurexs, Yvals: structureys}
	var contentdamagefunction = paireddata.PairedData{Xvals: contentxs, Yvals: contentys}
	return OccupancyType{Name: "RES1-SLWB", Structuredamfun: structuredamagefunction, Contentdamfun: contentdamagefunction}
}

func res2() OccupancyType {
	structurexs := []float64{-4, -3, -2, -1, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24}
	structureys := []float64{0, 0, 0, 0, 11, 44, 63, 73, 78, 79, 81, 82, 83, 84, 85, 86, 88, 89, 90, 91, 92, 94, 95, 96, 97, 98, 99, 100, 100}
	contentxs := []float64{-4, -3, -2, -1, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24}
	contentys := []float64{0, 0, 0, 0, 3, 27, 49, 64, 70, 76, 78, 79, 81, 83, 83, 83, 83, 83, 83, 83, 83, 83, 83, 83, 83, 83, 83, 83, 83}
	var structuredamagefunction = paireddata.PairedData{Xvals: structurexs, Yvals: structureys}
	var contentdamagefunction = paireddata.PairedData{Xvals: contentxs, Yvals: contentys}
	return OccupancyType{Name: "RES2", Structuredamfun: structuredamagefunction, Contentdamfun: contentdamagefunction}
}
func res3a() OccupancyType {
	structurexs := []float64{-4, -3, -2, -1, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24}
	structureys := []float64{0, 0, 0, 0, 5, 28, 29, 31, 36, 37, 39, 40, 41, 42, 44, 46, 48, 52, 55, 58, 61, 64, 68, 69, 70, 71, 72, 73, 74}
	contentxs := []float64{-4, -3, -2, -1, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24}
	contentys := []float64{0, 0, 0, 0, 4, 24, 34, 40, 47, 53, 56, 58, 58, 58, 61, 66, 68, 76, 81, 86, 91, 100, 100, 100, 100, 100, 100, 100, 100}
	var structuredamagefunction = paireddata.PairedData{Xvals: structurexs, Yvals: structureys}
	var contentdamagefunction = paireddata.PairedData{Xvals: contentxs, Yvals: contentys}
	return OccupancyType{Name: "RES3A", Structuredamfun: structuredamagefunction, Contentdamfun: contentdamagefunction}
}
func res3b() OccupancyType {
	structurexs := []float64{-4, -3, -2, -1, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24}
	structureys := []float64{0, 0, 0, 0, 5, 28, 29, 31, 36, 37, 39, 40, 41, 42, 44, 46, 48, 52, 55, 58, 61, 64, 68, 69, 70, 71, 72, 73, 74}
	contentxs := []float64{-4, -3, -2, -1, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24}
	contentys := []float64{0, 0, 0, 0, 4, 24, 34, 40, 47, 53, 56, 58, 58, 58, 61, 66, 68, 76, 81, 86, 91, 100, 100, 100, 100, 100, 100, 100, 100}
	var structuredamagefunction = paireddata.PairedData{Xvals: structurexs, Yvals: structureys}
	var contentdamagefunction = paireddata.PairedData{Xvals: contentxs, Yvals: contentys}
	return OccupancyType{Name: "RES3B", Structuredamfun: structuredamagefunction, Contentdamfun: contentdamagefunction}
}
func res3c() OccupancyType {
	structurexs := []float64{-4, -3, -2, -1, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24}
	structureys := []float64{0, 0, 0, 0, 5, 28, 29, 31, 36, 37, 39, 40, 41, 42, 44, 46, 48, 52, 55, 58, 61, 64, 68, 69, 70, 71, 72, 73, 74}
	contentxs := []float64{-4, -3, -2, -1, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24}
	contentys := []float64{0, 0, 0, 0, 4, 24, 34, 40, 47, 53, 56, 58, 58, 58, 61, 66, 68, 76, 81, 86, 91, 100, 100, 100, 100, 100, 100, 100, 100}
	var structuredamagefunction = paireddata.PairedData{Xvals: structurexs, Yvals: structureys}
	var contentdamagefunction = paireddata.PairedData{Xvals: contentxs, Yvals: contentys}
	return OccupancyType{Name: "RES3C", Structuredamfun: structuredamagefunction, Contentdamfun: contentdamagefunction}
}
func res3d() OccupancyType {
	structurexs := []float64{-4, -3, -2, -1, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24}
	structureys := []float64{0, 0, 0, 0, 5, 28, 29, 31, 36, 37, 39, 40, 41, 42, 44, 46, 48, 52, 55, 58, 61, 64, 68, 69, 70, 71, 72, 73, 74}
	contentxs := []float64{-4, -3, -2, -1, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24}
	contentys := []float64{0, 0, 0, 0, 4, 24, 34, 40, 47, 53, 56, 58, 58, 58, 61, 66, 68, 76, 81, 86, 91, 100, 100, 100, 100, 100, 100, 100, 100}
	var structuredamagefunction = paireddata.PairedData{Xvals: structurexs, Yvals: structureys}
	var contentdamagefunction = paireddata.PairedData{Xvals: contentxs, Yvals: contentys}
	return OccupancyType{Name: "RES3D", Structuredamfun: structuredamagefunction, Contentdamfun: contentdamagefunction}
}
func res3e() OccupancyType {
	structurexs := []float64{-4, -3, -2, -1, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24}
	structureys := []float64{0, 0, 0, 0, 5, 28, 29, 31, 36, 37, 39, 40, 41, 42, 44, 46, 48, 52, 55, 58, 61, 64, 68, 69, 70, 71, 72, 73, 74}
	contentxs := []float64{-4, -3, -2, -1, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24}
	contentys := []float64{0, 0, 0, 0, 4, 24, 34, 40, 47, 53, 56, 58, 58, 58, 61, 66, 68, 76, 81, 86, 91, 100, 100, 100, 100, 100, 100, 100, 100}
	var structuredamagefunction = paireddata.PairedData{Xvals: structurexs, Yvals: structureys}
	var contentdamagefunction = paireddata.PairedData{Xvals: contentxs, Yvals: contentys}
	return OccupancyType{Name: "RES3E", Structuredamfun: structuredamagefunction, Contentdamfun: contentdamagefunction}
}
func res3f() OccupancyType {
	structurexs := []float64{-4, -3, -2, -1, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24}
	structureys := []float64{0, 0, 0, 0, 5, 28, 29, 31, 36, 37, 39, 40, 41, 42, 44, 46, 48, 52, 55, 58, 61, 64, 68, 69, 70, 71, 72, 73, 74}
	contentxs := []float64{-4, -3, -2, -1, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24}
	contentys := []float64{0, 0, 0, 0, 4, 24, 34, 40, 47, 53, 56, 58, 58, 58, 61, 66, 68, 76, 81, 86, 91, 100, 100, 100, 100, 100, 100, 100, 100}
	var structuredamagefunction = paireddata.PairedData{Xvals: structurexs, Yvals: structureys}
	var contentdamagefunction = paireddata.PairedData{Xvals: contentxs, Yvals: contentys}
	return OccupancyType{Name: "RES3F", Structuredamfun: structuredamagefunction, Contentdamfun: contentdamagefunction}
}
func res4() OccupancyType {
	structurexs := []float64{-4, -3, -2, -1, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24}
	structureys := []float64{0, 0, 0, 0, 0, 3, 5, 6, 7, 9, 12, 14, 18, 21, 26, 31, 36, 41, 46, 50, 54, 58, 62, 66, 70, 74, 78, 82, 86}
	contentxs := []float64{-4, -3, -2, -1, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24}
	contentys := []float64{0, 0, 0, 0, 0, 11, 19, 25, 29, 34, 39, 44, 49, 56, 65, 74, 82, 88, 98, 98, 98, 98, 98, 98, 98, 98, 98, 98, 98}
	var structuredamagefunction = paireddata.PairedData{Xvals: structurexs, Yvals: structureys}
	var contentdamagefunction = paireddata.PairedData{Xvals: contentxs, Yvals: contentys}
	return OccupancyType{Name: "RES4", Structuredamfun: structuredamagefunction, Contentdamfun: contentdamagefunction}
}
func res5() OccupancyType {
	structurexs := []float64{-4, -3, -2, -1, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24}
	structureys := []float64{0, 0, 0, 0, 0, 7, 10, 14, 15, 15, 16, 18, 20, 23, 26, 30, 34, 38, 42, 47, 52, 57, 62, 67, 72, 77, 82, 87, 92}
	contentxs := []float64{-4, -3, -2, -1, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24}
	contentys := []float64{0, 0, 0, 0, 0, 38, 60, 73, 81, 88, 94, 100, 100, 100, 100, 100, 100, 100, 100, 100, 100, 100, 100, 100, 100, 100, 100, 100, 100}
	var structuredamagefunction = paireddata.PairedData{Xvals: structurexs, Yvals: structureys}
	var contentdamagefunction = paireddata.PairedData{Xvals: contentxs, Yvals: contentys}
	return OccupancyType{Name: "RES5", Structuredamfun: structuredamagefunction, Contentdamfun: contentdamagefunction}
}
func res6() OccupancyType {
	structurexs := []float64{-4, -3, -2, -1, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24}
	structureys := []float64{0, 0, 0, 0, 0, 7, 10, 14, 15, 15, 16, 18, 20, 23, 26, 30, 34, 38, 42, 47, 52, 57, 62, 67, 72, 77, 82, 87, 92}
	contentxs := []float64{-4, -3, -2, -1, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24}
	contentys := []float64{0, 0, 0, 0, 0, 38, 60, 73, 81, 88, 94, 100, 100, 100, 100, 100, 100, 100, 100, 100, 100, 100, 100, 100, 100, 100, 100, 100, 100}
	var structuredamagefunction = paireddata.PairedData{Xvals: structurexs, Yvals: structureys}
	var contentdamagefunction = paireddata.PairedData{Xvals: contentxs, Yvals: contentys}
	return OccupancyType{Name: "RES6", Structuredamfun: structuredamagefunction, Contentdamfun: contentdamagefunction}
}
