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
)

type ResultsWriterInfo struct {
	Type     ResultsWriterType `json:"results_writer_type"`
	FilePath string            `json:"output_file_path"`
}

func (info ResultsWriterInfo) CreateResultsWriter() (consequences.ResultsWriter, error) {
	switch info.Type {
	case JSON:
		return InitGeoJsonResultsWriterFromFile(info.FilePath)
	case GPKG:
		return InitGpkResultsWriter(info.FilePath, "results")
	case SHP:
		return InitShpResultsWriter(info.FilePath, "results")
	case PARQUET:
		return InitGeoparquetResultsWriter(info.FilePath, "results")
	default:
		return nil, errors.New("could not create a result writer of that type")
	}
}
