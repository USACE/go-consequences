package structureprovider

import (
	"errors"

	"github.com/USACE/go-consequences/consequences"
	"github.com/USACE/go-consequences/geography"
	"github.com/USACE/go-consequences/structures"
	"github.com/dewberry/gdal"
)

type shpDataSet struct {
	FilePath  string
	LayerName string
	schemaIDX []int
	ds        *gdal.DataSource
}

func InitSHP(filepath string) (shpDataSet, error) {
	ds := gdal.OpenDataSource(filepath, int(gdal.ReadOnly))
	if ds.LayerCount() > 1 {
		return shpDataSet{}, errors.New("Shapefile at path " + filepath + "Found more than one layer please specify one layer.")
	}
	/*if ds.LayerCount() < 1 {
		return shpDataSet{}, errors.New("Shapefile at path " + filepath + "Found no layers please specify one layer.")
	}*/
	l := ds.LayerByIndex(0)
	def := l.Definition()
	s := StructureSchema()
	sIDX := make([]int, len(s))
	for i, f := range s {
		idx := def.FieldIndex(f)
		if idx < 0 {
			return shpDataSet{}, errors.New("Shapefile at path " + filepath + " Expected field named " + f + " none was found")
		}
		sIDX[i] = idx
	}
	return shpDataSet{FilePath: filepath, LayerName: l.Name(), schemaIDX: sIDX, ds: &ds}, nil
}

//ByFips a streaming service for structure stochastic based on a bounding box
func (shp shpDataSet) ByFips(fipscode string, sp consequences.StreamProcessor) {
	shp.processFipsStream(fipscode, sp)
}
func (shp shpDataSet) processFipsStream(fipscode string, sp consequences.StreamProcessor) {
	m := structures.OccupancyTypeMap()
	//define a default occtype in case of emergancy
	defaultOcctype := m["RES1-1SNB"]
	idx := 0
	l := shp.ds.LayerByName(shp.LayerName)
	fc, _ := l.FeatureCount(true)
	for idx < fc { // Iterate and fetch the records from result cursor
		f := l.NextFeature()
		idx++
		if f != nil {
			cbfips := f.FieldAsString(shp.schemaIDX[1])
			//check if CBID matches from the start of the string
			if len(fipscode) <= len(cbfips) {
				comp := cbfips[0:len(fipscode)]
				if comp == fipscode {
					s, err := featuretoStructure(f, m, defaultOcctype, shp.schemaIDX)
					if err == nil {
						sp(s)
					}

				} //else no match, do not send structure.
			} //else error?
		}
	}
}

//ByBbox allows a shapefile to be streamed by bounding box
func (shp shpDataSet) ByBbox(bbox geography.BBox, sp consequences.StreamProcessor) {
	shp.processBboxStream(bbox, sp)
}
func (shp shpDataSet) processBboxStream(bbox geography.BBox, sp consequences.StreamProcessor) {
	m := structures.OccupancyTypeMap()
	//define a default occtype in case of emergancy
	defaultOcctype := m["RES1-1SNB"]
	idx := 0
	l := shp.ds.LayerByName(shp.LayerName)
	l.SetSpatialFilterRect(bbox.Bbox[0], bbox.Bbox[3], bbox.Bbox[2], bbox.Bbox[1])
	fc, _ := l.FeatureCount(true)
	for idx < fc { // Iterate and fetch the records from result cursor
		f := l.NextFeature()
		idx++
		if f != nil {
			s, err := featuretoStructure(f, m, defaultOcctype, shp.schemaIDX)
			if err == nil {
				sp(s)
			}
		}
	}
}
