package routes

import (
	"net/http"
    "context"
    "strings"

	"github.com/labstack/echo/v4"
    "google.golang.org/api/idtoken"
)

func GoogleAuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		authHeader := c.Request().Header.Get("Authorization")
		if authHeader == "" {
			return echo.NewHTTPError(http.StatusUnauthorized, "Missing Authorization header")
		}

		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			return echo.NewHTTPError(http.StatusUnauthorized, "Invalid Authorization format")
		}

		tokenString := tokenParts[1]

        payload, err := idtoken.Validate(context.Background(), tokenString, "")
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, "Invalid ID Token")
		}

		userId, ok := payload.Claims["sub"].(string)
		if !ok {
			return echo.NewHTTPError(http.StatusUnauthorized, "User ID not found in token")
		}

		c.Set("userId", userId)
		return next(c)
	}
}

func SetupRoutes(app *echo.Group) {
    app.Use(GoogleAuthMiddleware)

	app.GET("/expense-lists", HandlerGetExpenseLists)
	app.POST("/expense-list", HandlerCreateExpenseList)
	app.POST("/expense", HandlerCreateExpense)
}

