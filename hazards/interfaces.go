package hazards

import (
	"time"
)

//HazardEvent is an interface I am trying to make to describe all Hazard Events
type HazardEvent interface {
	//parameters?
	Depth() float64
	Velocity() float64
	ArrivalTime() time.Time
	ArrivalTime2ft() time.Time
	Duration() float64
	WaveHeight() float64
	Salinity() bool
	//values?
	//hazardType?
	Parameters() Parameter
	Has(p Parameter) bool
}

//Parameter is a bitflag enum https://github.com/yourbasic/bit a possible place to expand the set of hazards
type Parameter byte

//Parameter types describe different parameters for hazards
const (
	Default        Parameter = 0  //0
	Depth          Parameter = 1  //1
	Velocity       Parameter = 2  //2
	ArrivalTime    Parameter = 4  //3
	ArrivalTime2ft Parameter = 8  //4
	Duration       Parameter = 16 //5
	WaveHeight     Parameter = 32 //6
	Salinity       Parameter = 64 //7
	//next parameter
	//fin

)

//SetHasDepth turns on a bitflag for the Parameter Depth
func SetHasDepth(h Parameter) Parameter {
	return h | Depth
}

//SetHasVelocity turns on a bitflag for the Parameter Velocity
func SetHasVelocity(h Parameter) Parameter {
	return h | Velocity
}

//SetHasArrivalTime turns on a bitflag for the Parameter Arrival Time
func SetHasArrivalTime(h Parameter) Parameter {
	return h | ArrivalTime
}

//SetHasArrivalTime2ft turns on a bitflag for the Parameter ArrivalTime2ft
func SetHasArrivalTime2ft(h Parameter) Parameter {
	return h | ArrivalTime2ft
}

//SetHasDuration turns on a bitflag for the Parameter Duration
func SetHasDuration(h Parameter) Parameter {
	return h | Duration
}

//SetHasWaveHeight turns on a bitflag for the Parameter WaveHeight
func SetHasWaveHeight(h Parameter) Parameter {
	return h | WaveHeight
}

//SetHasSalinity turns on a bitflag for the Parameter Salinity
func SetHasSalinity(h Parameter) Parameter {
	return h | Salinity
}
func (p Parameter) String() string {
	s := ""
	count := 0
	if p&Depth != 0 {
		s += "Depth"
		count++
	}
	if p&ArrivalTime != 0 {
		if count > 0 {
			s += ", "
		}
		s += "Arrival Time"

		count++
	}
	if p&ArrivalTime2ft != 0 {
		if count > 0 {
			s += ", "
		}
		s += "Arrival Time 2ft"

		count++
	}
	if p&Velocity != 0 {
		if count > 0 {
			s += ", "
		}
		s += "Velocity"

		count++
	}
	if p&Duration != 0 {
		if count > 0 {
			s += ", "
		}
		s += "Duration"

		count++
	}
	if p&WaveHeight != 0 {
		if count > 0 {
			s += ", "
		}
		s += "WaveHeight"

		count++
	}
	if p&Salinity != 0 {
		if count > 0 {
			s += ", "
		}
		s += "Salinity"

		count++
	}
	return s
}
