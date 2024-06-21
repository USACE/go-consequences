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
func (c ComplianceBasedWarningSystem) WarningFunction() PopulationReductionFunction { //@TODO:Consider housing groups instead of simply assuming each individual is independent
	return func(s structures.StructureDeterministic, hazard hazards.HazardEvent) (structures.PopulationSet, consequences.Result) {
		var remainingpop2amo65 int32 = 0
		var remainingpop2amu65 int32 = 0
		var remainingpop2pmo65 int32 = 0
		var remainingpop2pmu65 int32 = 0
		for i := 0; i < int(s.Pop2amo65); i++ {
			if c.rng.Float64() > c.ComplianceRate {
				remainingpop2amo65 += 1
			}
		}
		for i := 0; i < int(s.Pop2amu65); i++ {
			if c.rng.Float64() > c.ComplianceRate {
				remainingpop2amu65 += 1
			}
		}
		for i := 0; i < int(s.Pop2pmo65); i++ {
			if c.rng.Float64() > c.ComplianceRate {
				remainingpop2pmo65 += 1
			}
		}
		for i := 0; i < int(s.Pop2pmu65); i++ {
			if c.rng.Float64() > c.ComplianceRate {
				remainingpop2pmu65 += 1
			}
		}
		remainingPopulation := structures.PopulationSet{
			Pop2pmo65: int32(remainingpop2amo65),
			Pop2pmu65: int32(remainingpop2amu65),
			Pop2amo65: int32(remainingpop2pmo65),
			Pop2amu65: int32(remainingpop2pmu65),
		}
		return remainingPopulation, consequences.Result{Headers: []string{"rem2amo65", "rem2amu65", "rem2pmo65", "rem2pmu65"}, Result: []interface{}{remainingpop2amo65, remainingpop2amu65, remainingpop2pmo65, remainingpop2pmu65}}
	}
}
