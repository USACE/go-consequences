package hazardproviders

import (
	"fmt"

	"github.com/USACE/go-consequences/geography"
	"github.com/USACE/go-consequences/hazards"
)

type HazardProviderParameterAndPath struct {
	Hazard   hazards.Parameter `json:"hazard_parameter_type"`
	FilePath string            `json:"hazard_provider_file_path"` //this should get fixed to be able to represent more complex information. e.g. what parameter?
}
type HazardProviderInfo struct {
	Hazards []HazardProviderParameterAndPath `json:"hazards"`
}

// HazardProvider provides hazards as a return for an argument input
type HazardProvider interface {
	Hazard(location geography.Location) (hazards.HazardEvent, error)
	//ProcessedHazard(location geography.Location, process HazardFunction) (hazards.HazardEvent, error)
	HazardBoundary() (geography.BBox, error)
	Close()
}
type HazardFunction func(valueIn hazards.HazardData, hazard hazards.HazardEvent) (hazards.HazardEvent, error)

func DepthHazardFunction() HazardFunction {
	return func(valueIn hazards.HazardData, hazard hazards.HazardEvent) (hazards.HazardEvent, error) {
		d := hazards.DepthEvent{}
		d.SetDepth(valueIn.Depth)
		return d, nil
	}
}
func ArrivalAndDurationHazardFunction() HazardFunction {
	return func(valueIn hazards.HazardData, hazard hazards.HazardEvent) (hazards.HazardEvent, error) {
		d := hazards.ArrivalDepthandDurationEvent{}
		d.SetDuration(valueIn.Duration)
		d.SetArrivalTime(valueIn.ArrivalTime)
		return d, nil
	}
}

// NoHazardFoundError is an error for a situation where no hazard could be computed for the given args
type NoHazardFoundError struct {
	Input string
}

// NoDataHazardError is an error for a situation where no hazard could be computed for the given args
type NoDataHazardError struct {
	Input string
}

// NoFrequencyFoundError is an error for a situation where no frequency could be associated for the hazard for the given args
type NoFrequencyFoundError struct {
	Input string
}

// HazardError is an error for a generic hazarderror for the given args
type HazardError struct {
	Input string
}

// Error implements the error interface for NoHazardFoundError
func (h NoHazardFoundError) Error() string {
	return fmt.Sprintf("No hazard Found for %s", h.Input)
}

// Error implements the error interface for NoDataHazardError
func (h NoDataHazardError) Error() string {
	return fmt.Sprintf("Location yeilded No Data for a hazard. %s", h.Input)
}

// Error implements the error interface for NoFrequencyFoundError
func (h NoFrequencyFoundError) Error() string {
	return fmt.Sprintf("No frequency Found for %s", h.Input)
}

// Error implements the error interface for HazardError
func (h HazardError) Error() string {
	return fmt.Sprintf("Could not compute because: %s", h.Input)
}
