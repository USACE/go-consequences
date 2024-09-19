package structures

import (
	"fmt"
	"testing"

	"github.com/HydrologicEngineeringCenter/go-statistics/paireddata"
	"github.com/HydrologicEngineeringCenter/go-statistics/statistics"
	"github.com/USACE/go-consequences/hazards"
)

// for testing with Github action
///const path = "./data/occtypes.json"

const path2 = "./data/erosion_trial8.json"

// for testing locally
//const path = "/workspaces/Go_Consequences/data/occtypes.json"

//const path2 = "/workspaces/Go_Consequences/data/erosion_trial8.json"

func Test_JsonReading(t *testing.T) {
	jotp := JsonOccupancyTypeProvider{}
	jotp.InitLocalPath(path2)
	m := jotp.occupancyTypesContainer.OccupancyTypes
	fmt.Println(m["COM1"].ComponentDamageFunctions["contents"].DamageFunctions[hazards.Erosion].Source)
	for k, v := range m {
		fmt.Println(k)
		for ck, cv := range v.ComponentDamageFunctions {
			fmt.Println(ck + " " + cv.DamageFunctions[hazards.Erosion].Source)

		}
	}
}
func Test_ModifyDefault(t *testing.T) {
	jotp := JsonOccupancyTypeProvider{}
	jotp.InitDefault()
	//modify the defaults to include FFRD curves.
	//all damage functions use the same depth range.
	depths := []float64{-4, -3, -2, -1, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14}

	//ONE-STORY SLAB ON GRADE
	res11snb := "RES1-1SNB"
	//depth only
	depthParameter := hazards.Depth
	depthDamage := []float64{0, 0, 0, 0, 6, 20, 27, 32, 37, 41, 46, 49, 52, 55, 57, 59, 60, 62, 64, 65, 66}
	//long duration
	durationParameter := hazards.Depth | hazards.Duration | hazards.LongDuration | hazards.Velocity
	durationDamage := []float64{0, 0, 0, 0, 7, 22, 30, 35, 41, 45, 51, 54, 57, 61, 62, 65, 66, 68, 70, 72, 73}
	//moderate velocity //2-5 f/s
	moderateVParameter := hazards.Depth | hazards.Duration | hazards.Velocity | hazards.ModerateVelocity
	moderateVelocityDamage := []float64{0, 0, 0, 0, 7, 22, 30, 35, 41, 45, 51, 54, 57, 61, 62, 65, 66, 68, 70, 72, 73}
	//high velocity //5+ f/s - note about 5-10...
	highVParameter := hazards.Depth | hazards.Duration | hazards.Velocity | hazards.HighVelocity
	highVelocityDamage := []float64{0, 0, 0, 0, 10, 29, 39, 45, 52, 58, 65, 69, 74, 78, 80, 83, 85, 88, 90, 92, 94}
	cdf := make(map[string]DamageFunctionFamilyStochastic)
	dffs := DamageFunctionFamilyStochastic{
		DamageFunctions: map[hazards.Parameter]DamageFunctionStochastic{},
	}
	dffs.DamageFunctions[depthParameter] = DamageFunctionStochastic{
		Source:       "FEMA Inland Damage Functions Report 20211001",
		DamageDriver: hazards.Depth,
		DamageFunction: paireddata.UncertaintyPairedData{
			Xvals: depths,
			Yvals: toContinuousDistribution(depthDamage),
		},
	}
	dffs.DamageFunctions[durationParameter] = DamageFunctionStochastic{
		Source:       "FEMA Inland Damage Functions Report 20211001",
		DamageDriver: hazards.Depth,
		DamageFunction: paireddata.UncertaintyPairedData{
			Xvals: depths,
			Yvals: toContinuousDistribution(durationDamage),
		},
	}
	dffs.DamageFunctions[moderateVParameter] = DamageFunctionStochastic{
		Source:       "FEMA Inland Damage Functions Report 20211001",
		DamageDriver: hazards.Depth,
		DamageFunction: paireddata.UncertaintyPairedData{
			Xvals: depths,
			Yvals: toContinuousDistribution(moderateVelocityDamage),
		},
	}
	dffs.DamageFunctions[highVParameter] = DamageFunctionStochastic{
		Source:       "FEMA Inland Damage Functions Report 20211001",
		DamageDriver: hazards.Depth,
		DamageFunction: paireddata.UncertaintyPairedData{
			Xvals: depths,
			Yvals: toContinuousDistribution(highVelocityDamage),
		},
	}
	cdf["structure"] = dffs
	cdf["content"] = dffs //make them the same? no clue. thats what ill do.
	o := OccupancyTypeStochastic{
		Name:                     res11snb,
		ComponentDamageFunctions: cdf,
	}
	//fmt.Println(o)
	jotp.occupancyTypesContainer.OccupancyTypes[o.Name] = o
	for OccupancyTypeName, v := range jotp.occupancyTypesContainer.OccupancyTypes {
		fmt.Println(OccupancyTypeName)
		for loss_component, component_damageFunctions := range v.ComponentDamageFunctions {
			s := loss_component + ":\n"
			for damageFunctionParameter, damageFunction := range component_damageFunctions.DamageFunctions {
				s += "\t" + damageFunctionParameter.String() + ", " + damageFunction.Source + "\n"
			}
			fmt.Println(s)
		}
	}
	/*err := jotp.Write("./data/Inland_FFRD_damageFunctions.json")
	if err!=nil{
		panic(err)
	}*/
}
func toContinuousDistribution(data []float64) []statistics.ContinuousDistribution {
	dists := make([]statistics.ContinuousDistribution, len(data))
	for i, v := range data {
		dists[i] = statistics.DeterministicDistribution{
			Value: v,
		}
	}
	return dists
}

/*
func Test_JsonMerging(t *testing.T) {
	jotp := JsonOccupancyTypeProvider{}
	jotp.InitDefault()
	jotp2 := JsonOccupancyTypeProvider{}
	jotp2.InitLocalPath(path2)
	m := jotp2.occupancyTypesContainer.OccupancyTypes
	err := jotp.occupancyTypesContainer.MergeMap(m)
	if err != nil {
		panic(err)
	}
	fmt.Println(jotp.occupancyTypesContainer.OccupancyTypes["COM1"].ComponentDamageFunctions["contents"].DamageFunctions[hazards.Erosion].Source)
}
*/
/*
func Test_JsonWriting(t *testing.T) {
	jotp := JsonOccupancyTypeProvider{}
	jotp.InitDefault()
	jotp2 := JsonOccupancyTypeProvider{}
	jotp2.InitLocalPath(path2)
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
