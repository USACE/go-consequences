package resultswriters

import (
	"errors"
	"reflect"

	"github.com/USACE/go-consequences/consequences"
	"github.com/dewberry/gdal"
)

var gdalTypes map[reflect.Kind]gdal.FieldType = map[reflect.Kind]gdal.FieldType{
	reflect.Float32: gdal.FieldType(gdal.FT_Real),
	reflect.Float64: gdal.FieldType(gdal.FT_Real),
	reflect.Int32:   gdal.FieldType(gdal.FT_Integer),
	reflect.String:  gdal.FieldType(gdal.FT_String),
}

type ResultsWriterType string

const (
	Unknown ResultsWriterType = "UNKNOWN" //0
	JSON    ResultsWriterType = "JSON"    //1
	GPKG    ResultsWriterType = "GPKG"    //2
	SHP     ResultsWriterType = "SHP"
	PARQUET ResultsWriterType = "Parquet"
	OGR     ResultsWriterType = "OGR"
)

type ResultsWriterInfo struct {
	Type     ResultsWriterType `json:"results_writer_type"`
	Driver   string            `json:"results_writer_driver,omitempty"`
	FilePath string            `json:"output_file_path"`
}

func (info ResultsWriterInfo) CreateResultsWriter() (consequences.ResultsWriter, error) {
	switch info.Type {
	case JSON:
		return InitSpatialResultsWriter(info.FilePath, "results", "GeoJSON")
	case GPKG:
		return InitSpatialResultsWriter(info.FilePath, "results", "GPKG")
	case SHP:
		return InitSpatialResultsWriter(info.FilePath, "results", "ESRI Shapefile")
	case PARQUET:
		return InitSpatialResultsWriter(info.FilePath, "results", "Parquet")
	case OGR:
		return InitSpatialResultsWriter(info.FilePath, "results", info.Driver)
	default:
		return nil, errors.New("could not create a result writer of that type")
	}
}
