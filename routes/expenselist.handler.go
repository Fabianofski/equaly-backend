package routes

import (
	"net/http"

	"github.com/fabianofski/equaly-backend/db"
	"github.com/labstack/echo/v4"
)

func HandlerGetExpenseLists(c echo.Context) error {
		userId := c.QueryParam("userId")
		if userId == "" {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Bad Request"})
		}


		expenseLists, err := db.GetExpenseLists(userId)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"error": "Error requesting data",
			})
		}

		return c.JSON(http.StatusOK, expenseLists)

	}
