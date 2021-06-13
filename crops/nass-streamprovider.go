package crops

import (
	"fmt"
	"strconv"

	"github.com/USACE/go-consequences/consequences"
	"github.com/USACE/go-consequences/geography"
	"github.com/dewberry/gdal"
)

type NassCropProvider struct {
	Year       string
	CropFilter map[string]Crop
}

func InitNassCropProvider(year string, cropFilter []string) NassCropProvider {
	cfilter := make(map[string]Crop, len(cropFilter))
	m := NASSCropMap()
	for _, f := range cropFilter {
		c, ok := m[f]
		if ok {
			cfilter[f] = c
		}
	}
	return NassCropProvider{Year: year, CropFilter: cfilter}
}

func (n NassCropProvider) ByFips(fipscode string, sp consequences.StreamProcessor) {
	result, err := GetCDLFileByFIPS(n.Year, fipscode)
	if err != nil {
		panic(err)
	}
	result.iterate(sp, n.CropFilter)
}
func (n NassCropProvider) ByBbox(bbox geography.BBox, sp consequences.StreamProcessor) {
	result, err := GetCDLFileByBbox(n.Year, fmt.Sprintf("%v", bbox.Bbox[0]), fmt.Sprintf("%v", bbox.Bbox[3]), fmt.Sprintf("%v", bbox.Bbox[2]), fmt.Sprintf("%v", bbox.Bbox[1]))
	if err != nil {
		panic(err)
	}
	result.iterate(sp, n.CropFilter)
}
func (n nassTiffReader) iterate(sp consequences.StreamProcessor, cfilter map[string]Crop) {
	rb := n.ds.RasterBand(1)
	nYs := rb.YSize()
	nXs := rb.XSize()
	offset := 0
	arr := make([]byte, nYs*nXs, nYs*nXs)
	err := rb.IO(gdal.Read, offset, offset, nXs, nYs, arr, nXs, nYs, 0, 0)
	if err != nil {
		panic(err)
	}
	gt := n.ds.GeoTransform()
	xval := gt[0] + (gt[1] / 2)
	yval := gt[3] + (gt[5] / 2)
	for i, b := range arr {
		if i%nXs == 0 {
			xval = gt[0] + (gt[1] / 2)
			if i != 0 {
				yval += gt[5]
			}
		}
		s := strconv.Itoa(int(b)) //not sure this is right
		c, ok := cfilter[s]
		//need to add location
		c.WithLocation(xval, yval)
		if ok {
			sp(c)
		}
		xval += gt[1]

	}
	/*
		nXBlocksize, nYBlocksize := rb.BlockSize()
		nXBlocks := (rb.XSize() + nXBlocksize - 1) / nXBlocksize
		nYBlocks := (rb.YSize() + nYBlocksize - 1) / nYBlocksize
		for iYBlock := 0; iYBlock < nYBlocks; iYBlock++ {
			for iXBlock := 0; iXBlock < nXBlocks; iXBlock++ {
				//int nXValid = 0
				//int nYValid = 0
				b := make([]byte, nXBlocksize*nYBlocksize)
				pabyData := unsafe.Pointer(&b)
				rb.ReadBlock(iXBlock, iYBlock, pabyData)
				for iY := 0; iY < iYBlock; iY++ {
					for iX := 0; iX < iXBlock; iX++ {
						s := strconv.Itoa(int(b[iX+iY*nXBlocksize])) //not sure this is right
						c, ok := n.converter[s]
						//need to add location
						if ok {
							sp(c)
						}
					}
				}
			}
		}
	*/
}
