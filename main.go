package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
    "log"
    "strconv"

	"github.com/labstack/echo/v4"
    "github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func connect_to_postgres() *sql.DB {
    host     := os.Getenv("POSTGRES_HOST")
    portStr  := os.Getenv("POSTGRES_PORT")
    user     := os.Getenv("POSTGRES_USER")
    password := os.Getenv("POSTGRES_PASSWORD")
    dbname   := os.Getenv("POSTGRES_DBNAME")

    port, err := strconv.Atoi(portStr)
	if err != nil {
		log.Fatalf("Error converting POSTGRES_PORT to integer: %v", err)
	}

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully connected!")
	return db
}

func main() {
    err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

    db := connect_to_postgres()
    err = db.Ping()
	if err != nil {
		panic(err)
	}

	e := echo.New()

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	e.Logger.Fatal(e.Start(":3000"))
}
