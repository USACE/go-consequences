package compute

import (
	"fmt"
	"log"

	"github.com/USACE/go-consequences/consequences"
	"github.com/USACE/go-consequences/geography"
	"github.com/USACE/go-consequences/hazardproviders"
	"github.com/USACE/go-consequences/resultswriters"
)

func Aggregated_StageDamage(hps []hazardproviders.HazardProvider, sp consequences.StreamProvider, indexlocation geography.Location, terrainelevation float64, coutputfilepath string, soutputfilepath string) {

	//need to allow provision of the user specified occupancy types
	//need to allow provision of impact area polygons

	cw, err := resultswriters.InitAggregatedStageDamageWriterFromFile(coutputfilepath)
	if err != nil {
		log.Panicf("unable to initialize the content output writer: %s", err)
	}
	sw, err := resultswriters.InitAggregatedStageDamageWriterFromFile(soutputfilepath)
	if err != nil {
		log.Panicf("unable to initialize the structure output writer: %s", err)
	}
	for _, hp := range hps {
		//update the writer to know what the index location elevation is
		e, err := hp.ProvideHazard(indexlocation)
		if err != nil {
			log.Panicf("unable to get the index location elevation: %s", err)
		}
		cw.SetAggregationElevation(e.Depth() + terrainelevation)
		sw.SetAggregationElevation(e.Depth() + terrainelevation)
		//get boundingbox
		bbox, err := hp.ProvideHazardBoundary()
		if err != nil {
			log.Panicf("unable to get the raster bounding box: %s", err)
		}
		fmt.Println(bbox.ToString())
		//loop here so that damages can be simply added across all structures
		i := 0
		maxiter := 100
		for i < maxiter {
			//modify structure random seed. this is probably an inefficent way to do this, but it should work.
			cm := make(map[string]float64)
			sm := make(map[string]float64)
			sp.ByBbox(bbox, func(f consequences.Receptor) {

				//ProvideHazard works off of a geography.Location
				d, err2 := hp.ProvideHazard(geography.Location{X: f.Location().X, Y: f.Location().Y})
				//compute damages based on hazard being able to provide depth
				if err2 == nil {
					r, err3 := f.Compute(d)
					if err3 == nil {
						damcati, err := r.Fetch("damage category")
						if err != nil {
							log.Fatal("couldnt find the damage category")
						}
						damcat := damcati.(string)
						cmtd := cm[damcat]
						smtd := sm[damcat]
						sd, err := r.Fetch("structure damage")
						if err != nil {
							log.Fatal("error fetching structure damage")
						}
						smtd += sd.(float64)
						cm[damcat] = smtd
						cd, err := r.Fetch("content damage")
						if err != nil {
							log.Fatal("error fetching content damage")
						}
						cmtd += cd.(float64)
						cm[damcat] = cmtd
					}
				}
			})
			i++
			//loop over damcat level dictionaries and write to the writer.
			for dc, td := range cm {
				header := []string{"damage category", "damage"}
				results := []interface{}{dc, td}
				r := consequences.Result{Headers: header, Result: results}
				cw.Write(r)
			}
			for dc, td := range sm {
				header := []string{"damage category", "damage"}
				results := []interface{}{dc, td}
				r := consequences.Result{Headers: header, Result: results}
				sw.Write(r)
			}
		}
	}

}
