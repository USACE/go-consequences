package compute

import (
	"testing"

	"github.com/USACE/go-consequences/hazardproviders"
	"github.com/USACE/go-consequences/resultswriters"
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
	//initialize the NSI API structure provider
	nsp := structureprovider.InitNSISP()
	//identify the depth grid to apply to the structures.
	root := "/workspaces/Go_Consequences/data/humbolt/S_NWS_STAGE_EFS_33_4326"
	filepath := root + ".tif"
	//w := consequences.InitGeoJsonResultsWriterFromFile(root + "_consequences.json")
	//w := consequences.InitSummaryResultsWriterFromFile(root + "_consequences_SUMMARY.json")
	//create a result writer based on the name of the depth grid.
	w, _ := resultswriters.InitGpkResultsWriter(root+"_consequences_nsi.gpkg", "nsi_result")
	defer w.Close()
	//initialize a hazard provider based on the depth grid.
	dfr, _ := hazardproviders.Init(filepath)
	//compute consequences.
	StreamAbstract(dfr, nsp, w)
}
func Test_StreamAbstract_FIPS_ECAM(t *testing.T) {
	nsp := structureprovider.InitNSISP()
	filepath := "/workspaces/Go_Consequences/data/Base.tif"
	w, _ := resultswriters.InitSummaryResultsWriterFromFile("/workspaces/Go_Consequences/data/base_directLosses.csv")
	defer w.Close()
	dfr, _ := hazardproviders.Init(filepath)
	StreamAbstractByFIPS_WithECAM("48201", dfr, nsp, w)
}
func Test_StreamAbstract_smallDataset(t *testing.T) {
	nsp := structureprovider.InitNSISP()
	root := "/workspaces/Go_Consequences/data/clipped_sample"
	filepath := root + ".tif"
	w, _ := resultswriters.InitGeoJsonResultsWriterFromFile(root + "_consequences.json")
	defer w.Close()
	dfr, _ := hazardproviders.Init(filepath)
	StreamAbstract(dfr, nsp, w)
}
