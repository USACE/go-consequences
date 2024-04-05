package structureprovider

import (
	"errors"
	"fmt"
	"math/rand"

	"github.com/USACE/go-consequences/consequences"
	"github.com/USACE/go-consequences/geography"
	"github.com/USACE/go-consequences/structures"
	"github.com/dewberry/gdal"
)

type gdalDataSet struct {
	FilePath              string
	LayerName             string
	schemaIDX             []int
	optionalSchemaIDX     []int
	ds                    *gdal.DataSource
	deterministic         bool
	seed                  int64
	OccTypeProvider       structures.OccupancyTypeProvider
	FoundationUncertainty *structures.FoundationUncertainty
}

func InitStructureProvider(filepath string, layername string, driver string) (*gdalDataSet, error) {
	//validation?
	gpk, err := initalizestructureprovider(filepath, layername, driver)
	gpk.setOcctypeProvider(false, "")
	gpk.UpdateFoundationHeightUncertainty(false, "")
	return &gpk, err
}
func InitStructureProviderwithOcctypePath(filepath string, layername string, driver string, occtypefp string) (*gdalDataSet, error) {
	//validation?
	gpk, err := initalizestructureprovider(filepath, layername, driver)
	gpk.setOcctypeProvider(true, occtypefp)
	return &gpk, err
}
func (ds *gdalDataSet) UpdateFoundationHeightUncertainty(useFile bool, foundationHeightUncertaintyJsonFilePath string) {
	if useFile {
		fh, err := structures.InitFoundationUncertaintyFromFile(foundationHeightUncertaintyJsonFilePath)
		if err != nil {
			fh, _ = structures.InitFoundationUncertainty()
		}
		ds.FoundationUncertainty = fh
	} else {
		fh, _ := structures.InitFoundationUncertainty()
		ds.FoundationUncertainty = fh
	}
}
func initalizestructureprovider(filepath string, layername string, driver string) (gdalDataSet, error) {
	driverOut := gdal.OGRDriverByName(driver)
	ds, dsok := driverOut.Open(filepath, int(gdal.ReadOnly))
	if !dsok {
		return gdalDataSet{}, errors.New("error opening structure provider of type " + driver)
	}

	hasNSITable := false
	for i := 0; i < ds.LayerCount(); i++ {
		if layername == ds.LayerByIndex(i).Name() {
			hasNSITable = true
		}
	}
	if !hasNSITable {
		return gdalDataSet{}, errors.New("gdal dataset at path " + filepath + "does not have a layer titled " + layername + ". ")
	}
	l := ds.LayerByName(layername)
	def := l.Definition()
	s := StructureSchema()
	sIDX := make([]int, len(s))
	for i, f := range s {
		idx := def.FieldIndex(f)
		if idx < 0 {
			return gdalDataSet{}, errors.New("gdal dataset at path " + filepath + " Expected field named " + f + " none was found")
		}
		sIDX[i] = idx
	}
	o := OptionalSchema()
	oIDX := make([]int, len(o))
	for i, f := range o {
		idx := def.FieldIndex(f)
		oIDX[i] = idx
	}
	gpk := gdalDataSet{FilePath: filepath, LayerName: layername, schemaIDX: sIDX, optionalSchemaIDX: oIDX, ds: &ds, seed: 1234}
	return gpk, nil
}
func (gpk *gdalDataSet) setOcctypeProvider(useFilepath bool, filepath string) {
	if useFilepath {
		otp := structures.JsonOccupancyTypeProvider{}
		otp.InitLocalPath(filepath)
		gpk.OccTypeProvider = otp
	} else {
		otp := structures.JsonOccupancyTypeProvider{}
		otp.InitDefault()
		gpk.OccTypeProvider = otp
	}
}
func (gpk *gdalDataSet) SetDeterministic(useDeterministic bool) {
	gpk.deterministic = useDeterministic
}
func (gpk *gdalDataSet) SetSeed(seed int64) {
	gpk.seed = seed
}

