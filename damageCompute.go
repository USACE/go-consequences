package main
import(
	"sort"
	"fmt"
)
type depthEvent struct{
	depth float64
}

type occupancyType struct{
	name string
	damfun pairedData
}

type pairedData struct{
	xvals []float64
	yvals []float64
}

type structure struct{
	occType occupancyType
	damCat string
	structVal, contVal, foundHt float64
}
func main(){

	//fake data to test
	xs := []float64{1.0,2.0,3.0,4.0}
	ys := []float64{10.0,20.0,30.0,40.0}
	var dfun = pairedData{xvals:xs, yvals:ys}
	var o = occupancyType{name:"test",damfun:dfun}
	var s = structure{occType:o,damCat:"category",structVal:100.0, contVal:100.0, foundHt:0.0}
	var d = depthEvent{depth:3.0}

	//simplified compute
	ret := computeStructureDamageAtStructure(s,d)
	fmt.Println(ret)
}
func computeStructureDamageAtStructure(s structure,d depthEvent) float64{
	lower := sort.SearchFloat64s(s.occType.damfun.xvals,d.depth)
	fmt.Println(lower)
	return (s.occType.damfun.yvals[lower]/100.0)*s.structVal
}