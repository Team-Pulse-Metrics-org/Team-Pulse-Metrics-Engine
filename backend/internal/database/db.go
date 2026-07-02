package database

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func ConnectDB() {
	//constr connection string
	conStr := "host=localhost port=5432 user=postgres password=root dbname=team_pulse sslmode=disable"

	var err error

	DB, err = sql.Open("postgres", conStr)
	if err != nil {
		log.Fatal(err)
	}

	if err = DB.Ping(); err != nil {
		log.Fatal(err)
	}

	log.Println("Connected to PostgreSQL")
}
