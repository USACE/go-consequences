package structures

import (
	"fmt"
	"testing"

	"github.com/USACE/go-consequences/hazards"
)

// for testing with Github action
const path = "./data/occtypes.json"

//const path2 = "./data/erosion_trial4.json"

// for testing locally
//const path = "/workspaces/Go_Consequences/data/occtypes.json"

//const path2 = "/workspaces/Go_Consequences/data/erosion_trial6.json"

func Test_JsonReading(t *testing.T) {
	jotp := JsonOccupancyTypeProvider{}
	jotp.Init(path)
	m := jotp.occupancyTypesContainer.OccupancyTypes
	fmt.Println(m["COM1"].ComponentDamageFunctions["contents"].DamageFunctions[hazards.Depth].Source)
}

/*
func Test_JsonMerging(t *testing.T) {
	jotp := JsonOccupancyTypeProvider{}
	jotp.Init(path)
	jotp2 := JsonOccupancyTypeProvider{}
	jotp2.Init(path2)
	m := jotp2.occupancyTypesContainer.OccupancyTypes
	err := jotp.occupancyTypesContainer.MergeMap(m)
	if err != nil {
		panic(err)
	}
	fmt.Println(jotp.occupancyTypesContainer.OccupancyTypes["COM1"].ComponentDamageFunctions["contents"].DamageFunctions[hazards.Erosion].Source)
}

func Test_JsonWriting(t *testing.T) {
	jotp := JsonOccupancyTypeProvider{}
	jotp.Init(path)
	jotp2 := JsonOccupancyTypeProvider{}
	jotp2.Init(path2)
	m := jotp2.occupancyTypesContainer.OccupancyTypes
	err := jotp.occupancyTypesContainer.MergeMap(m)
	if err != nil {
		panic(err)
	}
	err = jotp.Write("/workspaces/Go_Consequences/data/occtypes_merged.json")
	if err != nil {
		panic(err)
	}
}
*/
