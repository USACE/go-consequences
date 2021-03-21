package structureprovider

import (
	"log"
	"strings"

	"github.com/USACE/go-consequences/geography"
	"github.com/USACE/go-consequences/structures"
	"github.com/dewberry/gdal"
)

type gpkDataSet struct {
	FilePath  string
	LayerName string
	schemaIDX []int
	ds        *gdal.DataSource
}

func InitGPK(filepath string, layername string) gpkDataSet {
	ds := gdal.OpenDataSource(filepath, int(gdal.ReadOnly))
	//validation?
	hasNSITable := false
	for i := 0; i < ds.LayerCount(); i++ {
		if layername == ds.LayerByIndex(i).Name() {
			hasNSITable = true
		}
	}
	if !hasNSITable {
		log.Fatalln("GeoPpackage does not have a layer titled nsi.  Killing everything! ")
	}
	l := ds.LayerByName(layername)
	def := l.Definition()
	s := StructureSchema()
	sIDX := make([]int, len(s))
	for i, f := range s {
		idx := def.FieldIndex(f)
		if idx < 0 {
			log.Fatalln("Expected field named " + f + " none was found.  Killing everything! ")
		}
		sIDX[i] = idx
	}
	return gpkDataSet{FilePath: filepath, LayerName: layername, schemaIDX: sIDX, ds: &ds}
}

//StreamByFips a streaming service for structure stochastic based on a bounding box
func (gpk gpkDataSet) ByFips(fipscode string, sp StreamProcessor) error {
	return gpk.processFipsStream(fipscode, sp)
}
func (gpk gpkDataSet) processFipsStream(fipscode string, sp StreamProcessor) error {
	m := structures.OccupancyTypeMap()
	//define a default occtype in case of emergancy
	defaultOcctype := m["RES1-1SNB"]
	idx := 0
	l := gpk.ds.LayerByName(gpk.LayerName)
	fc, _ := l.FeatureCount(true)
	for idx < fc { // Iterate and fetch the records from result cursor
		f := l.NextFeature()
		idx++
		if f != nil {
			cbfips := f.FieldAsString(gpk.schemaIDX[1])
			//check if CBID matches?
			if strings.Contains(cbfips, fipscode) {
				sp(featuretoStructure(f, m, defaultOcctype, gpk.schemaIDX))
			}
		}
	}
	return nil

}
func (gpk gpkDataSet) ByBbox(bbox geography.BBox, sp StreamProcessor) error {
	return gpk.processBboxStream(bbox, sp)
}
func (gpk gpkDataSet) processBboxStream(bbox geography.BBox, sp StreamProcessor) error {
	m := structures.OccupancyTypeMap()
	//define a default occtype in case of emergancy
	defaultOcctype := m["RES1-1SNB"]
	idx := 0
	l := gpk.ds.LayerByName(gpk.LayerName)
	l.SetSpatialFilterRect(bbox.Bbox[0], bbox.Bbox[3], bbox.Bbox[2], bbox.Bbox[1])
	fc, _ := l.FeatureCount(true)
	for idx < fc { // Iterate and fetch the records from result cursor
		f := l.NextFeature()
		idx++
		if f != nil {
			sp(featuretoStructure(f, m, defaultOcctype, gpk.schemaIDX))
		}
	}
	return nil
}
