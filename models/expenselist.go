package models

type ExpenseList struct {
	ID           string        `json:"id"`
	Color        string        `json:"color"`
	Emoji        string        `json:"emoji"`
	Title        string        `json:"title"`
	TotalCost    float64       `json:"totalCost"`
	CreatorId    string        `json:"creatorId"`
	Currency     string        `json:"currency"`
	InviteCode   string        `json:"inviteCode"`
	Expenses     []Expense     `json:"expenses"`
	Participants []Participant `json:"participants"`
}
