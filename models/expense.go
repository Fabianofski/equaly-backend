package models

import "time"

type Expense struct {
	ID            string    `json:"id"`
	ExpenseListId string    `json:"expenseListId"`
	Buyer         string    `json:"buyer"`
	Amount        float64   `json:"amount"`
	Description   string    `json:"description"`
	Participants  string    `json:"participants"`
	Date          time.Time `json:"date"`
}
