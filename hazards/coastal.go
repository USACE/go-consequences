package hazards

import (
	"errors"
	"fmt"
	"sort"
	"time"
)

// CoastalEvent describes a coastal event
type CoastalEvent struct {
	depth         float64 `default:"-901.0"` //still depth
	waveHeight    float64 `default:"-901.0"` //continuous variable.
	salinity      bool    //default is false
	percentEroded float64
}

func NewCoastalEvent(c CoastalEvent) *CoastalEvent {

	e := c

	if c.depth == 0.0 {
		e.depth = -901.0
	}

	if c.waveHeight == 0.0 {
		e.waveHeight = -901.0
	}

	return &e
}

func (d CoastalEvent) MarshalJSON() ([]byte, error) {
	s := fmt.Sprintf("{\"coastalevent\":{\"depth\":%f, \"waveheight\":%f,\"salinity\":%t}}", d.Depth(), d.WaveHeight(), d.Salinity())
	return []byte(s), nil
}
func (h CoastalEvent) Depth() float64 {
	return h.depth
}
func (h *CoastalEvent) SetDepth(d float64) {
	h.depth = d
}
func (h *CoastalEvent) SetErosion(e float64) {
	h.percentEroded = e
}
func (h CoastalEvent) Velocity() float64 {
	return -901.0
}
func (h CoastalEvent) ArrivalTime() time.Time {
	return time.Time{}
}
func (h CoastalEvent) Erosion() float64 {
	return h.percentEroded
}
func (h CoastalEvent) Duration() float64 {
	return -901.0
}
func (h CoastalEvent) WaveHeight() float64 {
	return h.waveHeight
}
func (h *CoastalEvent) SetWaveHeight(d float64) {
	h.waveHeight = d
}
func (h CoastalEvent) Salinity() bool {
	return h.salinity
}
func (h *CoastalEvent) SetSalinity(d bool) {
	h.salinity = d
}
func (h CoastalEvent) Qualitative() string {
	return ""
}
func (h CoastalEvent) DV() float64 {
	return -901.0
}

// Parameters implements the HazardEvent interface
func (ad CoastalEvent) Parameters() Parameter {
	adp := Default

	// -901.0 is the float64 convention for no data
	if ad.Depth() > -901.0 {
		adp = SetHasDepth(adp)
	}

	if ad.WaveHeight() > 0.0 {
		adp = SetHasWaveHeight(adp)
		if ad.WaveHeight() < 3.0 {
			adp = SetHasMediumWaveHeight(adp)
		} else {
			adp = SetHasHighWaveHeight(adp)
		}
	}

	if ad.Salinity() {
		adp = SetHasSalinity(adp)
	}

	if ad.Erosion() > 0.0 {
		adp = SetHasErosion(adp)
	}

	return adp
}

// Has implements the HazardEvent Interface
func (ad CoastalEvent) Has(p Parameter) bool {
	adp := ad.Parameters()
	return adp&p != 0
}

type CoastalFrequencyEvent struct {
	CoastalEvent
	frequency float64
}

func NewCoastalFrequencyEvent(c CoastalFrequencyEvent) *CoastalFrequencyEvent {

	e := c

	if c.depth == 0.0 {
		e.depth = -901.0
	}

	if c.waveHeight == 0.0 {
		e.waveHeight = -901.0
	}

	return &e
}

func (d CoastalFrequencyEvent) MarshalJSON() ([]byte, error) {
	s := fmt.Sprintf("{\"coastalfrequencyevent\":{\"depth\":%f, \"waveheight\":%f,\"salinity\":%t,\"frequency\":%f}}", d.Depth(), d.WaveHeight(), d.Salinity(), d.Frequency())
	return []byte(s), nil
}

func (h CoastalFrequencyEvent) Frequency() float64 {
	return h.frequency
}
func (h *CoastalFrequencyEvent) SetFrequency(f float64) {
	h.frequency = f
}

type MultiFrequencyCoastalEvent struct {
	index int
	// Frequencies []float64
	Events []CoastalFrequencyEvent
}

func (h MultiFrequencyCoastalEvent) Depth() float64 {
	return h.Events[h.index].Depth()
}

