package hazardproviders

import (
	"log"

	"github.com/USACE/go-consequences/geography"
	"github.com/USACE/go-consequences/hazards"
	"github.com/dewberry/gdal"
)

type cogHazardProvider struct {
	FilePath string
	ds       *gdal.Dataset
}

//Init creates and produces an unexported cogHazardProvider
func Init(fp string) cogHazardProvider {
	//read the file path
	//make sure it is a tif
	ds, err := gdal.Open(fp, gdal.ReadOnly)
	if err != nil {
		log.Fatalln("Cannot connect to raster.  Killing everything! " + err.Error())
	}
	return cogHazardProvider{fp, &ds}
}
func (chp *cogHazardProvider) Close() {
	chp.ds.Close()
}
func (chp *cogHazardProvider) ProvideHazard(l geography.Location) (hazards.HazardEvent, error) {
	rb := chp.ds.RasterBand(1)
	igt := chp.ds.InvGeoTransform()
	px := int(igt[0] + l.X*igt[1] + l.Y*igt[2])
	py := int(igt[3] + l.X*igt[4] + l.Y*igt[5])
	buffer := make([]float32, 1*1)
	rb.IO(gdal.Read, px, py, 1, 1, buffer, 1, 1, 0, 0)
	depth := buffer[0]
	//fmt.Println(depth)
	h := hazards.DepthEvent{}
	h.SetDepth(float64(depth))
	//fmt.Println(h.Depth())
	return h, nil
}
func (chp *cogHazardProvider) ProvideHazardBoundary() (geography.BBox, error) {
	bbox := make([]float64, 4)
	gt := chp.ds.GeoTransform()
	dx := chp.ds.RasterXSize()
	dy := chp.ds.RasterYSize()
	bbox[0] = gt[0]                     //upper left x
	bbox[1] = gt[3]                     //upper left y
	bbox[2] = gt[0] + gt[1]*float64(dx) //lower right x
	bbox[3] = gt[3] + gt[5]*float64(dy) //lower right y
	return geography.BBox{Bbox: bbox}, nil
}
