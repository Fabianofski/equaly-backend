package models

type ExpenseListShare struct {
	ParticipantId    string  `json:"id"`
	NumberOfExpenses int     `json:"numberOfExpenses"`
	ExpenseAmount    float64 `json:"expenseAmount"`
	Share            float64 `json:"share"`
	Difference       float64 `json:"difference"`
}
