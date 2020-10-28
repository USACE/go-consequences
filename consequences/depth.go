package consequences

/*
import (
	"github.com/dewberry/gdal"
)

func getDepth(lon float64, lat float64, depthGrid string) (float32, error) {
	ds, err := gdal.Open(depthGrid, gdal.ReadOnly)
	if err != nil {
		return 0.0, err
	}
	rb := ds.RasterBand(1)
	igt := ds.InvGeoTransform()
	px := int(igt[0] + lon*igt[1] + lat*igt[2])
	py := int(igt[3] + lon*igt[4] + lat*igt[5])
	buffer := make([]float32, 1*1)
	rb.IO(gdal.Read, px, py, 1, 1, buffer, 1, 1, 0, 0)
	return buffer[0], nil
}
*/
