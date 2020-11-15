package compute

import (
	"fmt"
	"testing"

	"github.com/USACE/go-consequences/hazard_providers"
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
	ComputeMultiEvent_NSIStream(ds, "11")
}
func TestComputeEAD(t *testing.T) {
	d := []float64{1, 2, 3, 4}
	f := []float64{.75, .5, .25, 0}
	val := computeEAD(d, f)
	if val != 2.0 {
		t.Errorf("computeEAD() yeilded %f; expected %f", val, 2.0)
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
