package main

import (
	"log"
	"net/http"

	"github.com/fabianofski/equaly-backend/routes"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/swaggo/echo-swagger"
    _ "github.com/fabianofski/equaly-backend/docs"
)

//	@title			Equaly Backend 
//	@version		1.0
//	@description	This is the backend for the equaly cost management app

//	@host		localhost:3000 	
//	@BasePath	/v1
func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

    app := echo.New()

	app.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
    app.GET("/docs", func(c echo.Context) error {
        return c.Redirect(http.StatusMovedPermanently, "/docs/index.html")
    })
    app.GET("/docs/*", echoSwagger.WrapHandler)

    version := app.Group("/v1")
    routes.SetupRoutes(version)

	app.Logger.Fatal(app.Start(":3000"))
}

