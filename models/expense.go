package models

import (
    "time"
    "encoding/json"
    "strings"
)

type Expense struct {
	ID            string    `json:"id"`
	ExpenseListId string    `json:"expenseListId"`
	Buyer         string    `json:"buyer"`
	Amount        float64   `json:"amount"`
	Description   string    `json:"description"`
	Participants  []string  `json:"participants"`
	Date          time.Time `json:"date"`
}

func (e *Expense) UnmarshalJSON(data []byte) error {
	type Alias Expense
	aux := &struct {
		Participants string `json:"participants"`
		*Alias
	}{
		Alias: (*Alias)(e),
	}

	if err := json.Unmarshal(data, aux); err != nil {
		return err
	}

	if aux.Participants != "" {
		e.Participants = strings.Split(aux.Participants, ",")
	} else {
		e.Participants = []string{}
	}

	return nil
}
