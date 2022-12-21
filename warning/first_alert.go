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
	daycurve := ComputeCurve(BDay, CDay)
	//generate night curve
	nightcurve := ComputeCurve(BNight, CNight)
	//interpolate
	return interpolateCurves(daycurve, nightcurve, t)
}

////from PAI2022 Curve Gen -- 2022_09_06.xlsm
func ComputeCurve(B float64, C float64) paireddata.PairedData {
	cumulative := 0.0
	timeStep := 0.0
	times := make([]float64, 0)
	percentWarned := make([]float64, 0)
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
	dayLength := day.Xvals[len(day.Xvals)-1]
	nightLength := night.Xvals[len(night.Xvals)-1]
	maxDuration := math.Max(dayLength, nightLength)
	xvals := make([]float64, int(maxDuration))
	yvals := make([]float64, int(maxDuration))
	minOfDay := (time.Hour() * 60.0) + time.Minute() //+1?
	//since people are not on the same consistent schedule there is some time in the morning (dawn) and evening (dusk)
	//that a mixture of people are in "night" mode and a mixture of people are in "day" mode in terms of their ability
	//to be alerted to some disaster. In those time zones we interpolate between the curves.
	//timeframe of warning diffusion
	// 12PM - night - 5AM - *dawn* - 8AM - day - 8PM - *dusk* - 11PM - night -  12PM
	if minOfDay <= 300 && minOfDay > 1380 { //11pm to 5am
		//night
		return night
	} else if minOfDay >= 480 && minOfDay < 1200 { //8am to 8pm
		//day
		return day
	} else {
		//interpolate
		if minOfDay < 480 {
			//dawn
			slope := float64(minOfDay-300) / (480 - 300)
			for i := 0; i < int(maxDuration); i++ {
				floati := float64(i)
				xvals[i] = floati
				nighty := night.SampleValue(floati)
				yval := nighty + slope*(day.SampleValue(floati)-nighty)
				yvals = append(yvals, yval)
			}
			return paireddata.PairedData{Xvals: xvals, Yvals: yvals}
		} else {
			//dusk
			slope := float64(minOfDay-1200) / (1380 - 1200)
			for i := 0; i < int(maxDuration); i++ {
				floati := float64(i)
				xvals[i] = floati
				dayy := day.SampleValue(floati)
				yval := dayy + slope*(night.SampleValue(floati)-dayy)
				yvals = append(yvals, yval)
			}
			return paireddata.PairedData{Xvals: xvals, Yvals: yvals}
		}
	}

}
