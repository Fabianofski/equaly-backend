package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"

	_ "github.com/lib/pq"
)

var (
    db *sql.DB
)

func GetPostgresConnection() (*sql.DB, error) {
    if db != nil {
        return db, nil
    }

	host := os.Getenv("POSTGRES_HOST")
	portStr := os.Getenv("POSTGRES_PORT")
	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	dbname := os.Getenv("POSTGRES_DBNAME")

	port, err := strconv.Atoi(portStr)
	if err != nil {
		log.Fatalf("Error converting POSTGRES_PORT to integer: %v", err)
        return nil, err
	}

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err = sql.Open("postgres", psqlInfo)
	if err != nil {
        return nil, err
	}

	err = db.Ping()
	if err != nil {
        return nil, err
	}

	fmt.Println("Successfully connected!")
	return db, nil
}
