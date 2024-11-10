package routes 

import (
	"github.com/labstack/echo/v4"
)

func SetupRoutes(app *echo.Group) {
	app.GET("/expense-lists", HandlerGetExpenseLists)
	app.POST("/expense-list", HandlerCreateExpenseList)
	app.POST("/expense", HandlerCreateExpense)
}

