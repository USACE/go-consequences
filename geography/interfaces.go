package geography

import (
	"fmt"
)

type Location struct {
	X    float64
	Y    float64
	SRID string
}

type BBox struct {
	Bbox []float64
}

type GeoJsonGeometry struct {
	Type        string    `json:"type"`
	Coordinates []float64 `json:"coordinates"`
}
type Projected interface {
	SpatialReference() string
	UpdateSpatialReference(sr_wkt string)
}

func (bb BBox) ToString() string {
	return fmt.Sprintf("%f,%f,%f,%f,%f,%f,%f,%f,%f,%f",
		bb.Bbox[0], bb.Bbox[1],
		bb.Bbox[2], bb.Bbox[1],
		bb.Bbox[2], bb.Bbox[3],
		bb.Bbox[0], bb.Bbox[3],
		bb.Bbox[0], bb.Bbox[1])
}
func (bb BBox) Contains(p Location) bool {
	return bb.Bbox[0] <= p.X && p.X <= bb.Bbox[2] && bb.Bbox[1] <= p.Y && p.Y <= bb.Bbox[3]
}
func (gjg GeoJsonGeometry) ToLocation() Location {
	return Location{
		X:    gjg.Coordinates[0],
		Y:    gjg.Coordinates[1],
		SRID: "",
	}
}
