package main

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/USACE/go-consequences/census"
	"github.com/USACE/go-consequences/compute"
	"github.com/aws/aws-lambda-go/lambda"
)

type Config struct {
	SkipJWT       bool
	LambdaContext bool
	DBUser        string
	DBPass        string
	DBName        string
	DBHost        string
	DBSSLMode     string
}

func computeConcurrentEvent(r compute.Computable, args compute.RequestArgs) {
	f := census.StateToCountyFipsMap()
	a, ok := args.Args.(compute.FipsCodeCompute)
	if ok {
		fips := a.FIPS
		if len(fips) == 2 {
			counties, exists := f[a.FIPS]
			if exists {
				var wg sync.WaitGroup
				wg.Add(len(counties))
				var count int64
				var cdam float64
				var sdam float64
				var startTime = time.Now()
				var nsitime = time.Now()
				var computetime = time.Now()
				//header := []string{"Damage Category", "Structure Count", "Total Structure Damage", "Total Content Damage"}
				rowMap := make(map[string]compute.SimulationSummaryRow)
				if !args.Concurrent {
					rr := r.ComputeStream(args)
					for _, row := range rr.Rows {
						if val, ok := rowMap[row.RowHeader]; ok {
							//fmt.Println(fmt.Sprintf("FIPS %s Computing Damages %d of %d", fips.FIPS, idx, len(s.Structures)))
							val.StructureCount += row.StructureCount
							val.StructureDamage += row.StructureDamage
							val.ContentDamage += row.ContentDamage
							rowMap[row.RowHeader] = val
						} else {
							rowMap[row.RowHeader] = compute.SimulationSummaryRow{RowHeader: row.RowHeader, StructureCount: row.StructureCount, StructureDamage: row.StructureDamage, ContentDamage: row.ContentDamage}
						}
						count += row.StructureCount
						sdam += row.StructureDamage
						cdam += row.ContentDamage
					}
					nsitime = nsitime.Add(rr.NSITime)
					computetime = computetime.Add(rr.Computetime)
				} else {
					for _, ccc := range counties {
						go func(county string) {
							defer wg.Done()
							b := compute.FipsCodeCompute{FIPS: county, ID: a.ID, HazardArgs: a.HazardArgs}
							cargs := compute.RequestArgs{Args: b}
							rr := r.ComputeStream(cargs)
							for _, row := range rr.Rows {
								if val, ok := rowMap[row.RowHeader]; ok {
									//fmt.Println(fmt.Sprintf("FIPS %s Computing Damages %d of %d", fips.FIPS, idx, len(s.Structures)))
									val.StructureCount += row.StructureCount
									val.StructureDamage += row.StructureDamage
									val.ContentDamage += row.ContentDamage
									rowMap[row.RowHeader] = val
								} else {
									rowMap[row.RowHeader] = compute.SimulationSummaryRow{RowHeader: row.RowHeader, StructureCount: row.StructureCount, StructureDamage: row.StructureDamage, ContentDamage: row.ContentDamage}
								}
								count += row.StructureCount
								sdam += row.StructureDamage
								cdam += row.ContentDamage
							}
							nsitime = nsitime.Add(rr.NSITime)
							computetime = computetime.Add(rr.Computetime)
						}(ccc)
					}
					wg.Wait()
				}

				//}

				fmt.Println("COMPLETE FOR SIMULATION")
				elapsedNSI := startTime.Sub(nsitime)
				elapsedCompute := startTime.Sub(computetime)
				elapsedClock := time.Since(startTime)
				fmt.Println(fmt.Sprintf("NSI Took %s", -elapsedNSI))
				fmt.Println(fmt.Sprintf("Compute Took %s", -elapsedCompute))
				fmt.Println(fmt.Sprintf("Clock Time Taken was %s", elapsedClock))
				fmt.Println(fmt.Sprintf("Total Structure Count %d", count))
				fmt.Println(fmt.Sprintf("Total Structure Damage %f", sdam))
				fmt.Println(fmt.Sprintf("Total Content Damage %f", cdam))
				fmt.Println("*****************SUMMMARY*****************")
				rows := make([]compute.SimulationSummaryRow, len(rowMap))
				idx := 0
				for _, val := range rowMap {
					fmt.Println(fmt.Sprintf("for %s, there were %d structures with %f structure damages %f content damages for damage category %s", fips, val.StructureCount, val.StructureDamage, val.ContentDamage, val.RowHeader))
					rows[idx] = val
					idx++
				}
				//var ret = SimulationSummary{ColumnNames: header, Rows: rows, NSITime: elapsedNsi, Computetime: elapsed}
			} else {
				//2 characters but not a state?
				r.Compute(args) //should fail
			}
		} else {
			//not two characters
			r.Compute(args) //should work
		}
	} else {
		r.Compute(args)
	}
}
func computeEvent(r compute.Computable, args compute.RequestArgs) {
	r.Compute(args)
}
func HandleRequestArgs(args compute.RequestArgs) (string, error) {
	fmt.Print(args)
	switch t := args.Args.(type) {
	case compute.FipsCodeCompute:
		_, ok := args.Args.(compute.FipsCodeCompute)
		if ok {
			var r = compute.NSIStructureSimulation{}
			computeConcurrentEvent(r, args)
			return "computing", nil
		}

	case compute.BboxCompute:
		_, ok := args.Args.(compute.BboxCompute)
		if ok {
			var r = compute.NSIStructureSimulation{}
			go computeEvent(r, args)
			return "computing", nil
		}

	default:
		s := fmt.Sprintf("I am de fault of your request %T\n.", t)
		return s, nil //Error{Error: "cannot handle it any longer."}
	}
	return "umm. shouldnt get here.", nil

}
func main() {
	var cfg Config
	if cfg.LambdaContext {
		log.Print("starting server; Running On AWS LAMBDA")
		lambda.Start(HandleRequestArgs)
	} else {
		log.Print("Not on Lambda")
	}
}
