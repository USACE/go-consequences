package structureprovider

import (
	"fmt"
	"testing"
)

func TestSHPByFips(t *testing.T) {
	root := "/workspaces/Go_Consequences/data/hurricane-laura/ORNLcentroids_LBattributes"
	filepath := root + ".shp"
	nsp, _ := InitSHP(filepath) //, "ORNLcentroids_LBattributes")
	fmt.Println(nsp.FilePath)
}
