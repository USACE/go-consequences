package hazardproviders

// func TestInitCSV(t *testing.T) {
// 	file := "/workspaces/go-consequences/data/lifecycle/test_arrival-depth-duration_hazards.csv"

// 	expectedDepths := []float64{1.0, 1.0, 1.0, 2.0, 2.0}
// 	expectedDurations := []float64{0.0, 5.0, 0.0, 0.0, 0.0}
// 	et1 := time.Date(1984, time.Month(1), 1, 0, 0, 0, 0, time.UTC)
// 	et2 := time.Date(1984, time.Month(1), 11, 0, 0, 0, 0, time.UTC)
// 	et3 := time.Date(1984, time.Month(1), 21, 0, 0, 0, 0, time.UTC)
// 	et4 := time.Date(1985, time.Month(1), 1, 0, 0, 0, 0, time.UTC)
// 	et5 := time.Date(1985, time.Month(1), 11, 0, 0, 0, 0, time.UTC)
// 	expectedArrivals := []time.Time{et1, et2, et3, et4, et5}

// 	b := geography.BBox{
// 		// xMin, yMin, xMax, yMax
// 		Bbox: []float64{-71.1, 43, -71, 43.1},
// 	}
// 	ADDMHP, err := InitCSV(file, b)
// 	if err != nil {
// 		panic(err)
// 	}
// 	defer ADDMHP.Close()

// 	loc := geography.Location{
// 		X:    -71.05,
// 		Y:    43.05,
// 		SRID: "",
// 	}

// 	haz, err := ADDMHP.Hazard(loc)
// 	h := haz.(*hazards.ArrivalDepthandDurationEventMulti)
// 	if err != nil {
// 		panic(err)
// 	}

// 	for {
// 		edepth := expectedDepths[h.Index()]
// 		edur := expectedDurations[h.Index()]
// 		earr := expectedArrivals[h.Index()]

// 		if h.Depth() != edepth {
// 			t.Errorf("Event at index %d had Depth = %v. Expected: %3.2f", h.Index(), h.Depth(), edepth)
// 		}
// 		if h.Duration() != edur {
// 			t.Errorf("Event at index %d had Duration = %v. Expected: %3.2f", h.Index(), h.Duration(), edur)
// 		}
// 		if h.Depth() != edepth {
// 			t.Errorf("Event at index %d had ArrivalTime = %v. Expected: %v", h.Index(), h.ArrivalTime(), earr)
// 		}

// 		if h.HasNext() {
// 			h.Increment()
// 		} else {
// 			break
// 		}
// 	}
// }
