package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"
)

func connect_to_postgres() *sql.DB {
	host := os.Getenv("POSTGRES_HOST")
	portStr := os.Getenv("POSTGRES_PORT")
	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	dbname := os.Getenv("POSTGRES_DBNAME")

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

type ExpenseList struct {
	ID        string    `json:"id"`
	Color     string    `json:"color"`
	Emoji     string    `json:"emoji"`
	Title     string    `json:"title"`
	TotalCost float64   `json:"totalCost"`
	CreatorId string    `json:"creatorId"`
	Currency  string    `json:"currency"`
	Expenses  []Expense `json:"expenses"`
}

type Expense struct {
	Buyer        string  `json:"buyer"`
	Amount       float64 `json:"amount"`
	Description  string  `json:"description"`
	Participants string  `json:"participants"`
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

	e.GET("/user-expense-lists", func(c echo.Context) error {

		userId := c.QueryParam("userId")
		if userId == "" {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Bad Request"})
		}
		expenseLists, err := getExpenseLists(db, userId)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"error": "Error requesting data",
			})
		}
		return c.JSON(http.StatusOK, expenseLists)

	})

	e.Logger.Fatal(e.Start(":3000"))
}

func getExpenseLists(db *sql.DB, userId string) ([]ExpenseList, error) {
	rows, err := db.Query(`
            SELECT ExpenseLists.*,
                   COALESCE(SUM(Expenses.amount), 0)                                   AS totalCost,
                   COALESCE(json_agg(row_to_json(Expenses)) FILTER (WHERE Expenses IS NOT NULL), '[]'::json) AS expenses
            FROM ExpenseLists
                     LEFT JOIN Expenses ON ExpenseLists.id = Expenses.expenseListId
            WHERE ExpenseLists.creatorId = $1 
            GROUP BY ExpenseLists.id;
        `, userId)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer rows.Close()

	var expenseLists []ExpenseList
	for rows.Next() {
		var expenseList ExpenseList
		var expensesJSON string

		err := rows.Scan(&expenseList.ID, &expenseList.Color, &expenseList.Emoji, &expenseList.Title, &expenseList.CreatorId, &expenseList.Currency, &expenseList.TotalCost, &expensesJSON)
		if err != nil {
			log.Println(err)
			return nil, err
		}

		err = json.Unmarshal([]byte(expensesJSON), &expenseList.Expenses)
		if err != nil {
			log.Println(err)
			return nil, err
		}

		expenseLists = append(expenseLists, expenseList)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return expenseLists, nil
}
