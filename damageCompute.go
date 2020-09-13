package main
import(
	"sort"
	"fmt"
	"strings"
)
type depthEvent struct{
	depth float64
}
type occupancyType struct{
	name string
	damfun pairedData
	contentdamfun pairedData
}
type pairedData struct{
	xvals []float64
	yvals []float64
}
type valueSampler interface{
	sampleValue(inputValue float64) float64
}
func (p pairedData) sampleValue(xval float64) float64{
	if xval < p.xvals[0]{
		return 0.0 //xval is less than lowest x value
	}
	size := len(p.xvals)
	if xval >= p.xvals[size-1]{
		return p.yvals[size-1] //xval yeilds largest y value
	}
	if xval == p.xvals[0]{
		return p.yvals[0]
	}
	upper := sort.SearchFloat64s(p.xvals,xval)
	//interpolate
	lower := upper - 1 // safe because we trapped the 0 case earlier
	slope := (p.yvals[upper] - p.yvals[lower])/(p.xvals[upper] - p.xvals[lower])
	a := p.yvals[lower]
	return a + slope* (xval-p.xvals[lower])
}
type structure struct{
	occType occupancyType
	damCat string
	structVal, contVal, foundHt float64
}
type structureDamageResult struct{
	structureDamage, contentDamage float64
}
type consequenceDamageResult struct{
	headers []string
	results []interface{}
}
type consequenceReceptor interface{
	computeConsequences(event interface{}) consequenceDamageResult
}
func (s structure) computeConsequences(d interface{}) consequenceDamageResult {
	header := []string{"structure damage", "content damage"}
	results :=[]interface{}{0.0,0.0}
	var ret = consequenceDamageResult{headers:header, results:results}
	de, ok := d.(depthEvent)
	if ok{
		depth := de.depth
		depthAboveFFE := depth - s.foundHt
		damagePercent := s.occType.damfun.sampleValue(depthAboveFFE)/100 //assumes what type the damage array is in
		cdamagePercent := s.occType.contentdamfun.sampleValue(depthAboveFFE)/100
		ret.results[0] = damagePercent*s.structVal
		ret.results[1] = cdamagePercent*s.contVal
		return ret
	}else{
		return ret
	}
}
func (c consequenceDamageResult) String() string{
	if len(c.headers)!=len(c.results){
		return "mismatched lengths"
	}
	var ret string = "the consequences were:"
	for i, h := range c.headers{
		ret += " " + h + " = " + fmt.Sprintf("%f",c.results[i].(float64)) + ","
	}
	return strings.Trim(ret, ",")
}
func main(){

	//fake data to test
	xs := []float64{1.0,2.0,3.0,4.0}
	ys := []float64{10.0,20.0,30.0,40.0}
	cxs := []float64{1.0,2.0,3.0,4.0}
	cys := []float64{5.0,10.0,15.0,20.0}
	var dfun = pairedData{xvals:xs, yvals:ys}
	var cdfun = pairedData{xvals:cxs, yvals:cys}
	var o = occupancyType{name:"test",damfun:dfun,contentdamfun:cdfun}
	var s = structure{occType:o,damCat:"category",structVal:100.0, contVal:10.0, foundHt:0.0}
	var d = depthEvent{depth:3.0}

	//simplified compute
	ret := s.computeConsequences(d)
	fmt.Println("for a depth of", d.depth, ret)

	d.depth = 0.0 // test lower case
	ret = s.computeConsequences(d)
	fmt.Println("for a depth of", d.depth,ret)

	d.depth = .5 // should return 0
	ret = s.computeConsequences(d)
	fmt.Println("for a depth of", d.depth, ret)

	d.depth = 1.0 // test lowest valid case
	ret = s.computeConsequences(d)
	fmt.Println("for a depth of", d.depth, ret)

	d.depth = 1.0001 // test lowest interp case
	ret = s.computeConsequences(d)
	fmt.Println("for a depth of", d.depth, ret)

	d.depth = 2.25 //test interpolation case
	ret = s.computeConsequences(d)
	fmt.Println("for a depth of", d.depth, ret)
	
	d.depth = 2.5 //test interpolation case
	ret = s.computeConsequences(d)
	fmt.Println("for a depth of", d.depth, ret)

	d.depth = 2.75 //test interpolation case
	ret = s.computeConsequences(d)
	fmt.Println("for a depth of", d.depth, ret)

	d.depth = 3.99 // test highest interp case
	ret = s.computeConsequences(d)
	fmt.Println("for a depth of", d.depth, ret)

	d.depth = 4.0 // test highest valid case
	ret = s.computeConsequences(d)
	fmt.Println("for a depth of", d.depth, ret)
	
	d.depth = 5.0 //test upper case
	ret = s.computeConsequences(d)
	fmt.Println("for a depth of", d.depth, ret)

	s.foundHt = 1.1 //test interpolation due to foundation height putting depth back in range
	ret = s.computeConsequences(d)
	fmt.Println("for a depth of", d.depth, ret)
}
