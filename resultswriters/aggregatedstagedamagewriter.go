package resultswriters

import (
	"io"
	"log"
	"os"

	"github.com/HydrologicEngineeringCenter/go-statistics/data"
	"github.com/USACE/go-consequences/consequences"
)

type aggregatedStageDamageWriter struct {
	filepath         string
	w                io.Writer
	m                map[string]map[float64]*data.InlineHistogram //damage category
	currentElevation float64
}

func InitAggregatedStageDamageWriterFromFile(filepath string) (*aggregatedStageDamageWriter, error) {
	w, err := os.OpenFile(filepath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0600)
	if err != nil {
		return &aggregatedStageDamageWriter{}, err
	}
	//make the maps
	m := make(map[string]map[float64]*data.InlineHistogram, 1)
	return &aggregatedStageDamageWriter{filepath: filepath, w: w, m: m}, nil
}
func (srw *aggregatedStageDamageWriter) SetAggregationElevation(ele float64) {
	srw.currentElevation = ele
}
func (srw *aggregatedStageDamageWriter) Write(r consequences.Result) {
	//hardcoding for structures to experiment and think it through.
	damcati, err := r.Fetch("damage category")

	if err != nil {
		log.Fatal("couldnt find the damage category")
	}
	damcat := damcati.(string)
	aggregateddamage, err := r.Fetch("damage")
	if err != nil {
		log.Fatal("couldnt find the damage")
	}
	agdam := aggregateddamage.(float64)
	//use the damcat to select the rigtht damage function to aggregate, update the correct elevation inline histogram with the aggregated damage stored in header "damage".
	elehisto, ok := srw.m[damcat]
	if ok {
		histo, hok := elehisto[srw.currentElevation]
		if hok {
			histo.AddObservation(agdam)
		} else {
			histo = data.Init(1000.0, agdam-500.0, agdam+500.0)
			histo.AddObservation(agdam)
		}
		elehisto[srw.currentElevation] = histo
		srw.m[damcat] = elehisto
	} else {

		histo := data.Init(1000.0, agdam-500.0, agdam+500.0)
		histo.AddObservation(agdam)
		mi := make(map[float64]*data.InlineHistogram)
		mi[srw.currentElevation] = histo
		srw.m[damcat] = mi
	}
}
func (srw *aggregatedStageDamageWriter) Close() {
	//write out the information in the map.

	w2, ok := srw.w.(io.WriteCloser)
	if ok {
		w2.Close()
	}
}
