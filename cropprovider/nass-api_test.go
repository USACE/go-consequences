package cropprovider

import (
	"fmt"
	"testing"
	"time"

	"github.com/USACE/go-consequences/crops"
	"github.com/USACE/go-consequences/hazards"
)

func TestNassStatsByBbox(t *testing.T) {
	//https://nassgeodata.gmu.edu/axis2/services/CDLService/GetCDLStat?year=2018&bbox=130783,2203171,153923,2217961&format=csv"
	stats := GetStatsByBbox("2018", "130783", "2203171", "153923", "2217961")
	//diff := stats.Acreage - 953459824.285892
	if !stats.Success {
		t.Errorf("GetByBox() yeilded %v;", stats)
	}
	fmt.Println(stats)
}
func TestNassCDLValue(t *testing.T) {
	//https://nassgeodata.gmu.edu/axis2/services/CDLService/GetCDLValue?year=2018&x=1551565.363&y=1909363.537
	result := GetCDLValue("2018", "1551565.363", "1909363.537")
	if result.GetCropName() == "" {
		t.Error("GetCDLValue() yeilded nothing;")
	}
	//fmt.Println(result)
}
func TestNassCDLFile(t *testing.T) {
	//requires write access to /workspaces/Go_Consequences/data/
	result, err := GetCDLFileByFIPS("2018", "19017")
	if err != nil {
		t.Error(err)
	}
	fmt.Println(result.FilePath)
}
func TestNassCDLFileSampleValue(t *testing.T) {
	ncp := Init("/workspaces/Go_Consequences/data/CDL_2018_19015.tif")
	fmt.Println(ncp.getCropValue(174133, 2125229)) //should be 1, Corn
	fmt.Println(ncp.getCropValue(180913, 2115830)) //should be 5, Soybeans
	fmt.Println(ncp.getCropValue(156842, 2125731)) //should be 36, Alfalfa
}
func TestNassCDLFileFiltered(t *testing.T) {
	//requires write access to C:\\Temp\\agtesting
	result := GetCDLFileByFIPSFiltered("2018", "19015", "1,5")
	if !result {
		t.Error("GetCDLFileByFIPSFiltered() returned false;")
	}
	//fmt.Println(result)
}
func TestCropDamage(t *testing.T) {
	//get crop
	cropFromNass := GetCDLValue("2018", "1551565.363", "1909363.537")
	path := "./" + cropFromNass.GetCropName() + ".crop"
	c := crops.ReadFromXML(path)
	// construct hazard
	at := time.Date(1984, time.Month(7), 29, 0, 0, 0, 0, time.UTC)
	h := hazards.ArrivalandDurationEvent{}
	h.SetArrivalTime(at)
	h.SetDuration(10)
	//compute
	cd, _ := c.Compute(h)
	//expected results
	expectedcase := crops.Impacted
	expecteddamage := 1285.98 //Based on corn

	//test
	if cd.Result[1] != expectedcase {
		t.Errorf("ComputeConsequence() = %v; expected %v", cd.Result[3], expectedcase)
	}
	if cd.Result[2] != expecteddamage {
		t.Errorf("ComputeConsequence() = %v; expected %v", cd.Result[4], expecteddamage)
	}

}

func TestCropDamage_DelayedPlant(t *testing.T) {
	//get crop
	cropFromNass := GetCDLValue("2018", "1551565.363", "1909363.537")
	path := "./" + cropFromNass.GetCropName() + ".crop"
	c := crops.ReadFromXML(path)
	// construct hazard
	at := time.Date(1984, time.Month(4), 15, 0, 0, 0, 0, time.UTC)
	h := hazards.ArrivalandDurationEvent{}
	h.SetArrivalTime(at)
	h.SetDuration(15)
	//compute
	cd, _ := c.Compute(h)
	//expected results
	expectedcase := crops.PlantingDelayed
	expecteddamage := 0.0 //Based on corn

	//test
	if cd.Result[1] != expectedcase {
		t.Errorf("Compute() = %v; expected %v", cd.Result[3], expectedcase)
	}
	if cd.Result[2] != expecteddamage {
		t.Errorf("Compute() = %v; expected %v", cd.Result[4], expecteddamage)
	}

}
