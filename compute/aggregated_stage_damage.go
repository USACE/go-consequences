package compute

import (
	"fmt"
	"log"

	"github.com/USACE/go-consequences/consequences"
	"github.com/USACE/go-consequences/geography"
	"github.com/USACE/go-consequences/hazardproviders"
	"github.com/USACE/go-consequences/resultswriters"
)

func Aggregated_StageDamage(hps []hazardproviders.HazardProvider, sp consequences.StreamProvider, indexlocation geography.Location, terrainelevation float64, outputfilepath string) {

	w, err := resultswriters.InitAggregatedStageDamageWriterFromFile(outputfilepath)
	if err != nil {
		log.Panicf("unable to initialize the output writer: %s", err)
	}
	for _, hp := range hps {
		//update the writer to know what the index location elevation is
		e, err := hp.ProvideHazard(indexlocation)
		if err != nil {
			log.Panicf("nable to get the index location elevation: %s", err)
		}
		w.SetAggregationElevation(e.Depth() + terrainelevation)
		//get boundingbox
		bbox, err := hp.ProvideHazardBoundary()
		if err != nil {
			log.Panicf("Unable to get the raster bounding box: %s", err)
		}
		fmt.Println(bbox.ToString())
		sp.ByBbox(bbox, func(f consequences.Receptor) {
			//ProvideHazard works off of a geography.Location
			d, err2 := hp.ProvideHazard(geography.Location{X: f.Location().X, Y: f.Location().Y})
			//compute damages based on hazard being able to provide depth
			if err2 == nil {
				//iterate with uncertainty turned on
				i := 0
				maxiter := 100
				for i < maxiter {
					r, err3 := f.Compute(d)
					if err3 == nil {
						w.Write(r)
						i++
					}
				}
			}
		})
	}
	fmt.Println("Getting bbox")

}
