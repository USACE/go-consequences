package structureprovider

import (
	"fmt"
	"sync"
	"testing"

	"github.com/USACE/go-consequences/census"
	"github.com/USACE/go-consequences/consequences"
	"github.com/USACE/go-consequences/geography"
	"github.com/USACE/go-consequences/structures"
)

func TestNsiByFips(t *testing.T) {
	var fips string = "15005" //Kalawao county (smallest county in the us by population)
	structures := GetByFips(fips)
	if len(structures.Features) != 101 {
		t.Errorf("GetByFips(%s) yeilded %d structures; expected 101", fips, len(structures.Features))
	}
}
func TestNsiByFipsStream(t *testing.T) {
	var fips string = "15005" //Kalawao county (smallest county in the us by population)
	index := 0
	GetByFipsStream(fips, func(str NsiFeature) {
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
	for ss := range f {
		go func(sfips string) {
			defer wg.Done()
			GetByFipsStream(sfips, func(str NsiFeature) {
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
func TestNsi_FL_FoundationTypes(t *testing.T) {
	f := []string{"12"}
	foundationTypes := make(map[string]int64)
	for _, ss := range f {
		GetByFipsStream(ss, func(str NsiFeature) {
			val, ok := foundationTypes[str.Properties.FoundType]
			if ok {
				v2 := val + 1
				foundationTypes[str.Properties.FoundType] = v2
			} else {
				foundationTypes[str.Properties.FoundType] = 1
			}
		})
	}
	fmt.Println(foundationTypes)
}
func TestNsiByBbox(t *testing.T) {
	var bbox string = "-81.58418,30.25165,-81.58161,30.26939,-81.55898,30.26939,-81.55281,30.24998,-81.58418,30.25165"
	structures := GetByBbox(bbox)
	if len(structures.Features) != 1959 {
		t.Errorf("GetByBox(%s) yeilded %d structures; expected 1959", bbox, len(structures.Features))
	}
}
func TestNsiByBboxStream(t *testing.T) {
	var bbox string = "-81.58418,30.25165,-81.58161,30.26939,-81.55898,30.26939,-81.55281,30.24998,-81.58418,30.25165"
	index := 0
	GetByBboxStream(bbox, func(str NsiFeature) {
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
			if len(structures.Features) == 0 {
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
func Test_StructureProvider_NSI_BBOX(t *testing.T) {
	bbox := make([]float64, 4) //i might have these values inverted
	bbox[0] = -81.58418        //upper left x
	bbox[1] = 30.25165         //upper left y
	bbox[2] = -81.58161        //lower right x
	bbox[3] = 30.26939         //lower right y
	gbbx := geography.BBox{Bbox: bbox}
	nsp := InitNSISP()
	nsp.ByBbox(gbbx, func(c consequences.Receptor) {
		s, _ := c.(structures.StructureStochastic)
		fmt.Println(s.Name)
	})
}
