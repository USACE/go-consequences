package resultswriters

import (
	"path/filepath"
	"testing"
	"time"

	"github.com/USACE/go-consequences/consequences"
	"github.com/dewberry/gdal"
)

func TestSpatialResultsWriterWritesDateTime(t *testing.T) {
	fp := filepath.Join(t.TempDir(), "results.gpkg")
	w, err := InitSpatialResultsWriter(fp, "results", "GPKG")
	if err != nil {
		t.Fatal(err)
	}
	completion := time.Date(2040, time.November, 5, 4, 17, 23, 0, time.UTC)
	headers := []string{"x", "y", "completion_date"}
	w.Write(consequences.Result{Headers: headers, Result: []interface{}{-75.1, 38.7, completion}})
	w.Write(consequences.Result{Headers: headers, Result: []interface{}{-75.1, 38.7, time.Time{}}})
	w.Close()

	ds := gdal.OpenDataSource(fp, 0)
	defer ds.Destroy()
	layer := ds.LayerByName("results")
	idx := layer.Definition().FieldIndex("completion")
	if got := layer.Definition().FieldDefinition(idx).Type(); got != gdal.FT_DateTime {
		t.Errorf("completion field type = %v, want DateTime", got)
	}

	f := layer.NextFeature()
	if got, ok := f.FieldAsDateTime(idx); !ok || !got.Equal(completion) {
		t.Errorf("completion = %v, want %v", got, completion)
	}
	f.Destroy()

	f = layer.NextFeature()
	if f.IsFieldSetAndNotNull(idx) {
		t.Errorf("zero completion time was written, want null")
	}
	f.Destroy()
}
