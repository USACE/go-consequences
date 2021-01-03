package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/USACE/go-consequences/census"
	"github.com/USACE/go-consequences/compute"
	"github.com/USACE/go-consequences/hazards"
	"github.com/aws/aws-lambda-go/lambda"
)

//Config describes the configuration settings for go-consequences.
type Config struct {
	SkipJWT       bool
	LambdaContext bool
	DBUser        string
	DBPass        string
	DBName        string
	DBHost        string
	DBSSLMode     string
}

func computeConcurrentEvent(r compute.Computable, args compute.RequestArgs) string {
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
				s := "COMPLETE FOR SIMULATION" + "\n"
				fmt.Println("COMPLETE FOR SIMULATION")
				elapsedNSI := startTime.Sub(nsitime)
				elapsedCompute := startTime.Sub(computetime)
				elapsedClock := time.Since(startTime)
				//s += fmt.Sprintf("NSI Took %s", -elapsedNSI) + "\n"
				fmt.Println(fmt.Sprintf("NSI Took %s", -elapsedNSI))
				s += fmt.Sprintf("Compute Took %s computer time", -elapsedCompute) + "\n"
				fmt.Println(fmt.Sprintf("Compute Took %s", -elapsedCompute))
				s += fmt.Sprintf("Clock Time Taken was %s", elapsedClock) + "\n"
				fmt.Println(fmt.Sprintf("Clock Time Taken was %s", elapsedClock))
				s += fmt.Sprintf("Total Structure Count %d", count) + "\n"
				fmt.Println(fmt.Sprintf("Total Structure Count %d", count))
				s += fmt.Sprintf("Total Structure Damage %f", sdam) + "\n"
				fmt.Println(fmt.Sprintf("Total Structure Damage %f", sdam))
				s += fmt.Sprintf("Total Content Damage %f", cdam) + "\n"
				fmt.Println(fmt.Sprintf("Total Content Damage %f", cdam))
				fmt.Println("*****************SUMMMARY*****************")
				s += "*****************SUMMMARY*****************\n"
				rows := make([]compute.SimulationSummaryRow, len(rowMap))
				idx := 0
				for _, val := range rowMap {
					fmt.Println(fmt.Sprintf("for %s, there were %d structures with %f structure damages %f content damages for damage category %s", fips, val.StructureCount, val.StructureDamage, val.ContentDamage, val.RowHeader))
					s += fmt.Sprintf("for %s, there were %d structures with %f structure damages %f content damages for damage category %s", fips, val.StructureCount, val.StructureDamage, val.ContentDamage, val.RowHeader) + "\n"
					rows[idx] = val
					idx++
				}
				//var ret = SimulationSummary{ColumnNames: header, Rows: rows, NSITime: elapsedNsi, Computetime: elapsed}
				return s
			}
			//2 characters but not a state?
			r.Compute(args) //should fail
		} else {
			//not two characters
			r.Compute(args) //should work
		}
	} else {
		r.Compute(args)
	}
	return "didnt work"
}
func computeEvent(r compute.Computable, args compute.RequestArgs) {
	r.Compute(args)
}

//HandleRequestArgs handles request args and returns a string and an error
func HandleRequestArgs(args compute.RequestArgs) (string, error) {
	fmt.Print(args)
	switch t := args.Args.(type) {
	case compute.FipsCodeCompute:
		_, ok := args.Args.(compute.FipsCodeCompute)
		if ok {
			var r = compute.NSIStructureSimulation{}
			s := computeConcurrentEvent(r, args)
			return s, nil
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
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			params := r.URL.Query()
			fipsid, fipsPresent := params["FIPS"]
			if !fipsPresent {
				http.Error(w, "No FIPS argument", http.StatusNotFound)
			} else {
				if len(fipsid[0]) != 2 {
					http.Error(w, "Invalid FIPS argument", http.StatusNotFound)
				} else {
					depthParam, depthPresent := params["Depth"]
					if !depthPresent {
						http.Error(w, "Invalid Depth argument", http.StatusNotFound)
					} else {
						//cast to args
						depth, err := strconv.ParseFloat(depthParam[0], 64)
						var hazard = hazards.DepthEvent{Depth: 12.34}
						if err == nil {
							hazard = hazards.DepthEvent{Depth: depth}
						}
						var args = compute.FipsCodeCompute{ID: "123", FIPS: fipsid[0], HazardArgs: hazard}
						var rargs = compute.RequestArgs{Args: args, Concurrent: true}
						s, _ := HandleRequestArgs(rargs)
						fmt.Fprintf(w, s)
					}
				}
			}
		})
		log.Print("Not on Lambda")
		log.Print("starting local server")
		log.Fatal(http.ListenAndServe("localhost:3030", nil))
	}
}
