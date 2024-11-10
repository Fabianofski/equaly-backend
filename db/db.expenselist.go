package db

import (
	"encoding/json"
	"log"

	. "github.com/fabianofski/equaly-backend/models"
)


func GetExpenseLists(userId string) ([]ExpenseList, error) {
    db, err := GetPostgresConnection()
    if err != nil {
        return nil, err
    }

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

func CreateExpenseList(expenseList *ExpenseList) error {
    db, err := GetPostgresConnection()
    if err != nil {
        return err
    }

    query := `
        INSERT INTO ExpenseLists (id, color, emoji, title, creatorId, currency)
        VALUES ($1, $2, $3, $4, $5, $6)
    `

    _, err = db.Exec(query, 
        expenseList.ID, 
        expenseList.Color, 
        expenseList.Emoji, 
        expenseList.Title, 
        expenseList.CreatorId, 
        expenseList.Currency,
    )

    if err != nil {
        return err
    }

    return nil;
}

func CreateExpense(expense *Expense) error {
    db, err := GetPostgresConnection()
    if err != nil {
        return err
    }

    query := `
        INSERT INTO Expenses (id, expenseListId, buyer, amount, description, participants)
        VALUES ($1, $2, $3, $4, $5, $6)
    `

    _, err = db.Exec(query, 
        expense.ID,
        expense.ExpenseListId,
        expense.Buyer,
        expense.Amount,
        expense.Description,
        expense.Participants,
    )

    if err != nil {
        return err
    }

    return nil;
}
