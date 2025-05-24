package models

import (
	"encoding/json"
	"strings"
	"time"
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
		Participants json.RawMessage `json:"participants"`
		*Alias
	}{
		Alias: (*Alias)(e),
	}

	if err := json.Unmarshal(data, aux); err != nil {
		return err
	}

	if aux.Participants != nil {
		var participantsStr string
		if err := json.Unmarshal(aux.Participants, &participantsStr); err == nil {
			e.Participants = strings.Split(participantsStr, ",")
		} else {
			var participantsList []string
			if err := json.Unmarshal(aux.Participants, &participantsList); err == nil {
				e.Participants = participantsList
			} else {
				return err
			}
		}
	} else {
		e.Participants = []string{}
	}

	return nil
}
