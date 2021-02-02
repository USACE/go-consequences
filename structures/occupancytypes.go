package structures

import (
	"math/rand"

	"github.com/HenryGeorgist/go-statistics/statistics"
	"github.com/USACE/go-consequences/hazards"
	"github.com/USACE/go-consequences/paireddata"
)

//OccupancyType interface allows for multiple hazards that integrate with structures
type OccupancyType interface {
	GetStructureDamageFunctionForHazard(h hazards.HazardEvent) paireddata.ValueSampler
	GetContentDamageFunctionForHazard(h hazards.HazardEvent) paireddata.ValueSampler
}
type DamageFunctionFamily struct {
	DamageFunctions map[hazards.Parameter]paireddata.ValueSampler //parameter is a bitflag
}

//OccupancyTypeStochastic is used to describe an occupancy type with uncertainty in the damage relationships it produces an OccupancyTypeDeterministic through the UncertaintyOccupancyTypeSampler interface
type OccupancyTypeStochastic struct { //this is mutable
	Name            string
	Structuredamfun interface{} //if i make this an empty interface, it could be a value sampler, or an uncertainty valuesampler sampler...
	Contentdamfun   interface{} //if i make this an empty interface, it could be a value sampler, or an uncertainty valuesampler sampler...
}

//OccupancyTypeDeterministic is used to describe an occupancy type without uncertainty in the damage relationships
type OccupancyTypeDeterministic struct {
	Name                        string
	Structuredamfun             paireddata.ValueSampler
	Structuredamfun_MW_Salinity paireddata.ValueSampler
	Contentdamfun               paireddata.ValueSampler
}

//GetStructureDamageFunctionForHazard implements OccupancyType on OccupancyTypeDeterministic
func (o OccupancyTypeDeterministic) GetStructureDamageFunctionForHazard(h hazards.HazardEvent) paireddata.ValueSampler {
	ce, okc := h.(hazards.CoastalEvent)
	//placeholder code needs to be developed
	if okc {
		//determine which curve to supply
		if ce.Salinity {
			//saltwater
			if ce.WaveHeight < 1.0 {
				//stillwater
				return o.Structuredamfun
			} else {
				if ce.WaveHeight > 2.9 {
					//high wave
					return o.Structuredamfun
				}
				//medium wave
				return o.Structuredamfun
			}
		}
		//freshwater
		if ce.WaveHeight < 1.0 {
			//stillwater
			return o.Structuredamfun
		} else {
			if ce.WaveHeight > 2.9 {
				//high wave
				return o.Structuredamfun
			}
			//medium wave
			return o.Structuredamfun
		}
	}
	return o.Structuredamfun
}

//GetContentDamageFunctionForHazard implements OccupancyType on OccupancyTypeDeterministic
func (o OccupancyTypeDeterministic) GetContentDamageFunctionForHazard(h hazards.HazardEvent) paireddata.ValueSampler {
	ce, okc := h.(hazards.CoastalEvent)
	//placeholder code needs to be developed
	if okc {
		//determine which curve to supply
		if ce.Salinity {
			//saltwater
			if ce.WaveHeight < 1.0 {
				//stillwater
				return o.Contentdamfun
			} else {
				if ce.WaveHeight > 2.9 {
					//high wave
					return o.Contentdamfun
				}
				//medium wave
				return o.Contentdamfun
			}
		}
		//freshwater
		if ce.WaveHeight < 1.0 {
			//stillwater
			return o.Contentdamfun
		} else {
			if ce.WaveHeight > 2.9 {
				//high wave
				return o.Contentdamfun
			}
			//medium wave
			return o.Contentdamfun
		}
	}
	return o.Contentdamfun
}

//UncertaintyOccupancyTypeSampler provides the pattern for an OccupancyTypeStochastic to produce an OccupancyTypeDeterministic
type UncertaintyOccupancyTypeSampler interface {
	SampleOccupancyType(rand int64) OccupancyTypeDeterministic
	CentralTendencyOccupancyType() OccupancyTypeDeterministic
}

//SampleOccupancyType implements the UncertaintyOccupancyTypeSampler on the OccupancyTypeStochastic interface.
func (o OccupancyTypeStochastic) SampleOccupancyType(seed int64) OccupancyTypeDeterministic {
	sd, oks := o.Structuredamfun.(paireddata.ValueSampler)
	cd, okc := o.Contentdamfun.(paireddata.ValueSampler)
	if oks && okc {
		return OccupancyTypeDeterministic{Name: o.Name, Structuredamfun: sd, Contentdamfun: cd}
	}
	rand.Seed(seed)
	if oks {
		cd2, okc1 := o.Contentdamfun.(paireddata.UncertaintyValueSamplerSampler)
		if okc1 {
			cd = cd2.SampleValueSampler(rand.Float64())
		} else {
			//cd = nil
		}
	} else {
		sd2, oks1 := o.Structuredamfun.(paireddata.UncertaintyValueSamplerSampler)
		if oks1 {
			sd = sd2.SampleValueSampler(rand.Float64())
		} else {
			//sd = nil
		}
	}
	if okc {
		sd3, oks2 := o.Structuredamfun.(paireddata.UncertaintyValueSamplerSampler)
		if oks2 {
			sd = sd3.SampleValueSampler(rand.Float64())
		} else {
			//sd = nil
		}
	} else {
		cd3, okc2 := o.Contentdamfun.(paireddata.UncertaintyValueSamplerSampler)
		if okc2 {
			cd = cd3.SampleValueSampler(rand.Float64())
		} else {
			//cd = nil
		}
	}

	return OccupancyTypeDeterministic{Name: o.Name, Structuredamfun: sd, Contentdamfun: cd}
}

//CentralTendency implements the UncertaintyOccupancyTypeSampler on the OccupancyTypeStochastic interface.
func (o OccupancyTypeStochastic) CentralTendency() OccupancyTypeDeterministic {
	sd, oks := o.Structuredamfun.(paireddata.ValueSampler)
	cd, okc := o.Contentdamfun.(paireddata.ValueSampler)
	if oks && okc {
		return OccupancyTypeDeterministic{Name: o.Name, Structuredamfun: sd, Contentdamfun: cd}
	}
	//rand.Seed(seed)
	if oks {
		cd2, okc1 := o.Contentdamfun.(paireddata.UncertaintyValueSamplerSampler)
		if okc1 {
			cd = cd2.CentralTendency()
		} else {
			//cd = nil
		}
	} else {
		sd2, oks1 := o.Structuredamfun.(paireddata.UncertaintyValueSamplerSampler)
		if oks1 {
			sd = sd2.CentralTendency()
		} else {
			//sd = nil
		}
	}
	if okc {
		sd3, oks2 := o.Structuredamfun.(paireddata.UncertaintyValueSamplerSampler)
		if oks2 {
			sd = sd3.CentralTendency()
		} else {
			//sd = nil
		}
	} else {
		cd3, okc2 := o.Contentdamfun.(paireddata.UncertaintyValueSamplerSampler)
		if okc2 {
			cd = cd3.CentralTendency()
		} else {
			//cd = nil
		}
	}

	return OccupancyTypeDeterministic{Name: o.Name, Structuredamfun: sd, Contentdamfun: cd}
}

