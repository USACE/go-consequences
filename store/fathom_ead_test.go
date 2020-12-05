package store

import (
	"database/sql"
	"fmt"
	"log"
	"testing"
)

func TestCreateDatabase(t *testing.T) {
	CreateDatabase()
}
func TestPrintDatabase(t *testing.T) {
	fmt.Println("Reading Database")
	db, _ := sql.Open("sqlite3", "./fathom-results.db")
	defer db.Close()
	row, err := db.Query("SELECT * FROM fathom ORDER BY result_id")
	if err != nil {
		log.Fatal(err)
	}
	defer row.Close()
	for row.Next() { // Iterate and fetch the records from result cursor
		var rid int
		var fid string
		var x float64
		var y float64
		var fips string
		var year int
		var hazard string
		var frequency string
		var str float64
		var cont float64
		row.Scan(&rid, &fid, &x, &y, &fips, &year, &hazard, &frequency, &str, &cont)
		fmt.Println(fmt.Sprintf("result: %v, %v, %f, %f, %v, %v, %v, %v, %f, %f", rid, fid, x, y, fips, year, hazard, frequency, str, cont))
	}
}

func TestPrintDatabaseSize(t *testing.T) {
	fmt.Println("Reading Database")
	db, _ := sql.Open("sqlite3", "./fathom-results.db")
	defer db.Close()
	row, err := db.Query("SELECT COUNT(*) as count FROM fathom")
	if err != nil {
		log.Fatal(err)
	}
	defer row.Close()
	for row.Next() {
		var count float64
		row.Scan(&count)
		fmt.Println(fmt.Sprintf("result: %v", count))
	}

}
