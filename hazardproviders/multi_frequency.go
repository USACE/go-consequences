package hazardproviders

import (
	"errors"
	"time"

	"github.com/USACE/go-consequences/geography"
	"github.com/USACE/go-consequences/hazards"
	"github.com/dewberry/gdal"
)

type MultiFrequencyDepthProvider struct {
	depths      []float64
	frequencies []float64
	process     HazardFunction
	bbox        geography.BBox
}

func mfDepthSchema() []string {
	s := []string{"event_id", "aep", "depth"}
	return s
}

func InitMultiFrequencyDepthProvider(filepath string, driver string, layername string) (MultiFrequencyDepthProvider, error) {

	var mfp = MultiFrequencyDepthProvider{}

	d := gdal.OGRDriverByName(driver)
	ds, dsok := d.Open(filepath, int(gdal.ReadOnly))
	if !dsok {
		return mfp, errors.New("error opening file" + filepath + " of type " + driver)
	}

	var layer gdal.Layer
	if layername != "" {
		layer = ds.LayerByName(layername)
	} else {
		layer = ds.LayerByIndex(0)
	}
	def := layer.Definition()
	s := mfDepthSchema()
	sIDX := make([]int, len(s))

	for i, f := range s {
		idx := def.FieldIndex(f)

		if idx < 0 {
			return mfp, errors.New("gdal dataset at path " + filepath + " Expected field named " + f + " none was found")
		}
		sIDX[i] = idx
	}

	fc, _ := layer.FeatureCount(true)
	depths := make([]float64, fc)
	freqs := make([]float64, fc)

	for i := range fc {
		feat := layer.NextFeature()
		if feat != nil {
			freqs[i] = feat.FieldAsFloat64(sIDX[1])
			depths[i] = feat.FieldAsFloat64(sIDX[2])
		}
	}
	// This provider is only for testing purposes now.
	// TODO: add user ability to specify BBox in Init func
	bb := geography.BBox{Bbox: []float64{-71.5, 42.9, -71.4, 43}}
	mfp.depths = depths
	mfp.frequencies = freqs
	mfp.bbox = bb
	mfp.process = DepthFrequencyHazardFunction()
	return mfp, nil
}

func (p MultiFrequencyDepthProvider) Close() {
	// Is there anything to close when doing a read via gdal?
}

func (p MultiFrequencyDepthProvider) Hazard(l geography.Location) (hazards.HazardEvent, error) {
	var hm hazards.DepthEventMultiFrequency

	if p.bbox.Contains(l) {
		for i, d := range p.depths {
			hd := hazards.HazardData{
				Depth:       d,
				Velocity:    0,
				ArrivalTime: time.Time{},
				Erosion:     0,
				Duration:    0,
				WaveHeight:  0,
				Salinity:    false,
				Qualitative: "",
				Frequency:   p.frequencies[i],
			}
			var h hazards.HazardEvent
			h, err := p.process(hd, h)
			if err != nil {
				panic(err)
			}
			hm.Append(h)
		}
	}
	return &hm, nil
}

func (p MultiFrequencyDepthProvider) HazardBoundary() (geography.BBox, error) {
	return p.bbox, nil
}
