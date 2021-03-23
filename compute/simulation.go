package compute

import (
	"fmt"
	"io"
	"log"
	"time"

	"github.com/USACE/go-consequences/consequences"
	"github.com/USACE/go-consequences/geography"
	"github.com/USACE/go-consequences/hazardproviders"
	"github.com/USACE/go-consequences/hazards"
	"github.com/USACE/go-consequences/structureprovider"
	"github.com/USACE/go-consequences/structures"
)

//RequestArgs describes the request for a compute
type RequestArgs struct {
	Args       interface{}
	Concurrent bool
}

//FipsCodeCompute describes a fips code based compute with the hazardArgs
type FipsCodeCompute struct {
	ID         string      `json:"id"`
	FIPS       string      `json:"fips"`
	HazardArgs interface{} `json:"hazardargs"`
}

//BboxCompute describes a boundingbox based compute with an argument for the hazard args.
type BboxCompute struct {
	ID         string      `json:"id"`
	BBOX       string      `json:"bbox"`
	HazardArgs interface{} `json:"hazardargs"`
}

//NSIStructureSimulation is a structure that takes a requestargs and implements the computable interface.
type NSIStructureSimulation struct {
	RequestArgs
	//StructureSimulation
}

//Computable is an interface that describes the ability for an object to compute or compute by streaming to produce a simulation summary.
type Computable interface {
	Compute(args RequestArgs) SimulationSummary
	ComputeStream(args RequestArgs) SimulationSummary
}

//SimulationSummaryRow describes the result from a simulation for a row, the row header describes what the row means, and damages are provided in terms of count, and damage for structure and content
type SimulationSummaryRow struct {
	RowHeader       string  `json:"rowheader"`
	StructureCount  int64   `json:"structurecount"`
	StructureDamage float64 `json:"structuredamage"`
	ContentDamage   float64 `json:"contentdamage"`
}

//SimulationSummary is a struct that keeps a list of simulation rows and timing information about the compute.
type SimulationSummary struct {
	ColumnNames []string               `json:"columnnames"`
	Rows        []SimulationSummaryRow `json:"rows"`
	NSITime     time.Duration
	Computetime time.Duration
}

