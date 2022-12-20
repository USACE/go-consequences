package warning

import (
	"math"
	"time"

	"github.com/HydrologicEngineeringCenter/go-statistics/paireddata"
)

type warning_system struct {
	day   warning_system_parameters
	night warning_system_parameters
}
type warning_system_parameters struct {
	b float64 //effectiveness of system.
	c float64 //effectiveness of indirect system.
}
type warning_system_set struct {
	systems []warning_system
}

//from PAI2022 Curve Gen -- 2022_09_06.xlsm
func rayleighPDF(sigma float64, t float64) float64 {
	sigma_squared := math.Pow(sigma, 2.0)
	return (t / sigma_squared) * math.Exp(-(math.Pow(t, 2.0))/(2.0*sigma_squared))
}

//https://www.researchgate.net/publication/356109540_First_Alert_or_Warning_Diffusion_Time_Estimation_for_Dam_Breaches_Controlled_Dam_Releases_and_Levee_Breaches_or_Overtopping
//section 7.2  - not used by lifesim. it uses a rayleigh distribution instead.
func (wss warning_system_set) GenerateCurve(t time.Time) paireddata.PairedData {
	//combine system b parameters
	//combine system c parameters
	bDay := 0.0
	bNight := 0.0
	cDay := 0.0
	cNight := 0.0
	for _, ws := range wss.systems {
		bDay *= ws.day.b //assumes all systems are always on.
		bNight *= ws.night.b
		cDay *= ws.day.c //assumes all systems are always on. //not specified in paper to combine them, the original intent of the author was that c was representative of community cohesiveness.
		cNight *= ws.night.c
	}
	BDay := 1 - bDay
	BNight := 1 - bNight
	CDay := 1 - cDay
	CNight := 1 - cNight

	//generate day curve
	daycurve := computeCurve(BDay, CDay)
	//generate night curve
	nightcurve := computeCurve(BNight, CNight)
	//interpolate
	return interpolateCurves(daycurve, nightcurve, t)
}

////from PAI2022 Curve Gen -- 2022_09_06.xlsm
func computeCurve(B float64, C float64) paireddata.PairedData {
	cumulative := 0.0
	timeStep := 0.0
	times := make([]float64, 0)
	percentWarned := make([]float64, 0)
	//times = append(times, timeStep)
	//percentWarned = append(percentWarned, cumulative)
	for cumulative < .99999 { //check this for better epsilons
		PUt := 1 - cumulative // percent unwarned
		cumulative += rayleighPDF(B, timeStep)*PUt + (PUt * cumulative * C)
		times = append(times, timeStep)
		percentWarned = append(percentWarned, cumulative)
		timeStep += 1.0
	}
	return paireddata.PairedData{Xvals: times, Yvals: percentWarned}
}
func interpolateCurves(day paireddata.PairedData, night paireddata.PairedData, time time.Time) paireddata.PairedData {
	return paireddata.PairedData{}
}
