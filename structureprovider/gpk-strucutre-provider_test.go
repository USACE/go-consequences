package structureprovider

import (
	"fmt"
	"testing"

	"github.com/USACE/go-consequences/structures"
)

func TestGPKByFips(t *testing.T) {
	filepath := "/workspaces/Go_Consequences/data/nsiv2_11.gpkg"
	nsp := InitGPK(filepath)
	fmt.Println(nsp.FilePath)
	nsp.ByFips("11", func(s structures.StructureStochastic) {
		fmt.Println(s.Name)
	})
}