func nsiInventorytoStructures(i structureprovider.NsiInventory) []structures.StructureStochastic {
	m := structures.OccupancyTypeMap()
	defaultOcctype := m["RES1-1SNB"]
	structures := make([]structures.StructureStochastic, len(i.Features))
	for idx, feature := range i.Features {
		structures[idx] = structureprovider.NsiFeaturetoStructure(feature, m, defaultOcctype)
	}
	return structures
}

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
func FromFile(filepath string) (string, error) {
	//open a tif reader
	tiffReader := hazardproviders.Init(filepath)
	defer tiffReader.Close()
	return compute(&tiffReader)

}
func compute(hp hazardproviders.HazardProvider) (string, error) {
	//get boundingbox
	fmt.Println("Getting bbox")
	bbox, err := hp.ProvideHazardBoundary()
	if err != nil {
		log.Panicf("Unable to get the raster bounding box: %s", err)
	}
	fmt.Println(bbox.ToString())
	//get a map of all occupancy types
	m := structures.OccupancyTypeMap()
	//define a default occtype in case of emergancy
	defaultOcctype := m["RES1-1SNB"]
	//create a results store
	header := []string{"fd_id", "x", "y", "depth", "structure damage", "content damage", "Pop_2amo65", "Pop_2amu65", "Pop_2pmo65", "Pop_2pmu65"}
	var rows []interface{}
	result := consequences.Results{IsTable: true}
	result.Result.Headers = header
	result.Result.Result = rows
	structureprovider.GetByBboxStream(bbox.ToString(), func(f structureprovider.NsiFeature) {
		//convert nsifeature to structure
		str := structureprovider.NsiFeaturetoStructure(f, m, defaultOcctype)
		//query input tiff for xy location
		d, _ := hp.ProvideHazard(geography.Location{X: str.X, Y: str.Y})
		//compute damages based on provided depths
		if d.Has(hazards.Depth) {
			//fmt.Println(fmt.Sprintf("Depth was %f at structure %s", d.Depth(), f.Properties.Name))
			if d.Depth() > 0.0 {
				r := str.Compute(d)
				//keep a summmary of damages that adds the structure name
				row := []interface{}{r.Result[0], r.Result[1], r.Result[2], r.Result[3], r.Result[4], r.Result[5], f.Properties.Pop2amo65, f.Properties.Pop2amu65, f.Properties.Pop2pmo65, f.Properties.Pop2pmu65}
				structureResult := consequences.Result{Headers: header, Result: row}
				result.AddResult(structureResult)
			}
		}
	})
	b, _ := result.MarshalJSON() //json.Marshal(result)
	return string(b), nil
	//fmt.Println(string(b))
	//fmt.Println(result)
}
func StreamFromFile(filepath string, w io.Writer) { //enc json.Encoder) { //w http.ResponseWriter) {
	//open a tif reader
	tiffReader := hazardproviders.Init(filepath)
	defer tiffReader.Close()
	computeStream(&tiffReader, w)

}
func computeStream(hp hazardproviders.HazardProvider, w io.Writer) { //enc json.Encoder){//w http.ResponseWriter) {
	//get boundingbox
	fmt.Println("Getting bbox")
	bbox, err := hp.ProvideHazardBoundary()
	if err != nil {
		log.Panicf("Unable to get the raster bounding box: %s", err)
	}
	fmt.Println(bbox.ToString())
	//get a map of all occupancy types
	m := structures.OccupancyTypeMap()
	//define a default occtype in case of emergancy
	defaultOcctype := m["RES1-1SNB"]
	//create a results store
	header := []string{"fd_id", "x", "y", "depth", "structure damage", "content damage", "Pop_2amo65", "Pop_2amu65", "Pop_2pmo65", "Pop_2pmu65"}
	var rows []interface{}
	result := consequences.Results{IsTable: true}
	result.Result.Headers = header
	result.Result.Result = rows
	structureprovider.GetByBboxStream(bbox.ToString(), func(f structureprovider.NsiFeature) {
		//convert nsifeature to structure
		str := structureprovider.NsiFeaturetoStructure(f, m, defaultOcctype)
		//query input tiff for xy location
		d, _ := hp.ProvideHazard(geography.Location{X: str.X, Y: str.Y})
		//compute damages based on provided depths
		if d.Has(hazards.Depth) {
			//fmt.Println(fmt.Sprintf("Depth was %f at structure %s", d.Depth(), f.Properties.Name))
			if d.Depth() > 0.0 {
				r := str.Compute(d)
				//keep a summmary of damages that adds the structure name
				row := []interface{}{r.Result[0], r.Result[1], r.Result[2], r.Result[3], r.Result[4], r.Result[5], f.Properties.Pop2amo65, f.Properties.Pop2amu65, f.Properties.Pop2pmo65, f.Properties.Pop2pmu65}
				structureResult := consequences.Result{Headers: header, Result: row}
				b, _ := structureResult.MarshalJSON()
				s := string(b) + "\n"
				fmt.Fprintf(w, s)
			}
		}
	})
}
func StreamFromFileAbstract(filepath string, sp structureprovider.StreamProvider, w io.Writer) { //enc json.Encoder) { //w http.ResponseWriter) {
	//open a tif reader
	tiffReader := hazardproviders.Init(filepath)
	defer tiffReader.Close()
	computeStreamAbstract(&tiffReader, sp, w)

}
func computeStreamAbstract(hp hazardproviders.HazardProvider, sp structureprovider.StreamProvider, w io.Writer) {
	//get boundingbox
	fmt.Println("Getting bbox")
	bbox, err := hp.ProvideHazardBoundary()
	if err != nil {
		log.Panicf("Unable to get the raster bounding box: %s", err)
	}
	fmt.Println(bbox.ToString())
	sp.ByBbox(bbox, func(f consequences.Receptor) {
		//ProvideHazard works off of a geography.Location
		d, _ := hp.ProvideHazard(geography.Location{X: f.Location().X, Y: f.Location().Y})
		//compute damages based on hazard being able to provide depth
		if d.Has(hazards.Depth) {
			if d.Depth() > 0.0 {
				r := f.Compute(d)
				b, _ := r.MarshalJSON()
				s := string(b) + "\n"
				fmt.Fprintf(w, s)
			}
		}
	})
}
