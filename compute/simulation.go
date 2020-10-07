package compute

import (
	"fmt"

	"github.com/USACE/go-consequences/consequences"
	"github.com/USACE/go-consequences/hazards"
	"github.com/USACE/go-consequences/nsi"
)

type ComputeArgs struct {
	Args interface{}
}
type FipsCode struct {
	FIPS string
}
type Bbox struct {
	BBOX string
}

type StructureSimulation struct {
	Structures []consequences.StructureStochastic
	Status     string
	Result     consequences.ConsequenceDamageResult
}
type NSIStructureSimulation struct {
	ComputeArgs
	StructureSimulation
}
type Computeable interface {
	Compute(args ComputeArgs)
	GetResults() consequences.ConsequenceDamageResult
}

type ProgressReportable interface {
	ReportProgress() string
}
type SimulationSummary struct {
	RowHeader       string
	StructureCount  int64
	StructureDamage float64
	ContentDamage   float64
}

func (s NSIStructureSimulation) ReportProgress() string {
	return s.Status
}
func (s StructureSimulation) ReportProgress() string {
	return s.Status
}

/*
func (s StructureSimulation) Compute(args ComputeArgs){
	fips, okfips := args.Args.(FipsCode)
	if okfips {
		s.Structures = nsi.GetByFips(fips.FIPS)
	}
}
*/
func (s NSIStructureSimulation) Compute(args ComputeArgs) {
	fips, okfips := args.Args.(FipsCode)
	if okfips {
		s.Status = "Downloading NSI by fips " + fips.FIPS
		s.Structures = nsi.GetByFips(fips.FIPS)
	}
	bbox, okbbox := args.Args.(Bbox)
	if okbbox {
		s.Status = "Downloading NSI by bbox " + bbox.BBOX
		s.Structures = nsi.GetByBbox(bbox.BBOX)
	}
	s.Status = "Computing Depths"
	//depths
	d := hazards.DepthEvent{Depth: 5.32}
	//ideally get from some sort of source.
	rmap := make(map[string]SimulationSummary)
	s.Status = fmt.Sprintf("Computing Damages %d of %d", 0, len(s.Structures))
	for i, str := range s.Structures {
		r := str.ComputeConsequences(d)
		if val, ok := rmap[str.DamCat]; ok {
			val.StructureCount += 1
			val.StructureDamage += r.Results[0].(float64) //based on convention - super risky
			val.ContentDamage += r.Results[1].(float64)   //based on convention - super risky
		} else {
			rmap[str.DamCat] = SimulationSummary{RowHeader: str.DamCat, StructureCount: 1, StructureDamage: r.Results[0].(float64), ContentDamage: r.Results[1].(float64)}
		}
		s.Status = fmt.Sprintf("Computing Damages %d of %d", i, len(s.Structures))
	}
	header := []string{"Damage Category", "Structure Count", "Total Structure Damage", "Total Content Damage"}
	rows := make([]SimulationSummary, len(rmap))
	results := []interface{}{rows}
	var ret = consequences.ConsequenceDamageResult{Headers: header, Results: results}
	idx := 0
	for _, val := range rmap {
		ret.Results[idx] = val
		idx++
	}
	s.Result = ret
}
