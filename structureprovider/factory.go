package structureprovider

import (
	"errors"

	"github.com/USACE/go-consequences/consequences"
	"github.com/USACE/go-consequences/geography"
)

type StructureProviderType string

const (
	Unknown StructureProviderType = "UNKNOWN" //0
	NSIAPI  StructureProviderType = "NSIAPI"  //1
	GPKG    StructureProviderType = "GPKG"    //2
	SHP     StructureProviderType = "SHP"
	OGR     StructureProviderType = "OGR"
)

func (spt StructureProviderType) String() string {
	switch spt {
	case NSIAPI:
		return "the NSI API"
	case SHP:
		return "a local Shapefile"
	case GPKG:
		return "a local geopackage"
	case OGR:
		return "ogr dataset determined based on driver"
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
	StructureProviderDriver string                `json:"structure_provider_driver,omitempty"` // ESRI SHP, GPKG, PARQUET (OGR DRIVERS...)
	StructureProviderType   StructureProviderType `json:"structure_provider_type"`             // Provider_NSI or Provider_Local
	StructureFilePath       string                `json:"structure_file_path,omitempty"`       // Required if StructureProviderType == Provider_Local
	OccTypeFilePath         string                `json:"occtype_file_path,omitempty"`         // optional
	LayerName               string                `json:"layername,omitempty"`                 // required if specified a geopackage StructureFilePath
}

// NewStructureProvider generates a structure provider
func (spi StructureProviderInfo) CreateStructureProvider() (StructureProvider, error) {
	var p StructureProvider
	var err error
	switch spi.StructureProviderType {
	case NSIAPI: // nsi

		if len(spi.OccTypeFilePath) == 0 {
			p = InitNSISP()
		} else {
			p = InitNSISPwithOcctypeFilePath(spi.OccTypeFilePath)
		}

	case SHP:
		if len(spi.OccTypeFilePath) == 0 {
			p, err = InitSHP(spi.StructureFilePath)
		} else {
			p, err = InitSHPwithOcctypeFile(spi.StructureFilePath, spi.OccTypeFilePath)
		}
	case GPKG:
		if spi.LayerName == "" {
			return nil, errors.New("NewStructureProvider - LayerName must be specified in StructureProviderInfo for geopackage provider")
		}
		if len(spi.OccTypeFilePath) == 0 {
			p, err = InitGPK(spi.StructureFilePath, spi.LayerName)
		} else {
			p, err = InitGPKwithOcctypePath(spi.StructureFilePath, spi.LayerName, spi.OccTypeFilePath)
		}
	case OGR:
		if spi.LayerName == "" {
			return nil, errors.New("NewStructureProvider - LayerName must be specified in StructureProviderInfo for ogr provider")
		}
		if len(spi.OccTypeFilePath) == 0 {
			p, err = InitStructureProvider(spi.StructureFilePath, spi.LayerName, spi.StructureProviderDriver)
		} else {
			p, err = InitStructureProviderwithOcctypePath(spi.StructureFilePath, spi.LayerName, spi.StructureProviderDriver, spi.OccTypeFilePath)
		}
	case Unknown:
		return nil, errors.New("NewStructureProvider - unable to generate new structure provider from " + spi.StructureProviderType.String())

	default:
		return nil, errors.New("NewStructureProvider - unable to generate new structure provider for from " + spi.StructureProviderType.String())
	}
	return p, err
}
