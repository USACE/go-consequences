package warning

import (
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

func (wss warning_system_set) GenerateCurve(t time.Time) paireddata.PairedData {
	//combine system b parameters
	//combine system c parameters
	//generate day curve
	//generate night curve
	//interpolate
	return paireddata.PairedData{}
}
