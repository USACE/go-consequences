package crops

import (
	"sort"

	"github.com/USACE/go-consequences/hazards"
)

//DamageFunction manages the methodology for computing damages to crops based on a duration and seasonal damage curves
type DamageFunction struct {
	DurationDamageCurves map[float64][]float64 //[days]][]damages by month in percents
}

//ComputeDamagePercent takes an input Arrival time and duration based event and produces a duration damage based on the season of the start of the impact.
func (df DamageFunction) ComputeDamagePercent(h hazards.ArrivalandDurationEvent) float64 {
	//find the duration curves above and below the duration of the event
	previousKey := 0.0
	previousValue := make([]float64, 12)
	firstIteration := true
	//find the month of the event
	hazardMonth := h.ArrivalTime().Month() //iota "enum"
	hazardMonthIndex := int(hazardMonth) - 1
	keys := make([]float64, 0)
	for k := range df.DurationDamageCurves {
		keys = append(keys, k)
	}
	sort.Float64s(keys)
	for _, k := range keys {
		v := df.DurationDamageCurves[k]
		if k > h.Duration() {
			if firstIteration {
				//linearly interpolate to zero?
				factor := h.Duration() / k
				return v[hazardMonthIndex] * factor //return the damage percentage.
			}
			//interpolate
			factor := (k - h.Duration()) / (k - previousKey)
			return previousValue[hazardMonthIndex] + factor*(v[hazardMonthIndex]-previousValue[hazardMonthIndex]) //return the damage percentage.
		}
		previousKey = k
		previousValue = v
		firstIteration = false
	}
	return previousValue[hazardMonthIndex] //return the damage percentage.
}
