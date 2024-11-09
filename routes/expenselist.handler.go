package routes

import (
	"log"
	"net/http"

	"github.com/fabianofski/equaly-backend/db"
	_ "github.com/fabianofski/equaly-backend/models"
	"github.com/labstack/echo/v4"
)

// HandlerGetExpenseLists godoc
//
//	@Summary		Get Expense Lists
//	@Description	Retrieves the list of expenses for a specified user.
//	@Tags			Expenses
//	@Param			userId	query	string			true	"User ID to retrieve expenses for"
//	@Success		200		{array}	models.Expense	"List of user expenses"
//	@Failure		400		"Bad Request"
//	@Failure		500		"Internal Server Error"
//	@Router			/user-expense-lists [get]
func HandlerGetExpenseLists(c echo.Context) error {
    log.Println("GET FOR USER")

	userId := c.QueryParam("userId")
	if userId == "" {
		return c.JSON(http.StatusBadRequest, "Bad Request")
	}

	expenseLists, err := db.GetExpenseLists(userId)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "Error requesting data")
	}

	return c.JSON(http.StatusOK, expenseLists)

}
