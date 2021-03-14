package structureprovider

import (
	"database/sql"

	"github.com/USACE/go-consequences/consequences"
	"github.com/USACE/go-consequences/structures"
)

type gpkDataSet struct {
	db *sql.DB
}

func Init(filepath string) gpkDataSet {
	db, _ := sql.Open("sqlite3", filepath)
	db.SetMaxOpenConns(1)
	return gpkDataSet{db: db}
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

//NsiStreamProcessor is a function used to process an in memory NsiFeature through the NsiStreaming service endpoints
type StreamProcessor func(str consequences.Receptor)

/*
memory effecient structure compute methods
*/

//StreamByFips a streaming service for structure stochastic based on a bounding box
func (gpk gpkDataSet) StreamByFips(fipscode string, sp StreamProcessor) error {
	return gpk.processStream(fipscode, sp)
}
func (gpk gpkDataSet) processStream(fipscode string, sp StreamProcessor) error {
	//the query below is FIPS based - it defines the schema of the geopackage as well.
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
	}
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
