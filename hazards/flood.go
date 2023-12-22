package hazards

import (
	"fmt"
	"time"
)

// DepthEvent describes a Hazard with Depth Only
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
func (h DepthEvent) Erosion() float64 {
	return -901.0
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
func (h DepthEvent) DV() float64 {
	return -901.0
}

// Parameters implements the HazardEvent interface
func (h DepthEvent) Parameters() Parameter {
	dp := Default
	dp = SetHasDepth(dp)
	return dp
}

// Has implements the HazardEvent Interface
func (h DepthEvent) Has(p Parameter) bool {
	dp := h.Parameters()
	return dp&p != 0
}
func (d DepthEvent) MarshalJSON() ([]byte, error) {
	s := fmt.Sprintf("{\"depthevent\":{\"depth\":%f}}", d.Depth())
	return []byte(s), nil
}

// ArrivalandDurationEvent describes an event with an arrival time and a duration in days
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
func (h ArrivalandDurationEvent) Erosion() float64 {
	return -901.0
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
func (h ArrivalandDurationEvent) DV() float64 {
	return -901.0
}

// Parameters implements the HazardEvent interface
func (ad ArrivalandDurationEvent) Parameters() Parameter {
	adp := Default
	adp = SetHasDuration(adp)
	adp = SetHasArrivalTime(adp)
	return adp
}

// Has implements the HazardEvent Interface
func (ad ArrivalandDurationEvent) Has(p Parameter) bool {
	adp := ad.Parameters()
	return adp&p != 0
}

// ArrivalandDurationEvent describes an event with an arrival time, depth and a duration in days
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
func (h ArrivalDepthandDurationEvent) Erosion() float64 {
	return -901.0
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
func (h ArrivalDepthandDurationEvent) DV() float64 {
	return -901.0
}

// Parameters implements the HazardEvent interface
func (ad ArrivalDepthandDurationEvent) Parameters() Parameter {
	adp := Default
	adp = SetHasDuration(adp)
	adp = SetHasDepth(adp)
	adp = SetHasArrivalTime(adp)
	return adp
}

// Has implements the HazardEvent Interface
func (ad ArrivalDepthandDurationEvent) Has(p Parameter) bool {
	adp := ad.Parameters()
	return adp&p != 0
}

// ArrivalandDurationEvent describes an event with an arrival time, depth and a duration in days
type QualitativeEvent struct {
	qualitative string
}

func (d QualitativeEvent) MarshalJSON() ([]byte, error) {
	s := fmt.Sprintf("{\"qualitativeevent\":{\"qualitative\":%s}}", d.Qualitative())
	return []byte(s), nil
}
func (h QualitativeEvent) Depth() float64 {
	return -901.0
}
func (h QualitativeEvent) Velocity() float64 {
	return -901.0
}
func (h QualitativeEvent) ArrivalTime() time.Time {
	return time.Time{}
}
func (h QualitativeEvent) Erosion() float64 {
	return -901.0
}
func (h QualitativeEvent) Duration() float64 {
	return -901.0
}
func (h QualitativeEvent) WaveHeight() float64 {
	return -901.0
}
func (h QualitativeEvent) Salinity() bool {
	return false
}
func (h QualitativeEvent) Qualitative() string {
	return h.qualitative
}
func (h *QualitativeEvent) SetQualitative(message string) {
	h.qualitative = message
}

// Parameters implements the HazardEvent interface
func (q QualitativeEvent) Parameters() Parameter {
	qp := Default
	qp = SetHasQualitative(qp)
	return qp
}

// Has implements the HazardEvent Interface
func (q QualitativeEvent) Has(p Parameter) bool {
	qp := q.Parameters()
	return qp&p != 0
}

// DepthandDVEvent describes an event with an arrival time and a duration in days
type DepthandDVEvent struct {
	depth float64
	dv    float64
}

func (d DepthandDVEvent) MarshalJSON() ([]byte, error) {
	s := fmt.Sprintf("{\"depthanddvevent\":{\"depth\":%f,\"dv\":%f}}", d.Depth(), d.DV())
	return []byte(s), nil
}
func (h DepthandDVEvent) Depth() float64 {
	return h.depth
}
func (h *DepthandDVEvent) SetDepth(value float64) {
	h.depth = value
}
func (h DepthandDVEvent) Velocity() float64 {
	return -901.0
}

func (h DepthandDVEvent) ArrivalTime() time.Time {
	return time.Time{}
}
func (h DepthandDVEvent) Erosion() float64 {
	return -901.0
}
func (h DepthandDVEvent) Duration() float64 {
	return -901.0
}
func (h DepthandDVEvent) WaveHeight() float64 {
	return -901.0
}
func (h DepthandDVEvent) Salinity() bool {
	return false
}
func (h DepthandDVEvent) Qualitative() string {
	return ""
}
func (h DepthandDVEvent) DV() float64 {
	return h.dv
}
func (h *DepthandDVEvent) SetDV(value float64) {
	h.dv = value
}

// Parameters implements the HazardEvent interface
func (ad DepthandDVEvent) Parameters() Parameter {
	adp := Default
	adp = SetHasDepth(adp)
	adp = SetHasDV(adp)
	return adp
}

// Has implements the HazardEvent Interface
func (ad DepthandDVEvent) Has(p Parameter) bool {
	adp := ad.Parameters()
	return adp&p != 0
}