//OccupancyTypeMap produces a map of all occupancy types as OccupancyTypeStochastic so they can be joined to the structure inventory to compute damage
func OccupancyTypeMap() map[string]OccupancyTypeStochastic {
	m := make(map[string]OccupancyTypeStochastic)
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
func res11snb() OccupancyTypeStochastic {
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
	contentxs := []float64{-2.0, -1.0, 0.0, 1.0, 2.0, 3.0, 4.0, 5.0, 6.0, 7.0, 8.0, 9.0, 10.0, 11.0, 12.0, 13.0, 14.0, 15.0, 16.0}
	contentydists := make([]statistics.ContinuousDistribution, 19)
	contentydists[0] = statistics.NormalDistribution{Mean: 0, StandardDeviation: 0}
	contentydists[1] = statistics.NormalDistribution{Mean: 2.4000000953674316, StandardDeviation: 0.20000000298023224}
	contentydists[2] = statistics.NormalDistribution{Mean: 8.1000003814697266, StandardDeviation: 0.800000011920929}
	contentydists[3] = statistics.NormalDistribution{Mean: 13.300000190734863, StandardDeviation: 1.2999999523162842}
	contentydists[4] = statistics.NormalDistribution{Mean: 17.899999618530273, StandardDeviation: 1.7000000476837158}
	contentydists[5] = statistics.NormalDistribution{Mean: 22, StandardDeviation: 1.8999999761581421}
	contentydists[6] = statistics.NormalDistribution{Mean: 25.700000762939453, StandardDeviation: 2.1700000762939453}
	contentydists[7] = statistics.NormalDistribution{Mean: 28.799999237060547, StandardDeviation: 2.5}
	contentydists[8] = statistics.NormalDistribution{Mean: 31.5, StandardDeviation: 2.7999999523162842}
	contentydists[9] = statistics.NormalDistribution{Mean: 33.799999237060547, StandardDeviation: 2.9500000476837158}
	contentydists[10] = statistics.NormalDistribution{Mean: 35.700000762939453, StandardDeviation: 3.0999999046325684}
	contentydists[11] = statistics.NormalDistribution{Mean: 37.200000762939453, StandardDeviation: 3.2000000476837158}
	contentydists[12] = statistics.NormalDistribution{Mean: 38.400001525878906, StandardDeviation: 3.2999999523162842}
	contentydists[13] = statistics.NormalDistribution{Mean: 39.200000762939453, StandardDeviation: 3.4000000953674316}
	contentydists[14] = statistics.NormalDistribution{Mean: 39.700000762939453, StandardDeviation: 3.4000000953674316}
	contentydists[15] = statistics.NormalDistribution{Mean: 40, StandardDeviation: 3.4100000858306885}
	contentydists[16] = statistics.NormalDistribution{Mean: 40, StandardDeviation: 3.4100000858306885}
	contentydists[17] = statistics.NormalDistribution{Mean: 40, StandardDeviation: 3.4100000858306885}
	contentydists[18] = statistics.NormalDistribution{Mean: 40, StandardDeviation: 3.4100000858306885}
	var structuredamagefunctionStochastic = paireddata.UncertaintyPairedData{Xvals: structurexs, Yvals: structureydists}
	var contentdamagefunctionStochastic = paireddata.UncertaintyPairedData{Xvals: contentxs, Yvals: contentydists}
	return OccupancyTypeStochastic{Name: "RES1-1SNB", Structuredamfun: structuredamagefunctionStochastic, Contentdamfun: contentdamagefunctionStochastic}
}
func res11swb() OccupancyTypeStochastic {
	structurexs := []float64{-8, -7, -6, -5, -4, -3, -2, -1, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16}
	structureydists := make([]statistics.ContinuousDistribution, 25)
	structureydists[0] = statistics.NormalDistribution{Mean: 0, StandardDeviation: 0}
	structureydists[1] = statistics.NormalDistribution{Mean: 0.699999988079071, StandardDeviation: 0.0099999997764825821}
	structureydists[2] = statistics.NormalDistribution{Mean: 0.800000011920929, StandardDeviation: 0.019999999552965164}
	structureydists[3] = statistics.NormalDistribution{Mean: 2.4000000953674316, StandardDeviation: 0.10000000149011612}
	structureydists[4] = statistics.NormalDistribution{Mean: 5.1999998092651367, StandardDeviation: 0.30000001192092896}
	structureydists[5] = statistics.NormalDistribution{Mean: 9, StandardDeviation: 0.699999988079071}
	structureydists[6] = statistics.NormalDistribution{Mean: 13.800000190734863, StandardDeviation: 0.85000002384185791}
	structureydists[7] = statistics.NormalDistribution{Mean: 19.399999618530273, StandardDeviation: 0.82999998331069946}
	structureydists[8] = statistics.NormalDistribution{Mean: 25.5, StandardDeviation: 0.85000002384185791}
	structureydists[9] = statistics.NormalDistribution{Mean: 32, StandardDeviation: 0.95999997854232788}
	structureydists[10] = statistics.NormalDistribution{Mean: 38.700000762939453, StandardDeviation: 1.1399999856948853}
	structureydists[11] = statistics.NormalDistribution{Mean: 45.5, StandardDeviation: 1.3700000047683716}
	structureydists[12] = statistics.NormalDistribution{Mean: 52.200000762939453, StandardDeviation: 1.6299999952316284}
	structureydists[13] = statistics.NormalDistribution{Mean: 58.599998474121094, StandardDeviation: 1.8899999856948853}
	structureydists[14] = statistics.NormalDistribution{Mean: 64.5, StandardDeviation: 1.8999999761581421}
	structureydists[15] = statistics.NormalDistribution{Mean: 69.800003051757812, StandardDeviation: 2.0199999809265137}
	structureydists[16] = statistics.NormalDistribution{Mean: 74.199996948242188, StandardDeviation: 2.0399999618530273}
	structureydists[17] = statistics.NormalDistribution{Mean: 77.699996948242188, StandardDeviation: 2.130000114440918}
	structureydists[18] = statistics.NormalDistribution{Mean: 80.0999984741211, StandardDeviation: 2.2000000476837158}
	structureydists[19] = statistics.NormalDistribution{Mean: 81.0999984741211, StandardDeviation: 2.2999999523162842}
	structureydists[20] = statistics.NormalDistribution{Mean: 81.0999984741211, StandardDeviation: 2.2999999523162842}
	structureydists[21] = statistics.NormalDistribution{Mean: 81.0999984741211, StandardDeviation: 2.2999999523162842}
	structureydists[22] = statistics.NormalDistribution{Mean: 81.0999984741211, StandardDeviation: 2.2999999523162842}
	structureydists[23] = statistics.NormalDistribution{Mean: 81.0999984741211, StandardDeviation: 2.2999999523162842}
	structureydists[24] = statistics.NormalDistribution{Mean: 81.0999984741211, StandardDeviation: 2.2999999523162842}
	contentxs := []float64{-8, -7, -6, -5, -4, -3, -2, -1, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16}
	contentydists := make([]statistics.ContinuousDistribution, 25)
	contentydists[0] = statistics.NormalDistribution{Mean: 0.10000000149011612, StandardDeviation: 0}
	contentydists[1] = statistics.NormalDistribution{Mean: 0.800000011920929, StandardDeviation: 0.0099999997764825821}
	contentydists[2] = statistics.NormalDistribution{Mean: 2.0999999046325684, StandardDeviation: 0.10000000149011612}
	contentydists[3] = statistics.NormalDistribution{Mean: 3.7000000476837158, StandardDeviation: 0.30000001192092896}
	contentydists[4] = statistics.NormalDistribution{Mean: 5.6999998092651367, StandardDeviation: 0.5}
	contentydists[5] = statistics.NormalDistribution{Mean: 8, StandardDeviation: 0.60000002384185791}
	contentydists[6] = statistics.NormalDistribution{Mean: 10.5, StandardDeviation: 0.74000000953674316}
	contentydists[7] = statistics.NormalDistribution{Mean: 13.199999809265137, StandardDeviation: 0.72000002861022949}
	contentydists[8] = statistics.NormalDistribution{Mean: 16, StandardDeviation: 0.74000000953674316}
	contentydists[9] = statistics.NormalDistribution{Mean: 18.899999618530273, StandardDeviation: 0.82999998331069946}
	contentydists[10] = statistics.NormalDistribution{Mean: 21.799999237060547, StandardDeviation: 0.98000001907348633}
	contentydists[11] = statistics.NormalDistribution{Mean: 24.700000762939453, StandardDeviation: 1.1699999570846558}
	contentydists[12] = statistics.NormalDistribution{Mean: 27.399999618530273, StandardDeviation: 1.3899999856948853}
	contentydists[13] = statistics.NormalDistribution{Mean: 30, StandardDeviation: 1.6000000238418579}
	contentydists[14] = statistics.NormalDistribution{Mean: 32.400001525878906, StandardDeviation: 1.8400000333786011}
	contentydists[15] = statistics.NormalDistribution{Mean: 34.5, StandardDeviation: 2}
	contentydists[16] = statistics.NormalDistribution{Mean: 36.299999237060547, StandardDeviation: 2.1600000858306885}
	contentydists[17] = statistics.NormalDistribution{Mean: 37.700000762939453, StandardDeviation: 2.2999999523162842}
	contentydists[18] = statistics.NormalDistribution{Mean: 38.599998474121094, StandardDeviation: 2.4000000953674316}
	contentydists[19] = statistics.NormalDistribution{Mean: 39.099998474121094, StandardDeviation: 2.4500000476837158}
	contentydists[20] = statistics.NormalDistribution{Mean: 39.099998474121094, StandardDeviation: 2.4500000476837158}
	contentydists[21] = statistics.NormalDistribution{Mean: 39.099998474121094, StandardDeviation: 2.4500000476837158}
	contentydists[22] = statistics.NormalDistribution{Mean: 39.099998474121094, StandardDeviation: 2.4500000476837158}
	contentydists[23] = statistics.NormalDistribution{Mean: 39.099998474121094, StandardDeviation: 2.4500000476837158}
	contentydists[24] = statistics.NormalDistribution{Mean: 39.099998474121094, StandardDeviation: 2.4500000476837158}
	var structuredamagefunctionStochastic = paireddata.UncertaintyPairedData{Xvals: structurexs, Yvals: structureydists}
	var contentdamagefunctionStochastic = paireddata.UncertaintyPairedData{Xvals: contentxs, Yvals: contentydists}
	return OccupancyTypeStochastic{Name: "RES1-1SWB", Structuredamfun: structuredamagefunctionStochastic, Contentdamfun: contentdamagefunctionStochastic}
}
func agr1() OccupancyTypeStochastic {
	structurexs := []float64{-4, -3, -2, -1, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24}
	structureys := []float64{0, 0, 0, 0, 0, 6, 11, 15, 19, 25, 30, 35, 41, 46, 51, 57, 63, 70, 75, 79, 82, 84, 87, 89, 90, 92, 93, 95, 96}
	contentxs := []float64{-4, -3, -2, -1, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24}
	contentys := []float64{0, 0, 0, 0, 6, 20, 43, 58, 65, 66, 66, 67, 70, 75, 76, 76, 76, 77, 77, 77, 78, 78, 78, 79, 79, 79, 79, 80, 80}
	var structuredamagefunction = paireddata.PairedData{Xvals: structurexs, Yvals: structureys}
	var contentdamagefunction = paireddata.PairedData{Xvals: contentxs, Yvals: contentys}
	return OccupancyTypeStochastic{Name: "AGR1", Structuredamfun: structuredamagefunction, Contentdamfun: contentdamagefunction}
}
func com1() OccupancyTypeStochastic {
	structurexs := []float64{-4, -3, -2, -1, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24}
	structureys := []float64{0, 0, 0, 0, 1, 9, 14, 16, 18, 20, 23, 26, 30, 34, 38, 42, 47, 51, 55, 58, 61, 64, 67, 69, 71, 74, 76, 78, 80}
	contentxs := []float64{-4, -3, -2, -1, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24}
	contentys := []float64{0, 0, 0, 0, 2, 26, 42, 56, 68, 78, 83, 85, 87, 88, 89, 90, 91, 92, 92, 92, 93, 93, 94, 94, 94, 94, 94, 94, 94}
	var structuredamagefunction = paireddata.PairedData{Xvals: structurexs, Yvals: structureys}
	var contentdamagefunction = paireddata.PairedData{Xvals: contentxs, Yvals: contentys}
	return OccupancyTypeStochastic{Name: "COM1", Structuredamfun: structuredamagefunction, Contentdamfun: contentdamagefunction}
}
func com2() OccupancyTypeStochastic {
	structurexs := []float64{-4, -3, -2, -1, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24}
	structureys := []float64{0, 0, 0, 0, 0, 5, 8, 11, 13, 16, 19, 22, 25, 29, 32, 37, 41, 45, 49, 52, 55, 58, 61, 63, 66, 68, 70, 71, 73}
	contentxs := []float64{-4, -3, -2, -1, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24}
	contentys := []float64{0, 0, 0, 0, 3, 16, 27, 36, 49, 57, 63, 69, 72, 76, 80, 82, 84, 86, 87, 87, 88, 89, 89, 89, 89, 89, 89, 89, 89}
	var structuredamagefunction = paireddata.PairedData{Xvals: structurexs, Yvals: structureys}
	var contentdamagefunction = paireddata.PairedData{Xvals: contentxs, Yvals: contentys}
	return OccupancyTypeStochastic{Name: "COM2", Structuredamfun: structuredamagefunction, Contentdamfun: contentdamagefunction}
}
func com3() OccupancyTypeStochastic {
	structurexs := []float64{-4, -3, -2, -1, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24}
	structureys := []float64{0, 0, 0, 0, 0, 9, 12, 13, 16, 19, 22, 25, 28, 32, 35, 39, 43, 47, 50, 54, 57, 61, 64, 68, 71, 75, 78, 82, 85}
	contentxs := []float64{-4, -3, -2, -1, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24}
	contentys := []float64{0, 0, 0, 0, 4, 29, 46, 67, 79, 85, 91, 92, 92, 93, 94, 96, 96, 97, 97, 98, 98, 98, 98, 98, 98, 98, 98, 98, 98}
	var structuredamagefunction = paireddata.PairedData{Xvals: structurexs, Yvals: structureys}
	var contentdamagefunction = paireddata.PairedData{Xvals: contentxs, Yvals: contentys}
	return OccupancyTypeStochastic{Name: "COM3", Structuredamfun: structuredamagefunction, Contentdamfun: contentdamagefunction}
}
func com4() OccupancyTypeStochastic {
	structurexs := []float64{-4, -3, -2, -1, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24}
	structureys := []float64{0, 0, 0, 0, 2, 11, 16, 22, 28, 35, 38, 41, 44, 47, 50, 54, 57, 59, 62, 66, 68, 70, 72, 74, 76, 77, 78, 79, 80}
	contentxs := []float64{-4, -3, -2, -1, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24}
	contentys := []float64{0, 0, 0, 0, 2, 18, 25, 35, 43, 49, 52, 55, 57, 58, 60, 65, 67, 68, 69, 70, 71, 71, 72, 72, 72, 72, 72, 72, 72}
	var structuredamagefunction = paireddata.PairedData{Xvals: structurexs, Yvals: structureys}
	var contentdamagefunction = paireddata.PairedData{Xvals: contentxs, Yvals: contentys}
	return OccupancyTypeStochastic{Name: "COM4", Structuredamfun: structuredamagefunction, Contentdamfun: contentdamagefunction}
}
func com5() OccupancyTypeStochastic {
	structurexs := []float64{-4, -3, -2, -1, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24}
	structureys := []float64{0, 0, 0, 0, 0, 11, 11, 12, 13, 15, 17, 19, 22, 24, 28, 31, 34, 37, 40, 44, 48, 51, 55, 59, 63, 67, 71, 75, 79}
	contentxs := []float64{-4, -3, -2, -1, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24}
	contentys := []float64{0, 0, 0, 0, 0, 50, 74, 83, 100, 100, 100, 100, 100, 100, 100, 100, 100, 100, 100, 100, 100, 100, 100, 100, 100, 100, 100, 100, 100}
	var structuredamagefunction = paireddata.PairedData{Xvals: structurexs, Yvals: structureys}
	var contentdamagefunction = paireddata.PairedData{Xvals: contentxs, Yvals: contentys}
	return OccupancyTypeStochastic{Name: "COM5", Structuredamfun: structuredamagefunction, Contentdamfun: contentdamagefunction}
}
func com6() OccupancyTypeStochastic {
	structurexs := []float64{-4, -3, -2, -1, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24}
	structureys := []float64{0, 0, 0, 0, 0, 0, 0, 20, 25, 30, 35, 40, 43, 47, 50, 53, 55, 57, 60, 60, 60, 60, 60, 60, 60, 60, 60, 60, 60}
	contentxs := []float64{-4, -3, -2, -1, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24}
	contentys := []float64{0, 0, 0, 0, 0, 0, 0, 10, 20, 30, 65, 72, 78, 85, 95, 95, 95, 95, 96, 96, 96, 96, 96, 96, 96, 96, 96, 96, 96}
	var structuredamagefunction = paireddata.PairedData{Xvals: structurexs, Yvals: structureys}
	var contentdamagefunction = paireddata.PairedData{Xvals: contentxs, Yvals: contentys}
	return OccupancyTypeStochastic{Name: "COM6", Structuredamfun: structuredamagefunction, Contentdamfun: contentdamagefunction}
}
func com7() OccupancyTypeStochastic {
	structurexs := []float64{-4, -3, -2, -1, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24}
	structureys := []float64{0, 0, 0, 0, 2, 11, 12, 13, 14, 16, 17, 18, 20, 22, 24, 27, 30, 34, 37, 41, 44, 48, 51, 54, 56, 59, 61, 64, 66}
	contentxs := []float64{-4, -3, -2, -1, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24}
	contentys := []float64{0, 0, 0, 0, 0, 28, 51, 60, 63, 67, 71, 72, 74, 77, 81, 86, 92, 94, 97, 99, 100, 100, 100, 100, 100, 100, 100, 100, 100}
	var structuredamagefunction = paireddata.PairedData{Xvals: structurexs, Yvals: structureys}
	var contentdamagefunction = paireddata.PairedData{Xvals: contentxs, Yvals: contentys}
	return OccupancyTypeStochastic{Name: "COM7", Structuredamfun: structuredamagefunction, Contentdamfun: contentdamagefunction}
}
func com8() OccupancyTypeStochastic {
	structurexs := []float64{-4, -3, -2, -1, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24}
	structureys := []float64{0, 0, 0, 0, 1, 9, 11, 12, 14, 16, 18, 20, 22, 26, 29, 33, 37, 41, 45, 50, 53, 57, 60, 63, 66, 69, 73, 76, 78}
	contentxs := []float64{-4, -3, -2, -1, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24}
	contentys := []float64{0, 0, 0, 0, 13, 45, 55, 64, 73, 77, 80, 82, 83, 85, 87, 89, 90, 91, 92, 93, 94, 95, 96, 96, 96, 96, 96, 96, 96}
	var structuredamagefunction = paireddata.PairedData{Xvals: structurexs, Yvals: structureys}
	var contentdamagefunction = paireddata.PairedData{Xvals: contentxs, Yvals: contentys}
	return OccupancyTypeStochastic{Name: "COM8", Structuredamfun: structuredamagefunction, Contentdamfun: contentdamagefunction}
}
func com9() OccupancyTypeStochastic {
	structurexs := []float64{-4, -3, -2, -1, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24}
	structureys := []float64{0, 0, 0, 0, 0, 2, 4, 5, 5, 5, 6, 8, 10, 12, 15, 20, 24, 29, 35, 42, 49, 56, 62, 68, 74, 80, 86, 92, 98}
	contentxs := []float64{-4, -3, -2, -1, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24}
	contentys := []float64{0, 0, 0, 0, 0, 4, 6, 8, 9, 10, 12, 17, 22, 30, 41, 57, 66, 73, 79, 84, 90, 97, 98, 100, 100, 100, 100, 100, 100}
	var structuredamagefunction = paireddata.PairedData{Xvals: structurexs, Yvals: structureys}
	var contentdamagefunction = paireddata.PairedData{Xvals: contentxs, Yvals: contentys}
	return OccupancyTypeStochastic{Name: "COM9", Structuredamfun: structuredamagefunction, Contentdamfun: contentdamagefunction}
}
func com10() OccupancyTypeStochastic {
	structurexs := []float64{-4, -3, -2, -1, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24}
	structureys := []float64{0, 0, 0, 0, 0, 3, 5, 6, 7, 8, 10, 13, 17, 21, 25, 30, 35, 41, 47, 52, 58, 65, 71, 76, 81, 86, 91, 96, 100}
	contentxs := []float64{-4, -3, -2, -1, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24}
	contentys := []float64{0, 0, 0, 0, 0, 11, 17, 20, 23, 25, 29, 35, 42, 51, 63, 77, 93, 100, 100, 100, 100, 100, 100, 100, 100, 100, 100, 100, 100}
	var structuredamagefunction = paireddata.PairedData{Xvals: structurexs, Yvals: structureys}
	var contentdamagefunction = paireddata.PairedData{Xvals: contentxs, Yvals: contentys}
	return OccupancyTypeStochastic{Name: "COM10", Structuredamfun: structuredamagefunction, Contentdamfun: contentdamagefunction}
}
func edu1() OccupancyTypeStochastic {
	structurexs := []float64{-4, -3, -2, -1, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24}
	structureys := []float64{0, 0, 0, 0, 0, 5, 7, 9, 9, 10, 11, 13, 15, 17, 20, 24, 28, 33, 39, 45, 52, 59, 64, 69, 74, 79, 84, 89, 94}
	contentxs := []float64{-4, -3, -2, -1, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24}
	contentys := []float64{0, 0, 0, 0, 0, 27, 38, 53, 64, 68, 70, 72, 75, 79, 83, 88, 94, 100, 100, 100, 100, 100, 100, 100, 100, 100, 100, 100, 100}
	var structuredamagefunction = paireddata.PairedData{Xvals: structurexs, Yvals: structureys}
	var contentdamagefunction = paireddata.PairedData{Xvals: contentxs, Yvals: contentys}
	return OccupancyTypeStochastic{Name: "EDU1", Structuredamfun: structuredamagefunction, Contentdamfun: contentdamagefunction}
}
func edu2() OccupancyTypeStochastic {
	structurexs := []float64{-4, -3, -2, -1, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24}
	structureys := []float64{0, 0, 0, 0, 0, 5, 7, 9, 9, 10, 11, 13, 15, 17, 20, 24, 28, 33, 39, 45, 52, 59, 64, 69, 74, 79, 84, 89, 94}
	contentxs := []float64{-4, -3, -2, -1, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24}
	contentys := []float64{0, 0, 0, 0, 0, 27, 38, 53, 64, 68, 70, 72, 75, 79, 83, 88, 94, 100, 100, 100, 100, 100, 100, 100, 100, 100, 100, 100, 100}
	var structuredamagefunction = paireddata.PairedData{Xvals: structurexs, Yvals: structureys}
	var contentdamagefunction = paireddata.PairedData{Xvals: contentxs, Yvals: contentys}
	return OccupancyTypeStochastic{Name: "EDU2", Structuredamfun: structuredamagefunction, Contentdamfun: contentdamagefunction}
}
func gov1() OccupancyTypeStochastic {
	structurexs := []float64{-4, -3, -2, -1, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24}
	structureys := []float64{0, 0, 0, 0, 0, 5, 8, 13, 14, 14, 15, 17, 19, 22, 26, 31, 37, 44, 51, 59, 65, 70, 74, 79, 83, 87, 91, 95, 98}
	contentxs := []float64{-4, -3, -2, -1, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24}
	contentys := []float64{0, 0, 0, 0, 0, 30, 59, 74, 83, 90, 100, 100, 100, 100, 100, 100, 100, 100, 100, 100, 100, 100, 100, 100, 100, 100, 100, 100, 100}
	var structuredamagefunction = paireddata.PairedData{Xvals: structurexs, Yvals: structureys}
	var contentdamagefunction = paireddata.PairedData{Xvals: contentxs, Yvals: contentys}
	return OccupancyTypeStochastic{Name: "GOV1", Structuredamfun: structuredamagefunction, Contentdamfun: contentdamagefunction}
}
func gov2() OccupancyTypeStochastic {
	structurexs := []float64{-4, -3, -2, -1, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24}
	structureys := []float64{0, 0, 0, 0, 0, 7, 10, 11, 12, 15, 17, 20, 23, 27, 31, 35, 40, 44, 48, 52, 56, 60, 64, 68, 72, 76, 80, 84, 88}
	contentxs := []float64{-4, -3, -2, -1, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24}
	contentys := []float64{0, 0, 0, 0, 0, 8, 20, 38, 55, 70, 81, 89, 98, 100, 100, 100, 100, 100, 100, 100, 100, 100, 100, 100, 100, 100, 100, 100, 100}
	var structuredamagefunction = paireddata.PairedData{Xvals: structurexs, Yvals: structureys}
	var contentdamagefunction = paireddata.PairedData{Xvals: contentxs, Yvals: contentys}
	return OccupancyTypeStochastic{Name: "GOV2", Structuredamfun: structuredamagefunction, Contentdamfun: contentdamagefunction}
}
func ind1() OccupancyTypeStochastic {
	structurexs := []float64{-4, -3, -2, -1, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24}
	structureys := []float64{0, 0, 0, 0, 1, 10, 12, 15, 19, 22, 26, 30, 35, 39, 42, 48, 50, 51, 53, 54, 55, 55, 56, 56, 57, 57, 57, 58, 58}
	contentxs := []float64{-4, -3, -2, -1, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24}
	contentys := []float64{0, 0, 0, 0, 0, 15, 24, 34, 41, 47, 52, 57, 60, 63, 64, 66, 68, 69, 72, 73, 73, 73, 74, 74, 74, 74, 75, 75, 75}
	var structuredamagefunction = paireddata.PairedData{Xvals: structurexs, Yvals: structureys}
	var contentdamagefunction = paireddata.PairedData{Xvals: contentxs, Yvals: contentys}
	return OccupancyTypeStochastic{Name: "IND1", Structuredamfun: structuredamagefunction, Contentdamfun: contentdamagefunction}
}
func ind2() OccupancyTypeStochastic {
	structurexs := []float64{-4, -3, -2, -1, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24}
	structureys := []float64{0, 0, 0, 0, 1, 9, 14, 17, 22, 26, 30, 32, 35, 37, 39, 43, 46, 48, 50, 51, 54, 55, 57, 59, 60, 62, 63, 65, 66}
	contentxs := []float64{-4, -3, -2, -1, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24}
	contentys := []float64{0, 0, 0, 0, 0, 9, 23, 35, 44, 52, 58, 62, 65, 68, 70, 73, 74, 77, 78, 78, 79, 80, 80, 80, 80, 81, 81, 81, 81}
	var structuredamagefunction = paireddata.PairedData{Xvals: structurexs, Yvals: structureys}
	var contentdamagefunction = paireddata.PairedData{Xvals: contentxs, Yvals: contentys}
	return OccupancyTypeStochastic{Name: "IND2", Structuredamfun: structuredamagefunction, Contentdamfun: contentdamagefunction}
}
func ind3() OccupancyTypeStochastic {
	structurexs := []float64{-4, -3, -2, -1, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24}
	structureys := []float64{0, 0, 0, 0, 0, 13, 14, 19, 22, 25, 28, 30, 33, 34, 36, 39, 40, 42, 42, 43, 43, 44, 44, 44, 44, 44, 45, 45, 45}
	contentxs := []float64{-4, -3, -2, -1, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24}
	contentys := []float64{0, 0, 0, 0, 2, 20, 41, 51, 62, 67, 71, 73, 76, 78, 79, 82, 83, 84, 86, 87, 87, 88, 88, 88, 88, 88, 88, 88, 88}
	var structuredamagefunction = paireddata.PairedData{Xvals: structurexs, Yvals: structureys}
	var contentdamagefunction = paireddata.PairedData{Xvals: contentxs, Yvals: contentys}
	return OccupancyTypeStochastic{Name: "IND3", Structuredamfun: structuredamagefunction, Contentdamfun: contentdamagefunction}
}
func ind4() OccupancyTypeStochastic {
	structurexs := []float64{-4, -3, -2, -1, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24}
	structureys := []float64{0, 0, 0, 0, 0, 10, 14, 18, 22, 26, 34, 41, 42, 42, 45, 47, 49, 50, 50, 50, 50, 50, 50, 50, 50, 50, 50, 50, 50}
	contentxs := []float64{-4, -3, -2, -1, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24}
	contentys := []float64{0, 0, 0, 0, 15, 20, 26, 31, 37, 40, 44, 48, 53, 56, 57, 60, 62, 63, 63, 63, 64, 65, 65, 65, 65, 65, 65, 65, 65}
	var structuredamagefunction = paireddata.PairedData{Xvals: structurexs, Yvals: structureys}
	var contentdamagefunction = paireddata.PairedData{Xvals: contentxs, Yvals: contentys}
	return OccupancyTypeStochastic{Name: "IND4", Structuredamfun: structuredamagefunction, Contentdamfun: contentdamagefunction}
}
func ind5() OccupancyTypeStochastic {
	structurexs := []float64{-4, -3, -2, -1, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24}
	structureys := []float64{0, 0, 0, 0, 0, 13, 14, 19, 22, 25, 28, 30, 33, 34, 36, 39, 40, 42, 42, 43, 43, 44, 44, 44, 44, 44, 45, 45, 45}
	contentxs := []float64{-4, -3, -2, -1, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24}
	contentys := []float64{0, 0, 0, 0, 2, 20, 41, 51, 62, 67, 71, 73, 76, 78, 79, 82, 83, 84, 86, 87, 87, 88, 88, 88, 88, 88, 88, 88, 88}
	var structuredamagefunction = paireddata.PairedData{Xvals: structurexs, Yvals: structureys}
	var contentdamagefunction = paireddata.PairedData{Xvals: contentxs, Yvals: contentys}
	return OccupancyTypeStochastic{Name: "IND5", Structuredamfun: structuredamagefunction, Contentdamfun: contentdamagefunction}
}
func ind6() OccupancyTypeStochastic {
	structurexs := []float64{-4, -3, -2, -1, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24}
	structureys := []float64{0, 0, 0, 0, 0, 22, 31, 37, 43, 47, 50, 54, 57, 61, 63, 64, 65, 67, 68, 69, 70, 71, 72, 73, 74, 75, 76, 76, 77}
	contentxs := []float64{-4, -3, -2, -1, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24}
	contentys := []float64{0, 0, 0, 0, 0, 20, 35, 47, 56, 59, 66, 69, 71, 72, 78, 79, 80, 80, 81, 81, 81, 82, 82, 82, 83, 83, 83, 83, 83}
	var structuredamagefunction = paireddata.PairedData{Xvals: structurexs, Yvals: structureys}
	var contentdamagefunction = paireddata.PairedData{Xvals: contentxs, Yvals: contentys}
	return OccupancyTypeStochastic{Name: "IND6", Structuredamfun: structuredamagefunction, Contentdamfun: contentdamagefunction}
}
func rel1() OccupancyTypeStochastic {
	structurexs := []float64{-4, -3, -2, -1, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24}
	structureys := []float64{0, 0, 0, 0, 0, 10, 11, 11, 12, 12, 13, 14, 14, 15, 17, 19, 24, 30, 38, 45, 52, 58, 64, 69, 74, 78, 82, 85, 88}
	contentxs := []float64{-4, -3, -2, -1, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24}
	contentys := []float64{0, 0, 0, 0, 10, 52, 72, 85, 92, 95, 98, 99, 100, 100, 100, 100, 100, 100, 100, 100, 100, 100, 100, 100, 100, 100, 100, 100, 100}
	var structuredamagefunction = paireddata.PairedData{Xvals: structurexs, Yvals: structureys}
	var contentdamagefunction = paireddata.PairedData{Xvals: contentxs, Yvals: contentys}
	return OccupancyTypeStochastic{Name: "REL1", Structuredamfun: structuredamagefunction, Contentdamfun: contentdamagefunction}
}
func res12snb() OccupancyTypeStochastic {
	structurexs := []float64{-2, -1, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16}
	structureydists := make([]statistics.ContinuousDistribution, 19)
	structureydists[0] = statistics.NormalDistribution{Mean: 0, StandardDeviation: 0}
	structureydists[1] = statistics.NormalDistribution{Mean: 3, StandardDeviation: 0.30000001192092896}
	structureydists[2] = statistics.NormalDistribution{Mean: 9.3000001907348633, StandardDeviation: 1}
	structureydists[3] = statistics.NormalDistribution{Mean: 15.199999809265137, StandardDeviation: 1.5}
	structureydists[4] = statistics.NormalDistribution{Mean: 20.899999618530273, StandardDeviation: 2}
	structureydists[5] = statistics.NormalDistribution{Mean: 26.299999237060547, StandardDeviation: 2.4000000953674316}
	structureydists[6] = statistics.NormalDistribution{Mean: 31.399999618530273, StandardDeviation: 2.7000000476837158}
	structureydists[7] = statistics.NormalDistribution{Mean: 36.200000762939453, StandardDeviation: 3.0999999046325684}
	structureydists[8] = statistics.NormalDistribution{Mean: 40.700000762939453, StandardDeviation: 3.2999999523162842}
	structureydists[9] = statistics.NormalDistribution{Mean: 44.900001525878906, StandardDeviation: 3.4500000476837158}
	structureydists[10] = statistics.NormalDistribution{Mean: 48.799999237060547, StandardDeviation: 3.4900000095367432}
	structureydists[11] = statistics.NormalDistribution{Mean: 52.400001525878906, StandardDeviation: 3.5099999904632568}
	structureydists[12] = statistics.NormalDistribution{Mean: 55.700000762939453, StandardDeviation: 3.5499999523162842}
	structureydists[13] = statistics.NormalDistribution{Mean: 58.700000762939453, StandardDeviation: 3.5999999046325684}
	structureydists[14] = statistics.NormalDistribution{Mean: 61.400001525878906, StandardDeviation: 3.6500000953674316}
	structureydists[15] = statistics.NormalDistribution{Mean: 63.799999237060547, StandardDeviation: 3.7000000476837158}
	structureydists[16] = statistics.NormalDistribution{Mean: 65.9000015258789, StandardDeviation: 3.7200000286102295}
	structureydists[17] = statistics.NormalDistribution{Mean: 67.699996948242188, StandardDeviation: 3.75}
	structureydists[18] = statistics.NormalDistribution{Mean: 69.199996948242188, StandardDeviation: 3.7999999523162842}
	contentxs := []float64{-2, -1, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16}
	contentydists := make([]statistics.ContinuousDistribution, 19)
	contentydists[0] = statistics.NormalDistribution{Mean: 0, StandardDeviation: 0}
	contentydists[1] = statistics.NormalDistribution{Mean: 1, StandardDeviation: 0.10000000149011612}
	contentydists[2] = statistics.NormalDistribution{Mean: 5, StandardDeviation: 0.550000011920929}
	contentydists[3] = statistics.NormalDistribution{Mean: 8.6999998092651367, StandardDeviation: 1}
	contentydists[4] = statistics.NormalDistribution{Mean: 12.199999809265137, StandardDeviation: 1.3999999761581421}
	contentydists[5] = statistics.NormalDistribution{Mean: 15.5, StandardDeviation: 1.7999999523162842}
	contentydists[6] = statistics.NormalDistribution{Mean: 18.5, StandardDeviation: 2.0999999046325684}
	contentydists[7] = statistics.NormalDistribution{Mean: 21.299999237060547, StandardDeviation: 2.4000000953674316}
	contentydists[8] = statistics.NormalDistribution{Mean: 23.899999618530273, StandardDeviation: 2.7000000476837158}
	contentydists[9] = statistics.NormalDistribution{Mean: 26.299999237060547, StandardDeviation: 3}
	contentydists[10] = statistics.NormalDistribution{Mean: 28.399999618530273, StandardDeviation: 3.0999999046325684}
	contentydists[11] = statistics.NormalDistribution{Mean: 30.299999237060547, StandardDeviation: 3.2999999523162842}
	contentydists[12] = statistics.NormalDistribution{Mean: 32, StandardDeviation: 3.5}
	contentydists[13] = statistics.NormalDistribution{Mean: 33.400001525878906, StandardDeviation: 3.5}
	contentydists[14] = statistics.NormalDistribution{Mean: 34.700000762939453, StandardDeviation: 3.5}
	contentydists[15] = statistics.NormalDistribution{Mean: 35.599998474121094, StandardDeviation: 3.5}
	contentydists[16] = statistics.NormalDistribution{Mean: 36.400001525878906, StandardDeviation: 3.5999999046325684}
	contentydists[17] = statistics.NormalDistribution{Mean: 36.900001525878906, StandardDeviation: 3.5999999046325684}
	contentydists[18] = statistics.NormalDistribution{Mean: 37.200000762939453, StandardDeviation: 3.5999999046325684}
	var structuredamagefunctionStochastic = paireddata.UncertaintyPairedData{Xvals: structurexs, Yvals: structureydists}
	var contentdamagefunctionStochastic = paireddata.UncertaintyPairedData{Xvals: contentxs, Yvals: contentydists}
	return OccupancyTypeStochastic{Name: "RES1-2SNB", Structuredamfun: structuredamagefunctionStochastic, Contentdamfun: contentdamagefunctionStochastic}
}
func res12swb() OccupancyTypeStochastic {
	structurexs := []float64{-8, -7, -6, -5, -4, -3, -2, -1, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16}
	structureydists := make([]statistics.ContinuousDistribution, 25)
	structureydists[0] = statistics.NormalDistribution{Mean: 1.7000000476837158, StandardDeviation: 0}
	structureydists[1] = statistics.NormalDistribution{Mean: 1.7000000476837158, StandardDeviation: 0}
	structureydists[2] = statistics.NormalDistribution{Mean: 1.8999999761581421, StandardDeviation: 0.0099999997764825821}
	structureydists[3] = statistics.NormalDistribution{Mean: 2.9000000953674316, StandardDeviation: 0.10000000149011612}
	structureydists[4] = statistics.NormalDistribution{Mean: 4.6999998092651367, StandardDeviation: 0.30000001192092896}
	structureydists[5] = statistics.NormalDistribution{Mean: 7.1999998092651367, StandardDeviation: 0.60000002384185791}
	structureydists[6] = statistics.NormalDistribution{Mean: 10.199999809265137, StandardDeviation: 0.89999997615814209}
	structureydists[7] = statistics.NormalDistribution{Mean: 13.899999618530273, StandardDeviation: 1.1000000238418579}
	structureydists[8] = statistics.NormalDistribution{Mean: 17.899999618530273, StandardDeviation: 1.3200000524520874}
	structureydists[9] = statistics.NormalDistribution{Mean: 22.299999237060547, StandardDeviation: 1.3500000238418579}
	structureydists[10] = statistics.NormalDistribution{Mean: 27, StandardDeviation: 1.5}
	structureydists[11] = statistics.NormalDistribution{Mean: 31.899999618530273, StandardDeviation: 1.75}
	structureydists[12] = statistics.NormalDistribution{Mean: 36.900001525878906, StandardDeviation: 2.0399999618530273}
	structureydists[13] = statistics.NormalDistribution{Mean: 41.900001525878906, StandardDeviation: 2.3399999141693115}
	structureydists[14] = statistics.NormalDistribution{Mean: 46.900001525878906, StandardDeviation: 2.5999999046325684}
	structureydists[15] = statistics.NormalDistribution{Mean: 51.799999237060547, StandardDeviation: 2.7000000476837158}
	structureydists[16] = statistics.NormalDistribution{Mean: 56.400001525878906, StandardDeviation: 2.75}
	structureydists[17] = statistics.NormalDistribution{Mean: 60.799999237060547, StandardDeviation: 2.7599999904632568}
	structureydists[18] = statistics.NormalDistribution{Mean: 64.800003051757812, StandardDeviation: 2.7699999809265137}
	structureydists[19] = statistics.NormalDistribution{Mean: 68.4000015258789, StandardDeviation: 2.7799999713897705}
	structureydists[20] = statistics.NormalDistribution{Mean: 71.4000015258789, StandardDeviation: 2.7899999618530273}
	structureydists[21] = statistics.NormalDistribution{Mean: 73.699996948242188, StandardDeviation: 2.7999999523162842}
	structureydists[22] = statistics.NormalDistribution{Mean: 75.4000015258789, StandardDeviation: 2.8299999237060547}
	structureydists[23] = statistics.NormalDistribution{Mean: 76.4000015258789, StandardDeviation: 2.8599998950958252}
	structureydists[24] = statistics.NormalDistribution{Mean: 76.4000015258789, StandardDeviation: 2.8599998950958252}
	contentxs := []float64{-8, -7, -6, -5, -4, -3, -2, -1, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16}
	contentydists := make([]statistics.ContinuousDistribution, 25)
	contentydists[0] = statistics.NormalDistribution{Mean: 0, StandardDeviation: 0}
	contentydists[1] = statistics.NormalDistribution{Mean: 1, StandardDeviation: 0.0099999997764825821}
	contentydists[2] = statistics.NormalDistribution{Mean: 2.2999999523162842, StandardDeviation: 0.10000000149011612}
	contentydists[3] = statistics.NormalDistribution{Mean: 3.7000000476837158, StandardDeviation: 0.20000000298023224}
	contentydists[4] = statistics.NormalDistribution{Mean: 5.1999998092651367, StandardDeviation: 0.34999999403953552}
	contentydists[5] = statistics.NormalDistribution{Mean: 6.8000001907348633, StandardDeviation: 0.5}
	contentydists[6] = statistics.NormalDistribution{Mean: 8.3999996185302734, StandardDeviation: 0.699999988079071}
	contentydists[7] = statistics.NormalDistribution{Mean: 10.100000381469727, StandardDeviation: 0.86000001430511475}
	contentydists[8] = statistics.NormalDistribution{Mean: 11.899999618530273, StandardDeviation: 1}
	contentydists[9] = statistics.NormalDistribution{Mean: 13.800000190734863, StandardDeviation: 1.1100000143051148}
	contentydists[10] = statistics.NormalDistribution{Mean: 15.699999809265137, StandardDeviation: 1.2300000190734863}
	contentydists[11] = statistics.NormalDistribution{Mean: 17.700000762939453, StandardDeviation: 1.4299999475479126}
	contentydists[12] = statistics.NormalDistribution{Mean: 19.799999237060547, StandardDeviation: 1.6699999570846558}
	contentydists[13] = statistics.NormalDistribution{Mean: 22, StandardDeviation: 1.9199999570846558}
	contentydists[14] = statistics.NormalDistribution{Mean: 24.299999237060547, StandardDeviation: 2.1500000953674316}
	contentydists[15] = statistics.NormalDistribution{Mean: 26.700000762939453, StandardDeviation: 2.3599998950958252}
	contentydists[16] = statistics.NormalDistribution{Mean: 29.100000381469727, StandardDeviation: 2.559999942779541}
	contentydists[17] = statistics.NormalDistribution{Mean: 31.700000762939453, StandardDeviation: 2.7599999904632568}
	contentydists[18] = statistics.NormalDistribution{Mean: 34.400001525878906, StandardDeviation: 3.0399999618530273}
	contentydists[19] = statistics.NormalDistribution{Mean: 37.200000762939453, StandardDeviation: 3.2999999523162842}
	contentydists[20] = statistics.NormalDistribution{Mean: 40, StandardDeviation: 3.5999999046325684}
	contentydists[21] = statistics.NormalDistribution{Mean: 43, StandardDeviation: 3.9000000953674316}
	contentydists[22] = statistics.NormalDistribution{Mean: 46.099998474121094, StandardDeviation: 4.2600002288818359}
	contentydists[23] = statistics.NormalDistribution{Mean: 49.299999237060547, StandardDeviation: 4.5999999046325684}
	contentydists[24] = statistics.NormalDistribution{Mean: 52.599998474121094, StandardDeviation: 5}
	var structuredamagefunctionStochastic = paireddata.UncertaintyPairedData{Xvals: structurexs, Yvals: structureydists}
	var contentdamagefunctionStochastic = paireddata.UncertaintyPairedData{Xvals: contentxs, Yvals: contentydists}
	return OccupancyTypeStochastic{Name: "RES1-2SWB", Structuredamfun: structuredamagefunctionStochastic, Contentdamfun: contentdamagefunctionStochastic}
}
func res13snb() OccupancyTypeStochastic {
	structurexs := []float64{-2, -1, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16}
	structureydists := make([]statistics.ContinuousDistribution, 19)
	structureydists[0] = statistics.NormalDistribution{Mean: 0, StandardDeviation: 0}
	structureydists[1] = statistics.NormalDistribution{Mean: 3, StandardDeviation: 0.30000001192092896}
	structureydists[2] = statistics.NormalDistribution{Mean: 9.3000001907348633, StandardDeviation: 1}
	structureydists[3] = statistics.NormalDistribution{Mean: 15.199999809265137, StandardDeviation: 1.5}
	structureydists[4] = statistics.NormalDistribution{Mean: 20.899999618530273, StandardDeviation: 2.0999999046325684}
	structureydists[5] = statistics.NormalDistribution{Mean: 26.299999237060547, StandardDeviation: 2.5999999046325684}
	structureydists[6] = statistics.NormalDistribution{Mean: 31.399999618530273, StandardDeviation: 3}
	structureydists[7] = statistics.NormalDistribution{Mean: 36.200000762939453, StandardDeviation: 3.2000000476837158}
	structureydists[8] = statistics.NormalDistribution{Mean: 40.700000762939453, StandardDeviation: 3.5}
	structureydists[9] = statistics.NormalDistribution{Mean: 44.900001525878906, StandardDeviation: 3.5499999523162842}
	structureydists[10] = statistics.NormalDistribution{Mean: 48.799999237060547, StandardDeviation: 3.5999999046325684}
	structureydists[11] = statistics.NormalDistribution{Mean: 52.400001525878906, StandardDeviation: 3.6500000953674316}
	structureydists[12] = statistics.NormalDistribution{Mean: 55.700000762939453, StandardDeviation: 3.7000000476837158}
	structureydists[13] = statistics.NormalDistribution{Mean: 58.700000762939453, StandardDeviation: 3.7300000190734863}
	structureydists[14] = statistics.NormalDistribution{Mean: 61.400001525878906, StandardDeviation: 3.7699999809265137}
	structureydists[15] = statistics.NormalDistribution{Mean: 63.799999237060547, StandardDeviation: 3.7799999713897705}
	structureydists[16] = statistics.NormalDistribution{Mean: 65.9000015258789, StandardDeviation: 3.7899999618530273}
	structureydists[17] = statistics.NormalDistribution{Mean: 67.699996948242188, StandardDeviation: 3.7999999523162842}
	structureydists[18] = statistics.NormalDistribution{Mean: 69.199996948242188, StandardDeviation: 3.7999999523162842}
	contentxs := []float64{-2, -1, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16}
	contentydists := make([]statistics.ContinuousDistribution, 19)
	contentydists[0] = statistics.NormalDistribution{Mean: 0, StandardDeviation: 0}
	contentydists[1] = statistics.NormalDistribution{Mean: 1, StandardDeviation: 0.05000000074505806}
	contentydists[2] = statistics.NormalDistribution{Mean: 5, StandardDeviation: 0.5}
	contentydists[3] = statistics.NormalDistribution{Mean: 8.6999998092651367, StandardDeviation: 0.89999997615814209}
	contentydists[4] = statistics.NormalDistribution{Mean: 12.199999809265137, StandardDeviation: 1.2999999523162842}
	contentydists[5] = statistics.NormalDistribution{Mean: 15.5, StandardDeviation: 1.7000000476837158}
	contentydists[6] = statistics.NormalDistribution{Mean: 18.5, StandardDeviation: 2}
	contentydists[7] = statistics.NormalDistribution{Mean: 21.299999237060547, StandardDeviation: 2.2999999523162842}
	contentydists[8] = statistics.NormalDistribution{Mean: 23.899999618530273, StandardDeviation: 2.5999999046325684}
	contentydists[9] = statistics.NormalDistribution{Mean: 26.299999237060547, StandardDeviation: 2.9000000953674316}
	contentydists[10] = statistics.NormalDistribution{Mean: 28.399999618530273, StandardDeviation: 3.0999999046325684}
	contentydists[11] = statistics.NormalDistribution{Mean: 30.299999237060547, StandardDeviation: 3.2999999523162842}
	contentydists[12] = statistics.NormalDistribution{Mean: 32, StandardDeviation: 3.5}
	contentydists[13] = statistics.NormalDistribution{Mean: 33.400001525878906, StandardDeviation: 3.5999999046325684}
	contentydists[14] = statistics.NormalDistribution{Mean: 34.700000762939453, StandardDeviation: 3.7000000476837158}
	contentydists[15] = statistics.NormalDistribution{Mean: 35.599998474121094, StandardDeviation: 3.7999999523162842}
	contentydists[16] = statistics.NormalDistribution{Mean: 36.400001525878906, StandardDeviation: 3.9000000953674316}
	contentydists[17] = statistics.NormalDistribution{Mean: 36.900001525878906, StandardDeviation: 3.9000000953674316}
	contentydists[18] = statistics.NormalDistribution{Mean: 37.200000762939453, StandardDeviation: 3.9000000953674316}
	var structuredamagefunctionStochastic = paireddata.UncertaintyPairedData{Xvals: structurexs, Yvals: structureydists}
	var contentdamagefunctionStochastic = paireddata.UncertaintyPairedData{Xvals: contentxs, Yvals: contentydists}
	return OccupancyTypeStochastic{Name: "RES1-3SNB", Structuredamfun: structuredamagefunctionStochastic, Contentdamfun: contentdamagefunctionStochastic}
}
func res13swb() OccupancyTypeStochastic {
	structurexs := []float64{-8, -7, -6, -5, -4, -3, -2, -1, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16}
	structureydists := make([]statistics.ContinuousDistribution, 25)
	structureydists[0] = statistics.NormalDistribution{Mean: 1.7000000476837158, StandardDeviation: 0.0099999997764825821}
	structureydists[1] = statistics.NormalDistribution{Mean: 1.7000000476837158, StandardDeviation: 0.0099999997764825821}
	structureydists[2] = statistics.NormalDistribution{Mean: 1.8999999761581421, StandardDeviation: 0.029999999329447746}
	structureydists[3] = statistics.NormalDistribution{Mean: 2.9000000953674316, StandardDeviation: 0.05000000074505806}
	structureydists[4] = statistics.NormalDistribution{Mean: 4.6999998092651367, StandardDeviation: 0.20000000298023224}
	structureydists[5] = statistics.NormalDistribution{Mean: 7.1999998092651367, StandardDeviation: 0.5}
	structureydists[6] = statistics.NormalDistribution{Mean: 10.199999809265137, StandardDeviation: 0.699999988079071}
	structureydists[7] = statistics.NormalDistribution{Mean: 13.899999618530273, StandardDeviation: 1}
	structureydists[8] = statistics.NormalDistribution{Mean: 17.899999618530273, StandardDeviation: 1.3200000524520874}
	structureydists[9] = statistics.NormalDistribution{Mean: 22.299999237060547, StandardDeviation: 1.3500000238418579}
	structureydists[10] = statistics.NormalDistribution{Mean: 27, StandardDeviation: 1.5}
	structureydists[11] = statistics.NormalDistribution{Mean: 31.899999618530273, StandardDeviation: 1.75}
	structureydists[12] = statistics.NormalDistribution{Mean: 36.900001525878906, StandardDeviation: 2.0399999618530273}
	structureydists[13] = statistics.NormalDistribution{Mean: 41.900001525878906, StandardDeviation: 2.3399999141693115}
	structureydists[14] = statistics.NormalDistribution{Mean: 46.900001525878906, StandardDeviation: 2.630000114440918}
	structureydists[15] = statistics.NormalDistribution{Mean: 51.799999237060547, StandardDeviation: 2.8900001049041748}
	structureydists[16] = statistics.NormalDistribution{Mean: 56.400001525878906, StandardDeviation: 3.130000114440918}
	structureydists[17] = statistics.NormalDistribution{Mean: 60.799999237060547, StandardDeviation: 3.2000000476837158}
	structureydists[18] = statistics.NormalDistribution{Mean: 64.800003051757812, StandardDeviation: 3.2300000190734863}
	structureydists[19] = statistics.NormalDistribution{Mean: 68.4000015258789, StandardDeviation: 3.2699999809265137}
	structureydists[20] = statistics.NormalDistribution{Mean: 71.4000015258789, StandardDeviation: 3.2999999523162842}
	structureydists[21] = statistics.NormalDistribution{Mean: 73.699996948242188, StandardDeviation: 3.0999999046325684}
	structureydists[22] = statistics.NormalDistribution{Mean: 75.4000015258789, StandardDeviation: 3}
	structureydists[23] = statistics.NormalDistribution{Mean: 76.4000015258789, StandardDeviation: 2.9000000953674316}
	structureydists[24] = statistics.NormalDistribution{Mean: 76.4000015258789, StandardDeviation: 2.9000000953674316}
	contentxs := []float64{-8, -7, -6, -5, -4, -3, -2, -1, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16}
	contentydists := make([]statistics.ContinuousDistribution, 25)
	contentydists[0] = statistics.NormalDistribution{Mean: 0, StandardDeviation: 0}
	contentydists[1] = statistics.NormalDistribution{Mean: 1, StandardDeviation: 0.019999999552965164}
	contentydists[2] = statistics.NormalDistribution{Mean: 2.2999999523162842, StandardDeviation: 0.15000000596046448}
	contentydists[3] = statistics.NormalDistribution{Mean: 3.7000000476837158, StandardDeviation: 0.30000001192092896}
	contentydists[4] = statistics.NormalDistribution{Mean: 5.1999998092651367, StandardDeviation: 0.44999998807907104}
	contentydists[5] = statistics.NormalDistribution{Mean: 6.8000001907348633, StandardDeviation: 0.60000002384185791}
	contentydists[6] = statistics.NormalDistribution{Mean: 8.3999996185302734, StandardDeviation: 0.800000011920929}
	contentydists[7] = statistics.NormalDistribution{Mean: 10.100000381469727, StandardDeviation: 1}
	contentydists[8] = statistics.NormalDistribution{Mean: 11.899999618530273, StandardDeviation: 1.0900000333786011}
	contentydists[9] = statistics.NormalDistribution{Mean: 13.800000190734863, StandardDeviation: 1.1100000143051148}
	contentydists[10] = statistics.NormalDistribution{Mean: 15.699999809265137, StandardDeviation: 1.2300000190734863}
	contentydists[11] = statistics.NormalDistribution{Mean: 17.700000762939453, StandardDeviation: 1.4299999475479126}
	contentydists[12] = statistics.NormalDistribution{Mean: 19.799999237060547, StandardDeviation: 1.6699999570846558}
	contentydists[13] = statistics.NormalDistribution{Mean: 22, StandardDeviation: 1.9199999570846558}
	contentydists[14] = statistics.NormalDistribution{Mean: 24.299999237060547, StandardDeviation: 2.1500000953674316}
	contentydists[15] = statistics.NormalDistribution{Mean: 26.700000762939453, StandardDeviation: 2.3599998950958252}
	contentydists[16] = statistics.NormalDistribution{Mean: 29.100000381469727, StandardDeviation: 2.559999942779541}
	contentydists[17] = statistics.NormalDistribution{Mean: 31.700000762939453, StandardDeviation: 2.7599999904632568}
	contentydists[18] = statistics.NormalDistribution{Mean: 34.400001525878906, StandardDeviation: 3.0399999618530273}
	contentydists[19] = statistics.NormalDistribution{Mean: 37.200000762939453, StandardDeviation: 3.3499999046325684}
	contentydists[20] = statistics.NormalDistribution{Mean: 40, StandardDeviation: 3.7000000476837158}
	contentydists[21] = statistics.NormalDistribution{Mean: 43, StandardDeviation: 4}
	contentydists[22] = statistics.NormalDistribution{Mean: 46.099998474121094, StandardDeviation: 4.3000001907348633}
	contentydists[23] = statistics.NormalDistribution{Mean: 49.299999237060547, StandardDeviation: 4.5999999046325684}
	contentydists[24] = statistics.NormalDistribution{Mean: 52.599998474121094, StandardDeviation: 5}
	var structuredamagefunctionStochastic = paireddata.UncertaintyPairedData{Xvals: structurexs, Yvals: structureydists}
	var contentdamagefunctionStochastic = paireddata.UncertaintyPairedData{Xvals: contentxs, Yvals: contentydists}
	return OccupancyTypeStochastic{Name: "RES1-3SWB", Structuredamfun: structuredamagefunctionStochastic, Contentdamfun: contentdamagefunctionStochastic}
}
func res1slnb() OccupancyTypeStochastic {
	structurexs := []float64{-2, -1, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16}
	structureydists := make([]statistics.ContinuousDistribution, 19)
	structureydists[0] = statistics.NormalDistribution{Mean: 0, StandardDeviation: 0}
	structureydists[1] = statistics.NormalDistribution{Mean: 6.4000000953674316, StandardDeviation: 0}
	structureydists[2] = statistics.NormalDistribution{Mean: 7.1999998092651367, StandardDeviation: 0.039999999105930328}
	structureydists[3] = statistics.NormalDistribution{Mean: 9.3999996185302734, StandardDeviation: 0.20000000298023224}
	structureydists[4] = statistics.NormalDistribution{Mean: 12.899999618530273, StandardDeviation: 0.5}
	structureydists[5] = statistics.NormalDistribution{Mean: 17.399999618530273, StandardDeviation: 1}
	structureydists[6] = statistics.NormalDistribution{Mean: 22.799999237060547, StandardDeviation: 1.5}
	structureydists[7] = statistics.NormalDistribution{Mean: 28.899999618530273, StandardDeviation: 2}
	structureydists[8] = statistics.NormalDistribution{Mean: 35.5, StandardDeviation: 2.7000000476837158}
	structureydists[9] = statistics.NormalDistribution{Mean: 42.299999237060547, StandardDeviation: 3.2000000476837158}
	structureydists[10] = statistics.NormalDistribution{Mean: 49.200000762939453, StandardDeviation: 3.5}
	structureydists[11] = statistics.NormalDistribution{Mean: 56.099998474121094, StandardDeviation: 3.7999999523162842}
	structureydists[12] = statistics.NormalDistribution{Mean: 62.599998474121094, StandardDeviation: 4}
	structureydists[13] = statistics.NormalDistribution{Mean: 68.5999984741211, StandardDeviation: 3.5}
	structureydists[14] = statistics.NormalDistribution{Mean: 73.9000015258789, StandardDeviation: 3}
	structureydists[15] = statistics.NormalDistribution{Mean: 78.4000015258789, StandardDeviation: 2.5}
	structureydists[16] = statistics.NormalDistribution{Mean: 81.699996948242188, StandardDeviation: 2.0999999046325684}
	structureydists[17] = statistics.NormalDistribution{Mean: 83.800003051757812, StandardDeviation: 1.8999999761581421}
	structureydists[18] = statistics.NormalDistribution{Mean: 84.4000015258789, StandardDeviation: 1.8999999761581421}
	contentxs := []float64{-2, -1, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16}
	contentydists := make([]statistics.ContinuousDistribution, 19)
	contentydists[0] = statistics.NormalDistribution{Mean: 0, StandardDeviation: 0}
	contentydists[1] = statistics.NormalDistribution{Mean: 2.2000000476837158, StandardDeviation: 0}
	contentydists[2] = statistics.NormalDistribution{Mean: 2.9000000953674316, StandardDeviation: 0.059999998658895493}
	contentydists[3] = statistics.NormalDistribution{Mean: 4.6999998092651367, StandardDeviation: 0.25}
	contentydists[4] = statistics.NormalDistribution{Mean: 7.5, StandardDeviation: 0.60000002384185791}
	contentydists[5] = statistics.NormalDistribution{Mean: 11.100000381469727, StandardDeviation: 1}
	contentydists[6] = statistics.NormalDistribution{Mean: 15.300000190734863, StandardDeviation: 1.5}
	contentydists[7] = statistics.NormalDistribution{Mean: 20.100000381469727, StandardDeviation: 1.6000000238418579}
	contentydists[8] = statistics.NormalDistribution{Mean: 25.200000762939453, StandardDeviation: 1.7999999523162842}
	contentydists[9] = statistics.NormalDistribution{Mean: 30.5, StandardDeviation: 2.0999999046325684}
	contentydists[10] = statistics.NormalDistribution{Mean: 35.700000762939453, StandardDeviation: 2.5}
	contentydists[11] = statistics.NormalDistribution{Mean: 40.900001525878906, StandardDeviation: 3}
	contentydists[12] = statistics.NormalDistribution{Mean: 45.799999237060547, StandardDeviation: 3.5}
	contentydists[13] = statistics.NormalDistribution{Mean: 50.200000762939453, StandardDeviation: 4}
	contentydists[14] = statistics.NormalDistribution{Mean: 54.099998474121094, StandardDeviation: 4.4000000953674316}
	contentydists[15] = statistics.NormalDistribution{Mean: 57.200000762939453, StandardDeviation: 4.75}
	contentydists[16] = statistics.NormalDistribution{Mean: 59.400001525878906, StandardDeviation: 4.8000001907348633}
	contentydists[17] = statistics.NormalDistribution{Mean: 60.5, StandardDeviation: 4.8000001907348633}
	contentydists[18] = statistics.NormalDistribution{Mean: 60.5, StandardDeviation: 4.8000001907348633}
	var structuredamagefunctionStochastic = paireddata.UncertaintyPairedData{Xvals: structurexs, Yvals: structureydists}
	var contentdamagefunctionStochastic = paireddata.UncertaintyPairedData{Xvals: contentxs, Yvals: contentydists}
	return OccupancyTypeStochastic{Name: "RES1-SLNB", Structuredamfun: structuredamagefunctionStochastic, Contentdamfun: contentdamagefunctionStochastic}
}
func res1slwb() OccupancyTypeStochastic {
	structurexs := []float64{-8, -7, -6, -5, -4, -3, -2, -1, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16}
	structureydists := make([]statistics.ContinuousDistribution, 25)
	structureydists[0] = statistics.NormalDistribution{Mean: 0, StandardDeviation: 0}
	structureydists[1] = statistics.NormalDistribution{Mean: 0, StandardDeviation: 0}
	structureydists[2] = statistics.NormalDistribution{Mean: 2.5, StandardDeviation: 0.30000001192092896}
	structureydists[3] = statistics.NormalDistribution{Mean: 3.0999999046325684, StandardDeviation: 0.30000001192092896}
	structureydists[4] = statistics.NormalDistribution{Mean: 4.6999998092651367, StandardDeviation: 0.5}
	structureydists[5] = statistics.NormalDistribution{Mean: 7.1999998092651367, StandardDeviation: 0.699999988079071}
	structureydists[6] = statistics.NormalDistribution{Mean: 10.399999618530273, StandardDeviation: 1}
	structureydists[7] = statistics.NormalDistribution{Mean: 14.199999809265137, StandardDeviation: 1.2000000476837158}
	structureydists[8] = statistics.NormalDistribution{Mean: 18.5, StandardDeviation: 1.6000000238418579}
	structureydists[9] = statistics.NormalDistribution{Mean: 23.200000762939453, StandardDeviation: 1.7000000476837158}
	structureydists[10] = statistics.NormalDistribution{Mean: 28.200000762939453, StandardDeviation: 1.8999999761581421}
	structureydists[11] = statistics.NormalDistribution{Mean: 33.400001525878906, StandardDeviation: 2.0999999046325684}
	structureydists[12] = statistics.NormalDistribution{Mean: 38.599998474121094, StandardDeviation: 2.4000000953674316}
	structureydists[13] = statistics.NormalDistribution{Mean: 43.799999237060547, StandardDeviation: 2.5999999046325684}
	structureydists[14] = statistics.NormalDistribution{Mean: 48.799999237060547, StandardDeviation: 2.9000000953674316}
	structureydists[15] = statistics.NormalDistribution{Mean: 53.5, StandardDeviation: 3.2000000476837158}
	structureydists[16] = statistics.NormalDistribution{Mean: 57.799999237060547, StandardDeviation: 3.2999999523162842}
	structureydists[17] = statistics.NormalDistribution{Mean: 61.599998474121094, StandardDeviation: 3.4000000953674316}
	structureydists[18] = statistics.NormalDistribution{Mean: 64.800003051757812, StandardDeviation: 3.4500000476837158}
	structureydists[19] = statistics.NormalDistribution{Mean: 67.199996948242188, StandardDeviation: 3.5}
	structureydists[20] = statistics.NormalDistribution{Mean: 68.800003051757812, StandardDeviation: 3.5699999332427979}
	structureydists[21] = statistics.NormalDistribution{Mean: 69.300003051757812, StandardDeviation: 3.619999885559082}
	structureydists[22] = statistics.NormalDistribution{Mean: 69.300003051757812, StandardDeviation: 3.619999885559082}
	structureydists[23] = statistics.NormalDistribution{Mean: 69.300003051757812, StandardDeviation: 3.619999885559082}
	structureydists[24] = statistics.NormalDistribution{Mean: 69.300003051757812, StandardDeviation: 3.619999885559082}
	contentxs := []float64{-8, -7, -6, -5, -4, -3, -2, -1, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16}
	contentydists := make([]statistics.ContinuousDistribution, 25)
	contentydists[0] = statistics.NormalDistribution{Mean: 0.60000002384185791, StandardDeviation: 0}
	contentydists[1] = statistics.NormalDistribution{Mean: 0.699999988079071, StandardDeviation: 0}
	contentydists[2] = statistics.NormalDistribution{Mean: 1.3999999761581421, StandardDeviation: 0.05000000074505806}
	contentydists[3] = statistics.NormalDistribution{Mean: 2.4000000953674316, StandardDeviation: 0.15000000596046448}
	contentydists[4] = statistics.NormalDistribution{Mean: 3.7999999523162842, StandardDeviation: 0.30000001192092896}
	contentydists[5] = statistics.NormalDistribution{Mean: 5.4000000953674316, StandardDeviation: 0.5}
	contentydists[6] = statistics.NormalDistribution{Mean: 7.3000001907348633, StandardDeviation: 0.699999988079071}
	contentydists[7] = statistics.NormalDistribution{Mean: 9.3999996185302734, StandardDeviation: 0.89999997615814209}
	contentydists[8] = statistics.NormalDistribution{Mean: 11.600000381469727, StandardDeviation: 1.059999942779541}
	contentydists[9] = statistics.NormalDistribution{Mean: 13.800000190734863, StandardDeviation: 1.2000000476837158}
	contentydists[10] = statistics.NormalDistribution{Mean: 16.100000381469727, StandardDeviation: 1.3999999761581421}
	contentydists[11] = statistics.NormalDistribution{Mean: 18.200000762939453, StandardDeviation: 1.6000000238418579}
	contentydists[12] = statistics.NormalDistribution{Mean: 20.200000762939453, StandardDeviation: 1.7999999523162842}
	contentydists[13] = statistics.NormalDistribution{Mean: 22.100000381469727, StandardDeviation: 2}
	contentydists[14] = statistics.NormalDistribution{Mean: 23.600000381469727, StandardDeviation: 2.1800000667572021}
	contentydists[15] = statistics.NormalDistribution{Mean: 24.899999618530273, StandardDeviation: 2.2999999523162842}
	contentydists[16] = statistics.NormalDistribution{Mean: 25.799999237060547, StandardDeviation: 2.4000000953674316}
	contentydists[17] = statistics.NormalDistribution{Mean: 26.299999237060547, StandardDeviation: 2.440000057220459}
	contentydists[18] = statistics.NormalDistribution{Mean: 26.299999237060547, StandardDeviation: 2.440000057220459}
	contentydists[19] = statistics.NormalDistribution{Mean: 26.299999237060547, StandardDeviation: 2.440000057220459}
	contentydists[20] = statistics.NormalDistribution{Mean: 26.299999237060547, StandardDeviation: 2.440000057220459}
	contentydists[21] = statistics.NormalDistribution{Mean: 26.299999237060547, StandardDeviation: 2.440000057220459}
	contentydists[22] = statistics.NormalDistribution{Mean: 26.299999237060547, StandardDeviation: 2.440000057220459}
	contentydists[23] = statistics.NormalDistribution{Mean: 26.299999237060547, StandardDeviation: 2.440000057220459}
	contentydists[24] = statistics.NormalDistribution{Mean: 26.299999237060547, StandardDeviation: 2.440000057220459}
	var structuredamagefunctionStochastic = paireddata.UncertaintyPairedData{Xvals: structurexs, Yvals: structureydists}
	var contentdamagefunctionStochastic = paireddata.UncertaintyPairedData{Xvals: contentxs, Yvals: contentydists}
	return OccupancyTypeStochastic{Name: "RES1-SLWB", Structuredamfun: structuredamagefunctionStochastic, Contentdamfun: contentdamagefunctionStochastic}
}

func res2() OccupancyTypeStochastic {
	structurexs := []float64{-4, -3, -2, -1, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24}
	structureys := []float64{0, 0, 0, 0, 11, 44, 63, 73, 78, 79, 81, 82, 83, 84, 85, 86, 88, 89, 90, 91, 92, 94, 95, 96, 97, 98, 99, 100, 100}
	contentxs := []float64{-4, -3, -2, -1, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24}
	contentys := []float64{0, 0, 0, 0, 3, 27, 49, 64, 70, 76, 78, 79, 81, 83, 83, 83, 83, 83, 83, 83, 83, 83, 83, 83, 83, 83, 83, 83, 83}
	var structuredamagefunction = paireddata.PairedData{Xvals: structurexs, Yvals: structureys}
	var contentdamagefunction = paireddata.PairedData{Xvals: contentxs, Yvals: contentys}
	return OccupancyTypeStochastic{Name: "RES2", Structuredamfun: structuredamagefunction, Contentdamfun: contentdamagefunction}
}
func res3a() OccupancyTypeStochastic {
	structurexs := []float64{-4, -3, -2, -1, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24}
	structureys := []float64{0, 0, 0, 0, 5, 28, 29, 31, 36, 37, 39, 40, 41, 42, 44, 46, 48, 52, 55, 58, 61, 64, 68, 69, 70, 71, 72, 73, 74}
	contentxs := []float64{-4, -3, -2, -1, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24}
	contentys := []float64{0, 0, 0, 0, 4, 24, 34, 40, 47, 53, 56, 58, 58, 58, 61, 66, 68, 76, 81, 86, 91, 100, 100, 100, 100, 100, 100, 100, 100}
	var structuredamagefunction = paireddata.PairedData{Xvals: structurexs, Yvals: structureys}
	var contentdamagefunction = paireddata.PairedData{Xvals: contentxs, Yvals: contentys}
	return OccupancyTypeStochastic{Name: "RES3A", Structuredamfun: structuredamagefunction, Contentdamfun: contentdamagefunction}
}
func res3b() OccupancyTypeStochastic {
	structurexs := []float64{-4, -3, -2, -1, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24}
	structureys := []float64{0, 0, 0, 0, 5, 28, 29, 31, 36, 37, 39, 40, 41, 42, 44, 46, 48, 52, 55, 58, 61, 64, 68, 69, 70, 71, 72, 73, 74}
	contentxs := []float64{-4, -3, -2, -1, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24}
	contentys := []float64{0, 0, 0, 0, 4, 24, 34, 40, 47, 53, 56, 58, 58, 58, 61, 66, 68, 76, 81, 86, 91, 100, 100, 100, 100, 100, 100, 100, 100}
	var structuredamagefunction = paireddata.PairedData{Xvals: structurexs, Yvals: structureys}
	var contentdamagefunction = paireddata.PairedData{Xvals: contentxs, Yvals: contentys}
	return OccupancyTypeStochastic{Name: "RES3B", Structuredamfun: structuredamagefunction, Contentdamfun: contentdamagefunction}
}
func res3c() OccupancyTypeStochastic {
	structurexs := []float64{-4, -3, -2, -1, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24}
	structureys := []float64{0, 0, 0, 0, 5, 28, 29, 31, 36, 37, 39, 40, 41, 42, 44, 46, 48, 52, 55, 58, 61, 64, 68, 69, 70, 71, 72, 73, 74}
	contentxs := []float64{-4, -3, -2, -1, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24}
	contentys := []float64{0, 0, 0, 0, 4, 24, 34, 40, 47, 53, 56, 58, 58, 58, 61, 66, 68, 76, 81, 86, 91, 100, 100, 100, 100, 100, 100, 100, 100}
	var structuredamagefunction = paireddata.PairedData{Xvals: structurexs, Yvals: structureys}
	var contentdamagefunction = paireddata.PairedData{Xvals: contentxs, Yvals: contentys}
	return OccupancyTypeStochastic{Name: "RES3C", Structuredamfun: structuredamagefunction, Contentdamfun: contentdamagefunction}
}
func res3d() OccupancyTypeStochastic {
	structurexs := []float64{-4, -3, -2, -1, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24}
	structureys := []float64{0, 0, 0, 0, 5, 28, 29, 31, 36, 37, 39, 40, 41, 42, 44, 46, 48, 52, 55, 58, 61, 64, 68, 69, 70, 71, 72, 73, 74}
	contentxs := []float64{-4, -3, -2, -1, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24}
	contentys := []float64{0, 0, 0, 0, 4, 24, 34, 40, 47, 53, 56, 58, 58, 58, 61, 66, 68, 76, 81, 86, 91, 100, 100, 100, 100, 100, 100, 100, 100}
	var structuredamagefunction = paireddata.PairedData{Xvals: structurexs, Yvals: structureys}
	var contentdamagefunction = paireddata.PairedData{Xvals: contentxs, Yvals: contentys}
	return OccupancyTypeStochastic{Name: "RES3D", Structuredamfun: structuredamagefunction, Contentdamfun: contentdamagefunction}
}
func res3e() OccupancyTypeStochastic {
	structurexs := []float64{-4, -3, -2, -1, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24}
	structureys := []float64{0, 0, 0, 0, 5, 28, 29, 31, 36, 37, 39, 40, 41, 42, 44, 46, 48, 52, 55, 58, 61, 64, 68, 69, 70, 71, 72, 73, 74}
	contentxs := []float64{-4, -3, -2, -1, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24}
	contentys := []float64{0, 0, 0, 0, 4, 24, 34, 40, 47, 53, 56, 58, 58, 58, 61, 66, 68, 76, 81, 86, 91, 100, 100, 100, 100, 100, 100, 100, 100}
	var structuredamagefunction = paireddata.PairedData{Xvals: structurexs, Yvals: structureys}
	var contentdamagefunction = paireddata.PairedData{Xvals: contentxs, Yvals: contentys}
	return OccupancyTypeStochastic{Name: "RES3E", Structuredamfun: structuredamagefunction, Contentdamfun: contentdamagefunction}
}
func res3f() OccupancyTypeStochastic {
	structurexs := []float64{-4, -3, -2, -1, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24}
	structureys := []float64{0, 0, 0, 0, 5, 28, 29, 31, 36, 37, 39, 40, 41, 42, 44, 46, 48, 52, 55, 58, 61, 64, 68, 69, 70, 71, 72, 73, 74}
	contentxs := []float64{-4, -3, -2, -1, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24}
	contentys := []float64{0, 0, 0, 0, 4, 24, 34, 40, 47, 53, 56, 58, 58, 58, 61, 66, 68, 76, 81, 86, 91, 100, 100, 100, 100, 100, 100, 100, 100}
	var structuredamagefunction = paireddata.PairedData{Xvals: structurexs, Yvals: structureys}
	var contentdamagefunction = paireddata.PairedData{Xvals: contentxs, Yvals: contentys}
	return OccupancyTypeStochastic{Name: "RES3F", Structuredamfun: structuredamagefunction, Contentdamfun: contentdamagefunction}
}
func res4() OccupancyTypeStochastic {
	structurexs := []float64{-4, -3, -2, -1, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24}
	structureys := []float64{0, 0, 0, 0, 0, 3, 5, 6, 7, 9, 12, 14, 18, 21, 26, 31, 36, 41, 46, 50, 54, 58, 62, 66, 70, 74, 78, 82, 86}
	contentxs := []float64{-4, -3, -2, -1, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24}
	contentys := []float64{0, 0, 0, 0, 0, 11, 19, 25, 29, 34, 39, 44, 49, 56, 65, 74, 82, 88, 98, 98, 98, 98, 98, 98, 98, 98, 98, 98, 98}
	var structuredamagefunction = paireddata.PairedData{Xvals: structurexs, Yvals: structureys}
	var contentdamagefunction = paireddata.PairedData{Xvals: contentxs, Yvals: contentys}
	return OccupancyTypeStochastic{Name: "RES4", Structuredamfun: structuredamagefunction, Contentdamfun: contentdamagefunction}
}
func res5() OccupancyTypeStochastic {
	structurexs := []float64{-4, -3, -2, -1, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24}
	structureys := []float64{0, 0, 0, 0, 0, 7, 10, 14, 15, 15, 16, 18, 20, 23, 26, 30, 34, 38, 42, 47, 52, 57, 62, 67, 72, 77, 82, 87, 92}
	contentxs := []float64{-4, -3, -2, -1, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24}
	contentys := []float64{0, 0, 0, 0, 0, 38, 60, 73, 81, 88, 94, 100, 100, 100, 100, 100, 100, 100, 100, 100, 100, 100, 100, 100, 100, 100, 100, 100, 100}
	var structuredamagefunction = paireddata.PairedData{Xvals: structurexs, Yvals: structureys}
	var contentdamagefunction = paireddata.PairedData{Xvals: contentxs, Yvals: contentys}
	return OccupancyTypeStochastic{Name: "RES5", Structuredamfun: structuredamagefunction, Contentdamfun: contentdamagefunction}
}
func res6() OccupancyTypeStochastic {
	structurexs := []float64{-4, -3, -2, -1, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24}
	structureys := []float64{0, 0, 0, 0, 0, 7, 10, 14, 15, 15, 16, 18, 20, 23, 26, 30, 34, 38, 42, 47, 52, 57, 62, 67, 72, 77, 82, 87, 92}
	contentxs := []float64{-4, -3, -2, -1, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24}
	contentys := []float64{0, 0, 0, 0, 0, 38, 60, 73, 81, 88, 94, 100, 100, 100, 100, 100, 100, 100, 100, 100, 100, 100, 100, 100, 100, 100, 100, 100, 100}
	var structuredamagefunction = paireddata.PairedData{Xvals: structurexs, Yvals: structureys}
	var contentdamagefunction = paireddata.PairedData{Xvals: contentxs, Yvals: contentys}
	return OccupancyTypeStochastic{Name: "RES6", Structuredamfun: structuredamagefunction, Contentdamfun: contentdamagefunction}
}
