package structures

import (
	"fmt"
	"testing"

	"github.com/USACE/go-consequences/hazards"
)

// for testing with Github action
const path = "./data/occtypes.json"

// for testing locally
//const path = "/workspaces/Go_Consequences/data/occtypes.json"

func Test_JsonReading(t *testing.T) {
	jotp := JsonOccupancyTypeProvider{}
	jotp.Init(path)
	m := jotp.OccupancyTypeMap()
	fmt.Println(m["COM1"].ContentDFF.DamageFunctions[hazards.Erosion].Source)
}
