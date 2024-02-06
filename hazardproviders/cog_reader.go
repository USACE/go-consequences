package hazardproviders

import (
	"errors"
	"fmt"

	"github.com/USACE/go-consequences/geography"
	"github.com/dewberry/gdal"
)

type cogReader struct {
	FilePath         string
	ds               *gdal.Dataset
	nodata           float64
	verticalIsMeters bool //default false
	rb               gdal.RasterBand
	igt              [6]float64
}

func initCR_Meters(fp string) (cogReader, error) {
	cr, err := initCR(fp)
	cr.verticalIsMeters = true
	return cr, err
}

// init creates and produces an unexported cogReader
func initCR(fp string) (cogReader, error) {
	//read the file path
	//make sure it is a tif
	fmt.Println("Connecting to: " + fp)
	ds, err := gdal.Open(fp, gdal.Access(gdal.ReadOnly))
	if err != nil {
		return cogReader{}, errors.New("Cannot connect to raster at path " + fp + err.Error())
	}
	rb := ds.RasterBand(1)
	igt := ds.InvGeoTransform()
	v, valid := rb.NoDataValue()
	cr := cogReader{
		FilePath:         fp,
		ds:               &ds,
		nodata:           -9999,
		verticalIsMeters: false,
		rb:               rb,
		igt:              igt,
	}
	if valid {
		cr.nodata = v
	}
	return cr, nil
}
func (cr *cogReader) Close() {
	cr.ds.Close()
}
func (cr *cogReader) ProvideValue(l geography.Location) (float64, error) {
	igt := cr.igt
	px := int(igt[0] + l.X*igt[1] + l.Y*igt[2])
	py := int(igt[3] + l.X*igt[4] + l.Y*igt[5])
	buffer := make([]float32, 1*1)
	if px < 0 || px > cr.rb.XSize() {
		return cr.nodata, NoDataHazardError{Input: "X is out of range"}
	}
	if py < 0 || py > cr.rb.YSize() {
		return cr.nodata, NoDataHazardError{Input: "Y is out of range"}
	}
	err := cr.rb.IO(gdal.RWFlag(gdal.Read), px, py, 1, 1, buffer, 1, 1, 0, 0)
	if err != nil {
		return cr.nodata, NoDataHazardError{Input: err.Error()}
	}
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
