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

	depthParameter := hazards.Depth
	durationParameter := hazards.Depth | hazards.Duration | hazards.LongDuration | hazards.Velocity
	moderateVParameter := hazards.Depth | hazards.Duration | hazards.Velocity | hazards.ModerateVelocity
	highVParameter := hazards.Depth | hazards.Duration | hazards.Velocity | hazards.HighVelocity

	// 1: ONE-STORY SLAB ON GRADE
	res11snb := "RES1-1SNB"
	//depth only
	depthDamage1 := []float64{0, 0, 0, 0, 6, 20, 27, 32, 37, 41, 46, 49, 52, 55, 57, 59, 60, 62, 64, 65, 66}
	//long duration
	durationDamage1 := []float64{0, 0, 0, 0, 7, 22, 30, 35, 41, 45, 51, 54, 57, 61, 62, 65, 66, 68, 70, 72, 73}
	//moderate velocity //2-5 f/s
	moderateVelocityDamage1 := []float64{0, 0, 0, 0, 7, 22, 30, 35, 41, 45, 51, 54, 57, 61, 62, 65, 66, 68, 70, 72, 73}
	//high velocity //5+ f/s - note about 5-10...
	highVelocityDamage1 := []float64{0, 0, 0, 0, 10, 29, 39, 45, 52, 58, 65, 69, 74, 78, 80, 83, 85, 88, 90, 92, 94}
	cdf1 := make(map[string]DamageFunctionFamilyStochastic)
	dffs1 := DamageFunctionFamilyStochastic{
		DamageFunctions: map[hazards.Parameter]DamageFunctionStochastic{},
	}
	dffs1.DamageFunctions[depthParameter] = DamageFunctionStochastic{
		Source:       "FEMA Inland Damage Functions Report 20211001",
		DamageDriver: hazards.Depth,
		DamageFunction: paireddata.UncertaintyPairedData{
			Xvals: depths,
			Yvals: toContinuousDistribution(depthDamage1),
		},
	}
	dffs1.DamageFunctions[durationParameter] = DamageFunctionStochastic{
		Source:       "FEMA Inland Damage Functions Report 20211001",
		DamageDriver: hazards.Depth,
		DamageFunction: paireddata.UncertaintyPairedData{
			Xvals: depths,
			Yvals: toContinuousDistribution(durationDamage1),
		},
	}
	dffs1.DamageFunctions[moderateVParameter] = DamageFunctionStochastic{
		Source:       "FEMA Inland Damage Functions Report 20211001",
		DamageDriver: hazards.Depth,
		DamageFunction: paireddata.UncertaintyPairedData{
			Xvals: depths,
			Yvals: toContinuousDistribution(moderateVelocityDamage1),
		},
	}
	dffs1.DamageFunctions[highVParameter] = DamageFunctionStochastic{
		Source:       "FEMA Inland Damage Functions Report 20211001",
		DamageDriver: hazards.Depth,
		DamageFunction: paireddata.UncertaintyPairedData{
			Xvals: depths,
			Yvals: toContinuousDistribution(highVelocityDamage1),
		},
	}
	cdf1["structure"] = dffs1
	cdf1["content"] = dffs1 //make them the same? no clue. thats what ill do.
	o1 := OccupancyTypeStochastic{
		Name:                     res11snb,
		ComponentDamageFunctions: cdf1,
	}

	// 2: TWO-STORY SLAB ON GRADE
	res12snb := "RES1-2SNB"
	//depth only
	depthDamage2 := []float64{0, 0, 0, 0, 4, 14, 20, 23, 27, 30, 34, 36, 38, 40, 42, 43, 45, 46, 47, 48, 49}
	//long duration
	durationDamage2 := []float64{0, 0, 0, 0, 5, 16, 22, 26, 30, 33, 37, 40, 42, 45, 46, 48, 49, 51, 52, 53, 54}
	//moderate velocity //2-5 f/s
	moderateVelocityDamage2 := []float64{0, 0, 0, 0, 5, 18, 24, 28, 33, 36, 41, 43, 46, 49, 51, 52, 54, 55, 57, 58, 59}
	//high velocity //5+ f/s - note about 5-10...
	highVelocityDamage2 := []float64{0, 0, 0, 0, 7, 23, 31, 36, 42, 46, 52, 55, 59, 62, 64, 66, 68, 70, 72, 73, 75}
	cdf2 := make(map[string]DamageFunctionFamilyStochastic)
	dffs2 := DamageFunctionFamilyStochastic{
		DamageFunctions: map[hazards.Parameter]DamageFunctionStochastic{},
	}
	dffs2.DamageFunctions[depthParameter] = DamageFunctionStochastic{
		Source:       "FEMA Inland Damage Functions Report 20211001",
		DamageDriver: hazards.Depth,
		DamageFunction: paireddata.UncertaintyPairedData{
			Xvals: depths,
			Yvals: toContinuousDistribution(depthDamage2),
		},
	}
	dffs2.DamageFunctions[durationParameter] = DamageFunctionStochastic{
		Source:       "FEMA Inland Damage Functions Report 20211001",
		DamageDriver: hazards.Depth,
		DamageFunction: paireddata.UncertaintyPairedData{
			Xvals: depths,
			Yvals: toContinuousDistribution(durationDamage2),
		},
	}
	dffs2.DamageFunctions[moderateVParameter] = DamageFunctionStochastic{
		Source:       "FEMA Inland Damage Functions Report 20211001",
		DamageDriver: hazards.Depth,
		DamageFunction: paireddata.UncertaintyPairedData{
			Xvals: depths,
			Yvals: toContinuousDistribution(moderateVelocityDamage2),
		},
	}
	dffs2.DamageFunctions[highVParameter] = DamageFunctionStochastic{
		Source:       "FEMA Inland Damage Functions Report 20211001",
		DamageDriver: hazards.Depth,
		DamageFunction: paireddata.UncertaintyPairedData{
			Xvals: depths,
			Yvals: toContinuousDistribution(highVelocityDamage2),
		},
	}
	cdf2["structure"] = dffs2
	cdf2["content"] = dffs2 //make them the same? no clue. thats what ill do.
	o2 := OccupancyTypeStochastic{
		Name:                     res12snb,
		ComponentDamageFunctions: cdf2,
	}

	// 3: ONE-STORY w/ PIER/CRAWLSPACE w/ OPENINGS
	res11snbc := "RES1-1SNB-C"
	//depth only
	depthDamage3 := []float64{0, 0, 0, 5, 14, 27, 35, 39, 44, 48, 53, 56, 59, 62, 64, 66, 68, 69, 71, 73, 74}
	//long duration
	durationDamage3 := []float64{0, 0, 0, 5, 15, 30, 37, 43, 48, 53, 58, 61, 65, 68, 70, 72, 74, 76, 78, 79, 81}
	//moderate velocity //2-5 f/s
	moderateVelocityDamage3 := []float64{0, 0, 0, 5, 15, 30, 37, 43, 48, 53, 58, 61, 65, 68, 70, 72, 74, 76, 78, 79, 81}
	//high velocity //5+ f/s - note about 5-10...
	highVelocityDamage3 := []float64{0, 0, 0, 6, 17, 36, 46, 53, 60, 66, 73, 77, 81, 85, 88, 90, 93, 95, 97, 100, 100}
	cdf3 := make(map[string]DamageFunctionFamilyStochastic)
	dffs3 := DamageFunctionFamilyStochastic{
		DamageFunctions: map[hazards.Parameter]DamageFunctionStochastic{},
	}
	dffs3.DamageFunctions[depthParameter] = DamageFunctionStochastic{
		Source:       "FEMA Inland Damage Functions Report 20211001",
		DamageDriver: hazards.Depth,
		DamageFunction: paireddata.UncertaintyPairedData{
			Xvals: depths,
			Yvals: toContinuousDistribution(depthDamage3),
		},
	}
	dffs3.DamageFunctions[durationParameter] = DamageFunctionStochastic{
		Source:       "FEMA Inland Damage Functions Report 20211001",
		DamageDriver: hazards.Depth,
		DamageFunction: paireddata.UncertaintyPairedData{
			Xvals: depths,
			Yvals: toContinuousDistribution(durationDamage3),
		},
	}
	dffs3.DamageFunctions[moderateVParameter] = DamageFunctionStochastic{
		Source:       "FEMA Inland Damage Functions Report 20211001",
		DamageDriver: hazards.Depth,
		DamageFunction: paireddata.UncertaintyPairedData{
			Xvals: depths,
			Yvals: toContinuousDistribution(moderateVelocityDamage3),
		},
	}
	dffs3.DamageFunctions[highVParameter] = DamageFunctionStochastic{
		Source:       "FEMA Inland Damage Functions Report 20211001",
		DamageDriver: hazards.Depth,
		DamageFunction: paireddata.UncertaintyPairedData{
			Xvals: depths,
			Yvals: toContinuousDistribution(highVelocityDamage3),
		},
	}
	cdf3["structure"] = dffs3
	cdf3["content"] = dffs3 //make them the same? no clue. thats what ill do.
	o3 := OccupancyTypeStochastic{
		Name:                     res11snbc,
		ComponentDamageFunctions: cdf3,
	}

	// 4: TWO-STORY w/ PIER/CRAWLSPACE w/ OPENINGS
	res12snbc := "RES1-2SNB-C"
	//depth only
	depthDamage4 := []float64{0, 0, 0, 5, 11, 22, 27, 31, 34, 37, 41, 43, 46, 48, 49, 51, 52, 53, 55, 56, 57}
	//long duration
	durationDamage4 := []float64{0, 0, 0, 5, 12, 23, 29, 33, 37, 41, 45, 47, 50, 52, 54, 55, 57, 58, 60, 61, 62}
	//moderate velocity //2-5 f/s
	moderateVelocityDamage4 := []float64{0, 0, 0, 5, 13, 25, 32, 36, 40, 44, 48, 51, 54, 57, 58, 60, 61, 63, 64, 66, 67}
	//high velocity //5+ f/s - note about 5-10...
	highVelocityDamage4 := []float64{0, 0, 0, 5, 15, 30, 38, 43, 49, 54, 59, 62, 66, 69, 71, 74, 76, 77, 79, 81, 82}
	cdf4 := make(map[string]DamageFunctionFamilyStochastic)
	dffs4 := DamageFunctionFamilyStochastic{
		DamageFunctions: map[hazards.Parameter]DamageFunctionStochastic{},
	}
	dffs4.DamageFunctions[depthParameter] = DamageFunctionStochastic{
		Source:       "FEMA Inland Damage Functions Report 20211001",
		DamageDriver: hazards.Depth,
		DamageFunction: paireddata.UncertaintyPairedData{ 
			Xvals: depths,
			Yvals: toContinuousDistribution(depthDamage4),
		},
	}
	dffs4.DamageFunctions[durationParameter] = DamageFunctionStochastic{
		Source:       "FEMA Inland Damage Functions Report 20211001",
		DamageDriver: hazards.Depth,
		DamageFunction: paireddata.UncertaintyPairedData{
			Xvals: depths,
			Yvals: toContinuousDistribution(durationDamage4),
		},
	}
	dffs4.DamageFunctions[moderateVParameter] = DamageFunctionStochastic{
		Source:       "FEMA Inland Damage Functions Report 20211001",
		DamageDriver: hazards.Depth,
		DamageFunction: paireddata.UncertaintyPairedData{
			Xvals: depths,
			Yvals: toContinuousDistribution(moderateVelocityDamage4),
		},
	}
	dffs4.DamageFunctions[highVParameter] = DamageFunctionStochastic{
		Source:       "FEMA Inland Damage Functions Report 20211001",
		DamageDriver: hazards.Depth,
		DamageFunction: paireddata.UncertaintyPairedData{
			Xvals: depths,
			Yvals: toContinuousDistribution(highVelocityDamage4),
		},
	}
	cdf4["structure"] = dffs4
	cdf4["content"] = dffs4 //make them the same? no clue. thats what ill do.
	o4 := OccupancyTypeStochastic{
		Name:                     res12snbc,
		ComponentDamageFunctions: cdf4,
	}

	// 5: ONE-STORY w/ PIER/CRAWLSPACE w/o OPENINGS
	res11snbp := "RES1-1SNB-P"
	//depth only
	depthDamage5 := []float64{0, 0, 0, 5, 14, 27, 35, 39, 44, 48, 53, 56, 59, 62, 64, 66, 68, 69, 71, 73, 74}
	//long duration
	durationDamage5 := []float64{0, 0, 0, 5, 15, 30, 37, 43, 48, 53, 58, 61, 65, 68, 70, 72, 74, 76, 78, 79, 81}
	//moderate velocity //2-5 f/s
	moderateVelocityDamage5 := []float64{0, 0, 0, 5, 15, 30, 37, 43, 48, 53, 58, 61, 65, 68, 70, 72, 74, 76, 78, 79, 81}
	//high velocity //5+ f/s - note about 5-10...
	highVelocityDamage5 := []float64{0, 0, 0, 11, 22, 41, 51, 58, 65, 71, 78, 82, 86, 90, 93, 95, 98, 100, 100, 100, 100}
	cdf5 := make(map[string]DamageFunctionFamilyStochastic)
	dffs5 := DamageFunctionFamilyStochastic{
		DamageFunctions: map[hazards.Parameter]DamageFunctionStochastic{},
	}
	dffs5.DamageFunctions[depthParameter] = DamageFunctionStochastic{
		Source:       "FEMA Inland Damage Functions Report 20211001",
		DamageDriver: hazards.Depth,
		DamageFunction: paireddata.UncertaintyPairedData{
			Xvals: depths,
			Yvals: toContinuousDistribution(depthDamage5),
		},
	}
	dffs5.DamageFunctions[durationParameter] = DamageFunctionStochastic{
		Source:       "FEMA Inland Damage Functions Report 20211001",
		DamageDriver: hazards.Depth,
		DamageFunction: paireddata.UncertaintyPairedData{
			Xvals: depths,
			Yvals: toContinuousDistribution(durationDamage5),
		},
	}
	dffs5.DamageFunctions[moderateVParameter] = DamageFunctionStochastic{
		Source:       "FEMA Inland Damage Functions Report 20211001",
		DamageDriver: hazards.Depth,
		DamageFunction: paireddata.UncertaintyPairedData{
			Xvals: depths,
			Yvals: toContinuousDistribution(moderateVelocityDamage5),
		},
	}
	dffs5.DamageFunctions[highVParameter] = DamageFunctionStochastic{
		Source:       "FEMA Inland Damage Functions Report 20211001",
		DamageDriver: hazards.Depth,
		DamageFunction: paireddata.UncertaintyPairedData{
			Xvals: depths,
			Yvals: toContinuousDistribution(highVelocityDamage5),
		},
	}
	cdf5["structure"] = dffs5
	cdf5["content"] = dffs5 //make them the same? no clue. thats what ill do.
	o5 := OccupancyTypeStochastic{
		Name:                     res11snbp,
		ComponentDamageFunctions: cdf5,
	}

	// 6: TWO-STORY w/ PIER/CRAWLSPACE w/o OPENINGS
	res12snbp := "RES1-2SNB-P"
	//depth only
	depthDamage6 := []float64{0, 0, 0, 5, 11, 22, 27, 31, 34, 37, 41, 43, 46, 48, 49, 51, 52, 53, 55, 56, 57}
	//long duration
	durationDamage6 := []float64{0, 0, 0, 5, 12, 23, 29, 33, 37, 41, 45, 47, 50, 52, 54, 55, 57, 58, 60, 61, 62}
	//moderate velocity //2-5 f/s
	moderateVelocityDamage6 := []float64{0, 0, 0, 5, 13, 25, 32, 36, 40, 44, 48, 51, 54, 57, 58, 60, 61, 63, 64, 66, 67}
	//high velocity //5+ f/s - note about 5-10...
	highVelocityDamage6 := []float64{0, 0, 0, 10, 20, 35, 43, 48, 54, 59, 64, 67, 71, 74, 76, 79, 81, 82, 84, 86, 87}
	cdf6 := make(map[string]DamageFunctionFamilyStochastic)
	dffs6 := DamageFunctionFamilyStochastic{
		DamageFunctions: map[hazards.Parameter]DamageFunctionStochastic{},
	}
	dffs6.DamageFunctions[depthParameter] = DamageFunctionStochastic{
		Source:       "FEMA Inland Damage Functions Report 20211001",
		DamageDriver: hazards.Depth,
		DamageFunction: paireddata.UncertaintyPairedData{
			Xvals: depths,
			Yvals: toContinuousDistribution(depthDamage6),
		},
	}
	dffs6.DamageFunctions[durationParameter] = DamageFunctionStochastic{
		Source:       "FEMA Inland Damage Functions Report 20211001",
		DamageDriver: hazards.Depth,
		DamageFunction: paireddata.UncertaintyPairedData{
			Xvals: depths,
			Yvals: toContinuousDistribution(durationDamage6),
		},
	}
	dffs6.DamageFunctions[moderateVParameter] = DamageFunctionStochastic{
		Source:       "FEMA Inland Damage Functions Report 20211001",
		DamageDriver: hazards.Depth,
		DamageFunction: paireddata.UncertaintyPairedData{
			Xvals: depths,
			Yvals: toContinuousDistribution(moderateVelocityDamage6),
		},
	}
	dffs6.DamageFunctions[highVParameter] = DamageFunctionStochastic{
		Source:       "FEMA Inland Damage Functions Report 20211001",
		DamageDriver: hazards.Depth,
		DamageFunction: paireddata.UncertaintyPairedData{
			Xvals: depths,
			Yvals: toContinuousDistribution(highVelocityDamage6),
		},
	}
	cdf6["structure"] = dffs6
	cdf6["content"] = dffs6 //make them the same? no clue. thats what ill do.
	o6 := OccupancyTypeStochastic{
		Name:                     res12snbp,
		ComponentDamageFunctions: cdf6,
	}

	// 7: ONE-STORY w/ BASEMENT
	res11swb := "RES1-1SWB"
	//depth only
	depthDamage7 := []float64{0, 0, 10, 12, 20, 34, 41, 46, 51, 55, 60, 63, 66, 69, 71, 73, 74, 76, 78, 79, 80}
	//long duration
	durationDamage7 := []float64{0, 0, 11, 13, 22, 37, 45, 50, 56, 60, 66, 69, 73, 76, 78, 80, 82, 84, 86, 87, 88}
	//moderate velocity //2-5 f/s
	moderateVelocityDamage7 := []float64{0, 0, 10, 13, 21, 36, 43, 48, 53, 58, 63, 66, 69, 72, 74, 76, 78, 80, 82, 83, 84}
	//high velocity //5+ f/s - note about 5-10...
	highVelocityDamage7 := []float64{0, 0, 13, 16, 27, 46, 56, 62, 69, 74, 81, 85, 89, 93, 95, 98, 100, 100, 100, 100, 100}
	cdf7 := make(map[string]DamageFunctionFamilyStochastic)
	dffs7 := DamageFunctionFamilyStochastic{
		DamageFunctions: map[hazards.Parameter]DamageFunctionStochastic{},
	}
	dffs7.DamageFunctions[depthParameter] = DamageFunctionStochastic{
		Source:       "FEMA Inland Damage Functions Report 20211001",
		DamageDriver: hazards.Depth,
		DamageFunction: paireddata.UncertaintyPairedData{
			Xvals: depths,
			Yvals: toContinuousDistribution(depthDamage7),
		},
	}
	dffs7.DamageFunctions[durationParameter] = DamageFunctionStochastic{
		Source:       "FEMA Inland Damage Functions Report 20211001",
		DamageDriver: hazards.Depth,
		DamageFunction: paireddata.UncertaintyPairedData{
			Xvals: depths,
			Yvals: toContinuousDistribution(durationDamage7),
		},
	}
	dffs7.DamageFunctions[moderateVParameter] = DamageFunctionStochastic{
		Source:       "FEMA Inland Damage Functions Report 20211001",
		DamageDriver: hazards.Depth,
		DamageFunction: paireddata.UncertaintyPairedData{
			Xvals: depths,
			Yvals: toContinuousDistribution(moderateVelocityDamage7),
		},
	}
	dffs7.DamageFunctions[highVParameter] = DamageFunctionStochastic{
		Source:       "FEMA Inland Damage Functions Report 20211001",
		DamageDriver: hazards.Depth,
		DamageFunction: paireddata.UncertaintyPairedData{
			Xvals: depths,
			Yvals: toContinuousDistribution(highVelocityDamage7),
		},
	}
	cdf7["structure"] = dffs7
	cdf7["content"] = dffs7 //make them the same? no clue. thats what ill do.
	o7 := OccupancyTypeStochastic{
		Name:                     res11swb,
		ComponentDamageFunctions: cdf7,
	}

	// 8: TWO-STORY w/ BASEMENT
	res12swb := "RES1-2SWB"
	//depth only
	depthDamage8 := []float64{0, 0, 9, 10, 16, 26, 32, 35, 39, 42, 46, 48, 50, 53, 54, 55, 57, 58, 59, 60, 61}
	//long duration
	durationDamage8 := []float64{0, 0, 10, 11, 18, 29, 35, 39, 43, 46, 50, 53, 55, 58, 59, 61, 62, 64, 65, 66, 67}
	//moderate velocity //2-5 f/s
	moderateVelocityDamage8 := []float64{0, 0, 10, 11, 18, 29, 35, 39, 43, 46, 50, 53, 55, 58, 59, 61, 62, 64, 65, 66, 67}
	//high velocity //5+ f/s - note about 5-10...
	highVelocityDamage8 := []float64{0, 0, 12, 14, 23, 38, 45, 50, 56, 60, 66, 69, 72, 75, 77, 79, 81, 83, 85, 86, 88}
	cdf8 := make(map[string]DamageFunctionFamilyStochastic)
	dffs8 := DamageFunctionFamilyStochastic{
		DamageFunctions: map[hazards.Parameter]DamageFunctionStochastic{},
	}
	dffs8.DamageFunctions[depthParameter] = DamageFunctionStochastic{
		Source:       "FEMA Inland Damage Functions Report 20211001",
		DamageDriver: hazards.Depth,
		DamageFunction: paireddata.UncertaintyPairedData{
			Xvals: depths,
			Yvals: toContinuousDistribution(depthDamage8),
		},
	}
	dffs8.DamageFunctions[durationParameter] = DamageFunctionStochastic{
		Source:       "FEMA Inland Damage Functions Report 20211001",
		DamageDriver: hazards.Depth,
		DamageFunction: paireddata.UncertaintyPairedData{
			Xvals: depths,
			Yvals: toContinuousDistribution(durationDamage8),
		},
	}
	dffs8.DamageFunctions[moderateVParameter] = DamageFunctionStochastic{
		Source:       "FEMA Inland Damage Functions Report 20211001",
		DamageDriver: hazards.Depth,
		DamageFunction: paireddata.UncertaintyPairedData{
			Xvals: depths,
			Yvals: toContinuousDistribution(moderateVelocityDamage8),
		},
	}
	dffs8.DamageFunctions[highVParameter] = DamageFunctionStochastic{
		Source:       "FEMA Inland Damage Functions Report 20211001",
		DamageDriver: hazards.Depth,
		DamageFunction: paireddata.UncertaintyPairedData{
			Xvals: depths,
			Yvals: toContinuousDistribution(highVelocityDamage8),
		},
	}
	cdf8["structure"] = dffs8
	cdf8["content"] = dffs8 //make them the same? no clue. thats what ill do.
	o8 := OccupancyTypeStochastic{
		Name:                     res12swb,
		ComponentDamageFunctions: cdf8,
	}

	//fmt.Println(o)
	jotp.occupancyTypesContainer.OccupancyTypes[o1.Name] = o1
	jotp.occupancyTypesContainer.OccupancyTypes[o2.Name] = o2
	jotp.occupancyTypesContainer.OccupancyTypes[o3.Name] = o3
	jotp.occupancyTypesContainer.OccupancyTypes[o4.Name] = o4
	jotp.occupancyTypesContainer.OccupancyTypes[o5.Name] = o5
	jotp.occupancyTypesContainer.OccupancyTypes[o6.Name] = o6
	jotp.occupancyTypesContainer.OccupancyTypes[o7.Name] = o7
	jotp.occupancyTypesContainer.OccupancyTypes[o8.Name] = o8
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
