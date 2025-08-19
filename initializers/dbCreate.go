package initializers

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

func CreateDatabaseIfNotExists() {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	// Connect to default database (postgres)
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=postgres sslmode=disable",
		host, port, user, password)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal("Unable to connect to Postgres:", err)
	}
	defer db.Close()

	// Check if database exists
	var exists bool
	err = db.QueryRow("SELECT EXISTS(SELECT 1 FROM pg_database WHERE datname = $1)", dbname).Scan(&exists)
	if err != nil {
		log.Fatal("Failed to check if database exists:", err)
	}

	if !exists {
		_, err = db.Exec("CREATE DATABASE " + dbname)
		if err != nil {
			log.Fatal("Failed to create database:", err)
		}
		log.Println("✅ Database created:", dbname)
	} else {
		log.Println("✅ Database already exists:", dbname)
	}
}
