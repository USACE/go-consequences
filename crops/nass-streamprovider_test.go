package crops

import (
	"testing"
	"time"

	"github.com/USACE/go-consequences/consequences"
	"github.com/USACE/go-consequences/hazards"
)

func Test_StreamProcessor(t *testing.T) {
	filter := make([]string, 1)
	filter[0] = "1"                                //filter to corn
	nassSp := InitNassCropProvider("2017", filter) // choose a year
	at := time.Date(1984, time.Month(4), 15, 0, 0, 0, 0, time.UTC)
	h := hazards.ArrivalandDurationEvent{}
	h.SetArrivalTime(at)                                                                              //fake arrival time
	h.SetDuration(15)                                                                                 // fake duration
	rw := consequences.InitGpkResultsWriter("/workspaces/Go_Consequences/data/test.gpkg", "agdamage") // testing data output
	defer rw.Close()
	nassSp.ByFips("19017", func(r consequences.Receptor) { //iterate over a county for testing.
		c, ok := r.(Crop)
		if ok {
			r := c.Compute(h)
			rw.Write(r)
		}
	})
}
