package compute

import (
	"log"
	"math/rand"

	"github.com/USACE/go-consequences/consequences"
	"github.com/USACE/go-consequences/geography"
	"github.com/USACE/go-consequences/hazardproviders"
	"github.com/USACE/go-consequences/hazards"
	"github.com/USACE/go-consequences/lifeloss"
	"github.com/USACE/go-consequences/resultswriters"
	"github.com/USACE/go-consequences/structureprovider"
	"github.com/USACE/go-consequences/structures"
	"github.com/USACE/go-consequences/warning"
)

//todo add ag and indrect losses?

type Config struct {
	structureprovider.StructureProviderInfo `json:"structure_provider_info"`
	hazardproviders.HazardProviderInfo      `json:"hazard_provider_info"`
	resultswriters.ResultsWriterInfo        `json:"results_writer_info"`
	ComputeLifeloss                         bool    `json:"compute_lifeloss"`
	LifelossSeed                            int64   `json:"lifeloss_seed,omitempty"`
	ComplianceRate                          float64 `json:"warning_compliance_rate"`
	ComputeByFips                           bool    `json:"compute_by_fips"`
	FipsCode                                string  `json:"fips_code"`
}
type Computeable struct {
	structureprovider.StructureProvider
	hazardproviders.HazardProvider
	consequences.ResultsWriter
	ComputeLifeloss bool
	LifelossSeed    int64
	ComplianceRate  float64
	ComputeByFips   bool
	FipsCode        string
}

func (config Config) CreateComputable() (Computeable, error) {
	sp, err := config.CreateStructureProvider()
	if err != nil {
		return Computeable{}, err
	}
	hp, err := config.CreateHazardProvider()
	if err != nil {
		return Computeable{}, err
	}
	rw, err := config.CreateResultsWriter()
	if err != nil {
		return Computeable{}, err
	}
	return Computeable{
		StructureProvider: sp,
		HazardProvider:    hp,
		ResultsWriter:     rw,
		ComputeLifeloss:   config.ComputeLifeloss,
		LifelossSeed:      config.LifelossSeed,
		ComplianceRate:    config.ComplianceRate,
		ComputeByFips:     config.ComputeByFips,
		FipsCode:          config.FipsCode,
	}, nil
}
func (computable Computeable) Compute() error {
	if computable.ComputeLifeloss {
		if computable.ComputeByFips {
			return computable.computeWithLifelossByFips(computable.HazardProvider, computable.StructureProvider, computable.ResultsWriter)
		} else {
			return computable.computeWithLifelossByBbox(computable.HazardProvider, computable.StructureProvider, computable.ResultsWriter)
		}
	} else {
		StreamAbstract(computable.HazardProvider, computable.StructureProvider, computable.ResultsWriter) //bybbox. need to add logic for fips.
		computable.ResultsWriter.Close()
	}

	return nil
}
func (computable Computeable) computeWithLifelossByFips(hp hazardproviders.HazardProvider, sp consequences.StreamProvider, w consequences.ResultsWriter) error {

	rng := rand.New(rand.NewSource(computable.LifelossSeed))
	warningSystem := warning.InitComplianceBasedWarningSystem(rng.Int63(), computable.ComplianceRate)
	lle := lifeloss.Init(rng.Int63(), warningSystem)
	//err := errors.New("error")
	sp.ByFips(computable.FipsCode, func(f consequences.Receptor) {
		err := computeLifelossPerStructure(hp, f, rng, lle, w)
		if err != nil {
			log.Println(err)
		}
	})
	return nil
}
func (computable Computeable) computeWithLifelossByBbox(hp hazardproviders.HazardProvider, sp consequences.StreamProvider, w consequences.ResultsWriter) error {

	rng := rand.New(rand.NewSource(computable.LifelossSeed))
	warningSystem := warning.InitComplianceBasedWarningSystem(rng.Int63(), computable.ComplianceRate)
	lle := lifeloss.Init(rng.Int63(), warningSystem)
	bbox, err := hp.HazardBoundary()
	if err != nil {
		return err
	}
	sp.ByBbox(bbox, func(f consequences.Receptor) {
		err := computeLifelossPerStructure(hp, f, rng, lle, w)
		if err != nil {
			log.Println(err)
		}
	})
	return nil
}
func computeLifelossPerStructure(hp hazardproviders.HazardProvider, f consequences.Receptor, rng *rand.Rand, lle lifeloss.LifeLossEngine, w consequences.ResultsWriter) error {
	//ProvideHazard works off of a geography.Location
	d, err := hp.Hazard(geography.Location{X: f.Location().X, Y: f.Location().Y})
	if err == nil {
		r := consequences.Result{}
		//compute damages based on hazard being able to provide depth
		r, err = f.Compute(d)
		if err != nil {
			return err
		}

		//cast f to structure deterministic
		var sd structures.StructureDeterministic
		ss, ssok := f.(structures.StructureStochastic)
		if ssok {
			sd = ss.SampleStructure(rng.Int63())
		} else {
			sdtemp, sdok := f.(structures.StructureDeterministic)
			if sdok {
				sd = sdtemp
			} else {
				return err
			}
		}
		//if hazard provider does not have velocity or depth*velocity add it based on fema velocity zone
		llevent := hazards.DepthandDVEvent{}
		llevent.SetDepth(d.Depth())
		llevent.SetDV(0.0)
		if d.Has(hazards.Velocity) {
			llevent.SetDV(d.Depth() * d.Velocity())
		} else if d.Has(hazards.DV) {
			llevent.SetDV(d.DV())
		} else if d.Has(hazards.WaveHeight) {
			//if waveheight>3 => VE zone
			llevent.SetDV(d.Depth() * 6.5)
		} else {
			switch sd.FirmZone {
			case "VE", "V1-30":
				llevent.SetDV(d.Depth() * 6.5)
			}
		}
		//
		//compute lifeloss
		stability, err := lle.EvaluateStabilityCriteria(llevent, sd)
		if err != nil {
			return err
		}
		llr, err := lle.ComputeLifeLoss(llevent, sd, stability)
		if err != nil {
			return err
		}
		//append results
		r.Headers = append(r.Headers, llr.Headers...)
		r.Result = append(r.Result, llr.Result...)
		w.Write(r)
		return nil
	}
	return err
}
