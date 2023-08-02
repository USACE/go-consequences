package structureprovider

import (
	"errors"
	"strings"

	"github.com/USACE/go-consequences/consequences"
	"github.com/USACE/go-consequences/geography"
)

const (
	Provider_NSI   string = "Provider_NSI"
	Provider_Local        = "Provider_Local"
)

// StructureProvider can be a gpkDataSet, nsiStreamProvider, or shpDataSet
type StructureProvider interface {
	// TODO regulate more common methods
	ByFips(fipscode string, sp consequences.StreamProcessor)
	ByBbox(bbox geography.BBox, sp consequences.StreamProcessor)
}

type StructureProviderInfo struct {
	StructureProviderType *string // Provider_NSI or Provider_Local
	StructureFilePath     *string // Required if StructureProviderType == Provider_Local
	OccTypeFilePath       *string // optional
	LayerName             *string // required if specified a geopackage StructureFilePath
}

// NewStructureProvider generates a structure provider
func NewStructureProvider(spi StructureProviderInfo) (StructureProvider, error) {
	var p StructureProvider
	var err error
	switch *spi.StructureProviderType {
	case Provider_NSI: // nsi

		if *spi.OccTypeFilePath != "" {
			p = InitNSISP()
		} else {
			p = InitNSISPwithOcctypeFilePath(*spi.OccTypeFilePath)
		}

	case Provider_Local:

		if strings.Contains(*spi.StructureFilePath, ".shp") { // shapefile
			if *spi.OccTypeFilePath != "" {
				p, err = InitSHP(*spi.StructureFilePath)
			} else {
				p, err = InitSHPwithOcctypeFile(*spi.StructureFilePath, *spi.OccTypeFilePath)
			}

		} else if strings.Contains(*spi.StructureFilePath, ".gpkg") { // geopackage file
			if *spi.LayerName != "" {
				return nil, errors.New("NewStructureProvider - LayerName must be specified in StructureProviderInfo for geopackage provider")
			}
			if *spi.OccTypeFilePath != "" {
				p, err = InitGPK(*spi.StructureFilePath, *spi.LayerName)
			} else {
				p, err = InitGPKwithOcctypePath(*spi.StructureFilePath, *spi.LayerName, *spi.OccTypeFilePath)
			}

		} else {
			return nil, errors.New("NewStructureProvider - unable to generate new structure provider for type: " + *spi.StructureProviderType)
		}

	default:
		return nil, errors.New("NewStructureProvider - unable to generate new structure provider for type: " + *spi.StructureProviderType)
	}
	return p, err
}
