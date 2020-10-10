package main

import (
	"testing"

	"github.com/USACE/go-consequences/compute"
	"github.com/USACE/go-consequences/hazards"
)

func TestSampleSimulation(t *testing.T) {
	var hazard = hazards.DepthEvent{Depth: 12.34}
	var args = compute.FipsCodeCompute{ID: "123", FIPS: "06", HazardArgs: hazard}
	var rargs = compute.RequestArgs{Args: args}
	HandleRequestArgs(rargs)

}
