package compute

import (
	"github.com/USACE/go-consequences/hazardproviders"
	"github.com/USACE/go-consequences/resultswriters"
	"github.com/USACE/go-consequences/structureprovider"
)

type Config struct {
	structureprovider.StructureProviderInfo `json:"structure_provider_info"`
	hazardproviders.HazardProviderInfo      `json:"hazard_provider_info"`
	resultswriters.ResultsWriterInfo        `json:"results_writer_info"`
}
