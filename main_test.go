package main

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/USACE/go-consequences/census"
	"github.com/USACE/go-consequences/compute"
	"github.com/USACE/go-consequences/hazards"
)

func TestSampleSimulation(t *testing.T) {
	var hazard = hazards.DepthEvent{Depth: 12.34}
	var args = compute.FipsCodeCompute{ID: "123", FIPS: "11", HazardArgs: hazard}
	var rargs = compute.RequestArgs{Args: args, Concurrent: true}
	HandleRequestArgs(rargs)

}

/* Honestly, this is why i can't have nice things.
func TestNationalSimulationConcurrent(t *testing.T) {
	f := census.StateToCountyFipsMap()
	var hazard = hazards.DepthEvent{Depth: 12.34}
	start := time.Now()
	var wg sync.WaitGroup
	wg.Add(len(f))
	for key, _ := range f {
		go func(state string) {
			defer wg.Done()
			var args = compute.FipsCodeCompute{ID: "123", FIPS: state, HazardArgs: hazard}
			var rargs = compute.RequestArgs{Args: args}
			HandleRequestArgs(rargs)
		}(key)
	}
	wg.Wait()
	elapsed := time.Since(start)
	fmt.Println(fmt.Sprintf("Simulation complete, took: %s", elapsed))

}
*/
//TestNationalSimulationConccurentByStateOnly
func TestNationalSimulationConcurrentByStateOnly(t *testing.T) {
	f := census.StateToCountyFipsMap()
	var hazard = hazards.DepthEvent{Depth: 12.34}
	start := time.Now()
	var wg sync.WaitGroup
	wg.Add(len(f))
	for key, _ := range f {
		go func(state string) {
			defer wg.Done()
			var args = compute.FipsCodeCompute{ID: "123", FIPS: state, HazardArgs: hazard}
			var rargs = compute.RequestArgs{Args: args, Concurrent: false}
			HandleRequestArgs(rargs)
		}(key)
	}
	wg.Wait()
	elapsed := time.Since(start)
	fmt.Println(fmt.Sprintf("Simulation complete, took: %s", elapsed))

}

//TestNationalSimulation - runs sequential by state, concurrent by county within a state
func TestNationalSimulation(t *testing.T) {
	f := census.StateToCountyFipsMap()
	var hazard = hazards.DepthEvent{Depth: 12.34}
	start := time.Now()
	for key, _ := range f {
		var args = compute.FipsCodeCompute{ID: "123", FIPS: key, HazardArgs: hazard}
		var rargs = compute.RequestArgs{Args: args, Concurrent: true}
		HandleRequestArgs(rargs)
	}
	elapsed := time.Since(start)
	fmt.Println(fmt.Sprintf("Simulation complete, took: %s", elapsed))

}
