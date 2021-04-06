package consequences

import (
	"fmt"
	"io"
	"os"

	"github.com/HydrologicEngineeringCenter/go-statistics/data"
	"github.com/leekchan/accounting"
)

type summaryResultsWriter struct {
	filepath   string
	w          io.Writer
	grandTotal float64
	totals     map[string]float64
	m          map[string]*data.InlineHistogram
}

func InitSummaryResultsWriterFromFile(filepath string) *summaryResultsWriter {
	w, err := os.OpenFile(filepath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0600)
	if err != nil {
		panic(err)
	}
	//make the maps
	t := make(map[string]float64, 1)
	m := make(map[string]*data.InlineHistogram, 1)
	return &summaryResultsWriter{filepath: filepath, w: w, totals: t, m: m}
}
func InitSummaryResultsWriter(w io.Writer) *summaryResultsWriter {
	return &summaryResultsWriter{filepath: "not applicapble", w: w}
}
func (srw *summaryResultsWriter) Write(r Result) {
	//hardcoding for structures to experiment and think it through.
	var cat = "damage category"
	var structDam = "structure damage"
	var contDam = "content damage"
	var totDam = 0.0
	var damcat = ""
	h := r.Headers
	for i, v := range h {
		if v == cat {
			//add data to the map from this index in results
			damcat = r.Result[i].(string)
		}
		if v == structDam {
			totDam += r.Result[i].(float64)
		}
		if v == contDam {
			totDam += r.Result[i].(float64)
		}
	}
	srw.grandTotal += totDam
	ih, ok := srw.m[damcat]
	if ok {
		ih.AddObservation(totDam)
		srw.m[damcat] = ih
	} else {
		nih := data.Init(1000, 0, 10000)
		srw.m[damcat] = nih
	}
	//update damcat totals.
	t, ok := srw.totals[damcat]
	if ok {
		t += totDam
		srw.totals[damcat] = t
	} else {
		srw.totals[damcat] = totDam
	}
}
func (srw *summaryResultsWriter) Close() {
	ac := accounting.Accounting{Symbol: "$", Precision: 2}
	fmt.Fprintf(srw.w, fmt.Sprintf("Grand Total is %v\n", ac.FormatMoney(srw.grandTotal)))
	h := srw.totals
	for i, v := range h {
		fmt.Fprintf(srw.w, fmt.Sprintf("Damages for %v were %v\n", i, ac.FormatMoney(v)))
	}
	j := srw.m
	for i, v := range j {
		fmt.Fprintf(srw.w, fmt.Sprintf("Histogram for %v:\n%v", i, v.StringSparse()))
	}
	w2, ok := srw.w.(io.WriteCloser)
	if ok {
		w2.Close()
	}
}
