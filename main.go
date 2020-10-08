package main

import (
	"log"

	"github.com/USACE/go-consequences/compute"
	"github.com/USACE/go-consequences/consequences"
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

func HandleLambdaEvent(args compute.ComputeArgs) (consequences.ConsequenceDamageResult, error) {
	var r = compute.NSIStructureSimulation{}
	r.Compute(args)
	return r.Result, nil
}
func main() {
	var cfg Config
	if cfg.LambdaContext {
		log.Print("starting server; Running On AWS LAMBDA")
		lambda.Start(HandleLambdaEvent)
	} else {
		log.Print("Not on Lambda")
	}
	/*
		var s = consequences.BaseStructure()
		var d = hazards.DepthEvent{Depth: 3.0}
		depths := []float64{3.0, 0.0, 0.5, 1.0, 1.0001, 2.25, 2.5, 2.75, 3.99, 4, 5}
		for idx := range depths {
			d.Depth = depths[idx]
			fmt.Println("for a depth of", d.Depth, s.ComputeConsequences(d))
		}
		fmt.Println("*********Uncertainty************")
		var su = consequences.BaseStructureU()
		for i := 0; i < 10; i++ {
			fmt.Println("for a depth of", d.Depth, su.ComputeConsequences(d))
		}
		fmt.Println("*********Uncertainty************")
		s.FoundHt = 1.1 //test interpolation due to foundation height putting depth back in range
		ret := s.ComputeConsequences(d)
		fmt.Println("for a depth of", d.Depth, ret)

		var f = hazards.FireEvent{Intensity: hazards.Low}
		s = consequences.ConvertBaseStructureToFire(s)
		ret = s.ComputeConsequences(f)
		fmt.Println("for a fire intensity of", f.Intensity, ret)

		f.Intensity = hazards.Medium
		ret = s.ComputeConsequences(f)
		fmt.Println("for a fire intensity of", f.Intensity, ret)

		f.Intensity = hazards.High
		ret = s.ComputeConsequences(f)
		fmt.Println("for a fire intensity of", f.Intensity, ret)

		//var bbox string = "-81.58418,30.25165,-81.58161,30.26939,-81.55898,30.26939,-81.55281,30.24998,-81.58418,30.25165"

		startnsi := time.Now()
		var fips string = "06"
		d.Depth = 5.324 //testing cost of interpolation.
		structures := nsi.GetByFips(fips)
		//structures := nsi.GetByBbox(bbox)
		elapsedNsi := time.Since(startnsi)
		startcompute := time.Now()
		var count = 0
		for i, str := range structures {
			str.ComputeConsequences(d)
			//fmt.Println(i, "at structure", str.Name, "for a depth of", d.Depth, str.ComputeConsequences(d))
			count = i
		}
		count += 1
		elapsed := time.Since(startcompute)
		fmt.Println(fmt.Sprintf("NSI Fetching took %s Compute took %s for %d structures", elapsedNsi, elapsed, count))
	*/
}
