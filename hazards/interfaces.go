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
	Depth    Hazards = iota //0
	Velocity Hazards = iota // 1
	ArrivalTime Hazards = iota//2
	ArrivalTime2ft   Intensity = iota // 3
)

func (h *Parameter) SetHasDepth(){
	h | Depth
}
func (h *Parameter) SetHasVelocity(){
	h | Velocity
}
func (h *Parameter) SetHasArrivalTime(){
	h | ArrivalTime
}
func (h *Parameter) SetHasArrivalTime2ft(){
	h | Depth
}
func (h *Parameter) SetHasDuration(){
	h | Depth
}