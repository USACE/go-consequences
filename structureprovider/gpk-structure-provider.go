package structureprovider

import (
	"errors"
	"fmt"

	"github.com/USACE/go-consequences/consequences"
	"github.com/USACE/go-consequences/geography"
	"github.com/USACE/go-consequences/structures"
	"github.com/dewberry/gdal"
)

type gpkDataSet struct {
	FilePath      string
	LayerName     string
	schemaIDX     []int
	ds            *gdal.DataSource
	deterministic bool
}

func InitGPK(filepath string, layername string) (gpkDataSet, error) {
	ds := gdal.OpenDataSource(filepath, int(gdal.ReadOnly))
	//validation?
	hasNSITable := false
	for i := 0; i < ds.LayerCount(); i++ {
		if layername == ds.LayerByIndex(i).Name() {
			hasNSITable = true
		}
	}
	if !hasNSITable {
		return gpkDataSet{}, errors.New("GeoPpackage at path " + filepath + "does not have a layer titled nsi. ")
	}
	l := ds.LayerByName(layername)
	def := l.Definition()
	s := StructureSchema()
	sIDX := make([]int, len(s))
	for i, f := range s {
		idx := def.FieldIndex(f)
		if idx < 0 {
			return gpkDataSet{}, errors.New("GeoPpackage at path " + filepath + " Expected field named " + f + " none was found")
		}
		sIDX[i] = idx
	}
	return gpkDataSet{FilePath: filepath, LayerName: layername, schemaIDX: sIDX, ds: &ds}, nil
}
func (gpk *gpkDataSet) SetDeterministic(useDeterministic bool) {
	gpk.deterministic = useDeterministic
}

//StreamByFips a streaming service for structure stochastic based on a bounding box
func (gpk gpkDataSet) ByFips(fipscode string, sp consequences.StreamProcessor) {
	if gpk.deterministic {
		gpk.processFipsStreamDeterministic(fipscode, sp)
	} else {
		gpk.processFipsStream(fipscode, sp)
	}

}
func (gpk gpkDataSet) processFipsStream(fipscode string, sp consequences.StreamProcessor) {
	m := structures.OccupancyTypeMap()
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
	for idx < fc { // Iterate and fetch the records from result cursor
		f := l.NextFeature()
		idx++
		if f != nil {
			s, err := featuretoStructure(f, m, defaultOcctype, gpk.schemaIDX)
			if err == nil {
				sp(s)
			}
		}
	}
}
func (gpk gpkDataSet) processFipsStreamDeterministic(fipscode string, sp consequences.StreamProcessor) {
	m := structures.OccupancyTypeMap()
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
			s, err := featuretoDeterministicStructure(f, m2, defaultOcctype, gpk.schemaIDX)
			if err == nil {
				sp(s)
			}
		}
	}
}
func (gpk gpkDataSet) ByBbox(bbox geography.BBox, sp consequences.StreamProcessor) {
	if gpk.deterministic {
		gpk.processBboxStreamDeterministic(bbox, sp)
	} else {
		gpk.processBboxStream(bbox, sp)
	}

}
func (gpk gpkDataSet) processBboxStream(bbox geography.BBox, sp consequences.StreamProcessor) {
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
			s, err := featuretoStructure(f, m, defaultOcctype, gpk.schemaIDX)
			if err == nil {
				sp(s)
			}
		}
	}
}

func (gpk gpkDataSet) processBboxStreamDeterministic(bbox geography.BBox, sp consequences.StreamProcessor) {
	m := structures.OccupancyTypeMap()
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
			s, err := featuretoDeterministicStructure(f, m2, defaultOcctype, gpk.schemaIDX)
			if err == nil {
				sp(s)
			}
		}
	}
}
