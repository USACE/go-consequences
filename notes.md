# Reconcstruction/Life Cycle Implementation

## 1. Implement basic reconstruction as a component damage function (COMPLETE)

- current component damage functions exist for "structure" and "contents"

- the existing format can be leveraged to return a reconstruction time in days for a given input
    - rather than using the hazard intensity (e.g. flood depth) to return a percent loss, we use the pct_loss as the input. Return value is number of days to return to 0 damage. e.g.:
        | damcat | dmg_pct | t_rebuild_min | t_rebuild_mostlikely | t_rebuild_max |
        |---------|-----------|----------------|------------------------|-----------------|
        | RES | 0 | 0 | 0 | 0 |
        | RES | .5 | 13.5 | 15 | 16.5 |
        | RES | 1 | 27 | 30 | 33 |



## 2. Implement MultiHazard

- `computeConsequencesMulti(events []hazards.HazardEvent, s StructureDeterministic) ([]consequences.Result, error){}`
    - this function implements the logic to compute losses and reconstruction across a series of Hazards based on the implementation in g2crm
    - Downside is that the use of slices means we can't go through the `Compute()` method on the structure

- `computeConsequencesMultiHazard(event hazards.MultiHazardEvent, s StructureDeterministic) (consequences.Result, error)`
    - this function implements the same logic as the function above, but takes a new  `MultiHazardEvent` interface. 
    - the `MultiHazardEvent` interface also satisfies the `HazardEvent` interface so both types can be passed to `Compute()`

- added logic to Compute to call `computeConsequencesMultiHazard` if the `HazardEvent` paramter can be asserted as a `MultiHazardEvent` 
    - what happens if we pass `MultiHazardEvent` to the `Compute` method of a different receptor that we haven't provided multi-hazard logic for?
    - It should work fine, but Compute will only run on the first `HazardEvent` in the `MultiHazardEvent`.

## G2CRM Implementation - Structure.cs

### Rebuild

- param "rebuildFactor" - amount of structure to rebuild. 
    - range: 0.0-1.0
    - rebuildFactor < 1.0 ==> partial rebuild due to storm occuring during rebuild
    - e.g. rebuildFactor = 0.75 means rebuild 75% of the structure damage
        - **Where does this comefrom?**

- structureAmountToRebuild = CalculateStructureDepreciatedReplacementValue(year) - CalculateCurrentStructureValue(year)
    - = CSDRV(year) - CCSV(year)
    - = CSDRV(year) - (CSDRV(year) * (1 - CurrentStructureDamageFactor)
    - **structureAmountToRebuild = Structure Value * pct_damage**
    - ==> So rebuild the damage that occurred from the event

- `if (rebuildFactor < 1.0) { structureAmountToRebuild *= rebuildFactor }`
    - When function is called we know rebuildFactor. e.g. We know we want to repair 75% of the damage.
    - but why is it only partial rebuild? where is rebuildFactor calculated? What does this represent?
    - 

### Damage

- structureModifier ==> pct_damage
    - From damage function triangular distribution (`G2CRM.Core.Math.TriangularDistribution.triangularDegen()`)

- structureAmountToDamage = CalculateCurrentStructureValue(year) * structureModifier
    - ==> like go-consequences `sval * sdampercent`

## 3. Write results to file

- Assumption: we don't want to write several hundred columns to store results for every HazardEvent in a MultiHazardEvent
    - Depending on the number of hazards and variables tracked, it may not be possible: "The default setting for SQLITE_MAX_COLUMN is 2000" (https://sqlite.org/limits.html)

### Proposal 1: dedicate 1 column in the `consequences.Result` to results for specific hazard events
    - Result.Header[i] = "hazard results"
    - subResult = consequences.Result{}
    - Result.Results[i] = subResult
        - for event in events:
            1. create a new result with the fields we want to report for each event
            2. subResult.Header = append(subResult.Header, fmt.Sprintf("%d", event.Index())
            3. subResult.Results = append(subResult.Results, consequences.Result{Header: subHeader, Results: event_results}

Pros:
    1. This works and we can utilize the Fetch() method on the nested results to get results for specific events

Cons:
    1. the string headers for the results of individual events are repeated each time wasting space
    2. cannot save the results to gpkg unless the entire result column is converted to a string which removes our ability to parse as structs.
        - instead of just casting as a string, should I be using the MarshalJSON method?


### Proposal 2: Use consequences.Results{IsTable: true}
    - 



## General Brainstorming

- In `compute-configuration.go` the `Computable` struct includes a `ComputeLifeloss` bool. We could add another bool for `ComputeReconstruction` that would enable that calculation without making it a default.

# TODO today

1. [x] json_multi_hazard_test.go - replace printline with asserts that check expected values
2. [x] flood_test.go - when checking if event has given parameters, also assert that the values are as expected
3. [x] make a csv hazard provider that simply provides a constant arrival, depth, and duration for all locations. 

