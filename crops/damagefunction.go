package crops

import (
	"github.com/USACE/go-consequences/hazards"
)

type DamageFunction struct {
	DurationDamageCurves map[float64][]float64 //[days]][]damages by month in percents
}

func (df DamageFunction) ComputeDamagePercent(h hazards.ArrivalandDurationEvent) float64 {
	//find the duration curves above and below the duration of the event
	previousKey := 0.0
	previousValue := make([]float64, 12)
	firstIteration := true
	//find the month of the event
	hazardMonth := h.ArrivalTime.Month() //iota "enum"
	hazardMonthIndex := int(hazardMonth) - 1
	for k, v := range df.DurationDamageCurves {
		if k > h.DurationInDays {
			if firstIteration {
				//linearly interpolate to zero?
				factor := h.DurationInDays / k
				return v[hazardMonthIndex] * factor //return the damage percentage.
			} else {
				//interpolate
				factor := (k - h.DurationInDays) / (k - previousKey)
				return previousValue[hazardMonthIndex] + factor*(v[hazardMonthIndex]-previousValue[hazardMonthIndex]) //return the damage percentage.
			}
		}
		previousKey = k
		previousValue = v
		firstIteration = false
	}
	return previousValue[hazardMonthIndex] //return the damage percentage.
}
