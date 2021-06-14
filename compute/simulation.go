package compute

import (
	"fmt"
	"log"

	"github.com/USACE/go-consequences/consequences"
	"github.com/USACE/go-consequences/geography"
	"github.com/USACE/go-consequences/hazardproviders"
)

//ComputeEAD takes an array of damages and frequencies and integrates the curve. we should probably refactor this into paired data as a function.
func ComputeEAD(damages []float64, freq []float64) float64 {
	triangle := 0.0
	square := 0.0
	x1 := 1.0 // create a triangle to the first probability space - linear interpolation is probably a problem, maybe use log linear interpolation for the triangle
	y1 := 0.0
	eadT := 0.0
	for i := 0; i < len(freq); i++ {
		xdelta := x1 - freq[i]
		square = xdelta * y1
		triangle = ((xdelta) * (damages[i] - y1)) / 2.0
		eadT += square + triangle
		x1 = freq[i]
		y1 = damages[i]
	}
	if x1 != 0.0 {
		xdelta := x1 - 0.0
		eadT += xdelta * y1 //no extrapolation, just continue damages out as if it were truth for all remaining probability.

	}
	return eadT
}

//ComputeSpecialEAD integrates under the damage frequency curve but does not calculate the first triangle between 1 and the first frequency.
func ComputeSpecialEAD(damages []float64, freq []float64) float64 {
	//this differs from computeEAD in that it specifically does not calculate the first triangle between 1 and the first frequency to interpolate damages to zero.
	triangle := 0.0
	square := 0.0
	x1 := freq[0]
	y1 := damages[0]
	eadT := 0.0
	for i := 1; i < len(freq); i++ {
		xdelta := x1 - freq[i]
		square = xdelta * y1
		triangle = ((xdelta) * -(y1 - damages[i])) / 2.0
		eadT += square + triangle
		x1 = freq[i]
		y1 = damages[i]
	}
	if x1 != 0.0 {
		xdelta := x1 - 0.0
		eadT += xdelta * y1 //no extrapolation, just continue damages out as if it were truth for all remaining probability.

	}
	return eadT
}
func StreamFromFileAbstract(filepath string, sp consequences.StreamProvider, w consequences.ResultsWriter) { //enc json.Encoder) { //w http.ResponseWriter) {
	//open a tif reader
	tiffReader := hazardproviders.Init(filepath)
	defer tiffReader.Close()
	StreamAbstract(&tiffReader, sp, w)
	w.Close()
}
func StreamAbstract(hp hazardproviders.HazardProvider, sp consequences.StreamProvider, w consequences.ResultsWriter) {
	//get boundingbox
	fmt.Println("Getting bbox")
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
			r, err3 := f.Compute(d)
			if err3 == nil {
				w.Write(r)
			}
		}
	})
}
func StreamAbstractByFIPS(FIPSCODE string, hp hazardproviders.HazardProvider, sp consequences.StreamProvider, w consequences.ResultsWriter) {
	fmt.Println("FIPS Code is " + FIPSCODE)
	sp.ByFips(FIPSCODE, func(f consequences.Receptor) {
		//ProvideHazard works off of a geography.Location
		d, err := hp.ProvideHazard(geography.Location{X: f.Location().X, Y: f.Location().Y})
		//compute damages based on hazard being able to provide depth
		if err == nil {
			r, err3 := f.Compute(d)
			if err3 == nil {
				w.Write(r)
			}
		}
	})
}