func (h MultiFrequencyCoastalEvent) Velocity() float64 {
	return h.Events[h.index].Velocity()
}

func (h MultiFrequencyCoastalEvent) ArrivalTime() time.Time {
	return h.Events[h.index].ArrivalTime()
}

func (h MultiFrequencyCoastalEvent) Erosion() float64 {
	return h.Events[h.index].Erosion()
}

func (h MultiFrequencyCoastalEvent) Duration() float64 {
	return h.Events[h.index].Duration()
}

func (h MultiFrequencyCoastalEvent) WaveHeight() float64 {
	return h.Events[h.index].WaveHeight()
}

func (h MultiFrequencyCoastalEvent) Salinity() bool {
	return h.Events[h.index].Salinity()
}

func (h MultiFrequencyCoastalEvent) Qualitative() string {
	return h.Events[h.index].Qualitative()
}

func (h MultiFrequencyCoastalEvent) DV() float64 {
	return h.Events[h.index].DV()
}

func (h MultiFrequencyCoastalEvent) Frequency() float64 {
	return h.Events[h.index].Frequency()
}

func (h MultiFrequencyCoastalEvent) Parameters() Parameter {
	return h.Events[h.index].Parameters()
}

func (h MultiFrequencyCoastalEvent) Has(p Parameter) bool {
	return h.Events[h.index].Has(p)
}

func (h MultiFrequencyCoastalEvent) Index() int {
	return h.index
}

func (h MultiFrequencyCoastalEvent) HasNext() bool {
	return h.index < (len(h.Events) - 1)
}

func (h MultiFrequencyCoastalEvent) HasPrevious() bool {
	return h.index > 0
}

func (h MultiFrequencyCoastalEvent) This() HazardEvent {
	return h.Events[h.index]
}

func (h MultiFrequencyCoastalEvent) Next() (HazardEvent, error) {
	var err error = nil
	if h.HasNext() {
		return h.Events[h.index+1], err
	} else {
		return ArrivalDepthandDurationEvent{}, errors.New("hazards: MultiFrequencyCoastalEvent does not have Next event")
	}
}

func (h MultiFrequencyCoastalEvent) Previous() (HazardEvent, error) {
	var err error = nil
	if h.HasPrevious() {
		return h.Events[h.index-1], err
	} else {
		return ArrivalDepthandDurationEvent{}, errors.New("hazards: MultiFrequencyCoastalEvent does not have Previous event")
	}
}

func (h *MultiFrequencyCoastalEvent) Increment() {
	if h.HasNext() {
		h.index++
	}
}

func (h *MultiFrequencyCoastalEvent) ResetIndex() {
	h.index = 0
}

func (h *MultiFrequencyCoastalEvent) Append(n HazardEvent) {
	newEvent := n.(CoastalFrequencyEvent)
	h.Events = append(h.Events, newEvent)
}

func (h MultiFrequencyCoastalEvent) Sort() { // ensure the hazard events are in order of arrival time
	sort.Sort(h)
}

func (h MultiFrequencyCoastalEvent) IsSorted() bool {
	return sort.IsSorted(h)
}

// Len is part of sort.Interface.
func (h MultiFrequencyCoastalEvent) Len() int {
	return len(h.Events)
}

// Swap is part of sort.Interface.
func (h MultiFrequencyCoastalEvent) Swap(i, j int) {
	h.Events[i], h.Events[j] = h.Events[j], h.Events[i]
}

func (h MultiFrequencyCoastalEvent) Less(i, j int) bool {
	return h.Events[i].Frequency() < h.Events[j].Frequency() // This means the 500-year flood is "Less" than the 100-year event because we are sorting on frequency
}

func (h *MultiFrequencyCoastalEvent) SetIndex(i int) error {
	if i < (len(h.Events)) {
		h.index = i
		return nil
	}
	return errors.New("hazards: Attempted to set out of bounds index on DepthEventMultiFrequency event.")
}
func (h MultiFrequencyCoastalEvent) HazardFrequencies() []float64 {
	f := make([]float64, len(h.Events))
	for i, e := range h.Events {
		f[i] = e.Frequency()
	}
	return f
}
