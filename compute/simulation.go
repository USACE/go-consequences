package compute

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"

	"github.com/USACE/go-consequences/consequences"
	"github.com/USACE/go-consequences/geography"
	"github.com/USACE/go-consequences/hazardproviders"
	"github.com/USACE/go-consequences/hazards"
	"github.com/USACE/go-consequences/indirecteconomics"
	"github.com/USACE/go-consequences/structures"
)

// ComputeEAD takes an array of damages and frequencies and integrates the curve. we should probably refactor this into paired data as a function.
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

// ComputeSpecialEAD integrates under the damage frequency curve but does not calculate the first triangle between 1 and the first frequency.
func ComputeSpecialEAD(damages []float64, freq []float64) float64 {
	//this differs from computeEAD in that it specifically does not calculate the first triangle between 1 and the first frequency to interpolate damages to zero.
	if len(damages) != len(freq) {
		panic("frequency curve is unbalanced")
	}
	triangle := 0.0
	square := 0.0
	x1 := freq[0]
	y1 := damages[0]
	eadT := 0.0
	if len(damages) > 1 {
		for i := 1; i < len(freq); i++ {
			xdelta := x1 - freq[i]
			square = xdelta * y1
			if square != 0.0 { //we dont know where damage really begins until we see it. we can guess it is inbetween ordinates, but who knows.
				triangle = ((xdelta) * -(y1 - damages[i])) / 2.0
			} else {
				triangle = 0.0
			}
			eadT += square + triangle
			x1 = freq[i]
			y1 = damages[i]
		}
	}
	if x1 != 0.0 {
		xdelta := x1 - 0.0
		eadT += xdelta * y1 //no extrapolation, just continue damages out as if it were truth for all remaining probability.
	}
	return eadT
}
func StreamAbstract(hp hazardproviders.HazardProvider, sp consequences.StreamProvider, w consequences.ResultsWriter) {
	//get boundingbox
	fmt.Println("Getting bbox")
	bbox, err := hp.HazardBoundary()
	if err != nil {
		log.Panicf("Unable to get the raster bounding box: %s", err)
	}
	fmt.Println(bbox.ToString())
	sp.ByBbox(bbox, func(f consequences.Receptor) {
		//ProvideHazard works off of a geography.Location
		d, err2 := hp.Hazard(geography.Location{X: f.Location().X, Y: f.Location().Y})
		//compute damages based on hazard being able to provide depth
		if err2 == nil {
			r, err3 := f.Compute(d)
			if err3 == nil {
				w.Write(r)
			}
		}
	})
}
func StreamAbstractMultiFrequency(hps []hazardproviders.HazardProvider, freqs []float64, sp consequences.StreamProvider, w consequences.ResultsWriter) {
	fmt.Printf("Computing %v frequencies\n", len(hps))
	//ASSUMPTION hazard providers and frequencies are in the same order
	//ASSUMPTION ordered by most frequent to least frequent event
	//ASSUMPTION! get bounding box from largest frequency.

	largestHp := hps[len(hps)-1]
	bbox, err := largestHp.HazardBoundary()
	if err != nil {
		fmt.Print(err)
		return
	}
	//set up output tables for all frequencies.
	header := []string{"fd_id", "x", "y", "damcat", "occtype", "s EAD", "c EAD", "pop2amu65", "pop2amo65", "pop2pmu65", "pop2pmo65", "found_ht", "cb_id"}

	for _, f := range freqs {
		header = append(header, fmt.Sprintf("%2.6fS", f))
		header = append(header, fmt.Sprintf("%2.6fC", f))
		header = append(header, fmt.Sprintf("%2.6fH", f))
	}

	sp.ByBbox(bbox, func(f consequences.Receptor) {
		s, sok := f.(structures.StructureStochastic)
		if !sok {
			return
		}
		results := []interface{}{s.Name, s.X, s.Y, s.DamCat, s.OccType.Name, 0.0, 0.0, s.Pop2amu65, s.Pop2amo65, s.Pop2pmu65, s.Pop2pmo65, s.FoundHt.CentralTendency(), s.CBFips}

		sEADs := make([]float64, len(freqs))
		cEADs := make([]float64, len(freqs))
		hazarddata := make([]hazards.HazardEvent, len(freqs))
		//ProvideHazard works off of a geography.Location
		gotWet := false
		for index, hp := range hps {
			d, err := hp.Hazard(geography.Location{X: f.Location().X, Y: f.Location().Y})
			hazarddata = append(hazarddata, d)
			//compute damages based on hazard being able to provide depth

			if err == nil {
				r, err3 := f.Compute(d)
				if err3 == nil {
					gotWet = true
					sdam, err := r.Fetch("structure damage")
					if err != nil {
						//panic?
						sEADs[index] = 0.0
					} else {
						damage := sdam.(float64)
						sEADs[index] = damage
					}
					cdam, err := r.Fetch("content damage")
					if err != nil {
						//panic?
						cEADs[index] = 0.0
					} else {
						damage := cdam.(float64)
						cEADs[index] = damage
					}
				}
				results = append(results, sEADs[index])
				results = append(results, cEADs[index])
				b, _ := json.Marshal(d)
				shaz := string(b)
				results = append(results, shaz)
			} else {
				//record zeros?
				results = append(results, 0.0)
				results = append(results, 0.0)
				results = append(results, "no hazard")
			}
		}
		sEAD := ComputeSpecialEAD(sEADs, freqs)
		results[5] = sEAD
		cEAD := ComputeEAD(cEADs, freqs)
		results[6] = cEAD
		var ret = consequences.Result{Headers: header, Result: results}
		if gotWet {
			w.Write(ret)
		}

	})

}
func StreamAbstractByFIPS(FIPSCODE string, hp hazardproviders.HazardProvider, sp consequences.StreamProvider, w consequences.ResultsWriter) {
	fmt.Println("FIPS Code is " + FIPSCODE)
	sp.ByFips(FIPSCODE, func(f consequences.Receptor) {
		//ProvideHazard works off of a geography.Location
		d, err := hp.Hazard(geography.Location{X: f.Location().X, Y: f.Location().Y})
		//compute damages based on hazard being able to provide depth
		if err == nil {
			r, err3 := f.Compute(d)
			if err3 == nil {
				w.Write(r)
			}
		}
	})
}

