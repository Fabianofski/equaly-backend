package routes

import (
	"log"
	"net/http"

	"github.com/fabianofski/equaly-backend/db"
	"github.com/fabianofski/equaly-backend/lib"
	"github.com/fabianofski/equaly-backend/models"
	"github.com/labstack/echo/v4"
)

// HandlerGetExpenseLists godoc
//
//	@Summary		Get Expense Lists
//	@Description	Retrieves the list of expenses for a specified user.
//	@Tags			Expenses
//	@Param			Authorization	header	string						true	"Bearer Token"
//	@Success		200				{array}	models.ExpenseListWrapper	"List of user expense lists"
//	@Failure		400				"Bad Request"
//	@Failure		500				"Internal Server Error"
//	@Router			/expense-lists [get]
func HandlerGetExpenseLists(c echo.Context) error {
    userId := c.Get("userId").(string)
    log.Println("GET Expense Lists for", userId)
    
	expenseLists, err := db.GetExpenseLists(userId)
	if err != nil {
		log.Println(err)
		return c.String(http.StatusInternalServerError, "Error requesting data")
	}

	var expenseListWrappers = make([]models.ExpenseListWrapper, 0)
	for _, expenseList := range expenseLists {
		expenseListWrappers = append(expenseListWrappers, lib.Calculate_shares_and_compensations(*expenseList))
	}

	return c.JSON(http.StatusOK, expenseListWrappers)

}

// HandlerCreateExpenseLists godoc
//
//	@Summary		Create Expense List
//	@Description	Creates a new Expense list with given data for a specified user
//	@Tags			Expenses
//	@Param			expenseList		body		models.ExpenseListWrapper	true	"Expense List Data"
//	@Param			Authorization	header		string						true	"Bearer Token"
//	@Success		200				{object}	models.ExpenseListWrapper	"Created Expense List"
//	@Failure		400				"Bad Request"
//	@Failure		500				"Internal Server Error"
//	@Router			/expense-list [post]
func HandlerCreateExpenseList(c echo.Context) error {
    userId := c.Get("userId").(string)
	log.Println("POST Create new Expense List")
	expenseList := new(models.ExpenseList)
	if err := c.Bind(expenseList); err != nil {
		log.Println("400 Bad Request")
		log.Println(err)
		return c.String(http.StatusBadRequest, "Bad Request")
	}

    expenseList.CreatorId = userId;

	if expenseList.Color == "" || expenseList.Emoji == "" || expenseList.Title == "" || expenseList.Currency == "" {
		log.Println("400 Bad Request")
		return c.String(http.StatusBadRequest, "Bad Request")
	}

	id, err := db.CreateExpenseList(expenseList)
	if err != nil {
		log.Println("500 Internal Server Error")
		log.Println(err)
		return c.String(http.StatusInternalServerError, "Error posting data")
	}

    expenseList.ID = id
    expenseListWrapper := lib.Calculate_shares_and_compensations(*expenseList) 

	log.Println("200 Success")
	return c.JSON(http.StatusOK, expenseListWrapper)

}

// HandlerCreateExpense godoc
//
//	@Summary		Create Expense
//	@Description	Creates a new Expense for an Expense List
//	@Tags			Expenses
//	@Param			expense			body		models.Expense				true	"Expense Data"
//	@Param			Authorization	header		string						true	"Bearer Token"
//	@Success		200				{object}	models.ExpenseListWrapper	"Updated Expense List with new Expense, Compensation and Shares"
//	@Failure		400				"Bad Request"
//	@Failure		500				"Internal Server Error"
//	@Router			/expense [post]
func HandlerCreateExpense(c echo.Context) error {
	log.Println("POST Create new Expense")
	expense := new(models.Expense)
	if err := c.Bind(expense); err != nil {
		log.Println("400 Bad Request")
        log.Println(err)
		return c.String(http.StatusBadRequest, "Bad Request")
	}

	err := db.CreateExpense(expense)
	if err != nil {
		log.Println("500 Internal Server Error")
		log.Println(err)
		return c.String(http.StatusInternalServerError, "Error posting data")
	}

    expenseList, err := db.GetExpenseList(expense.ExpenseListId)
    if err != nil {
		return c.String(http.StatusInternalServerError, "Error posting data")
    }

    expenseListWrapper := lib.Calculate_shares_and_compensations(*expenseList)

	log.Println("200 Sucess")
	return c.JSON(http.StatusOK, expenseListWrapper)

}
