package database

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

var DB *sql.DB

// функция подключения к базе данных
func Connect(connStr string) error {
	// connStr := "user=username dbname=catfoodstore sslmode=disable password=yourpassword"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return fmt.Errorf("error DB connection: %v", err)
	}

	if err := db.Ping(); err != nil {
		return fmt.Errorf("error DB Ping: %v", err)
	}

	DB = db
	return nil
}
