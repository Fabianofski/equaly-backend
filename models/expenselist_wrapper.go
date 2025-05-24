package models

type ExpenseListWrapper struct {
	ExpenseList   ExpenseList               `json:"expenseList"`
	Shares        []*ExpenseListShare       `json:"shares"`
	Compensations []ExpenseListCompensation `json:"compensations"`
}
