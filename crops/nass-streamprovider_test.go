package crops

import (
	"testing"
	"time"

	"github.com/USACE/go-consequences/compute"
	"github.com/USACE/go-consequences/consequences"
	"github.com/USACE/go-consequences/hazardproviders"
	"github.com/USACE/go-consequences/hazards"
)

func Test_StreamProcessor(t *testing.T) {
	filter := make([]string, 11)
	filter[0] = "1" //filter to corn
	filter[1] = "5"
	filter[2] = "6"
	filter[3] = "22"
	filter[4] = "23"
	filter[5] = "24"
	filter[6] = "28"
	filter[7] = "36"
	filter[8] = "42"
	filter[9] = "52"
	filter[10] = "21"
	nassSp := InitNassCropProvider("2016", filter) // choose a year
	at := time.Date(1984, time.Month(4), 15, 0, 0, 0, 0, time.UTC)
	h := hazards.ArrivalandDurationEvent{}
	h.SetArrivalTime(at)                                                                                                     //fake arrival time
	h.SetDuration(15)                                                                                                        // fake duration
	rw, _ := consequences.InitGpkResultsWriter_Projected("/workspaces/Go_Consequences/data/test2016.gpkg", "agdamage", 5070) // testing data output
	defer rw.Close()
	nassSp.ByFips("19017", func(r consequences.Receptor) { //iterate over a county for testing.
		c, ok := r.(Crop)
		if ok {
			r, err := c.Compute(h)
			if err == nil {
				rw.Write(r)
			}
		}
	})
}
func Test_StreamAbstract(t *testing.T) {
	filter := make([]string, 11)
	filter[0] = "1" //filter to corn
	filter[1] = "5"
	filter[2] = "6"
	filter[3] = "22"
	filter[4] = "23"
	filter[5] = "24"
	filter[6] = "28"
	filter[7] = "36"
	filter[8] = "42"
	filter[9] = "52"
	filter[10] = "21"
	nassSp := InitNassCropProvider("2016", filter)                                                                           // choose a year                                                                                                 // fake duration
	rw, _ := consequences.InitGpkResultsWriter_Projected("/workspaces/Go_Consequences/data/abstract.gpkg", "agdamage", 5070) // testing data output
	defer rw.Close()
	at := time.Date(1984, time.Month(4), 15, 0, 0, 0, 0, time.UTC)
	hp := hazardproviders.InitDaAHP("/workspaces/Go_Consequences/data/Duration5070.tif", "/workspaces/Go_Consequences/data/Arrival5070.tif", at)
	defer hp.Close()
	compute.StreamAbstract(hp, nassSp, rw)
}
