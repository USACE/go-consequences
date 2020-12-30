package crops

import (
	"fmt"
	"testing"
)

func TestNassStatsByBbox(t *testing.T) {
	//https://nassgeodata.gmu.edu/axis2/services/CDLService/GetCDLStat?year=2018&bbox=130783,2203171,153923,2217961&format=json"
	stats := GetStatsByBbox("2018", "130783", "2203171", "153923", "2217961")
	//diff := stats.Acreage - 953459824.285892
	if !stats.Success {
		t.Errorf("GetByBox() yeilded %v;", stats)
	}
	fmt.Println(stats)
}
