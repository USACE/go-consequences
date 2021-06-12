package crops

import (
	"fmt"
	"strconv"
	"unsafe"

	"github.com/USACE/go-consequences/consequences"
	"github.com/USACE/go-consequences/geography"
)

type NassCropProvider struct {
	Year string
}

func (n NassCropProvider) ByFips(fipscode string, sp consequences.StreamProcessor) {
	result, err := GetCDLFileByFIPS(n.Year, fipscode)
	if err != nil {
		panic(err)
	}
	result.iterate(sp)
}
func (n NassCropProvider) ByBbox(bbox geography.BBox, sp consequences.StreamProcessor) {
	result, err := GetCDLFileByBbox(n.Year, fmt.Sprintf("%v", bbox.Bbox[0]), fmt.Sprintf("%v", bbox.Bbox[3]), fmt.Sprintf("%v", bbox.Bbox[2]), fmt.Sprintf("%v", bbox.Bbox[1]))
	if err != nil {
		panic(err)
	}
	result.iterate(sp)
}
func (n nassTiffReader) iterate(sp consequences.StreamProcessor) {
	rb := n.ds.RasterBand(1)
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
					if ok {
						sp(c)
					}
				}
			}
		}
	}
}
