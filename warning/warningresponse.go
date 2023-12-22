package warning

import (
	"math/rand"

	"github.com/USACE/go-consequences/structures"
)

type PopulationReductionFunction func(s *structures.StructureDeterministic)

type WarningResponseSystem interface {
	WarningFunction() PopulationReductionFunction
}
type ComplianceBasedWarningSystem struct {
	ComplianceRate float64
}

func (c ComplianceBasedWarningSystem) WarningFunction() PopulationReductionFunction {
	return func(s *structures.StructureDeterministic) {
		rng := rand.New(rand.NewSource(12345)) //TODO:Fix this
		pop2amo65 := 0
		pop2amu65 := 0
		pop2pmo65 := 0
		pop2pmu65 := 0
		for i := 0; i < int(s.Pop2amo65); i++ {
			if rng.Float64() > c.ComplianceRate {
				pop2amo65 += 1
			}
		}
		for i := 0; i < int(s.Pop2amu65); i++ {
			if rng.Float64() > c.ComplianceRate {
				pop2amu65 += 1
			}
		}
		for i := 0; i < int(s.Pop2pmo65); i++ {
			if rng.Float64() > c.ComplianceRate {
				pop2pmo65 += 1
			}
		}
		for i := 0; i < int(s.Pop2pmu65); i++ {
			if rng.Float64() > c.ComplianceRate {
				pop2pmu65 += 1
			}
		}
		s.Pop2amo65 = int32(pop2amo65)
		s.Pop2amu65 = int32(pop2amu65)
		s.Pop2pmo65 = int32(pop2pmo65)
		s.Pop2pmu65 = int32(pop2pmu65)
	}
}
