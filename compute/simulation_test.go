package compute

import (
	"testing"

	"github.com/USACE/go-consequences/consequences"
	"github.com/USACE/go-consequences/hazardproviders"
	"github.com/USACE/go-consequences/structureprovider"
)

func TestComputeEAD(t *testing.T) {
	d := []float64{1, 2, 3, 4}
	f := []float64{.75, .5, .25, 0}
	val := ComputeEAD(d, f)
	if val != 2.0 {
		t.Errorf("computeEAD() yielded %f; expected %f", val, 2.0)
	}
}

func TestComputeEAD2(t *testing.T) {
	d := []float64{1, 10, 30, 45, 59, 78, 89, 102, 140, 180, 240, 330, 350, 370}
	f := []float64{.99, .95, .9, .8, .7, .6, .5, .4, .3, .2, .1, .01, .002, .001}
	val := ComputeEAD(d, f)
	if val != 113.125 {
		t.Errorf("computeEAD() yielded %f; expected %f", val, 113.125)
	}
}
func TestComputeSpecialEAD(t *testing.T) {
	d := []float64{1, 2, 3, 4}
	f := []float64{.75, .5, .25, 0}
	val := ComputeSpecialEAD(d, f)
	if val != 1.875 {
		t.Errorf("computeEAD() yeilded %f; expected %f", val, 1.875)
	}
}
func Test_StreamAbstract(t *testing.T) {
	nsp := structureprovider.InitSHP("/workspaces/Go_Consequences/data/hurricane-laura/ORNLcentroids_LBattributes.shp")
	//nsp := structureprovider.InitNSISP()
	root := "/workspaces/Go_Consequences/data/CERA_Adv29_maxwaterelev_4326_90m"
	filepath := root + ".tif"
	w := consequences.InitGeoJsonResultsWriterFromFile(root + "_consequences.json")
	defer w.Close()
	dfr := hazardproviders.Init(filepath)
	StreamAbstract(dfr, nsp, w)
}
