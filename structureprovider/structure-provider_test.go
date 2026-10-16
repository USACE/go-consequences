package structureprovider

import (
	"fmt"
	"testing"

	"github.com/USACE/go-consequences/consequences"
	"github.com/USACE/go-consequences/geography"
	"github.com/USACE/go-consequences/structures"
)

func Test_InitStructureProvider(t *testing.T) {
	//TODO: use finalized testing data
	fp := "/workspaces/go-consequences/data/lifecycle/nsi_2026_test.gpkg"

	sp, err := InitStructureProvider(fp, "nsi2022", "GPKG")
	if err != nil {
		t.Error(err)
	}
	bbox := geography.BBox{Bbox: []float64{-180, 30, 180, 40}}
	i := 0
	sp.ByBbox(bbox, func(f consequences.Receptor) {
		if i == 0 {
			s := f.(structures.StructureDeterministic)
			//Check all BaseStructure attributes are loaded
			fmt.Printf("Name: %s\n", s.BaseStructure.Name)
			fmt.Printf("DamCat: %s\n", s.BaseStructure.DamCat)
			fmt.Printf("CBFips: %s\n", s.BaseStructure.CBFips)
			fmt.Printf("X, Y: %v, %v\n", s.BaseStructure.X, s.BaseStructure.Y)
			fmt.Printf("GroundElevation: %v\n", s.BaseStructure.GroundElevation)
		}
		i++
	})
}
