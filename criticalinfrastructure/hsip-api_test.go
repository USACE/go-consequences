package criticalinfrastructure

import (
	"testing"

	"github.com/USACE/go-consequences/consequences"
	"github.com/USACE/go-consequences/geography"
	"github.com/USACE/go-consequences/hazards"
)

func TestHSIP(t *testing.T) {
	list := []Layer{Hospitals}
	provider := InitHsipProvider(list)
	bbox := geography.BBox{
		Bbox: []float64{-95, 40, -94.32, 41},
	}
	provider.ByBbox(bbox, func(ci consequences.Receptor) {
		ci.Compute(hazards.DepthEvent{})
	})

}
