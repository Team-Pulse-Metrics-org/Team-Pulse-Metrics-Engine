package database

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib" // 👈 The new driver registration!
	"github.com/joho/godotenv"
)

var DB *sql.DB

func ConnectDB() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	conStr := os.Getenv("DATABASE_URL")

	// 👈 Change "postgres" to "pgx" right here
	DB, err = sql.Open("pgx", conStr)
	if err != nil {
		log.Fatal(err)
	}

	if err = DB.Ping(); err != nil {
		log.Fatal(err)
	}

	log.Println("Connected to PostgreSQL using pgx driver!")
}
