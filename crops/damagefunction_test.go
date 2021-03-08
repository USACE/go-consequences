package crops

import (
	"testing"
	"time"

	"github.com/USACE/go-consequences/hazards"
)

func TestDamageFunctionCompute_one(t *testing.T) {

	one := []float64{1.1, 2.1, 3.1, 4.1, 5.1, 6.1, 7.1, 8.1, 9.1, 10.1, 11.1, 12.1}
	two := []float64{1.2, 2.2, 3.2, 4.2, 5.2, 6.2, 7.2, 8.2, 9.2, 10.2, 11.2, 12.2}
	three := []float64{1.3, 2.3, 3.3, 4.3, 5.3, 6.3, 7.3, 8.3, 9.3, 10.3, 11.3, 12.3}
	four := []float64{1.4, 2.4, 3.4, 4.4, 5.4, 6.4, 7.4, 8.4, 9.4, 10.4, 11.4, 12.4}

	m := make(map[float64][]float64)
	m[1.0] = one
	m[2.0] = two
	m[3.0] = three
	m[4.0] = four

	df := DamageFunction{DurationDamageCurves: m}

	at := time.Date(1984, time.Month(1), 22, 0, 0, 0, 0, time.UTC)
	h := hazards.ArrivalandDurationEvent{}
	h.SetArrivalTime(at)
	h.SetDuration(1.5)
	got := df.ComputeDamagePercent(h)
	expected := 1.15
	if got != expected {
		t.Errorf("ComputeDamagePercent() = %f; expected %f", got, expected)
	}
}
func TestDamageFunctionCompute_two(t *testing.T) {

	one := []float64{1.1, 2.1, 3.1, 4.1, 5.1, 6.1, 7.1, 8.1, 9.1, 10.1, 11.1, 12.1}
	two := []float64{1.2, 2.2, 3.2, 4.2, 5.2, 6.2, 7.2, 8.2, 9.2, 10.2, 11.2, 12.2}
	three := []float64{1.3, 2.3, 3.3, 4.3, 5.3, 6.3, 7.3, 8.3, 9.3, 10.3, 11.3, 12.3}
	four := []float64{1.4, 2.4, 3.4, 4.4, 5.4, 6.4, 7.4, 8.4, 9.4, 10.4, 11.4, 12.4}

	m := make(map[float64][]float64)
	m[1.0] = one
	m[2.0] = two
	m[3.0] = three
	m[4.0] = four

	df := DamageFunction{DurationDamageCurves: m}

	at := time.Date(1984, time.Month(1), 22, 0, 0, 0, 0, time.UTC)
	h := hazards.ArrivalandDurationEvent{}
	h.SetArrivalTime(at)
	h.SetDuration(.5)
	got := df.ComputeDamagePercent(h)
	expected := .55
	if got != expected {
		t.Errorf("ComputeDamagePercent() = %f; expected %f", got, expected)
	}
}
func TestDamageFunctionCompute_three(t *testing.T) {

	one := []float64{1.1, 2.1, 3.1, 4.1, 5.1, 6.1, 7.1, 8.1, 9.1, 10.1, 11.1, 12.1}
	two := []float64{1.2, 2.2, 3.2, 4.2, 5.2, 6.2, 7.2, 8.2, 9.2, 10.2, 11.2, 12.2}
	three := []float64{1.3, 2.3, 3.3, 4.3, 5.3, 6.3, 7.3, 8.3, 9.3, 10.3, 11.3, 12.3}
	four := []float64{1.4, 2.4, 3.4, 4.4, 5.4, 6.4, 7.4, 8.4, 9.4, 10.4, 11.4, 12.4}

	m := make(map[float64][]float64)
	m[1.0] = one
	m[2.0] = two
	m[3.0] = three
	m[4.0] = four

	df := DamageFunction{DurationDamageCurves: m}

	at := time.Date(1984, time.Month(2), 22, 0, 0, 0, 0, time.UTC)
	h := hazards.ArrivalandDurationEvent{}
	h.SetArrivalTime(at)
	h.SetDuration(2.5)
	got := df.ComputeDamagePercent(h)
	expected := 2.25
	if got != expected {
		t.Errorf("ComputeDamagePercent() = %f; expected %f", got, expected)
	}
}
