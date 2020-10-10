package compute

import (
	"fmt"
	"time"

	"github.com/USACE/go-consequences/consequences"
	"github.com/USACE/go-consequences/hazards"
	"github.com/USACE/go-consequences/nsi"
)

type RequestArgs struct {
	Args interface{}
}
type FipsCodeCompute struct {
	ID         string      `json:"id"`
	FIPS       string      `json:"fips"`
	HazardArgs interface{} `json:"hazardargs"`
}
type BboxCompute struct {
	ID         string      `json:"id"`
	BBOX       string      `json:"bbox"`
	HazardArgs interface{} `json:"hazardargs"`
}
type StructureSimulation struct {
	Structures []consequences.StructureStochastic
}
type NSIStructureSimulation struct {
	RequestArgs
	StructureSimulation
}
type Computable interface {
	Compute(args RequestArgs) SimulationSummary
}

type SimulationSummaryRow struct {
	RowHeader       string  `json:"rowheader"`
	StructureCount  int64   `json:"structurecount"`
	StructureDamage float64 `json:"structuredamage"`
	ContentDamage   float64 `json:"contentdamage"`
}
type SimulationSummary struct {
	ColumnNames []string               `json:"columnnames"`
	Rows        []SimulationSummaryRow `json:"rows"`
	NSITime     time.Duration
	Computetime time.Duration
}

func (s NSIStructureSimulation) Compute(args RequestArgs) SimulationSummary {
	var depthevent = hazards.DepthEvent{Depth: 5.32}
	okd := false
	fips, okfips := args.Args.(FipsCodeCompute)
	startnsi := time.Now()
	if okfips {
		//s.Status = "Downloading NSI by fips " + fips.FIPS
		fmt.Println("Downloading NSI by fips " + fips.FIPS)
		s.Structures = nsi.GetByFips(fips.FIPS)
		depthevent, okd = fips.HazardArgs.(hazards.DepthEvent)
	} else {
		bbox, okbbox := args.Args.(BboxCompute)
		if okbbox {
			//s.Status = "Downloading NSI by bbox " + bbox.BBOX
			fmt.Println("Downloading NSI by bbox " + bbox.BBOX)
			s.Structures = nsi.GetByBbox(bbox.BBOX)
			depthevent, okd = bbox.HazardArgs.(hazards.DepthEvent)
		}
	}
	elapsedNsi := time.Since(startnsi)
	fmt.Println(fmt.Sprintf("FIPS %s retrieved %d structures from the NSI in %d", fips.FIPS, len(s.Structures), elapsedNsi))
	//s.Status = "Computing Depths"
	//depths
	fmt.Println("Computing depths for" + fips.FIPS)
	var d = hazards.DepthEvent{Depth: 5.32}
	if okd {
		d = depthevent
	}
	startcompute := time.Now()
	//ideally get from some sort of source.
	rmap := make(map[string]SimulationSummaryRow)
	//s.Status = fmt.Sprintf("Computing Damages %d of %d", 0, len(s.Structures))
	for _, str := range s.Structures {
		r := str.ComputeConsequences(d)
		if val, ok := rmap[str.DamCat]; ok {
			//fmt.Println(fmt.Sprintf("FIPS %s Computing Damages %d of %d", fips.FIPS, idx, len(s.Structures)))
			val.StructureCount += 1
			val.StructureDamage += r.Results[0].(float64) //based on convention - super risky
			val.ContentDamage += r.Results[1].(float64)   //based on convention - super risky
		} else {
			rmap[str.DamCat] = SimulationSummaryRow{RowHeader: str.DamCat, StructureCount: 1, StructureDamage: r.Results[0].(float64), ContentDamage: r.Results[1].(float64)}
		}
		//s.Status = fmt.Sprintf("Computing Damages %d of %d", i, len(s.Structures))
	}
	header := []string{"Damage Category", "Structure Count", "Total Structure Damage", "Total Content Damage"}
	rows := make([]SimulationSummaryRow, len(rmap))
	idx := 0
	for _, val := range rmap {
		fmt.Println(fmt.Sprintf("for %s, there were %d structures with %f structure damages %f content damages for damage category %s", fips.FIPS, val.StructureCount, val.StructureDamage, val.ContentDamage, val.RowHeader))
		rows[idx] = val
		idx++
	}
	elapsed := time.Since(startcompute)
	var ret = SimulationSummary{ColumnNames: header, Rows: rows, NSITime: elapsedNsi, Computetime: elapsed}
	//s.Status = "Complete"
	fmt.Println("Complete for" + fips.FIPS)
	//s.Result = ret
	return ret
}
