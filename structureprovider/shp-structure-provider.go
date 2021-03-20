package structureprovider

import (
	"fmt"
	"log"

	"github.com/USACE/go-consequences/geography"
	"github.com/dewberry/gdal"
)

type shpDataSet struct {
	FilePath string
	ds       *gdal.Dataset
}

func InitSHP(filepath string) shpDataSet {
	ds, err := gdal.Open(filepath, gdal.ReadOnly)
	if err != nil {
		log.Fatalln("Cannot connect to shapefile.  Killing everything! " + err.Error())
	}
	fmt.Println(ds.Driver().ShortName())
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
