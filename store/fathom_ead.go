package store

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

type fathom_result struct {
	fd_id                 string  `db:"fd_id"`
	x                     float64 `db:"x"`
	y                     float64 `db:"y"`
	fips                  string  `db:"fips"`
	hazard_Year           int     `db:"hazard_year"` //2020, 2050
	hazard_Type           string  `db:"hazard_type"` //fluvial, pluvial
	frequency             string  `db:"frequency"`   // 5, 20, 100, 250, 500
	structure_Consequence float64 `db:"structure_consequence"`
	content_Consequence   float64 `db:"content_consequence"`
}

func CreateResult(fd_id string, x float64, y float64, sfips string, year int, hazard string, frequency string, structure_damage float64, content_damage float64) fathom_result {
	result := fathom_result{fd_id: fd_id, x: x, y: y, fips: sfips, hazard_Year: year, hazard_Type: hazard, frequency: frequency, structure_Consequence: structure_damage, content_Consequence: content_damage}
	return result
}
func CreateDatabase() *sql.DB {
	os.Remove("fathom-results.db")
	fmt.Println("Creating fathom-results.db...")
	file, err := os.Create("fathom-results.db")
	if err != nil {
		fmt.Println("error")
	}
	file.Close()
	fmt.Println("fathom-results.db created")

	db, _ := sql.Open("sqlite3", "./fathom-results.db")
	//defer db.Close()
	createTable(db)
	return db
}
func CreateWALDatabase() *sql.DB {
	os.Remove("fathom-results.db")
	fmt.Println("Creating fathom-results.db...")
	file, err := os.Create("fathom-results.db")
	if err != nil {
		fmt.Println("error")
	}
	file.Close()
	fmt.Println("fathom-results.db created")

	db, _ := sql.Open("sqlite3", "./fathom-results.db?_journal_mode=WAL") //https://stackoverflow.com/questions/35804884/sqlite-concurrent-writing-performance/35805826
	db.SetMaxOpenConns(1)
	//defer db.Close()
	createTable(db)
	return db
}
func createTable(db *sql.DB) {
	createfathom := `CREATE TABLE fathom (
		"result_id" integer NOT NULL PRIMARY KEY AUTOINCREMENT,
		"fd_id" string,
		"x" float,
		"y" float,
		"fips" string,	
		"hazard_year" integer,
		"hazard_type" string,
		"frequency" string,
		"structure_consequence" float,
		"content_consequence" float
	  );`

	statement, err := db.Prepare(createfathom) // Prepare SQL Statement
	if err != nil {
		log.Fatal(err.Error())
	}
	statement.Exec() // Execute SQL Statements
	log.Println("fathom table created")
}
func CreateStatement(db *sql.DB) *sql.Stmt {
	//https://golangbot.com/mysql-create-table-insert-row/
	insertresult := `INSERT INTO fathom(fd_id, x, y, fips, hazard_year, hazard_type, frequency, structure_consequence, content_consequence) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`
	statement, err := db.Prepare(insertresult)
	if err != nil {
		log.Fatalln(err.Error())
	} else {
		return statement
	}
	return nil
}
func WriteArrayToDatabase(db *sql.DB, results []interface{}) {
	insertresult := `INSERT INTO fathom(fd_id, x, y, fips, hazard_year, hazard_type, frequency, structure_consequence, content_consequence) VALUES `
	var inserts []string
	var params []interface{}
	somethingtoadd := false
	for _, result := range results {
		res, ok := result.(fathom_result)
		if ok {
			somethingtoadd = true
			inserts = append(inserts, "(?, ?, ?, ?, ?, ?, ?, ?, ?)")
			params = append(params, res.fd_id, res.x, res.y, res.fips, res.hazard_Year, res.hazard_Type, res.frequency, res.structure_Consequence, res.content_Consequence)
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
func WriteToDatabase(stmt *sql.Stmt, fd_id string, year int, hazard string, frequency string, structure_damage float64, content_damage float64) {
	_, err := stmt.Exec(fd_id, year, hazard, frequency, structure_damage, content_damage)
	if err != nil {
		log.Fatalln(err.Error())
	}
}
