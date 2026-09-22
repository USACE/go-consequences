package compute

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
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
func DepthHazardFunctionModified() hazardproviders.HazardFunction {
	return func(valueIn hazards.HazardData, hazard hazards.HazardEvent) (hazards.HazardEvent, error) {
		d := hazards.DepthEvent{}
		if valueIn.Depth > 1000 {
			return d, hazardproviders.NoDataHazardError{}
		}

		d.SetDepth(valueIn.Depth)
		return d, nil
	}
}
func Test_StreamAbstract_MultiFrequency(t *testing.T) {
	//initialize the NSI API structure provider
	dataset := "BourbonCo_Depth"
	nsp := structureprovider.InitNSISP()

	//initialize a set of frequencies
	frequencies := []float64{.10, .04, .02, .01, .002}
	//specify a working directory for data
	//root := fmt.Sprintf("/vsis3/mmc-storage-6/nsi/Kansas_Silver_Jackets/kansas_ble/%v/", dataset)
	root := fmt.Sprintf("/workspaces/Go_Consequences/data/kc_silverjackets/%v/", dataset)
	//identify the depth grids to represent the frequencies.
	hazardProviders := make([]hazardproviders.HazardProvider, len(frequencies))

	hp1, err := hazardproviders.Init_CustomFunction(fmt.Sprint(root, "Depth_10pct_4326.tif"), DepthHazardFunctionModified())
	if err != nil {
		t.Fail()
	}
	hazardProviders[0] = hp1

	hp2, err := hazardproviders.Init_CustomFunction(fmt.Sprint(root, "Depth_04pct_4326.tif"), DepthHazardFunctionModified())
	if err != nil {
		t.Fail()
	}
	hazardProviders[1] = hp2

	hp3, err := hazardproviders.Init_CustomFunction(fmt.Sprint(root, "Depth_02pct_4326.tif"), DepthHazardFunctionModified())
	if err != nil {
		t.Fail()
	}
	hazardProviders[2] = hp3

	hp4, err := hazardproviders.Init_CustomFunction(fmt.Sprint(root, "Depth_01pct_4326.tif"), DepthHazardFunctionModified())
	if err != nil {
		t.Fail()
	}
	hazardProviders[3] = hp4

	hp5, err := hazardproviders.Init_CustomFunction(fmt.Sprint(root, "Depth_0_2pct_4326.tif"), DepthHazardFunctionModified())
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

// func Test_StreamAbstract_MultiHazard(t *testing.T) {
// 	//initialize the NSI API structure provider
// 	nsp := structureprovider.InitNSISPwithOcctypeFilePath("/workspaces/go-consequences/data/lifecycle/occtypes_reconstruction.json")
// 	now := time.Now()
// 	fmt.Println(now)

// 	root := "/workspaces/go-consequences/data/lifecycle/"
// 	filepath := root + "test_arrival-depth-duration_hazards.json"
// 	w, _ := resultswriters.InitSpatialResultsWriter(root+"multihazardtest_consequences.gpkg", "results", "GPKG")
// 	//w := consequences.InitSummaryResultsWriterFromFile(root + "_consequences_SUMMARY.json")
// 	//create a result writer based on the name of the depth grid.
// 	//w, _ := resultswriters.InitGpkResultsWriter(root+"_consequences_nsi.gpkg", "nsi_result")
// 	defer w.Close()
// 	//initialize a hazard provider based on the depth grid.
// 	dfr, err := hazardproviders.InitADDMHP(filepath)
// 	if err != nil {
// 		panic(err)
// 	}
// 	//compute consequences.
// 	StreamAbstract(dfr, nsp, w)
// 	fmt.Println(time.Since(now))
// }

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
	nsp := structureprovider.InitNSISP()
	now := time.Now()
	fmt.Println(now)
	//nsp, _ := structureprovider.InitStructureProvider("/workspaces/Go_Consequences/data/ffrd/Lower Kanawha-Elk Lower.gpkg", "Lower Kanawha-Elk Lower", "GPKG")
	//nsp.SetDeterministic(true)
	//identify the depth grid to apply to the structures.
	root := "/workspaces/Go_Consequences/data/kc_silverjackets/Douglas_Co_depth/DG_Depth_01pct"
	filepath := root + ".tif"
	w, _ := resultswriters.InitSpatialResultsWriter(root+"_consequences3.gpkg", "results", "GPKG")
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

// ---- EDIT THESE for the one run -----------------------------------------
const (
	rootDir    = "/workspaces/Go_Consequences/data/kc_silverjackets/"                               // root directory containing the folders
	foldersCSV = "BourbonCo_Depth"                                                                  // folder names under root
	filesCSV   = "depth_0_2pct.tif,depth_01pct.tif,depth_02pct.tif,depth_04pct.tif,depth_10pct.tif" // geotiff names expected in each folder
	workers    = 4                                                                                  // max concurrent gdalwarp processes
)

// with4326Suffix inserts _4326 before the extension: dem.tif -> dem_4326.tif.
// Idempotent: a name already ending in _4326 is returned unchanged.
func with4326Suffix(name string) string {
	ext := filepath.Ext(name)
	base := strings.TrimSuffix(name, ext)
	if strings.HasSuffix(base, "_4326") {
		return name
	}
	return base + "_4326" + ext
}

func splitCSV(s string) []string {
	var out []string
	for _, p := range strings.Split(s, ",") {
		if p = strings.TrimSpace(p); p != "" {
			out = append(out, p)
		}
	}
	return out
}

type job struct {
	in  string // absolute path to input geotiff
	out string // absolute path to output (<in>_4326)
}

func TestReproject(t *testing.T) {
	if _, err := exec.LookPath("gdalwarp"); err != nil {
		t.Fatalf("gdalwarp not found on PATH: %v", err)
	}

	// Build the worklist: every folder x every expected file.
	var jobs []job
	for _, f := range splitCSV(foldersCSV) {
		dir := filepath.Join(rootDir, f)
		for _, fn := range splitCSV(filesCSV) {
			jobs = append(jobs, job{
				in:  filepath.Join(dir, fn),
				out: filepath.Join(dir, with4326Suffix(fn)),
			})
		}
	}

	sem := make(chan struct{}, workers) // bound concurrent gdalwarp processes
	var wg sync.WaitGroup
	var mu sync.Mutex
	var failures []string
	var doneCount, skipCount int

	for _, j := range jobs {
		wg.Add(1)
		go func(j job) {
			defer wg.Done()
			sem <- struct{}{}        // acquire a worker slot
			defer func() { <-sem }() // release it

			// Skip if the reprojected file already exists (safe to re-run).
			if _, err := os.Stat(j.out); err == nil {
				mu.Lock()
				skipCount++
				t.Logf("skip (already present): %s", j.out)
				mu.Unlock()
				return
			}
			// Fail the job if the input is missing.
			if _, err := os.Stat(j.in); err != nil {
				mu.Lock()
				failures = append(failures, fmt.Sprintf("%s: input not found: %v", j.in, err))
				mu.Unlock()
				return
			}

			// Source SRS is read from the tiff itself, so only -t_srs is needed.
			// -co COMPRESS=DEFLATE writes a deflate-compressed GeoTIFF.
			cmd := exec.Command("gdalwarp",
				"-t_srs", "EPSG:4326",
				"-of", "GTiff",
				"-r", "bilinear",
				"-co", "COMPRESS=DEFLATE",
				j.in,
				j.out,
			)
			var stderr strings.Builder
			cmd.Stderr = &stderr
			if err := cmd.Run(); err != nil {
				mu.Lock()
				failures = append(failures, fmt.Sprintf("%s: gdalwarp failed: %v\n%s", j.in, err, stderr.String()))
				mu.Unlock()
				return
			}

			mu.Lock()
			doneCount++
			t.Logf("ok: %s -> %s", j.in, j.out)
			mu.Unlock()
		}(j)
	}
	wg.Wait()

	t.Logf("summary: %d reprojected, %d skipped, %d failed", doneCount, skipCount, len(failures))
	for _, f := range failures {
		t.Error(f) // marks the test failed and prints the reason
	}
}
