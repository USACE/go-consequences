package store

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

type fathom_result struct {
	fd_id                 string  `db:"fd_id"`
	hazard_Year           int     `db:"hazard_year"` //2020, 2050
	hazard_Type           string  `db:"hazard_type"` //fluvial, pluvial
	frequency             string  `db:"frequency"`   // 5, 20, 100, 250, 500
	structure_Consequence float64 `db:"structure_consequence"`
	content_Consequence   float64 `db:"content_consequence"`
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
func createTable(db *sql.DB) {
	createfathom := `CREATE TABLE fathom (
		"result_id" integer NOT NULL PRIMARY KEY AUTOINCREMENT,
		"fd_id" string,		
		"hazard_year" integer,
		"hazard_type" TEXT,
		"frequency" TEXT,
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

func WriteToDatabase(db *sql.DB, fd_id string, year int, hazard string, frequency string, structure_damage float64, content_damage float64) {
	//f := fathom_result{fd_id: fd_id, hazard_Year: year, hazard_Type: hazard, frequency: frequency, structure_Consequence: structure_damage, content_Consequence: content_damage}
	//db, _ := sql.Open("sqlite3", "./fathom-results.db")
	//defer db.Close()

	insertresult := `INSERT INTO fathom(fd_id, hazard_year, hazard_type, frequency, structure_consequence, content_consequence) VALUES (?, ?, ?, ?, ?, ?)`
	statement, err := db.Prepare(insertresult)
	if err != nil {
		log.Fatalln(err.Error())
	}
	_, err = statement.Exec(fd_id, year, hazard, frequency, structure_damage, content_damage)
	if err != nil {
		log.Fatalln(err.Error())
	}
	//write to database.

}
