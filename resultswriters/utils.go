package resultswriters

import (
	"reflect"

	"github.com/dewberry/gdal"
)

var gdalTypes map[reflect.Kind]gdal.FieldType = map[reflect.Kind]gdal.FieldType{
	reflect.Float32: gdal.FT_Real,
	reflect.Float64: gdal.FT_Real,
	reflect.Int32:   gdal.FT_Integer,
	reflect.String:  gdal.FT_String,
}
