package hazard_providers

import "fmt"

type HazardProvider interface {
	ProvideHazard(args interface{}) (interface{}, error)
}
type NoHazardFoundError struct {
	Input string
}
type NoFrequencyFoundError struct {
	Input string
}
type HazardError struct {
	Input string
}

func (h NoHazardFoundError) Error() string {
	return fmt.Sprintf("No hazard Found for %s", h.Input)
}
func (h NoFrequencyFoundError) Error() string {
	return fmt.Sprintf("No frequency Found for %s", h.Input)
}
func (h HazardError) Error() string {
	return fmt.Sprintf("Could not compute because: %s", h.Input)
}
