package compute

import (
	"fmt"
	"testing"

	"github.com/USACE/go-consequences/hazard_providers"
	"github.com/USACE/go-consequences/store"
)

func TestSingleEvent(t *testing.T) {
	fmt.Println("Reading Depths")
	ds := hazard_providers.ReadFeetFile("C:\\Users\\Q0HECWPL\\Documents\\NSI\\NSI_Fathom_depths\\NSI_Fathom_depths_Filtered_Feet.csv")
	fmt.Println("Finished Reading Depths")
	fe := hazard_providers.FathomEvent{Year: 2050, Frequency: 5, Fluvial: true}
	ComputeSingleEvent_NSIStream(ds, "11", fe)
}
func TestMultiEvent(t *testing.T) {
	fmt.Println("Reading Depths")
	ds := hazard_providers.ReadFeetFile("C:\\Users\\Q0HECWPL\\Documents\\NSI\\NSI_Fathom_depths\\NSI_Fathom_depths_Filtered_Feet.csv")
	fmt.Println("Finished Reading Depths")
	db := store.CreateDatabase()
	defer db.Close()
	ComputeMultiEvent_NSIStream(ds, "11", db)
}
func TestMultiEvent_MultiState(t *testing.T) {
	fmt.Println("Reading Depths")
	ds := hazard_providers.ReadFeetFile("C:\\Users\\Q0HECWPL\\Documents\\NSI\\NSI_Fathom_depths\\NSI_Fathom_depths_Filtered_Feet.csv")
	fmt.Println("Finished Reading Depths")
	ComputeMultiFips_MultiEvent(ds)
}
func TestSQLMultiEvent_SingleState(t *testing.T) {
	fmt.Println("Reading Depths")
	ds := hazard_providers.OpenSQLDepthDataSet()
	fmt.Println("Finished Reading Depths")
	db := store.CreateDatabase()
	defer db.Close()
	ComputeMultiEvent_NSIStream(ds, "11", db)
}
func TestSQL_MultiEvent_MultiState(t *testing.T) {
	fmt.Println("Reading Depths")
	ds := hazard_providers.OpenSQLDepthDataSet()
	fmt.Println("Finished Reading Depths")
	ComputeMultiFips_MultiEvent(ds)
}
func TestComputeEAD(t *testing.T) {
	d := []float64{1, 2, 3, 4}
	f := []float64{.75, .5, .25, 0}
	val := computeEAD(d, f)
	if val != 2.0 {
		t.Errorf("computeEAD() yielded %f; expected %f", val, 2.0)
	}
}

func TestComputeEAD2(t *testing.T) {
	d := []float64{1, 10, 30, 45, 59, 78, 89, 102, 140, 180, 240, 330, 350, 370}
	f := []float64{.99, .95, .9, .8, .7, .6, .5, .4, .3, .2, .1, .01, .002, .001}
	val := computeEAD(d, f)
	if val != 113.125 {
		t.Errorf("computeEAD() yielded %f; expected %f", val, 113.125)
	}
}
func TestComputeSpecialEAD(t *testing.T) {
	d := []float64{1, 2, 3, 4}
	f := []float64{.75, .5, .25, 0}
	val := computeSpecialEAD(d, f)
	if val != 1.875 {
		t.Errorf("computeEAD() yeilded %f; expected %f", val, 1.875)
	}
}
