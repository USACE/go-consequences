package geography

import "fmt"

type Location struct {
	X    float64
	Y    float64
	SRID string
}

type BBox struct {
	bbox []float64
}

func (bb BBox) ToString() string {
	return fmt.Sprintf("%f,%f,%f,%f,%f,%f,%f,%f,%f,%f",
		bb.bbox[0], bb.bbox[1],
		bb.bbox[2], bb.bbox[1],
		bb.bbox[2], bb.bbox[3],
		bb.bbox[0], bb.bbox[3],
		bb.bbox[0], bb.bbox[1])
}