// StreamByFips a streaming service for structure stochastic based on a bounding box
func (gpk gdalDataSet) ByFips(fipscode string, sp consequences.StreamProcessor) {
	if gpk.deterministic {
		gpk.processFipsStreamDeterministic(fipscode, sp)
	} else {
		gpk.processFipsStream(fipscode, sp)
	}

}
func (gpk gdalDataSet) processFipsStream(fipscode string, sp consequences.StreamProcessor) {
	m := gpk.OccTypeProvider.OccupancyTypeMap()
	//define a default occtype in case of emergancy
	defaultOcctype := m["RES1-1SNB"]
	idx := 0
	l := gpk.ds.LayerByName(gpk.LayerName)
	fdef := l.Definition().FieldDefinition(gpk.schemaIDX[1])
	filterstring := "SUBSTR(" + fdef.Name() + ",1," + fmt.Sprint(len(fipscode)) + ") = '" + fipscode + "'"
	err := l.SetAttributeFilter(filterstring)
	if err != nil {
		panic(err)
	}
	fc, _ := l.FeatureCount(true)
	r := rand.New(rand.NewSource(gpk.seed))
	for idx < fc { // Iterate and fetch the records from result cursor
		f := l.NextFeature()
		idx++
		if f != nil {
			s, err := featuretoStructure(f, m, defaultOcctype, gpk.schemaIDX, gpk.optionalSchemaIDX)
			s.ApplyFoundationHeightUncertanty(gpk.FoundationUncertainty)
			s.UseUncertainty = true
			sd := s.SampleStructure(r.Int63())
			if err == nil {
				sp(sd)
			}
		}
	}
}
func (gpk gdalDataSet) processFipsStreamDeterministic(fipscode string, sp consequences.StreamProcessor) {
	m := gpk.OccTypeProvider.OccupancyTypeMap()
	m2 := swapOcctypeMap(m)
	//define a default occtype in case of emergancy
	defaultOcctype := m2["RES1-1SNB"]
	idx := 0
	l := gpk.ds.LayerByName(gpk.LayerName)
	fdef := l.Definition().FieldDefinition(gpk.schemaIDX[1])
	filterstring := "SUBSTR(" + fdef.Name() + ",1," + fmt.Sprint(len(fipscode)) + ") = '" + fipscode + "'"
	err := l.SetAttributeFilter(filterstring)
	if err != nil {
		panic(err)
	}
	fc, _ := l.FeatureCount(true)
	for idx < fc { // Iterate and fetch the records from result cursor
		f := l.NextFeature()
		idx++
		if f != nil {
			s, err := featuretoDeterministicStructure(f, m2, defaultOcctype, gpk.schemaIDX, gpk.optionalSchemaIDX)
			if err == nil {
				sp(s)
			}
		}
	}
}
func (gpk gdalDataSet) ByBbox(bbox geography.BBox, sp consequences.StreamProcessor) {
	if gpk.deterministic {
		gpk.processBboxStreamDeterministic(bbox, sp)
	} else {
		gpk.processBboxStream(bbox, sp)
	}

}
func (gpk gdalDataSet) processBboxStream(bbox geography.BBox, sp consequences.StreamProcessor) {
	m := gpk.OccTypeProvider.OccupancyTypeMap()
	//define a default occtype in case of emergancy
	defaultOcctype := m["RES1-1SNB"]
	idx := 0
	l := gpk.ds.LayerByName(gpk.LayerName)
	l.SetSpatialFilterRect(bbox.Bbox[0], bbox.Bbox[3], bbox.Bbox[2], bbox.Bbox[1])
	fc, _ := l.FeatureCount(true)
	r := rand.New(rand.NewSource(gpk.seed))
	for idx < fc { // Iterate and fetch the records from result cursor
		f := l.NextFeature()
		idx++
		if f != nil {
			s, err := featuretoStructure(f, m, defaultOcctype, gpk.schemaIDX, gpk.optionalSchemaIDX)
			s.ApplyFoundationHeightUncertanty(gpk.FoundationUncertainty)
			s.UseUncertainty = true
			sd := s.SampleStructure(r.Int63())
			if err == nil {
				sp(sd)
			}
		}
	}
}

func (gpk gdalDataSet) processBboxStreamDeterministic(bbox geography.BBox, sp consequences.StreamProcessor) {
	m := gpk.OccTypeProvider.OccupancyTypeMap()
	m2 := swapOcctypeMap(m)
	//define a default occtype in case of emergancy
	defaultOcctype := m2["RES1-1SNB"]
	idx := 0
	l := gpk.ds.LayerByName(gpk.LayerName)
	l.SetSpatialFilterRect(bbox.Bbox[0], bbox.Bbox[3], bbox.Bbox[2], bbox.Bbox[1])
	fc, _ := l.FeatureCount(true)
	for idx < fc { // Iterate and fetch the records from result cursor
		f := l.NextFeature()
		idx++
		if f != nil {
			s, err := featuretoDeterministicStructure(f, m2, defaultOcctype, gpk.schemaIDX, gpk.optionalSchemaIDX)
			if err == nil {
				sp(s)
			}
		}
	}
}
