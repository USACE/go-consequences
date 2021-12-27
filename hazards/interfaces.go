package hazards

import (
	"bytes"
	"encoding/json"
	"strings"
	"time"
)

//HazardEvent is an interface I am trying to make to describe all Hazard Events
type HazardEvent interface {
	//parameters?
	Depth() float64
	Velocity() float64
	ArrivalTime() time.Time
	Erosion() float64
	Duration() float64
	WaveHeight() float64
	Salinity() bool
	Qualitative() string
	//values?
	//hazardType?
	Parameters() Parameter
	Has(p Parameter) bool
}

//Parameter is a bitflag enum
type Parameter uint //switch to uint64 if we hit 32 slots and need another.

//Parameter types describe different parameters for hazards
const (
	Default          Parameter = 0    //0
	Depth            Parameter = 1    //1
	Velocity         Parameter = 2    //2
	ArrivalTime      Parameter = 4    //3
	ArrivalTime2ft   Parameter = 8    //4
	Erosion          Parameter = 16   //5
	Duration         Parameter = 32   //6
	WaveHeight       Parameter = 64   //7
	MediumWaveHeight Parameter = 128  //8
	HighWaveHeight   Parameter = 256  //9
	Salinity         Parameter = 512  //10
	Qualitative      Parameter = 1024 //11
)

var parametersToStrings = map[Parameter]string{
	Default:          "default",
	Depth:            "depth",
	Velocity:         "velocity",
	ArrivalTime:      "arrivaltime",
	Erosion:          "erosion",
	Duration:         "duration",
	WaveHeight:       "waveheight",
	MediumWaveHeight: "mediumwaveheight",
	HighWaveHeight:   "highwaveheight",
	Salinity:         "salinity",
	Qualitative:      "qualitative",
}

var stringsToParameters = map[string]Parameter{
	"default":          Default,
	"depth":            Depth,
	"velocity":         Velocity,
	"arrivaltime":      ArrivalTime,
	"erosion":          Erosion,
	"duration":         Duration,
	"waveheight":       WaveHeight,
	"highwaveheight":   HighWaveHeight,
	"mediumwaveheight": MediumWaveHeight,
	"salinity":         Salinity,
	"qualitative":      Qualitative,
}

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

//SetHasErosion turns on a bitflag for the Parameter Erosion
func SetHasErosion(h Parameter) Parameter {
	return h | Erosion
}

//SetHasDuration turns on a bitflag for the Parameter Duration
func SetHasDuration(h Parameter) Parameter {
	return h | Duration
}

//SetHasWaveHeight turns on a bitflag for the Parameter WaveHeight
func SetHasWaveHeight(h Parameter) Parameter {
	return h | WaveHeight
}
func SetHasMediumWaveHeight(h Parameter) Parameter {
	return h | MediumWaveHeight
}
func SetHasHighWaveHeight(h Parameter) Parameter {
	return h | HighWaveHeight
}

//SetHasSalinity turns on a bitflag for the Parameter Salinity
func SetHasSalinity(h Parameter) Parameter {
	return h | Salinity
}

//SetHasSalinity turns on a bitflag for the Parameter Salinity
func SetHasQualitative(h Parameter) Parameter {
	return h | Qualitative
}
func (p Parameter) String() string {
	s := ""
	count := 0
	if p < 1 {
		return "default"
	}
	for key, val := range parametersToStrings {
		if p&key != 0 {
			if count > 0 {
				s += ", "
			}
			s += val
			count++
		}
	}
	return s
}
func toParameter(s string) Parameter {
	parts := strings.Split(s, ", ")
	var p Parameter
	for _, sp := range parts {
		pval, found := stringsToParameters[sp]
		if found {
			p = p | pval
		}
	}
	return p
}

// MarshalJSON marshals the enum as a quoted json string
func (p Parameter) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString(`"`)
	buffer.WriteString(p.String())
	buffer.WriteString(`"`)
	return buffer.Bytes(), nil
}

// UnmarshalJSON unmashals a quoted, comma separated string to the parameter value
func (p *Parameter) UnmarshalJSON(b []byte) error {
	var s string
	err := json.Unmarshal(b, &s)
	if err != nil {
		return err
	}
	// Note that if the string cannot be found then it will be set to the zero value, 'default' in this case.
	*p = toParameter(s)
	return nil
}
