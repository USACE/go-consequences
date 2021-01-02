package paireddata

//ValueSampler is an interface that enforces SampleValue which returns a float64 this is intended to represent any parameter in a struct if it can be deterministic or stochastic.
type ValueSampler interface {
	SampleValue(inputValue interface{}) float64
}

//UncertaintyValueSamplerSampler provides an input random number and produces a value sampler (such as a paired data)
type UncertaintyValueSamplerSampler interface {
	SampleValueSampler(randomValue float64) ValueSampler
}
