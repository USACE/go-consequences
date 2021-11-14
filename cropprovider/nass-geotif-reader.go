package cropprovider

import (
	"fmt"
	"log"
	"strconv"

	"github.com/USACE/go-consequences/crops"
	"github.com/dewberry/gdal"
)

type nassTiffReader struct {
	FilePath  string
	ds        *gdal.Dataset
	converter map[string]crops.Crop
}

// Init creates a streaming crop provider
func Init(fp string) nassTiffReader {
	//read the file path
	//make sure it is a tif
	ds, err := gdal.Open(fp, gdal.ReadOnly)
	if err != nil {
		log.Fatalln("Cannot connect to NASS GeoTiff.  Killing everything! " + err.Error())
	}
	//m := NASSCropMap()
	spatialRef := gdal.CreateSpatialReference("")
	spatialRef.FromEPSG(5070)
	srString, err := spatialRef.ToWKT()
	if err != nil {
		panic(err)
	}
	ds.SetProjection(srString)
	return nassTiffReader{fp, &ds, nil}
}
func (ncp *nassTiffReader) getCropValue(y float64, x float64) (crops.Crop, error) {
	rb := ncp.ds.RasterBand(1)
	igt := ncp.ds.InvGeoTransform()
	px := int(igt[0] + y*igt[1] + x*igt[2])
	py := int(igt[3] + y*igt[4] + x*igt[5])
	buffer := make([]uint8, 1*1)
	rb.IO(gdal.Read, px, py, 1, 1, buffer, 1, 1, 0, 0)
	s := strconv.Itoa(int(buffer[0]))
	if ncp.converter == nil {
		ncp.converter = crops.NASSCropMap()
	}
	c, ok := ncp.converter[s]

	if ok {
		return c, nil
	} else {
		return crops.Crop{}, NoCropFoundError{fmt.Sprintf("requested %f, %f, %s and no crop was found.", y, x, ncp.FilePath)}
	}

}

type NoCropFoundError struct {
	Input string
}

func (c NoCropFoundError) Error() string {
	return c.Input
}
