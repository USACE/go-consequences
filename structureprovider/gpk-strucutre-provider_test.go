package structureprovider

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/USACE/go-consequences/hazards"
	"github.com/USACE/go-consequences/structures"
)

func TestGPKByFips(t *testing.T) {
	filepath := "/workspaces/Go_Consequences/data/nsiv2_11.gpkg"
	nsp := InitGPK(filepath)
	fmt.Println(nsp.FilePath)
	d := hazards.DepthEvent{}
	d.SetDepth(2.4)
	nsp.ByFips("11", func(s structures.StructureStochastic) {
		r := s.Compute(d)
		b, _ := json.Marshal(r)
		fmt.Println(string(b))
	})
}
