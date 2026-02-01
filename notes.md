# Reconcstruction/Life Cycle Implementation

## 1. Implement basic reconstruction as a component damage function

- current component damage functions exist for "structure" and "contents"

- the existing format can be leveraged to return a reconstruction time in days for a given input
    - rather than using the hazard intensity (e.g. flood depth) to return a percent loss, we use the pct_loss as the input. Return value is number of days to return to 0 damage. e.g.:
        | damcat | dmg_pct | t_rebuild_min | t_rebuild_mostlikely | t_rebuild_max |
        |---------|-----------|----------------|------------------------|-----------------|
        | RES | 0 | 0 | 0 | 0 |
        | RES | .5 | 13.5 | 15 | 16.5 |
        | RES | 1 | 27 | 30 | 33 |

- Could add this functionality to computeConsequences() func, but this could break existing workflows (e.g. user has custom occtypes.json)
    - however, if we create a separate computeConsequencesWithReconstruction() func, there is no pathway for the receptor's Compute() method to call that func
    - commented out code shows a prototype of this where computeConsequencesWithReconstruction() was called for hazard events with arrival and duration, but I think it is reasonable that a user may want reconstruction time for a standard depth event.


## 2. Implement MultiHazard

### 2.1 Unit testing functionality

- test capability to calculate consequence and reconstruction for a generic series of events using existing hazard type

- What happens when event occurs while structure is still in reconstruction?
    - **Is there literature on this?**
    - Proposal: Time to reconstruction unnaffected by current value.
        - Structure components that were damaged in the previous event but have not been replaced will still require replacement after the next event.
        - 

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



## General Brainstorming

- In `compute-configuration.go` the `Computable` struct includes a `ComputeLifeloss` bool. We could add another bool for `ComputeReconstruction` that would enable that calculation without making it a default.


