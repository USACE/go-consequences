package paireddata

type ValueSampler interface {
	SampleValue(inputValue interface{}) float64
}

type UncertaintyValueSamplerSampler interface {
	SampleValueSampler(randomValue float64) ValueSampler
}
