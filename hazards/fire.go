package hazards

type FireEvent struct {
	Intensity
}
type Intensity int

const (
	Low    Intensity = iota //0
	Medium Intensity = iota // 1
	High   Intensity = iota // 2
)

type FireDamageFunction struct {
}

func (f FireDamageFunction) SampleValue(inputValue interface{}) float64 {
	input, ok := inputValue.(Intensity)
	if !ok {
		return 0.0
	}
	if input == Low {
		return 33.3
	}
	if input == Medium {
		return 50.0
	}
	if input == High {
		return 100.0
	}
	return 0.0
}
