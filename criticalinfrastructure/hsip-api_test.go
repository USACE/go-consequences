package criticalinfrastructure

import (
	"testing"

	"github.com/USACE/go-consequences/consequences"
	"github.com/USACE/go-consequences/geography"
	"github.com/USACE/go-consequences/hazards"
	"github.com/USACE/go-consequences/resultswriters"
)

func TestHSIP(t *testing.T) {
	list := []Layer{Hospitals, PowerPlants, FireStations, WasteWater, LawEnforcement, EmergencyMedicalServices, BRSandEBSTransmitters, CellularTowers, DialysisCenters, EPAandFRSPowerPlants, FacilityInterests, GeneratingUnits, HurricaneEvacuationRoutes, LandMobileBroadcastTowers, LandMobileCommercialTransmissionTowers, LocalEmergencyOperationsCenterEOC, LocalLawEnforcementLocations, MicrowaveServiceTowers, NursingHomes, PagingTransmissionTowers, Pharmacies, PublicHealthDepartments, PublicRefrigeratedWarehouses, UrgentCareFacilities, VeteransHealthAdministrationFacilities}
	provider := InitHsipProvider(list)
	bbox := geography.BBox{
		// {top-left longitude, top-left latitude, bottom-right longitude, bottom-right latitude}
		Bbox: []float64{-106.6456, 47.4567, -66.9499, 24.5231},
	}
	rw, _ := resultswriters.InitSpatialResultsWriter_EPSG_Projected("/workspaces/Go_Consequences/data/test6.GPKG", "criticalInfrastructure", string(resultswriters.GPKG), 4326)
	defer rw.Close()
	provider.ByBbox(bbox, func(ci consequences.Receptor) {
		result, _ := ci.Compute(hazards.DepthEvent{})
		rw.Write(result)
	})
	//fmt.Println(string(rw.Bytes()))

}
