package compute

import (
	"fmt"

	"github.com/USACE/go-consequences/consequences"
	"github.com/USACE/go-consequences/hazard_providers"
	"github.com/USACE/go-consequences/hazards"
	"github.com/USACE/go-consequences/nsi"
)

func ComputeMultiEvent_NSIStream(ds hazard_providers.DataSet, fips string) {
	//rmapMap := make(map[string]map[string]SimulationRow)
	fmt.Println("Downloading NSI by fips " + fips)
	years := [2]int{2020, 2050}
	frequencies := [5]int{5, 20, 100, 250, 500}
	fluvial := [2]bool{true, false}
	nsi.GetByFipsStream(fips, func(str consequences.StructureStochastic) {
		//check to see if the structure exists for a first "default event"
		fe := hazard_providers.FathomEvent{Year: 2050, Frequency: 500, Fluvial: true}
		fq := hazard_providers.FathomQuery{Fd_id: str.Name, FathomEvent: fe}
		_, err := ds.ProvideHazard(fq)
		if err == nil {
			//structure presumably exists?
			cfdam := make([]float64, 5)
			cpdam := make([]float64, 5)
			ffdam := make([]float64, 5)
			fpdam := make([]float64, 5)
			for _, flu := range fluvial {
				for _, y := range years {
					for _, f := range frequencies {
						fe = hazard_providers.FathomEvent{Year: y, Frequency: f, Fluvial: flu}
						fq.FathomEvent = fe
						result, _ := ds.ProvideHazard(fq)
						depthevent, okd := result.(hazards.DepthEvent)
						if okd {
							if depthevent.Depth <= 0 {
								//skip
								recordDamage(flu, y, f, 0, ffdam, fpdam, cfdam, cpdam)
							} else {
								r := str.ComputeConsequences(depthevent)
								StructureDamage := r.Results[0].(float64) //based on convention - super risky
								ContentDamage := r.Results[1].(float64)   //based on convention - super risky
								recordDamage(flu, y, f, StructureDamage+ContentDamage, ffdam, fpdam, cfdam, cpdam)
							}
						}
					}

				}
			}
			//compute ead's for each of the 4 caases.
			fmt.Println(fmt.Sprintf("FD_ID: %v has EADs: %f, %f, %f, %f", str.Name, computeEAD(cfdam), computeEAD(cpdam), computeEAD(ffdam), computeEAD(fpdam)))
		}
	})

}
func frequencyIndex(frequency int) int {
	switch frequency {
	case 5:
		return 0
	case 20:
		return 1
	case 100:
		return 2
	case 250:
		return 3
	case 500:
		return 4
	default:
		return -1 //bad frequency
	}
}
func recordDamage(fluvial bool, year int, frequency int, damage float64, ffdam []float64, fpdam []float64, cfdam []float64, cpdam []float64) {
	if fluvial {
		if year == 2020 {
			cfdam[frequencyIndex(frequency)] = damage
		} else if year == 2050 {
			ffdam[frequencyIndex(frequency)] = damage
		}
	} else {
		if year == 2020 {
			cpdam[frequencyIndex(frequency)] = damage
		} else if year == 2050 {
			fpdam[frequencyIndex(frequency)] = damage
		}
	}

}
func computeEAD(damages []float64) float64 {
	frequencyWeight := [5]float64{.2, .05, .01, .004, .002} //5,20,100,250,500
	ead := 0.0
	for i, d := range damages {
		ead += d * frequencyWeight[i] //not trapazoidal.
	}
	return ead
}
func ComputeSingleEvent_NSIStream(ds hazard_providers.DataSet, fips string, fe hazard_providers.FathomEvent) {
	rmap := make(map[string]SimulationSummaryRow)
	fmt.Println("Downloading NSI by fips " + fips)
	nsi.GetByFipsStream(fips, func(str consequences.StructureStochastic) {
		fq := hazard_providers.FathomQuery{Fd_id: str.Name, FathomEvent: fe}
		result, err := ds.ProvideHazard(fq)
		if err == nil {
			//structure presumably exists?
			depthevent, okd := result.(hazards.DepthEvent)
			if okd {
				if depthevent.Depth <= 0 {
					//skip
				} else {
					r := str.ComputeConsequences(depthevent)
					if val, ok := rmap[str.DamCat]; ok {
						val.StructureCount += 1
						val.StructureDamage += r.Results[0].(float64) //based on convention - super risky
						val.ContentDamage += r.Results[1].(float64)   //based on convention - super risky
						rmap[str.DamCat] = val
					} else {
						rmap[str.DamCat] = SimulationSummaryRow{RowHeader: str.DamCat, StructureCount: 1, StructureDamage: r.Results[0].(float64), ContentDamage: r.Results[1].(float64)}
					}
				}
			}

		}

	})
	rows := make([]SimulationSummaryRow, len(rmap))
	idx := 0
	//s := "COMPLETE FOR SIMULATION" + "\n"
	for _, val := range rmap {
		fmt.Println(fmt.Sprintf("for %s, there were %d structures with %f structure damages %f content damages for damage category %s", fips, val.StructureCount, val.StructureDamage, val.ContentDamage, val.RowHeader))
		//s += fmt.Sprintf("for %s, there were %d structures with %f structure damages %f content damages for damage category %s", fips, val.StructureCount, val.StructureDamage, val.ContentDamage, val.RowHeader) + "\n"
		rows[idx] = val
		idx++
	}
	//elapsed := time.Since(startnsi)
	fmt.Println("Complete for" + fips)
}