func StreamAbstractByFIPS_WithECAM(FIPSCODE string, hp hazardproviders.HazardProvider, sp consequences.StreamProvider, w consequences.ResultsWriter) {
	fmt.Println("FIPS Code is " + FIPSCODE)
	totalCounty := make(map[string]indirecteconomics.CapitalAndLabor)
	lossCounty := make(map[string]indirecteconomics.CapitalAndLabor)
	sp.ByFips(FIPSCODE, func(f consequences.Receptor) {
		//ProvideHazard works off of a geography.Location
		s, sok := f.(structures.StructureStochastic)
		if sok {
			//parse to get county level
			d := s.SampleStructure(rand.Int63()) //this is not a good idea, it will advance the seed and change results beween ECAM and non ecam computes.
			cbfips := s.CBFips[0:5]
			c, cok := totalCounty[cbfips]
			if cok {
				c.Capital += d.ContVal + d.StructVal
				c.Labor += float64(d.Pop2pmu65) //day workers (summing labor - could sum night under as an alternative, this assumes that people cant go to work because work is damaged, if we summed night, we would be saying people cant go to work because they are displaced.)
				totalCounty[cbfips] = c
			} else {
				newc := indirecteconomics.CapitalAndLabor{Capital: d.ContVal + d.StructVal, Labor: float64(d.Pop2pmu65)}
				totalCounty[cbfips] = newc
			}
		}
		d, err := hp.Hazard(geography.Location{X: f.Location().X, Y: f.Location().Y})
		//compute damages based on hazard being able to provide depth
		if err == nil {
			r, err3 := f.Compute(d)
			if err3 == nil {
				//we know it is a structure, so just jump to the values (unsafe operation, data structure of results subject to change)
				cbfips := r.Result[12].(string)[0:5]
				c, cok := lossCounty[cbfips]
				duration := .5
				if d.Has(hazards.Duration) {
					duration = d.Duration()
				}
				if cok {
					c.Capital += r.Result[7].(float64) + r.Result[6].(float64)
					c.Labor += float64(r.Result[10].(int32)) * duration //day workers (summing labor - could sum night under as an alternative, this assumes that people cant go to work because work is damaged, if we summed night, we would be saying people cant go to work because they are displaced.)
					lossCounty[cbfips] = c
				} else {
					newc := indirecteconomics.CapitalAndLabor{Capital: r.Result[7].(float64) + r.Result[6].(float64), Labor: float64(r.Result[10].(int32)) * duration}
					lossCounty[cbfips] = newc
				}
				w.Write(r)
			}
		}
	})
	//create loss ratios!
	for k, v := range lossCounty {
		statefips := k[0:2]
		countyfips := k[2:5]
		tc, tcok := totalCounty[k]
		if tcok {

			llr := float64(float64(v.Labor) / float64(tc.Labor))
			clr := v.Capital / tc.Capital
			if llr > 0 {
				if clr > 0 {
					fmt.Println(fmt.Sprintf("capital loss ratio %f", clr))
					fmt.Println(fmt.Sprintf("labor loss ratio %f", llr))
					er, err := indirecteconomics.ComputeEcam(statefips, countyfips, clr, llr)
					if err != nil {
						fmt.Println("Couldnt compute ECAM for " + k)
					} else {
						fmt.Println(er) //computed ecam!
					}
				} else {
					fmt.Printf("Couldnt compute ECAM for %v, Capital loss ratio was %f\n", k, clr)
				}
			} else {
				fmt.Printf("Couldnt compute ECAM for %v, labor loss ratio was %f\n", k, llr)
			}
		}
	}
}
