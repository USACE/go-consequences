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
type Parameter byte

const (
	Default        Parameter = 0  //0
	Depth          Parameter = 1  //1
	Velocity       Parameter = 2  //2
	ArrivalTime    Parameter = 4  //3
	ArrivalTime2ft Parameter = 8  //4
	Duration       Parameter = 16 //5
	//next parameter
	//next parameter
	//next parameter
	//fin

)

func SetHasDepth(h Parameter) Parameter {
	return h | Depth
}
func SetHasVelocity(h Parameter) Parameter {
	return h | Velocity
}
func SetHasArrivalTime(h Parameter) Parameter {
	return h | ArrivalTime
}
func SetHasArrivalTime2ft(h Parameter) Parameter {
	return h | ArrivalTime2ft
}
func SetHasDuration(h Parameter) Parameter {
	return h | Duration
}

func (p Parameter) String() string {
	s := ""
	count := 0
	if p&Depth != 0 {
		s += "Depth"
		count += 1
	}
	if p&ArrivalTime != 0 {
		if count > 0 {
			s += ", "
		}
		s += "Arrival Time"

		count += 1
	}
	if p&ArrivalTime2ft != 0 {
		if count > 0 {
			s += ", "
		}
		s += "Arrival Time 2ft"

		count += 1
	}
	if p&Velocity != 0 {
		if count > 0 {
			s += ", "
		}
		s += "Velocity"

		count += 1
	}
	if p&Duration != 0 {
		if count > 0 {
			s += ", "
		}
		s += "Duration"

		count += 1
	}
	return s
}
