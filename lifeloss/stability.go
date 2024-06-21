package lifeloss

import (
	"log"

	"github.com/USACE/go-consequences/hazards"
)

type Stability uint

const (
	Stable    Stability = 1 //1
	Collapsed Stability = 2 //2
)

func (s Stability) String() string {
	return []string{"Stable", "Collapsed"}[s-1]
}

type StabilityCriteria struct {
	//curve paireddata.PairedData //assume x is depth and y is velocity?
	//curve paireddata.UncertaintyPairedData
	MinimumVelocity    float64
	MinimumDepth       float64
	DepthTimesVelocity float64
}

// @TODO: add interface for Stability as well as adding Engineered Stability mobile homes stability which will require paired data, max velocity, uncertainty and probably some other features.
var RescDamWoodUnanchored = StabilityCriteria{
	MinimumVelocity:    0.0,
	MinimumDepth:       0.0,
	DepthTimesVelocity: 32.3,
}
var RescDamWoodAnchored = StabilityCriteria{
	MinimumVelocity:    0.0,
	MinimumDepth:       0.0,
	DepthTimesVelocity: 75.3,
}
var RescDamMasonryConcreteBrick = StabilityCriteria{
	MinimumVelocity:    6.6,
	MinimumDepth:       0.0,
	DepthTimesVelocity: 75.3,
}

func (sc StabilityCriteria) Evaluate(e hazards.HazardEvent) Stability {
	depth := 0.0    //expected to be maximum depth across all time
	velocity := 0.0 //expected to be maximum velocity
	dv := 0.0
	if e.Has(hazards.DV) {
		dv = e.DV()
		if e.Has(hazards.Depth) {
			//get depth from the hazard
			depth = e.Depth()
			if e.Has(hazards.Velocity) {
				//get velocity from the hazard
				velocity = e.Velocity()
			} else {
				if depth == 0 {
					velocity = dv
				} else {
					velocity = dv / depth // assumes max depth happened at the same time as max velocity - which is not true, so this is an underestimate of max velocity
				}

				//ergo velocity was strictly greater than this value at some point, if i compare this velocity to the threshold of minimum velocity that must be exceeded
				//for concrete strucures, it will yeild the correct result less frequently (because velocity could be greater at some point)
			}
		} else {
			log.Fatal("no depth to evaluate stability criteria")
		}

	} else {
		log.Fatal("no dv to evaluate stability criteria")
	}
	if depth > sc.MinimumDepth {
		if velocity > sc.MinimumVelocity {
			if dv >= sc.DepthTimesVelocity {
				return Collapsed
			}
		}
	}
	return Stable
	//the below psuedo code would work great if i have the instantaneous depth and velocity at the time of the maximum of their product - if i dont have that, using d*V calculated externally is a better overall approach
	/*
		depth := 0.0    //this should be depth at the time of the peak velocity*depth calculation.
		velocity := 0.0 //this should be the velocity at the time of the peak velocity*depth calculation
		if e.Has(hazards.Depth) {
			//get depth from the hazard
			depth = e.Depth()
		}
		if e.Has(hazards.Velocity) {
			//get velocity from the hazard
			velocity = e.Velocity()
		}

		//shortcut if velocity is less than the min value...
		if sc.curve.Xvals[0] > velocity {
			return false
		}
		comparisonDepth := sc.curve.SampleValue(velocity)

		return depth >= comparisonDepth
	*/
}
