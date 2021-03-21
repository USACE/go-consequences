package structureprovider

import (
	"fmt"
	"log"
	"strings"

	"github.com/USACE/go-consequences/consequences"
	"github.com/USACE/go-consequences/geography"
	"github.com/USACE/go-consequences/structures"
	"github.com/dewberry/gdal"
)

type gpkDataSet struct {
	FilePath string
	ds       *gdal.DataSource
}

func InitGPK(filepath string) gpkDataSet {
	ds := gdal.OpenDataSource(filepath, int(gdal.ReadOnly))
	//validation?
	hasNSITable := false
	for i := 0; i < ds.LayerCount(); i++ {
		if "nsi" == ds.LayerByIndex(i).Name() {
			hasNSITable = true
		}
	}
	if hasNSITable {
		return gpkDataSet{FilePath: filepath, ds: &ds}
	}
	log.Fatalln("GeoPpackage does not have a layer titled nsi.  Killing everything! ")
	return gpkDataSet{FilePath: filepath}
}

//StreamByFips a streaming service for structure stochastic based on a bounding box
func (gpk gpkDataSet) ByFips(fipscode string, sp StreamProcessor) error {
	return gpk.processFipsStream(fipscode, sp)
}
func (gpk gpkDataSet) processFipsStream(fipscode string, sp StreamProcessor) error {
	m := structures.OccupancyTypeMap()
	//define a default occtype in case of emergancy
	defaultOcctype := m["RES1-1SNB"]
	idx := 0
	l := gpk.ds.LayerByName("nsi")
	fc, _ := l.FeatureCount(true)
	for idx < fc { // Iterate and fetch the records from result cursor
		f := l.NextFeature()
		cbfips := f.FieldAsString(3)
		//check if CBID matches?
		if strings.Contains(cbfips, fipscode) {
			sp(featuretoStructure(f, m, defaultOcctype))
		}
	}
	return nil

}
func (gpk gpkDataSet) ByBbox(bbox geography.BBox, sp StreamProcessor) error {
	return gpk.processBboxStream(bbox, sp)
}
func (gpk gpkDataSet) processBboxStream(bbox geography.BBox, sp StreamProcessor) error {
	/*//the query below is FIPS based - it defines the schema of the geopackage as well.
	rows, err := gpk.db.Query("SELECT fd_id, x, y, occtype, found_ht, found_type, st_damcat, val_struct, val_cont, pop2amu65, pop2amo65, pop2pmu65, pop2pmo65 FROM nsi WHERE cbfips LIKE '" + fipscode + "%'")
	m := structures.OccupancyTypeMap()
	//define a default occtype in case of emergancy
	defaultOcctype := m["RES1-1SNB"]
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() { // Iterate and fetch the records from result cursor
		s := SimpleStructure{}
		err := rows.Scan(&s.Name, &s.X, &s.Y, &s.OcctypeName, &s.Found_ht, &s.Found_type, &s.DamCat, &s.Val_struct, &s.Val_cont, &s.Pop2amu65, &s.Pop2amo65, &s.Pop2pmu65, &s.Pop2pmo65)
		if err != nil {
			return err
		}
		sp(toStructure(s, m, defaultOcctype))
	}*/
	return nil

}
func featuretoStructure(f *gdal.Feature, m map[string]structures.OccupancyTypeStochastic, defaultOcctype structures.OccupancyTypeStochastic) structures.StructureStochastic {
	s := structures.StructureStochastic{}
	s.Name = fmt.Sprintf("%v", f.FieldAsInteger(0))
	OccTypeName := f.FieldAsString(4)
	var occtype = defaultOcctype
	if ot, ok := m[OccTypeName]; ok {
		occtype = ot
	} else {
		occtype = defaultOcctype
		msg := "Using default " + OccTypeName + " not found"
		panic(msg)
	}
	s.OccType = occtype
	s.X = f.FieldAsFloat64(1)
	s.Y = f.FieldAsFloat64(2)
	s.DamCat = f.FieldAsString(18)
	s.StructVal = consequences.ParameterValue{Value: f.FieldAsFloat64(23)}
	s.ContVal = consequences.ParameterValue{Value: f.FieldAsFloat64(24)}
	s.FoundHt = consequences.ParameterValue{Value: f.FieldAsFloat64(21)}
	return s
}
