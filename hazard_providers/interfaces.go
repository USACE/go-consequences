package hazard_providers

type HazardProvider interface {
	ProvideHazard(args interface{}) interface{}
}
