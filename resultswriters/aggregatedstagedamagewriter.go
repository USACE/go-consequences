package resultswriters

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/USACE/go-consequences/consequences"
	"github.com/USACE/go-consequences/paireddata"
)

type aggregatedStageDamageWriter struct {
	filepath         string
	w                io.Writer
	m                map[string]paireddata.UncertaintyPairedData //damage category
	currentElevation float64
}

func InitAggregatedStageDamageWriterFromFile(filepath string) (*aggregatedStageDamageWriter, error) {
	w, err := os.OpenFile(filepath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0600)
	if err != nil {
		return &aggregatedStageDamageWriter{}, err
	}
	//make the maps
	m := make(map[string]paireddata.UncertaintyPairedData, 1)
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

	//use the damcat to select the rigtht damage function to aggregate, update the correct elevation inline histogram with the aggregated damage stored in header "damage".
	fmt.Println(damcat)

}
func (srw *aggregatedStageDamageWriter) Close() {
	//write out the information in the map.

	w2, ok := srw.w.(io.WriteCloser)
	if ok {
		w2.Close()
	}
}
