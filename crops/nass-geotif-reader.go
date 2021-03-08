package crops
/*
import (
	"github.com/dewberry/gdal"
)

func getCropValue(y float64, x float64, filepath string) (float32, error) {
	ds, err := gdal.Open(filepath, gdal.ReadOnly)
	if err != nil {
		return 0.0, err
	}
	rb := ds.RasterBand(1)
	igt := ds.InvGeoTransform()
	px := int(igt[0] + y*igt[1] + x*igt[2])
	py := int(igt[3] + y*igt[4] + x*igt[5])
	buffer := make([]float32, 1*1)
	rb.IO(gdal.Read, px, py, 1, 1, buffer, 1, 1, 0, 0)
	return buffer[0], nil
}
*/
