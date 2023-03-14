package structureprovider

import (
	"github.com/USACE/go-consequences/consequences"
	"github.com/USACE/go-consequences/structures"
)

type Stats struct {
	TotalCount          int
	TotalValue          float64
	CountByCategory     map[string]int
	ValueByCategory     map[string]float64
	TotalResPop2AM      int32   //total population to calculate percapita losses
	WorkingResPop2AM    int32   //working population to calculate percapita losses
	NonResidentialValue float64 //for capital loss ratio for indirect economics
}

func StatsByFips(fips string, sp StructureProvider) Stats {
	cbc := make(map[string]int)
	vbc := make(map[string]float64)
	stats := Stats{CountByCategory: cbc, ValueByCategory: vbc}
	sp.ByFips(fips, func(f consequences.Receptor) {
		s, ok := f.(structures.StructureStochastic)
		if ok {
			stats.TotalCount += 1
			stats.TotalValue += s.StructVal.CentralTendency()
			stats.TotalValue += s.ContVal.CentralTendency()
			stats.NonResidentialValue += s.StructVal.CentralTendency() //residential value gets removed below
			stats.NonResidentialValue += s.ContVal.CentralTendency()
			occtype := s.OccType.Name
			assetType := ""
			switch occtype {
			case "REL1":
				assetType = "Assembly"
			case "AGR1":
				assetType = "Agriculture"
			case "EDU1", "EDU2":
				assetType = "Education"
			case "GOV1", "GOV2", "COM6": //why are hospitals encoded to governement?
				assetType = "Government"
			case "IND1", "IND2", "IND3", "IND4", "IND5", "IND6":
				assetType = "Industrial"
			case "COM1", "COM2", "COM3", "COM4", "COM5", "COM7", "COM8", "COM9", "COM10":
				assetType = "Commercial"
			default:
				assetType = "Residential"
				stats.TotalResPop2AM += s.Pop2amo65 + s.Pop2amu65
				stats.WorkingResPop2AM += s.Pop2amu65
				//remove residential value...
				stats.NonResidentialValue -= s.StructVal.CentralTendency()
				stats.NonResidentialValue -= s.ContVal.CentralTendency()
			}
			//add to maps.
			count, cok := stats.CountByCategory[assetType]
			value, _ := stats.ValueByCategory[assetType]
			if cok {
				count += 1
				value += s.StructVal.CentralTendency()
				value += s.ContVal.CentralTendency()
			} else {
				count = 1
				value = s.StructVal.CentralTendency()
				value += s.ContVal.CentralTendency()
			}
			stats.CountByCategory[assetType] = count
			stats.ValueByCategory[assetType] = value

		}
	})
	return stats
}
