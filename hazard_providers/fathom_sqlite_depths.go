package hazard_providers

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

func (ds DataSet) WriteToSqlite() {
	db := CreateDepthDatabase()
	index := 0
	maxTransaction := 500
	transaction := make([]interface{}, 500)
	for _, val := range ds.Data {
		transaction[index] = val
		index++
		if index >= maxTransaction {
			WriteArrayToDepthDatabase(db, transaction)
			index = 0
		}
	}
	fmt.Println("Done.")
}
func CreateDepthDatabase() *sql.DB {
	os.Remove("fathom-depths.db")
	fmt.Println("Creating fathom-depths.db...")
	file, err := os.Create("fathom-depths.db")
	if err != nil {
		fmt.Println("error")
	}
	file.Close()
	fmt.Println("fathom-depths.db created")

	db, _ := sql.Open("sqlite3", "./fathom-depths.db")
	//defer db.Close()
	createDepthTable(db)
	return db
}
func CreateDepthWALDatabase() *sql.DB {
	os.Remove("fathom-depths.db")
	fmt.Println("Creating fathom-depths.db...")
	file, err := os.Create("fathom-depths.db")
	if err != nil {
		fmt.Println("error")
	}
	file.Close()
	fmt.Println("fathom-depths.db created")

	db, _ := sql.Open("sqlite3", "./fathom-depths.db?_journal_mode=WAL") //https://stackoverflow.com/questions/35804884/sqlite-concurrent-writing-performance/35805826
	db.SetMaxOpenConns(1)
	//defer db.Close()
	createDepthTable(db)
	return db
}
func createDepthTable(db *sql.DB) {
	createfathom := `CREATE TABLE fathom_depths (
		"fd_id" string NOT NULL PRIMARY KEY,	
		"fluv_2020_5yr" float,
		"pluv_2020_5yr" float,
		"fluv_2020_20yr" float,
		"pluv_2020_20yr" float,
		"fluv_2020_100yr" float,
		"pluv_2020_100yr" float,
		"fluv_2020_250yr" float,
		"pluv_2020_250yr" float,
		"fluv_2020_500yr" float,
		"pluv_2020_500yr" float,
		"fluv_2050_5yr" float,
		"pluv_2050_5yr" float,
		"fluv_2050_20yr" float,
		"pluv_2050_20yr" float,
		"fluv_2050_100yr" float,
		"pluv_2050_100yr" float,
		"fluv_2050_250yr" float,
		"pluv_2050_250yr" float,
		"fluv_2050_500yr" float,
		"pluv_2050_500yr" float
	  );`

	statement, err := db.Prepare(createfathom) // Prepare SQL Statement
	if err != nil {
		log.Fatal(err.Error())
	}
	statement.Exec() // Execute SQL Statements
	log.Println("fathom_depths table created")
}
func WriteArrayToDepthDatabase(db *sql.DB, dset []interface{}) {
	insertresult := `INSERT INTO fathom_depths(fd_id, fluv_2020_5yr, pluv_2020_5yr, fluv_2020_20yr, pluv_2020_20yr, fluv_2020_100yr, pluv_2020_100yr, fluv_2020_250yr, pluv_2020_250yr, fluv_2020_500yr, pluv_2020_500yr, fluv_2050_5yr, pluv_2050_5yr, fluv_2050_20yr, pluv_2050_20yr, fluv_2050_100yr, pluv_2050_100yr, fluv_2050_250yr, pluv_2050_250yr, fluv_2050_500yr, pluv_2050_500yr) VALUES `
	var inserts []string
	var params []interface{}
	somethingtoadd := false
	for _, data := range dset {
		record, ok := data.(Record)
		if ok {
			somethingtoadd = true
			inserts = append(inserts, "(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
			params = append(params, record.Fd_id, record.CurrentFluvial.Values[0], record.CurrentPluvial.Values[0], record.CurrentFluvial.Values[1], record.CurrentPluvial.Values[1], record.CurrentFluvial.Values[2], record.CurrentPluvial.Values[2], record.CurrentFluvial.Values[3], record.CurrentPluvial.Values[3], record.CurrentFluvial.Values[4], record.CurrentPluvial.Values[4], record.FutureFluvial.Values[0], record.FuturePluvial.Values[0], record.FutureFluvial.Values[1], record.FuturePluvial.Values[1], record.FutureFluvial.Values[2], record.FuturePluvial.Values[2], record.FutureFluvial.Values[3], record.FuturePluvial.Values[3], record.FutureFluvial.Values[4], record.FuturePluvial.Values[4])
		}

	}
	if somethingtoadd {
		queryVals := strings.Join(inserts, ",")
		insertresult += queryVals
		statement, err := db.Prepare(insertresult)
		if err != nil {
			fmt.Println(insertresult)
			log.Fatalln("ERROR WITH DB PREPARE " + err.Error())
		}
		_, err = statement.Exec(params...)
		if err != nil {
			fmt.Println(params)
			log.Fatalln("ERROR WITH EXECUTION " + err.Error())
		}
	}

}
