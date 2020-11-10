package compute

import (
	"fmt"

	"github.com/USACE/go-consequences/consequences"
	"github.com/USACE/go-consequences/hazard_providers"
	"github.com/USACE/go-consequences/hazards"
	"github.com/USACE/go-consequences/nsi"
)

func ComputeMultiEvent_NSIStream(ds hazard_providers.DataSet, fips string) {
	//rmapMap := make(map[string]SimulationSummaryRow)
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
			for _, flu := range fluvial {
				for _, y := range years {
					for _, f := range frequencies {
						rmap := make(map[string]SimulationSummaryRow)
						fe = hazard_providers.FathomEvent{Year: y, Frequency: f, Fluvial: flu}
						fq.FathomEvent = fe
						result, _ := ds.ProvideHazard(fq)
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

				}
			}

		}
	})

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
