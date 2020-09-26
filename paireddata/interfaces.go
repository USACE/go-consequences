package paireddata

type ValueSampler interface {
	SampleValue(inputValue interface{}) float64
}

type UncertaintyValueSampler interface {
	SampleValue(inputValue interface{}, randomValue float64) float64
}
type UncertaintyValueSamplerSampler interface {
	SampleValueSampler(randomValue float64) ValueSampler
}
