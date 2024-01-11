package warning

import (
	"math/rand"

	"github.com/USACE/go-consequences/consequences"
	"github.com/USACE/go-consequences/hazards"
	"github.com/USACE/go-consequences/structures"
)

type PopulationReductionFunction func(s structures.StructureDeterministic, hazard hazards.HazardEvent) (structures.PopulationSet, consequences.Result)

type WarningResponseSystem interface {
	WarningFunction() PopulationReductionFunction
}
type ComplianceBasedWarningSystem struct {
	rng            *rand.Rand
	ComplianceRate float64
}

func InitComplianceBasedWarningSystem(seed int64, complianceRate float64) ComplianceBasedWarningSystem {
	return ComplianceBasedWarningSystem{rng: rand.New(rand.NewSource(seed)), ComplianceRate: complianceRate}
}
func (c ComplianceBasedWarningSystem) WarningFunction() PopulationReductionFunction {
	return func(s structures.StructureDeterministic, hazard hazards.HazardEvent) (structures.PopulationSet, consequences.Result) {
		pop2amo65 := 0
		pop2amu65 := 0
		pop2pmo65 := 0
		pop2pmu65 := 0
		for i := 0; i < int(s.Pop2amo65); i++ {
			if c.rng.Float64() > c.ComplianceRate {
				pop2amo65 += 1
			}
		}
		for i := 0; i < int(s.Pop2amu65); i++ {
			if c.rng.Float64() > c.ComplianceRate {
				pop2amu65 += 1
			}
		}
		for i := 0; i < int(s.Pop2pmo65); i++ {
			if c.rng.Float64() > c.ComplianceRate {
				pop2pmo65 += 1
			}
		}
		for i := 0; i < int(s.Pop2pmu65); i++ {
			if c.rng.Float64() > c.ComplianceRate {
				pop2pmu65 += 1
			}
		}
		warnedPopulation := structures.PopulationSet{
			Pop2pmo65: int32(pop2amo65),
			Pop2pmu65: int32(pop2amu65),
			Pop2amo65: int32(pop2pmo65),
			Pop2amu65: int32(pop2pmu65),
		}
		return warnedPopulation, consequences.Result{Headers: []string{"rem2amo65", "rem2amu65", "rem2pmo65", "rem2pmu65"}, Result: []interface{}{pop2amo65, pop2amu65, pop2pmo65, pop2pmu65}}
	}
}
