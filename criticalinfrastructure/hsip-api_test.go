package criticalinfrastructure

import (
	"testing"

	"github.com/USACE/go-consequences/consequences"
	"github.com/USACE/go-consequences/geography"
	"github.com/USACE/go-consequences/hazards"
)

func TestHSIP(t *testing.T) {
	list := []Layer{Hospitals, PowerPlants, FireStations}
	provider := InitHsipProvider(list)
	bbox := geography.BBox{
		Bbox: []float64{-80, 36, -79.5, 35.5},
	}
	provider.ByBbox(bbox, func(ci consequences.Receptor) {
		ci.Compute(hazards.DepthEvent{})
	})

}
