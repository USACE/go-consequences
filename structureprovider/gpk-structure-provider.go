package structureprovider

import (
	"fmt"

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
	fmt.Println(ds.Driver().Name())
	for i := 0; i < ds.LayerCount(); i++ {
		fmt.Println(ds.LayerByIndex(i).Name())
		layer := ds.LayerByIndex(i)
		fieldDef := layer.Definition()

		for j := 0; j < fieldDef.FieldCount(); j++ {
			fieldName := fieldDef.FieldDefinition(j).Name()
			fieldType := fieldDef.FieldDefinition(j).Type().Name()
			fmt.Println(fmt.Sprintf("%s, %s", fieldName, fieldType))
		}
	}
	return gpkDataSet{FilePath: filepath, ds: &ds}
}

type SimpleStructure struct {
	Name        string
	X           float64
	Y           float64
	DamCat      string
	OcctypeName string
	Found_ht    float64
	Found_type  string
	Val_struct  float64
	Val_cont    float64
	Pop2amo65   int16
	Pop2amu65   int16
	Pop2pmo65   int16
	Pop2pmu65   int16
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
	fc, _ := gpk.ds.LayerByName("nsi").FeatureCount(true)
	for idx < fc { // Iterate and fetch the records from result cursor
		s := SimpleStructure{}
		f := gpk.ds.LayerByName("nsi").NextFeature()
		s.Name = fmt.Sprintf("%v", f.FieldAsInteger(0))
		s.OcctypeName = f.FieldAsString(4)
		//check if CBID matches?
		sp(toStructure(s, m, defaultOcctype))
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
func toStructure(s SimpleStructure, m map[string]structures.OccupancyTypeStochastic, defaultOcctype structures.OccupancyTypeStochastic) structures.StructureStochastic {
	var occtype = defaultOcctype
	if ot, ok := m[s.OcctypeName]; ok {
		occtype = ot
	} else {
		occtype = defaultOcctype
		msg := "Using default " + s.OcctypeName + " not found"
		panic(msg)
	}
	return structures.StructureStochastic{
		OccType:   occtype,
		StructVal: consequences.ParameterValue{Value: s.Val_struct},
		ContVal:   consequences.ParameterValue{Value: s.Val_cont},
		FoundHt:   consequences.ParameterValue{Value: s.Found_ht},
		BaseStructure: structures.BaseStructure{
			Name:   s.Name,
			DamCat: s.DamCat,
			X:      s.X,
			Y:      s.Y,
		},
	}
}
