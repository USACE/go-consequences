package crops

import (
	"fmt"
	"testing"
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
	//requires write access to C:\\Temp\\agtesting
	result := GetCDLFileByFIPS("2018", "19015")
	if !result {
		t.Error("GetCDLFile() returned false;")
	}
	//fmt.Println(result)
}
func TestNassCDLFileFiltered(t *testing.T) {
	//requires write access to C:\\Temp\\agtesting
	result := GetCDLFileByFIPSFiltered("2018", "19015", "1,5")
	if !result {
		t.Error("GetCDLFileByFIPSFiltered() returned false;")
	}
	//fmt.Println(result)
}
