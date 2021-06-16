package hazardproviders

import (
	"fmt"
	"log"

	"github.com/USACE/go-consequences/geography"
	"github.com/dewberry/gdal"
)

type cogReader struct {
	FilePath         string
	ds               *gdal.Dataset
	nodata           float64
	verticalIsMeters bool //default false
}

func initCR_Meters(fp string) cogReader {
	cr := initCR(fp)
	cr.verticalIsMeters = true
	return cr
}

//init creates and produces an unexported cogReader
func initCR(fp string) cogReader {
	//read the file path
	//make sure it is a tif
	fmt.Println("Connecting to: " + fp)
	ds, err := gdal.Open(fp, gdal.ReadOnly)
	if err != nil {
		log.Fatalln("Cannot connect to raster.  Killing everything! " + err.Error())
	}
	v, valid := ds.RasterBand(1).NoDataValue()
	cr := cogReader{FilePath: fp, ds: &ds, verticalIsMeters: false}
	if valid {
		cr.nodata = v
	}
	return cr
}
func (cr *cogReader) Close() {
	cr.ds.Close()
}
func (cr *cogReader) ProvideValue(l geography.Location) (float64, error) {
	rb := cr.ds.RasterBand(1)
	igt := cr.ds.InvGeoTransform()
	px := int(igt[0] + l.X*igt[1] + l.Y*igt[2])
	py := int(igt[3] + l.X*igt[4] + l.Y*igt[5])
	buffer := make([]float32, 1*1)
	rb.IO(gdal.Read, px, py, 1, 1, buffer, 1, 1, 0, 0)
	depth := buffer[0]
	d := float64(depth)
	if d == cr.nodata {
		return cr.nodata, NoDataHazardError{Input: fmt.Sprintf("COG reader had the no data value observed, setting to %v", cr.nodata)}
	}
	if cr.verticalIsMeters {
		d *= 3.28084
	}
	return d, nil
}
func (cr *cogReader) GetBoundingBox() (geography.BBox, error) {
	bbox := make([]float64, 4)
	gt := cr.ds.GeoTransform()
	dx := cr.ds.RasterXSize()
	dy := cr.ds.RasterYSize()
	bbox[0] = gt[0]                     //upper left x
	bbox[1] = gt[3]                     //upper left y
	bbox[2] = gt[0] + gt[1]*float64(dx) //lower right x
	bbox[3] = gt[3] + gt[5]*float64(dy) //lower right y
	return geography.BBox{Bbox: bbox}, nil
}
