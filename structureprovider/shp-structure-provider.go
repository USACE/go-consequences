package structureprovider

import (
	"fmt"

	"github.com/USACE/go-consequences/geography"
	"github.com/dewberry/gdal"
)

type shpDataSet struct {
	FilePath string
	ds       *gdal.DataSource
}

func InitSHP(filepath string) shpDataSet {
	ds := gdal.OpenDataSource(filepath, int(gdal.ReadOnly))
	fmt.Println(ds.Driver().Name())
	for i := 0; i < ds.LayerCount(); i++ {
		fmt.Println(ds.LayerByIndex(i).Name())
		layer := ds.LayerByIndex(i)
		fieldDef := layer.Definition()

		for j := 0; j < fieldDef.FieldCount(); j++ {
			fieldName := fieldDef.FieldDefinition(j).Name()
			fieldType := fieldDef.FieldDefinition(j).Type().Name()
			fmt.Println(fmt.Sprintf("%s, %s", fieldName, fieldType))
		}
	}
	return shpDataSet{filepath, &ds}
}

//StreamByFips a streaming service for structure stochastic based on a bounding box
func (gpk shpDataSet) ByFips(fipscode string, sp StreamProcessor) error {
	return gpk.processFipsStream(fipscode, sp)
}
func (gpk shpDataSet) processFipsStream(fipscode string, sp StreamProcessor) error {

	return nil

}
func (gpk shpDataSet) ByBbox(bbox geography.BBox, sp StreamProcessor) error {
	return gpk.processBboxStream(bbox, sp)
}
func (gpk shpDataSet) processBboxStream(bbox geography.BBox, sp StreamProcessor) error {
	return nil

}
