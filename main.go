package main

import (
	"log"

	"github.com/fabianofski/equaly-backend/routes"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

    app := echo.New()

    routes.SetupRoutes(app)

	app.Logger.Fatal(app.Start(":3000"))
}

