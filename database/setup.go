package setup

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

var DB *sql.DB

// Establishing the connection
func InitDB() {
	connStr := "postgres://ankush:2004@localhost:5432/mytop100?sslmode=disable"
	var err error
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Error Connecting the datatbase", err)
	}

	//Verifying the connection
	err = DB.Ping()
	if err != nil {
		log.Fatal("Cannot connect the database", err)
	}
	fmt.Println("✅ Connected to PostgreSql Succesfully")
}

func GetDB() *sql.DB {
	return DB
}
