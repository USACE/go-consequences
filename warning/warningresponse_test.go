package warning

import (
	"fmt"
	"testing"

	"github.com/USACE/go-consequences/hazards"
	"github.com/USACE/go-consequences/structures"
)

func Test_ComplianceBasedWarnings(t *testing.T) {
	ws := InitComplianceBasedWarningSystem(1234, .75)
	e := hazards.DepthEvent{}
	e.SetDepth(12345.678) //not used in this warning system - design is to ultimately allow for arrival time to be used in warning if needed.
	s := createStructureDeterministicForTesting(1000)
	remainingpop, result := ws.WarningFunction()(s, e)
	fmt.Println(result)
	fmt.Println(remainingpop)
}
func ExampleComplianceBasedWarnings() {
	ws := InitComplianceBasedWarningSystem(1234, .75)
	e := hazards.DepthEvent{}
	e.SetDepth(12345.678) //not used in this warning system - design is to ultimately allow for arrival time to be used in warning if needed.
	s := createStructureDeterministicForTesting(1000)
	remainingpop, result := ws.WarningFunction()(s, e)
	fmt.Println(result)
	fmt.Println(remainingpop)
}
func createStructureDeterministicForTesting(population int32) structures.StructureDeterministic {
	s := structures.StructureDeterministic{
		BaseStructure: structures.BaseStructure{
			Name:            "blank",
			DamCat:          "blank",
			CBFips:          "blank",
			X:               0,
			Y:               0,
			GroundElevation: 0,
		},
		OccType: structures.OccupancyTypeDeterministic{
			Name:                     "blank",
			ComponentDamageFunctions: map[string]structures.DamageFunctionFamily{},
		},
		FoundType:        "blank",
		FirmZone:         "blank",
		ConstructionType: "blank",
		StructVal:        0,
		ContVal:          0,
		FoundHt:          0,
		NumStories:       0,
		PopulationSet: structures.PopulationSet{
			Pop2pmo65: population,
			Pop2pmu65: population,
			Pop2amo65: population,
			Pop2amu65: population,
		},
	}
	return s
}
