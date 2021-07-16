package hazards

import (
	"fmt"
	"time"
)

//DepthEvent describes a Hazard with Depth Only
type DepthEvent struct {
	depth float64
}

func (h DepthEvent) Depth() float64 {
	return h.depth
}
func (h *DepthEvent) SetDepth(d float64) {
	//fmt.Println(d)
	h.depth = d
}
func (h DepthEvent) Velocity() float64 {
	return -901.0
}
func (h DepthEvent) ArrivalTime() time.Time {
	return time.Time{}
}
func (h DepthEvent) ArrivalTime2ft() time.Time {
	return time.Time{}
}
func (h DepthEvent) Duration() float64 {
	return -901.0
}
func (h DepthEvent) WaveHeight() float64 {
	return -901.0
}
func (h DepthEvent) Salinity() bool {
	return false
}
func (h DepthEvent) Qualitative() string {
	return ""
}

//Parameters implements the HazardEvent interface
func (h DepthEvent) Parameters() Parameter {
	dp := Default
	dp = SetHasDepth(dp)
	return dp
}

//Has implements the HazardEvent Interface
func (h DepthEvent) Has(p Parameter) bool {
	dp := h.Parameters()
	return dp&p != 0
}
func (d DepthEvent) MarshalJSON() ([]byte, error) {
	s := fmt.Sprintf("{\"depthevent\":{\"depth\":%f}}", d.Depth())
	return []byte(s), nil
}

//ArrivalandDurationEvent describes an event with an arrival time and a duration in days
type ArrivalandDurationEvent struct {
	arrivalTime    time.Time
	durationInDays float64
}

func (d ArrivalandDurationEvent) MarshalJSON() ([]byte, error) {
	s := fmt.Sprintf("{\"arrivalanddurationevent\":{\"arrivaltime\":%s,\"duration\":%f}}", d.ArrivalTime().Format("Jan _2 15:04"), d.Duration())
	return []byte(s), nil
}
func (h ArrivalandDurationEvent) Depth() float64 {
	return -901.0
}
func (h ArrivalandDurationEvent) Velocity() float64 {
	return -901.0
}
func (h *ArrivalandDurationEvent) SetArrivalTime(t time.Time) {
	h.arrivalTime = t
}
func (h ArrivalandDurationEvent) ArrivalTime() time.Time {
	return h.arrivalTime
}
func (h ArrivalandDurationEvent) ArrivalTime2ft() time.Time {
	return time.Time{}
}
func (h ArrivalandDurationEvent) Duration() float64 {
	return h.durationInDays
}
func (h *ArrivalandDurationEvent) SetDuration(d float64) {
	h.durationInDays = d
}
func (h ArrivalandDurationEvent) WaveHeight() float64 {
	return -901.0
}
func (h ArrivalandDurationEvent) Salinity() bool {
	return false
}
func (h ArrivalandDurationEvent) Qualitative() string {
	return ""
}

//Parameters implements the HazardEvent interface
func (ad ArrivalandDurationEvent) Parameters() Parameter {
	adp := Default
	adp = SetHasDuration(adp)
	adp = SetHasArrivalTime(adp)
	return adp
}

//Has implements the HazardEvent Interface
func (ad ArrivalandDurationEvent) Has(p Parameter) bool {
	adp := ad.Parameters()
	return adp&p != 0
}

//ArrivalandDurationEvent describes an event with an arrival time and a duration in days
type ArrivalDepthandDurationEvent struct {
	arrivalTime    time.Time
	depth          float64
	durationInDays float64
}

func (d ArrivalDepthandDurationEvent) MarshalJSON() ([]byte, error) {
	s := fmt.Sprintf("{\"arrivaldepthanddurationevent\":{\"arrivaltime\":%s,\"depth\":%f,\"duration\":%f}}", d.ArrivalTime().Format("Jan _2 15:04"), d.Depth(), d.Duration())
	return []byte(s), nil
}
func (h *ArrivalDepthandDurationEvent) SetDepth(d float64) {
	h.depth = d
}
func (h ArrivalDepthandDurationEvent) Depth() float64 {
	return h.depth
}
func (h ArrivalDepthandDurationEvent) Velocity() float64 {
	return -901.0
}
func (h *ArrivalDepthandDurationEvent) SetArrivalTime(t time.Time) {
	h.arrivalTime = t
}
func (h ArrivalDepthandDurationEvent) ArrivalTime() time.Time {
	return h.arrivalTime
}
func (h ArrivalDepthandDurationEvent) ArrivalTime2ft() time.Time {
	return time.Time{}
}
func (h ArrivalDepthandDurationEvent) Duration() float64 {
	return h.durationInDays
}
func (h *ArrivalDepthandDurationEvent) SetDuration(d float64) {
	h.durationInDays = d
}
func (h ArrivalDepthandDurationEvent) WaveHeight() float64 {
	return -901.0
}
func (h ArrivalDepthandDurationEvent) Salinity() bool {
	return false
}
func (h ArrivalDepthandDurationEvent) Qualitative() string {
	return ""
}

//Parameters implements the HazardEvent interface
func (ad ArrivalDepthandDurationEvent) Parameters() Parameter {
	adp := Default
	adp = SetHasDuration(adp)
	adp = SetHasDepth(adp)
	adp = SetHasArrivalTime(adp)
	return adp
}

//Has implements the HazardEvent Interface
func (ad ArrivalDepthandDurationEvent) Has(p Parameter) bool {
	adp := ad.Parameters()
	return adp&p != 0
}
