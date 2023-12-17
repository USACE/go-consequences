package structureprovider

import (
	"errors"

	"github.com/USACE/go-consequences/consequences"
	"github.com/USACE/go-consequences/geography"
)

type StructureProviderType int

const (
	Unknown StructureProviderType = 0 //0
	NSIAPI  StructureProviderType = 1 //1
	GPKG    StructureProviderType = 2 //2
	SHP     StructureProviderType = 3
)

func (spt StructureProviderType) String() string {
	switch spt {
	case NSIAPI:
		return "the NSI API"
	case SHP:
		return "a local Shapefile"
	case GPKG:
		return "a local geopackage"
	case Unknown:
		return "an unknown structure provider type"

	default:
		return "an unspecified structure provider type"
	}
}

// StructureProvider can be a gpkDataSet, nsiStreamProvider, or shpDataSet
type StructureProvider interface {
	// TODO regulate more common methods
	ByFips(fipscode string, sp consequences.StreamProcessor)
	ByBbox(bbox geography.BBox, sp consequences.StreamProcessor)
}

type StructureProviderInfo struct {
	StructureProviderType StructureProviderType `json:"structure_provider_type"`     // Provider_NSI or Provider_Local
	StructureFilePath     string                `json:"structure_file_path"`         // Required if StructureProviderType == Provider_Local
	OccTypeFilePath       string                `json:"occtype_file_path,omitempty"` // optional
	LayerName             string                `json:"layername,omitempty"`         // required if specified a geopackage StructureFilePath
}

// NewStructureProvider generates a structure provider
func NewStructureProvider(spi StructureProviderInfo) (StructureProvider, error) {
	var p StructureProvider
	var err error
	switch spi.StructureProviderType {
	case NSIAPI: // nsi

		if spi.OccTypeFilePath != "" {
			p = InitNSISP()
		} else {
			p = InitNSISPwithOcctypeFilePath(spi.OccTypeFilePath)
		}

	case SHP:
		if spi.OccTypeFilePath != "" {
			p, err = InitSHP(spi.StructureFilePath)
		} else {
			p, err = InitSHPwithOcctypeFile(spi.StructureFilePath, spi.OccTypeFilePath)
		}
	case GPKG:
		if spi.LayerName != "" {
			return nil, errors.New("NewStructureProvider - LayerName must be specified in StructureProviderInfo for geopackage provider")
		}
		if spi.OccTypeFilePath != "" {
			p, err = InitGPK(spi.StructureFilePath, spi.LayerName)
		} else {
			p, err = InitGPKwithOcctypePath(spi.StructureFilePath, spi.LayerName, spi.OccTypeFilePath)
		}
	case Unknown:
		return nil, errors.New("NewStructureProvider - unable to generate new structure provider from " + spi.StructureProviderType.String())

	default:
		return nil, errors.New("NewStructureProvider - unable to generate new structure provider for from " + spi.StructureProviderType.String())
	}
	return p, err
}
