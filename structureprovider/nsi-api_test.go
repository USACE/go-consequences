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
			index := 0
			n.ByFips(sfips, func(s consequences.Receptor) {
				index++
			})
			fmt.Println(fmt.Sprintf("Completed %s with %v structures", sfips, index))
			if countByState(sfips) == index {
				fmt.Println(fmt.Sprintf("For state %s the count matched", sfips))
			} else {
				fmt.Println(fmt.Sprintf("For state %s the count did NOT match!", sfips))
			}
		}(ss)
	}
	wg.Wait()
	if index != 109406858 {
		t.Errorf("ByFips(%s) yeilded %d structures; expected 109,406,858", "all states", index)
	} else {
		fmt.Println("Completed 109,406,858 structures")
	}
}
func TestNsiByFipsStream_MultiState_Sequential(t *testing.T) {
	f := census.StateToCountyFipsMap()
	n := InitNSISP()
	index := 0
	for ss := range f {
		func(sfips string) {
			index := 0
			n.ByFips(sfips, func(s consequences.Receptor) {
				index++
			})
			fmt.Println(fmt.Sprintf("Completed %s with %v structures", sfips, index))
			if countByState(sfips) == index {
				fmt.Println(fmt.Sprintf("For state %s the count matched", sfips))
			} else {
				fmt.Println(fmt.Sprintf("For state %s the count did NOT match!", sfips))
			}
		}(ss)
	}
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

func countByState(ss string) int {
	m := make(map[string]int)
	m["01"] = 2110979
	m["02"] = 293688
	m["04"] = 2668111
	m["05"] = 1263237
	m["06"] = 11111929
	m["08"] = 2045575
	m["09"] = 1234236
	m["10"] = 381352
	m["11"] = 153062
	m["12"] = 7694670
	m["13"] = 3810113
	m["15"] = 344996
	m["16"] = 666303
	m["17"] = 4066775
	m["18"] = 2613847
	m["19"] = 1304639
	m["20"] = 1205171
	m["21"] = 1607171
	m["22"] = 1872003
	m["23"] = 607351
	m["24"] = 2155044
	m["25"] = 2042137
	m["26"] = 4329378
	m["27"] = 2167795
	m["28"] = 1371025
	m["29"] = 2544779
	m["30"] = 501699
	m["31"] = 785530
	m["32"] = 1017555
	m["33"] = 542509
	m["34"] = 2982407
	m["35"] = 800302
	m["36"] = 4737239
	m["37"] = 4266543
	m["38"] = 312468
	m["39"] = 4510863
	m["40"] = 1617867
	m["41"] = 1539773
	m["42"] = 4870271
	m["44"] = 375415
	m["45"] = 1932100
	m["46"] = 371769
	m["47"] = 2578884
	m["48"] = 9083959
	m["49"] = 938757
	m["50"] = 229969
	m["51"] = 3466599
	m["53"] = 2647194
	m["54"] = 773467
	m["55"] = 1989245
	m["56"] = 275446

	return m[ss]
}
