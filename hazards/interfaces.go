package hazards

type Hazard_Event interface {
	//parameters?
	//values?
	//hazardType?
	HazardType() Hazards
}

type Hazards interface {
	Has(p Parameter) bool
}

//https://github.com/yourbasic/bit a possible place to expand the set of hazards
type Parameter uint8

const (
	Depth          Parameter = iota //0
	Velocity       Parameter = iota // 1
	ArrivalTime    Parameter = iota //2
	ArrivalTime2ft Parameter = iota // 3
	Duration       Parameter = iota
)

func (h Parameter) SetHasDepth() {
	h = h | Depth
}
func (h Parameter) SetHasVelocity() {
	h = h | Velocity
}
func (h Parameter) SetHasArrivalTime() {
	h = h | ArrivalTime
}
func (h Parameter) SetHasArrivalTime2ft() {
	h = h | ArrivalTime2ft
}
func (h Parameter) SetHasDuration() {
	h = h | Duration
}
