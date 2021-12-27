package structureprovider

import (
	"errors"

	"github.com/USACE/go-consequences/consequences"
	"github.com/USACE/go-consequences/geography"
	"github.com/USACE/go-consequences/structures"
	"github.com/dewberry/gdal"
)

type shpDataSet struct {
	FilePath          string
	LayerName         string
	schemaIDX         []int
	optionalSchemaIDX []int
	ds                *gdal.DataSource
	OccTypeProvider   structures.OccupancyTypeProvider
}

func InitSHP(filepath string) (shpDataSet, error) {
	shp, err := initialize(filepath)
	shp.setOcctypeProvider(false, "")
	return shp, err
}
func InitSHPwithOcctypeFile(filepath string, occtypefp string) (shpDataSet, error) {
	shp, err := initialize(filepath)
	shp.setOcctypeProvider(true, occtypefp)
	return shp, err
}

func initialize(filepath string) (shpDataSet, error) {
	ds := gdal.OpenDataSource(filepath, int(gdal.ReadOnly))
	if ds.LayerCount() > 1 {
		return shpDataSet{}, errors.New("Shapefile at path " + filepath + "Found more than one layer please specify one layer.")
	}

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
	o := OptionalSchema()
	oIDX := make([]int, len(o))
	for i, f := range o {
		idx := def.FieldIndex(f)
		oIDX[i] = idx
	}
	shp := shpDataSet{FilePath: filepath, LayerName: l.Name(), schemaIDX: sIDX, optionalSchemaIDX: oIDX, ds: &ds}
	return shp, nil
}
func (shp *shpDataSet) setOcctypeProvider(useFilepath bool, filepath string) {
	if useFilepath {
		otp := structures.JsonOccupancyTypeProvider{}
		otp.InitLocalPath(filepath)
		shp.OccTypeProvider = otp
	} else {
		otp := structures.JsonOccupancyTypeProvider{}
		otp.InitDefault()
		shp.OccTypeProvider = otp
	}
}

//ByFips a streaming service for structure stochastic based on a bounding box
func (shp shpDataSet) ByFips(fipscode string, sp consequences.StreamProcessor) {
	shp.processFipsStream(fipscode, sp)
}
func (shp shpDataSet) processFipsStream(fipscode string, sp consequences.StreamProcessor) {
	m := shp.OccTypeProvider.OccupancyTypeMap()
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
					s, err := featuretoStructure(f, m, defaultOcctype, shp.schemaIDX, shp.optionalSchemaIDX)
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
	m := shp.OccTypeProvider.OccupancyTypeMap()
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
			s, err := featuretoStructure(f, m, defaultOcctype, shp.schemaIDX, shp.optionalSchemaIDX)
			if err == nil {
				sp(s)
			}
		}
	}
}
