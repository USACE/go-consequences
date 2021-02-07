package compute

import (
	"fmt"
	"testing"

	"github.com/USACE/go-consequences/nsi"
	"github.com/USACE/go-consequences/structures"
)

func TestConvertNSIFeatureToStructure(t *testing.T) {
	bbox := "-81.58418,30.25165,-81.58161,30.26939,-81.55898,30.26939,-81.55281,30.24998,-81.58418,30.25165"
	//get a map of all occupancy types
	m := structures.OccupancyTypeMap()
	//define a default occtype in case of emergancy
	defaultOcctype := m["RES1-1SNB"]
	nsi.GetByBboxStream(bbox, func(f nsi.NsiFeature) {
		//convert nsifeature to structure
		str := NsiFeaturetoStructure(f, m, defaultOcctype)
		fmt.Println(str.OccType.Name)
	})
}
func TestComputeEAD(t *testing.T) {
	d := []float64{1, 2, 3, 4}
	f := []float64{.75, .5, .25, 0}
	val := ComputeEAD(d, f)
	if val != 2.0 {
		t.Errorf("computeEAD() yielded %f; expected %f", val, 2.0)
	}
}

func TestComputeEAD2(t *testing.T) {
	d := []float64{1, 10, 30, 45, 59, 78, 89, 102, 140, 180, 240, 330, 350, 370}
	f := []float64{.99, .95, .9, .8, .7, .6, .5, .4, .3, .2, .1, .01, .002, .001}
	val := ComputeEAD(d, f)
	if val != 113.125 {
		t.Errorf("computeEAD() yielded %f; expected %f", val, 113.125)
	}
}
func TestComputeSpecialEAD(t *testing.T) {
	d := []float64{1, 2, 3, 4}
	f := []float64{.75, .5, .25, 0}
	val := ComputeSpecialEAD(d, f)
	if val != 1.875 {
		t.Errorf("computeEAD() yeilded %f; expected %f", val, 1.875)
	}
}
