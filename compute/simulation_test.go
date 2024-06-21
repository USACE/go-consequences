package compute

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/USACE/go-consequences/hazardproviders"
	"github.com/USACE/go-consequences/hazards"
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
func Test_StreamAbstract_MultiFrequency(t *testing.T) {
	//initialize the NSI API structure provider
	dataset := "Rice_CowCreek_depth"
	nsp := structureprovider.InitNSISP()

	//initialize a set of frequencies
	frequencies := []float64{.10, .04, .02, .01, .002}
	//specify a working directory for data
	//root := fmt.Sprintf("/vsis3/mmc-storage-6/nsi/Kansas_Silver_Jackets/kansas_ble/%v/", dataset)
	root := fmt.Sprintf("/workspaces/Go_Consequences/data/kc_silverjackets/%v/", dataset)
	//identify the depth grids to represent the frequencies.
	hazardProviders := make([]hazardproviders.HazardProvider, len(frequencies))

	hp1, err := hazardproviders.Init(fmt.Sprint(root, "Depth_10pct.tif"))
	if err != nil {
		t.Fail()
	}
	hazardProviders[0] = hp1

	hp2, err := hazardproviders.Init(fmt.Sprint(root, "Depth_04pct.tif"))
	if err != nil {
		t.Fail()
	}
	hazardProviders[1] = hp2

	hp3, err := hazardproviders.Init(fmt.Sprint(root, "Depth_02pct.tif"))
	if err != nil {
		t.Fail()
	}
	hazardProviders[2] = hp3

	hp4, err := hazardproviders.Init(fmt.Sprint(root, "Depth_01pct.tif"))
	if err != nil {
		t.Fail()
	}
	hazardProviders[3] = hp4

	hp5, err := hazardproviders.Init(fmt.Sprint(root, "Depth_0_2pct.tif"))
	if err != nil {
		t.Fail()
	}
	hazardProviders[4] = hp5

	//create a result writer based on the name of the depth grid.
	//write local
	path := fmt.Sprintf("/workspaces/Go_Consequences/data/kc_silverjackets/%v/%v_consequences_nsi.gpkg", dataset, dataset)
	w, _ := resultswriters.InitSpatialResultsWriter(path, "nsi_result", "GPKG")
	defer w.Close()
	//compute consequences.
	StreamAbstractMultiFrequency(hazardProviders, frequencies, nsp, w)
}
func Test_Config(t *testing.T) {
	config := Config{
		StructureProviderInfo: structureprovider.StructureProviderInfo{
			StructureProviderType:   structureprovider.OGR,
			StructureProviderDriver: "PARQUET",
			LayerName:               "lower_kanawha_lower_elk",
			StructureFilePath:       "/workspaces/Go_Consequences/data/ffrd/lower_kanawha_lower_elk.parquet",
		},
		HazardProviderInfo: hazardproviders.HazardProviderInfo{
			Hazards: []hazardproviders.HazardProviderParameterAndPath{
				hazardproviders.HazardProviderParameterAndPath{
					Hazard:   hazards.Depth,
					FilePath: "/workspaces/Go_Consequences/data/ffrd/LowKanLowElk/depth_grid.vrt",
				},
			},
		},
		ResultsWriterInfo: resultswriters.ResultsWriterInfo{
			Type:     resultswriters.JSON,
			FilePath: "/workspaces/Go_Consequences/data/ffrd/LowKanLowElk/depth_grid_consequences.json",
		},
	}
	b, err := json.Marshal(config)
	if err != nil {
		t.Fail()
	}
	configPath := "/workspaces/Go_Consequences/data/ffrd/configexample.json"
	if _, err := os.Stat(configPath); errors.Is(err, os.ErrNotExist) {
		//does not exist
	} else {
		os.Remove(configPath)
	}
	os.WriteFile(configPath, b, os.ModeAppend)
	computable, err := config.CreateComputable()
	if err != nil {
		t.Fail()
	}
	err = computable.Compute()
	if err != nil {
		t.Fail()
	}

}
func Test_StreamAbstract(t *testing.T) {
	//initialize the NSI API structure provider
	//nsp := structureprovider.InitNSISP()
	now := time.Now()
	fmt.Println(now)
	nsp, _ := structureprovider.InitStructureProvider("/workspaces/Go_Consequences/data/ffrd/Lower Kanawha-Elk Lower.gpkg", "Lower Kanawha-Elk Lower", "GPKG")
	nsp.SetDeterministic(true)
	//identify the depth grid to apply to the structures.
	root := "/workspaces/Go_Consequences/data/ffrd/LowKanLowElk/depth_grid"
	filepath := root + ".vrt"
	w, _ := resultswriters.InitSpatialResultsWriter(root+"_consequences.json", "results", "GeoJSON")
	//w := consequences.InitSummaryResultsWriterFromFile(root + "_consequences_SUMMARY.json")
	//create a result writer based on the name of the depth grid.
	//w, _ := resultswriters.InitGpkResultsWriter(root+"_consequences_nsi.gpkg", "nsi_result")
	defer w.Close()
	//initialize a hazard provider based on the depth grid.
	dfr, _ := hazardproviders.Init_CustomFunction(filepath, func(valueIn hazards.HazardData, hazard hazards.HazardEvent) (hazards.HazardEvent, error) {
		if valueIn.Depth == 0 {
			return hazard, hazardproviders.NoHazardFoundError{}
		}
		process := hazardproviders.DepthHazardFunction()
		return process(valueIn, hazard)
	})
	//compute consequences.
	StreamAbstract(dfr, nsp, w)
	fmt.Println(time.Since(now))
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
	w, _ := resultswriters.InitSpatialResultsWriter(root+"_consequences.json", "results", "GeoJSON")
	defer w.Close()
	dfr, _ := hazardproviders.Init(filepath)
	StreamAbstract(dfr, nsp, w)
}
