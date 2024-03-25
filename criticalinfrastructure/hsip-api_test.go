package criticalinfrastructure

import (
	"testing"

	"github.com/USACE/go-consequences/consequences"
	"github.com/USACE/go-consequences/geography"
	"github.com/USACE/go-consequences/hazards"
)

func TestHSIP(t *testing.T) {
	list := []Layer{Hospitals, PowerPlants}
	provider := InitHsipProvider(list)
	bbox := geography.BBox{
		Bbox: []float64{-80, 35.5, -79.5, 36},
	}
	provider.ByBbox(bbox, func(ci consequences.Receptor) {
		ci.Compute(hazards.DepthEvent{})
	})

}
