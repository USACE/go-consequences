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
type valueSampler interface{
	sampleValue(xval float64) float64
}
func (p pairedData) sampleValue(xval float64) float64{
	if xval < p.xvals[0]{
		return 0.0 //xval is less than lowest x value
	}
	size := len(p.xvals)
	if xval >= p.xvals[size-1]{
		return p.yvals[size-1] //xval yeilds largest y value
	}
	lower := sort.SearchFloat64s(p.xvals,xval)
	//interpolate
	return p.yvals[lower]
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
	fmt.Println("for a depth of", d.depth, "the damage is",ret)

	d.depth = 0.0 // test lower case
	ret = computeStructureDamageAtStructure(s,d)
	fmt.Println("for a depth of", d.depth, "the damage is",ret)

	d.depth = 1.0 // test lowest valid case
	ret = computeStructureDamageAtStructure(s,d)
	fmt.Println("for a depth of", d.depth, "the damage is",ret)

	d.depth = 2.5 //test interpolation case (not passing currently)
	ret = computeStructureDamageAtStructure(s,d)
	fmt.Println("for a depth of", d.depth, "the damage is",ret)

	d.depth = 4.0 // test highest valid case
	ret = computeStructureDamageAtStructure(s,d)
	fmt.Println("for a depth of", d.depth, "the damage is",ret)
	
	d.depth = 5.0 //test upper case
	ret = computeStructureDamageAtStructure(s,d)
	fmt.Println("for a depth of", d.depth, "the damage is",ret)
}
func computeStructureDamageAtStructure(s structure,d depthEvent) float64{
	damagePercent := s.occType.damfun.sampleValue(d.depth)/100
	return damagePercent*s.structVal
}