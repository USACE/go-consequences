package crops

import (
	"fmt"
	"log"

	"github.com/dewberry/gdal"
)

type nassCropProvider struct {
	FilePath  string
	ds        *gdal.Dataset
	converter map[string]Crop
}

//Init creates and produces an unexported cogHazardProvider
func Init(fp string) nassCropProvider {
	//read the file path
	//make sure it is a tif
	ds, err := gdal.Open(fp, gdal.ReadOnly)
	if err != nil {
		log.Fatalln("Cannot connect to NASS GeoTiff.  Killing everything! " + err.Error())
	}
	m := NASSCropMap()
	return nassCropProvider{fp, &ds, m}
}
func (ncp *nassCropProvider) getCropValue(y float64, x float64) (Crop, error) {
	rb := ncp.ds.RasterBand(1)
	igt := ncp.ds.InvGeoTransform()
	px := int(igt[0] + y*igt[1] + x*igt[2])
	py := int(igt[3] + y*igt[4] + x*igt[5])
	buffer := make([]uint8, 1*1)
	rb.IO(gdal.Read, px, py, 1, 1, buffer, 1, 1, 0, 0)
	c, ok := ncp.converter[string(buffer[0])]
	if ok {
		return c, nil
	} else {
		return Crop{}, NoCropFoundError{fmt.Sprintf("requested %f, %f, %s and no crop was found.", y, x, ncp.FilePath)}
	}

}

type NoCropFoundError struct {
	Input string
}

func (c NoCropFoundError) Error() string {
	return c.Input
}
