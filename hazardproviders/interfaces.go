package hazardproviders

import "fmt"

//HazardProvider provides hazards as a return for an argument input
type HazardProvider interface {
	ProvideHazard(args interface{}) (interface{}, error)
}

//NoHazardFoundError is an error for a situation where no hazard could be computed for the given args
type NoHazardFoundError struct {
	Input string
}

//NoFrequencyFoundError is an error for a situation where no frequency could be associated for the hazard for the given args
type NoFrequencyFoundError struct {
	Input string
}

//HazardError is an error for a generic hazarderror for the given args
type HazardError struct {
	Input string
}

//Error implements the error interface for NoHazardFoundError
func (h NoHazardFoundError) Error() string {
	return fmt.Sprintf("No hazard Found for %s", h.Input)
}

//Error implements the error interface for NoFrequencyFoundError
func (h NoFrequencyFoundError) Error() string {
	return fmt.Sprintf("No frequency Found for %s", h.Input)
}

//Error implements the error interface for HazardError
func (h HazardError) Error() string {
	return fmt.Sprintf("Could not compute because: %s", h.Input)
}
