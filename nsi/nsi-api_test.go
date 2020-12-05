package nsi

import (
	"fmt"
	"sync"
	"testing"

	"github.com/USACE/go-consequences/census"
	"github.com/USACE/go-consequences/consequences"
)

func TestNsiByFips(t *testing.T) {
	var fips string = "15005" //Kalawao county (smallest county in the us by population)
	structures := GetByFips(fips)
	if len(structures) != 101 {
		t.Errorf("GetByFips(%s) yeilded %d structures; expected 101", fips, len(structures))
	}
}
func TestNsiByFipsStream(t *testing.T) {
	var fips string = "15005" //Kalawao county (smallest county in the us by population)
	index := 0
	GetByFipsStream(fips, func(str consequences.StructureStochastic) {
		index++
	})
	if index != 101 {
		t.Errorf("GetByFipsStream(%s) yeilded %d structures; expected 101", fips, index)
	}
}
func TestNsiByFipsStream_MultiState(t *testing.T) {
	f := census.StateToCountyFipsMap()
	var wg sync.WaitGroup
	wg.Add(len(f))
	index := 0
	for ss, _ := range f {
		go func(sfips string) {
			defer wg.Done()
			GetByFipsStream(sfips, func(str consequences.StructureStochastic) {
				index++
			})
			fmt.Println("Completed " + sfips)
		}(ss)
	}
	wg.Wait()
	if index != 109406858 {
		t.Errorf("GetByFipsStream(%s) yeilded %d structures; expected 109,406,858", "all states", index)
	} else {
		fmt.Println("Completed 109,406,858 structures")
	}
}
func TestNsiByBbox(t *testing.T) {
	var bbox string = "-81.58418,30.25165,-81.58161,30.26939,-81.55898,30.26939,-81.55281,30.24998,-81.58418,30.25165"
	structures := GetByBbox(bbox)
	if len(structures) != 1959 {
		t.Errorf("GetByBox(%s) yeilded %d structures; expected 1959", bbox, len(structures))
	}
}
func TestNsiByBboxStream(t *testing.T) {
	var bbox string = "-81.58418,30.25165,-81.58161,30.26939,-81.55898,30.26939,-81.55281,30.24998,-81.58418,30.25165"
	index := 0
	GetByBboxStream(bbox, func(str consequences.StructureStochastic) {
		index++
	})
	if index != 1959 {
		t.Errorf("GetByBoxStream(%s) yeilded %d structures; expected 1959", bbox, index)
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
		t.Errorf("There were %d failures of %d total counties, failed counties were: %s", len(fails), len(counties), s)
	}
}
