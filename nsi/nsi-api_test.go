package nsi

import (
	"testing"
)

func TestNsiByFips(t *testing.T) {
	var fips string = "15005" //Kalawao county (smallest county in the us by population)
	structures := GetByFips(fips)
	if len(structures) != 101 {
		t.Errorf("GetByFips(%s) yeilded %d structures; expected 101", fips, len(structures))
	}
}
func TestNsiByBbox(t *testing.T) {
	var bbox string = "-81.58418,30.25165,-81.58161,30.26939,-81.55898,30.26939,-81.55281,30.24998,-81.58418,30.25165"
	structures := GetByBbox(bbox)
	if len(structures) != 1939 {
		t.Errorf("GetByBox(%s) yeilded %d structures; expected 1939", bbox, len(structures))
	}
}
