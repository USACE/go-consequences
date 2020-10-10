package nsi

import (
	"sync"
	"testing"

	"github.com/USACE/go-consequences/census"
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
func TestNSI_FIPS_CA_ERRORS(t *testing.T) {
	f := census.StateToCountyFipsMap()
	var wg sync.WaitGroup
	counties := f["06"]
	fails := make([]string, 0)
	wg.Add(len(counties))
	for _, ccc := range counties {
		go func(county string) {
			defer wg.Done()
			structures := GetByFips(county)
			if len(structures) == 0 {
				fails = append(fails, county)
				//t.Errorf("GetByFips(%s) yeilded %d structures; expected more than zero", county, len(structures))
			}
		}(ccc)
	}
	wg.Wait()
	if len(fails) > 0 {
		s := "Counties: "
		for _, f := range fails {
			s += f + ", "
		}
		t.Errorf("There were %d failures %s", len(fails), s)
	}
}
