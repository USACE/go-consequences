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

## 3. Implement LifeCycle consequence calculation

- Most likely as StreamAbstractLifecycle()

## General Brainstorming

- In `compute-configuration.go` the `Computable` struct includes a `ComputeLifeloss` bool. We could add another bool for `ComputeReconstruction` that would enable that calculation without making it a default.


