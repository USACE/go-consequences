package compute

import (
	"fmt"
	"time"

	"github.com/USACE/go-consequences/consequences"
	"github.com/USACE/go-consequences/hazards"
	"github.com/USACE/go-consequences/nsi"
	"github.com/USACE/go-consequences/structures"
)

type RequestArgs struct {
	Args       interface{}
	Concurrent bool
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
	//Structures []consequences.StructureStochastic
}
type NSIStructureSimulation struct {
	RequestArgs
	//StructureSimulation
}
type Computable interface {
	Compute(args RequestArgs) SimulationSummary
	ComputeStream(args RequestArgs) SimulationSummary
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

func NsiFeaturetoStructure(f nsi.NsiFeature, m map[string]structures.OccupancyTypeStochastic, defaultOcctype structures.OccupancyTypeStochastic) structures.StructureStochastic {
	var occtype = defaultOcctype
	if ot, ok := m[f.Properties.Occtype]; ok {
		occtype = ot
	} else {
		occtype = defaultOcctype
		msg := "Using default " + f.Properties.Occtype + " not found"
		panic(msg)
	}
	return structures.StructureStochastic{
		OccType:   occtype,
		StructVal: consequences.ParameterValue{Value: f.Properties.StructVal},
		ContVal:   consequences.ParameterValue{Value: f.Properties.ContVal},
		FoundHt:   consequences.ParameterValue{Value: f.Properties.FoundHt},
		BaseStructure: structures.BaseStructure{
			Name:   f.Properties.Name,
			DamCat: f.Properties.DamCat,
			X:      f.Properties.X,
			Y:      f.Properties.Y,
		},
	}
}
func nsiInventorytoStructures(i nsi.NsiInventory) []structures.StructureStochastic {
	m := structures.OccupancyTypeMap()
	defaultOcctype := m["RES1-1SNB"]
	structures := make([]structures.StructureStochastic, len(i.Features))
	for idx, feature := range i.Features {
		structures[idx] = NsiFeaturetoStructure(feature, m, defaultOcctype)
	}
	return structures
}
func (s NSIStructureSimulation) Compute(args RequestArgs) SimulationSummary {
	var depthevent = hazards.DepthEvent{Depth: 5.32}
	okd := false
	fips, okfips := args.Args.(FipsCodeCompute)
	startnsi := time.Now()
	var structures []structures.StructureStochastic
	if okfips {
		//s.Status = "Downloading NSI by fips " + fips.FIPS
		fmt.Println("Downloading NSI by fips " + fips.FIPS)
		structures = nsiInventorytoStructures(nsi.GetByFips(fips.FIPS))
		depthevent, okd = fips.HazardArgs.(hazards.DepthEvent)
	} else {
		bbox, okbbox := args.Args.(BboxCompute)
		if okbbox {
			//s.Status = "Downloading NSI by bbox " + bbox.BBOX
			fmt.Println("Downloading NSI by bbox " + bbox.BBOX)
			structures = nsiInventorytoStructures(nsi.GetByBbox(bbox.BBOX))
			depthevent, okd = bbox.HazardArgs.(hazards.DepthEvent)
		}
	}
	elapsedNsi := time.Since(startnsi)
	fmt.Println(fmt.Sprintf("FIPS %s retrieved %d structures from the NSI in %d", fips.FIPS, len(structures), elapsedNsi))
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
	for _, str := range structures {
		r := str.ComputeConsequences(d)
		if val, ok := rmap[str.DamCat]; ok {
			//fmt.Println(fmt.Sprintf("FIPS %s Computing Damages %d of %d", fips.FIPS, idx, len(s.Structures)))
			val.StructureCount += 1
			val.StructureDamage += r.Results[0].(float64) //based on convention - super risky
			val.ContentDamage += r.Results[1].(float64)   //based on convention - super risky
			rmap[str.DamCat] = val
		} else {
			rmap[str.DamCat] = SimulationSummaryRow{RowHeader: str.DamCat, StructureCount: 1, StructureDamage: r.Results[0].(float64), ContentDamage: r.Results[1].(float64)}
		}
		//s.Status = fmt.Sprintf("Computing Damages %d of %d", i, len(s.Structures))
	}
	header := []string{"Damage Category", "Structure Count", "Total Structure Damage", "Total Content Damage"}
	rows := make([]SimulationSummaryRow, len(rmap))
	idx := 0
	for _, val := range rmap {
		//fmt.Println(fmt.Sprintf("for %s, there were %d structures with %f structure damages %f content damages for damage category %s", fips.FIPS, val.StructureCount, val.StructureDamage, val.ContentDamage, val.RowHeader))
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

/*
 more memory efficient version of compute
 each point is processed as it is received
 from the server
*/
func (s NSIStructureSimulation) ComputeStream(args RequestArgs) SimulationSummary {
	var depthevent = hazards.DepthEvent{Depth: 5.32}
	okd := false
	fips, okfips := args.Args.(FipsCodeCompute)
	startnsi := time.Now()
	rmap := make(map[string]SimulationSummaryRow)
	if okfips {
		fmt.Println("Downloading NSI by fips " + fips.FIPS)
		depthevent, okd = fips.HazardArgs.(hazards.DepthEvent)
		if !okd {
			depthevent = hazards.DepthEvent{Depth: 5.32}
		}
		fmt.Println("Computing depths for" + fips.FIPS)
		m := structures.OccupancyTypeMap()
		defaultOcctype := m["RES1-1SNB"]
		nsi.GetByFipsStream(fips.FIPS, func(f nsi.NsiFeature) {
			str := NsiFeaturetoStructure(f, m, defaultOcctype)
			r := str.ComputeConsequences(depthevent)
			if val, ok := rmap[str.DamCat]; ok {
				val.StructureCount += 1
				val.StructureDamage += r.Results[0].(float64) //based on convention - super risky
				val.ContentDamage += r.Results[1].(float64)   //based on convention - super risky
				rmap[str.DamCat] = val
			} else {
				rmap[str.DamCat] = SimulationSummaryRow{RowHeader: str.DamCat, StructureCount: 1, StructureDamage: r.Results[0].(float64), ContentDamage: r.Results[1].(float64)}
			}
		})
	}
	elapsedNsi := time.Since(startnsi)
	header := []string{"Damage Category", "Structure Count", "Total Structure Damage", "Total Content Damage"}
	rows := make([]SimulationSummaryRow, len(rmap))
	idx := 0
	for _, val := range rmap {
		rows[idx] = val
		idx++
	}
	elapsed := time.Since(startnsi)
	var ret = SimulationSummary{ColumnNames: header, Rows: rows, NSITime: elapsedNsi, Computetime: elapsed}
	fmt.Println("Complete for" + fips.FIPS)
	return ret
}
func ComputeEAD(damages []float64, freq []float64) float64 {
	triangle := 0.0
	square := 0.0
	x1 := 1.0 // create a triangle to the first probability space - linear interpolation is probably a problem, maybe use log linear interpolation for the triangle
	y1 := 0.0
	eadT := 0.0
	for i := 0; i < len(freq); i++ {
		xdelta := x1 - freq[i]
		square = xdelta * y1
		triangle = ((xdelta) * (damages[i] - y1)) / 2.0
		eadT += square + triangle
		x1 = freq[i]
		y1 = damages[i]
	}
	if x1 != 0.0 {
		xdelta := x1 - 0.0
		eadT += xdelta * y1 //no extrapolation, just continue damages out as if it were truth for all remaining probability.

	}
	return eadT
}
func ComputeSpecialEAD(damages []float64, freq []float64) float64 {
	//this differs from computeEAD in that it specifically does not calculate the first triangle between 1 and the first frequency to interpolate damages to zero.
	triangle := 0.0
	square := 0.0
	x1 := freq[0]
	y1 := damages[0]
	eadT := 0.0
	for i := 1; i < len(freq); i++ {
		xdelta := x1 - freq[i]
		square = xdelta * y1
		triangle = ((xdelta) * -(y1 - damages[i])) / 2.0
		eadT += square + triangle
		x1 = freq[i]
		y1 = damages[i]
	}
	if x1 != 0.0 {
		xdelta := x1 - 0.0
		eadT += xdelta * y1 //no extrapolation, just continue damages out as if it were truth for all remaining probability.

	}
	return eadT
}
