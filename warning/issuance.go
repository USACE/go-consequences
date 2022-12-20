package warning

import (
	"math"
)

type lindellDist struct {
	a float64
	b float64
}

/*
   Public Overrides Function Validate() As String
       If _A < 0.1 Then Return "Lindell Distribution Error: A value must be greater than or equal to 0.1."
       If _A > 1.8 Then Return "Lindell Distribution Error: A value must be less than or equal to 1.8."
       If _B < 0.5 Then Return "Lindell Distribution Error: B value must be greater than or equal to 0.5."
       If _B > 2.1 Then Return "Lindell Distribution Error: B value must be less than or equal to 2.1."
       Return Nothing
   End Function
*/
func (dist lindellDist) CDF(value float64) float64 {
	if value <= 0 {
		return 0.0
	}
	if value >= 360 {
		return 1.0
	} //seems dangerous.
	return 1 - (math.Exp(math.Pow(-dist.a*(value/60), dist.b)))
}
func (dist lindellDist) Sample(probability float64) float64 {
	if probability <= 0 {
		return 0.0
	}
	if probability >= 1 {
		return 360.0
	} //fix this.
	return math.Min(-60*math.Pow(math.Log(1-probability)/dist.a, 1/dist.b), 360) //test. //360 seems arbitrary, but it is what we have in our notes
}
