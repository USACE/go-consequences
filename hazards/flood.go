package hazards

import "time"

type DepthEvent struct {
	Depth     float64
	parameter Parameter
}

func (h DepthEvent) Has(p Parameter) bool {
	return h.parameter&p != 0
}

type ArrivalandDurationEvent struct {
	ArrivalTime    time.Time
	DurationInDays int
}
