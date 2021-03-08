package geography

import "fmt"

type Location struct {
	X    float64
	Y    float64
	SRID string
}

type BBox struct {
	Bbox []float64
}

func (bb BBox) ToString() string {
	return fmt.Sprintf("%f,%f,%f,%f,%f,%f,%f,%f,%f,%f",
		bb.Bbox[0], bb.Bbox[1],
		bb.Bbox[2], bb.Bbox[1],
		bb.Bbox[2], bb.Bbox[3],
		bb.Bbox[0], bb.Bbox[3],
		bb.Bbox[0], bb.Bbox[1])
}
