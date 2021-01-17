package models

import (
	"gihub.com/jastribl/balancedot/entities"
)

// AccountActivity is the interface for account activities
type AccountActivity interface {
	ToAccountActivitiyEntity(account *entities.Account) *entities.AccountActivity
}

// ChaseAccountActivity holds line items from the Chase Chequing Account Activity Report file
type ChaseAccountActivity struct {
	Details           string      `csv:"Details"`
	PostingDate       ChaseDate   `csv:"Posting Date"`
	Description       string      `csv:"Description"`
	Amount            MoneyAmount `csv:"Amount"`
	Type              string      `csv:"Type"`
	Balance           MoneyAmount `csv:"Balance"`
	CheckOrSlipNumber string      `csv:"Check or Slip #"`
}

// ToAccountActivitiyEntity converts to an AccountActivity entity
func (m *ChaseAccountActivity) ToAccountActivitiyEntity(account *entities.Account) *entities.AccountActivity {
	return &entities.AccountActivity{
		AccountUUID: account.UUID,
		Details:     m.Details,
		PostingDate: m.PostingDate.Time,
		Description: m.Description,
		Amount:      m.Amount.ToFloat64(),
		Type:        m.Type,
	}
}

// todo: re-arrange this folder structure to reflect that it contains both bofa and chase things
// something like:
// banks
//     models.go
//     bofa
//         cards.go
//         accounts.go
//     chase
//         cards.go
//         accounts.go

// BofAAccountActivity holds line items from the BofA Chequing Account Activity Report file
type BofAAccountActivity struct {
	Date           BofADate    `csv:"Date"`
	Description    string      `csv:"Description"`
	Amount         MoneyAmount `csv:"Amount"`
	RunningBalance MoneyAmount `csv:"Running Bal."`
}

// ToAccountActivitiyEntity converts to an AccountActivity entity
func (m *BofAAccountActivity) ToAccountActivitiyEntity(account *entities.Account) *entities.AccountActivity {
	return &entities.AccountActivity{
		AccountUUID: account.UUID,
		PostingDate: m.Date.Time,
		Description: m.Description,
		Amount:      m.Amount.ToFloat64(),
	}
}
