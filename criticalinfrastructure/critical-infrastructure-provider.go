package criticalinfrastructure

import (
	"errors"
	"log"

	"github.com/USACE/go-consequences/consequences"
	"github.com/USACE/go-consequences/geography"
	"github.com/dewberry/gdal"
)

type gdalDataSet struct {
	FilePath  string
	LayerName string
	ds        *gdal.DataSource
	schemaIDX []int
}

func InitCriticalInfrastructureProvider(filepath string, layername string, driver string) (*gdalDataSet, error) {
	//validation?
	gpk, err := initalizeprovider(filepath, layername, driver)
	return &gpk, err
}

func initalizeprovider(filepath string, layername string, driver string) (gdalDataSet, error) {
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
	s := header[:len(header)-1]
	sIDX := make([]int, len(s))
	for i, f := range s {
		idx := def.FieldIndex(f)
		if idx < 0 {
			return gdalDataSet{}, errors.New("gdal dataset at path " + filepath + " Expected field named " + f + " none was found")
		}
		sIDX[i] = idx
	}
	gpk := gdalDataSet{FilePath: filepath, LayerName: layername, schemaIDX: sIDX, ds: &ds}
	return gpk, nil
}

// StreamByFips a streaming service for structure stochastic based on a bounding box
func (gpk gdalDataSet) ByFips(fipscode string, sp consequences.StreamProcessor) {
	gpk.processFipsStream(fipscode, sp)
}
func (gpk gdalDataSet) processFipsStream(fipscode string, sp consequences.StreamProcessor) {
	log.Fatal("no fips codes provided on critical infrastructure")
}
func (gpk gdalDataSet) ByBbox(bbox geography.BBox, sp consequences.StreamProcessor) {

	gpk.processBboxStream(bbox, sp)
}
func (gpk gdalDataSet) processBboxStream(bbox geography.BBox, sp consequences.StreamProcessor) {

	idx := 0
	l := gpk.ds.LayerByName(gpk.LayerName)
	l.SetSpatialFilterRect(bbox.Bbox[0], bbox.Bbox[3], bbox.Bbox[2], bbox.Bbox[1])
	fc, _ := l.FeatureCount(true)
	for idx < fc { // Iterate and fetch the records from result cursor
		f := l.NextFeature()
		idx++
		if f != nil {
			s := featuretoCI(f, gpk.schemaIDX)
			sp(s)
		}
	}
}
func featuretoCI(f *gdal.Feature, idxs []int) CriticalInfrastructureFeature {
	defer f.Destroy()
	var x = 0.0
	var y = 0.0
	g := f.Geometry()
	if g.IsNull() || g.IsEmpty() {
		x = f.FieldAsFloat64(idxs[1])
		y = f.FieldAsFloat64(idxs[2])
	} else {
		x = f.Geometry().X(0)
		y = f.Geometry().Y(0)
	}
	return CriticalInfrastructureFeature{
		Attributes: CriticalInfrastructureAttributes{
			Name:           f.FieldAsString(idxs[0]),
			DamageCategory: f.FieldAsString(idxs[3]),
			OccupancyType:  f.FieldAsString(idxs[4]),
		},
		Geometry: geography.GeoJsonGeometry{
			Type:        "point",
			Coordinates: []float64{x, y},
		},
	}
}
