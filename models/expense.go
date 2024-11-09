package models

type Expense struct {
	Buyer        string  `json:"buyer"`
	Amount       float64 `json:"amount"`
	Description  string  `json:"description"`
	Participants string  `json:"participants"`
}
