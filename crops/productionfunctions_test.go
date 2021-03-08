package crops

import (
	"testing"
	"time"
)

func TestCreateProductionFunction(t *testing.T) {
	st := time.Date(1984, time.Month(1), 22, 0, 0, 0, 0, time.UTC)
	et := time.Date(1984, time.Month(1), 28, 0, 0, 0, 0, time.UTC)
	cs := CropSchedule{StartPlantingDate: st, LastPlantingDate: et, DaysToMaturity: 330}

	mcfp := []float64{1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0}
	mclp := []float64{1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0}
	fc := []float64{1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0}
	expected := []float64{2.0, 4.0, 6.0, 8.0, 10.0, 12.0, 14.0, 16.0, 18.0, 20.0, 22.0, 24.0}
	pf := NewProductionFunction(mcfp, mclp, fc, cs, 1.0, .1)

	for idx := range mcfp {
		got := pf.GetCumulativeMonthlyProductionCostsEarly()[idx]
		if got != expected[idx] {
			t.Errorf("productionFunction.GetCumulativeMonthlyProductionCosts()[%v] = %f; expected %f", idx, got, expected[idx])
		}
	}
}
func TestCreateProductionFunctionWrapYear(t *testing.T) {
	st := time.Date(1984, time.Month(2), 22, 0, 0, 0, 0, time.UTC)
	et := time.Date(1984, time.Month(2), 28, 0, 0, 0, 0, time.UTC)
	cs := CropSchedule{StartPlantingDate: st, LastPlantingDate: et, DaysToMaturity: 330}

	mcfp := []float64{1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0}
	mclp := []float64{1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0}
	fc := []float64{1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0}
	expected := []float64{24.0, 2.0, 4.0, 6.0, 8.0, 10.0, 12.0, 14.0, 16.0, 18.0, 20.0, 22.0}
	pf := NewProductionFunction(mcfp, mclp, fc, cs, 1.0, .1)

	for idx := range mcfp {
		got := pf.GetCumulativeMonthlyProductionCostsEarly()[idx]
		if got != expected[idx] {
			t.Errorf("productionFunction.GetCumulativeMonthlyProductionCosts()[%v] = %f; expected %f", idx, got, expected[idx])
		}
	}
}
func TestCreateProductionFunction_ShorterCropSchedule(t *testing.T) {
	st := time.Date(1984, time.Month(2), 22, 0, 0, 0, 0, time.UTC)
	et := time.Date(1984, time.Month(2), 28, 0, 0, 0, 0, time.UTC)
	cs := CropSchedule{StartPlantingDate: st, LastPlantingDate: et, DaysToMaturity: 299}

	mcfp := []float64{1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0}
	mclp := []float64{1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0}
	fc := []float64{1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0}
	expected := []float64{0.0, 2.0, 4.0, 6.0, 8.0, 10.0, 12.0, 14.0, 16.0, 18.0, 20.0, 22.0}
	pf := NewProductionFunction(mcfp, mclp, fc, cs, 1.0, .1)

	for idx := range mcfp {
		got := pf.GetCumulativeMonthlyProductionCostsEarly()[idx]
		if got != expected[idx] {
			t.Errorf("productionFunction.GetCumulativeMonthlyProductionCosts()[%v] = %f; expected %f", idx, got, expected[idx])
		}
	}
}
