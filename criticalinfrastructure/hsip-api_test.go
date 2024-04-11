package criticalinfrastructure

import (
	"testing"

	"github.com/USACE/go-consequences/consequences"
	"github.com/USACE/go-consequences/geography"
	"github.com/USACE/go-consequences/hazards"
	"github.com/USACE/go-consequences/resultswriters"
)

func TestHSIP(t *testing.T) {
	list := []Layer{Hospitals, PowerPlants, FireStations, WasteWater, LawEnforcement, EmergencyMedicalServices, AmtrakStations}
	provider := InitHsipProvider(list)
	bbox := geography.BBox{
		Bbox: []float64{-80, 36, -79.5, 35.5},
	}
	rw, _ := resultswriters.InitSpatialResultsWriter_EPSG_Projected("/workspaces/Go_Consequences/data/test5", "criticalInfrastructure", string(resultswriters.PARQUET), 4326)
	defer rw.Close()
	provider.ByBbox(bbox, func(ci consequences.Receptor) {
		result, _ := ci.Compute(hazards.DepthEvent{})
		rw.Write(result)
	})
	//fmt.Println(string(rw.Bytes()))

}
