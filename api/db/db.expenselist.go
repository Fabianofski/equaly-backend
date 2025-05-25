package db

import (
	"database/sql"
	"encoding/json"
	"log"
	"slices"
	"strings"

	. "github.com/fabianofski/equaly-backend/models"
)

func GetExpenseLists(userId string) ([]*ExpenseList, error) {
	db, err := GetPostgresConnection()
	if err != nil {
		return nil, err
	}

	rows, err := db.Query(`
            SELECT ExpenseLists.*,
							COALESCE(SUM(Expenses.amount), 0) AS totalCost,
							COALESCE(json_agg(row_to_json(Expenses)) FILTER (WHERE Expenses IS NOT NULL), '[]'::json) AS expenses
            FROM ExpenseLists
							LEFT JOIN Expenses ON ExpenseLists.id = Expenses.expenseListId
            WHERE ExpenseLists.creatorId = $1 
							OR ExpenseLists.id IN (
								SELECT expenseListId FROM ExpenseListsUsers WHERE userId = $1
							)
            GROUP BY ExpenseLists.id;
        `, userId)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer rows.Close()

	var expenseLists []*ExpenseList
	for rows.Next() {
		expenseList, err := RowToExpenseList(rows)
		if err != nil {
			return nil, err
		}
		expenseLists = append(expenseLists, expenseList)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return expenseLists, nil
}

func GetExpenseList(listId string) (*ExpenseList, error) {
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
            WHERE ExpenseLists.id = $1 
            GROUP BY ExpenseLists.id;
        `, listId)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer rows.Close()

	rows.Next()
	expenseList, err := RowToExpenseList(rows)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return expenseList, nil
}

func RowToExpenseList(rows *sql.Rows) (*ExpenseList, error) {
	expenseList := &ExpenseList{}
	var expensesJSON string
	var participantsJSON string

	err := rows.Scan(
		&expenseList.ID,
		&expenseList.Color,
		&expenseList.Emoji,
		&expenseList.Title,
		&expenseList.CreatorId,
		&expenseList.Currency,
		&participantsJSON,
		&expenseList.InviteCode,
		&expenseList.TotalCost,
		&expensesJSON,
	)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	err = json.Unmarshal([]byte(expensesJSON), &expenseList.Expenses)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	slices.SortFunc(expenseList.Expenses, func(a, b *Expense) int {
		return a.Date.Compare(b.Date)
	})

	err = json.Unmarshal([]byte(participantsJSON), &expenseList.Participants)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return expenseList, nil
}

func CreateExpenseList(expenseList *ExpenseList) (string, error) {
	db, err := GetPostgresConnection()
	if err != nil {
		return "", err
	}

	query := `
        INSERT INTO ExpenseLists (color, emoji, title, creatorId, currency, participants)
        VALUES ($1, $2, $3, $4, $5, $6)
        RETURNING id 
    `

	participants, err := json.Marshal(expenseList.Participants)
	if err != nil {
		return "", err
	}

	var id string
	err = db.QueryRow(query,
		expenseList.Color,
		expenseList.Emoji,
		expenseList.Title,
		expenseList.CreatorId,
		expenseList.Currency,
		participants,
	).Scan(&id)
	log.Println(id)

	if err != nil {
		return "", err
	}

	return id, nil
}

func CreateExpense(expense *Expense) error {
	db, err := GetPostgresConnection()
	if err != nil {
		return err
	}

	query := `
        INSERT INTO Expenses (expenseListId, buyer, amount, description, participants, date)
        VALUES ($1, $2, $3, $4, $5, $6)
    `

	_, err = db.Exec(query,
		expense.ExpenseListId,
		expense.Buyer,
		expense.Amount,
		expense.Description,
		strings.Join(expense.Participants, ","),
		expense.Date,
	)

	if err != nil {
		return err
	}

	return nil
}

func IsUserOwner(expenseListId string, userId string) (bool, error) {
	db, err := GetPostgresConnection()
	if err != nil {
		return false, err
	}

	var isOwner bool

	err = db.QueryRow(`
			SELECT EXISTS (
					SELECT 1 FROM ExpenseLists
					WHERE id = $1 AND creatorId = $2
			)
	`, expenseListId, userId).Scan(&isOwner)

	if err != nil {
		return false, err
	}

	return isOwner, nil
}

func IsUserAuthorized(expenseListId string, userId string) (bool, error) {
	db, err := GetPostgresConnection()
	if err != nil {
		return false, err
	}

	var isAuthorized bool

	err = db.QueryRow(`
			SELECT EXISTS (
					SELECT 1 FROM ExpenseLists
					WHERE id = $1 AND creatorId = $2
					UNION
					SELECT 1 FROM ExpenseListsUsers
					WHERE expenseListId = $1 AND userId = $2
			)
	`, expenseListId, userId).Scan(&isAuthorized)

	if err != nil {
		return false, err
	}

	return isAuthorized, nil
}

func IsInviteCodeValid(expenseListId string, inviteCode string) (bool, error) {
	db, err := GetPostgresConnection()
	if err != nil {
		return false, err
	}

	var inviteCodeValid bool

	err = db.QueryRow(`
			SELECT EXISTS (
					SELECT 1 FROM ExpenseLists
					WHERE id = $1 AND inviteCode = $2
			)
	`, expenseListId, inviteCode).Scan(&inviteCodeValid)

	if err != nil {
		return false, err
	}

	return inviteCodeValid, nil

}

func JoinExpenseList(expenseListId string, userId string) error {
	db, err := GetPostgresConnection()
	if err != nil {
		return err
	}

	query := `
        INSERT INTO ExpenseListsUsers (expenseListId, userId)
        VALUES ($1, $2)
    `

	_, err = db.Exec(query, expenseListId, userId)

	if err != nil {
		return err
	}

	return nil

}
