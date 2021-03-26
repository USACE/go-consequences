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

func TestNsiByFipsStream(t *testing.T) {
	var fips string = "15005" //Kalawao county (smallest county in the us by population)
	n := InitNSISP()
	counter := 0
	n.ByFips(fips, func(s consequences.Receptor) {
		counter++
	})
	if counter != 101 {
		t.Errorf("GetByFips(%s) yeilded %d structures; expected 101", fips, counter)
	}
}

func TestNsiByFipsStream_MultiState(t *testing.T) {
	f := census.StateToCountyFipsMap()
	var wg sync.WaitGroup
	wg.Add(len(f))
	n := InitNSISP()
	index := 0
	for ss := range f {
		go func(sfips string) {
			defer wg.Done()
			n.ByFips(sfips, func(s consequences.Receptor) {
				index++
			})
			fmt.Println("Completed " + sfips)
		}(ss)
	}
	wg.Wait()
	if index != 109406858 {
		t.Errorf("ByFips(%s) yeilded %d structures; expected 109,406,858", "all states", index)
	} else {
		fmt.Println("Completed 109,406,858 structures")
	}
}
func TestNsiByBboxStream(t *testing.T) {
	bbox := make([]float64, 4) //i might have these values inverted
	bbox[0] = -81.58418        //upper left x
	bbox[1] = 30.25165         //upper left y
	bbox[2] = -81.58161        //lower right x
	bbox[3] = 30.26939         //lower right y
	gbbx := geography.BBox{Bbox: bbox}
	n := InitNSISP()
	counter := 0
	n.ByBbox(gbbx, func(s consequences.Receptor) {
		counter++
	})
	if counter != 1959 {
		t.Errorf("ByBox(%s) yeilded %d structures; expected 1959", gbbx.ToString(), counter)
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
