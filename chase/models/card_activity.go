package models

import (
	"fmt"

	"gihub.com/jastribl/balancedot/entities"
)

// CardActivity is the interface for card activities
type CardActivity interface {
	ToCardActivityEntity(card *entities.Card) *entities.CardActivity
	ToCardActivityUniqueMatcher(card *entities.Card) *entities.CardActivity
}

// ChaseCardActivity holds line items from the Chase Card Activity Report file
type ChaseCardActivity struct {
	TransactionDate ChaseDate   `csv:"Transaction Date"`
	PostDate        ChaseDate   `csv:"Post Date"`
	Description     string      `csv:"Description"`
	Category        string      `csv:"Category"`
	Type            string      `csv:"Type"`
	Amount          MoneyAmount `csv:"Amount"`
	Memo            string      `csv:"Memo"`
}

// ToCardActivityEntity converts to an CardActivity entity
func (m *ChaseCardActivity) ToCardActivityEntity(card *entities.Card) *entities.CardActivity {
	return &entities.CardActivity{
		CardUUID:        card.UUID,
		TransactionDate: m.TransactionDate.Time,
		PostDate:        m.PostDate.Time,
		Description:     m.Description,
		Category:        m.Category,
		Type:            m.Type,
		Amount:          m.Amount.ToFloat64(),
	}
}

// ToCardActivityUniqueMatcher converts to an CardActivity entity matcher
func (m *ChaseCardActivity) ToCardActivityUniqueMatcher(card *entities.Card) *entities.CardActivity {
	return &entities.CardActivity{
		CardUUID:        card.UUID,
		TransactionDate: m.TransactionDate.Time,
		PostDate:        m.PostDate.Time,
		Description:     m.Description,
		Type:            m.Type,
		Amount:          m.Amount.ToFloat64(),
	}
}

// BofACardActivity holds line items from the BofA Card Activity Report file
type BofACardActivity struct {
	TransactionDate BofADate
	PostingDate     BofADate
	Description     string
	ReferenceNumber string
	AccountNumber   string
	Amount          MoneyAmount
}

// ToCardActivityEntity converts to an CardActivity entity
func (m *BofACardActivity) ToCardActivityEntity(card *entities.Card) *entities.CardActivity {
	return &entities.CardActivity{
		CardUUID:        card.UUID,
		TransactionDate: m.TransactionDate.Time,
		PostDate:        m.PostingDate.Time,
		Description:     m.Description,
		Type:            fmt.Sprintf("Ref: %s, Account: %s", m.ReferenceNumber, m.AccountNumber),
		Amount:          m.Amount.ToFloat64(),
	}
}

// ToCardActivityUniqueMatcher  converts to an CardActivity entity macher
func (m *BofACardActivity) ToCardActivityUniqueMatcher(card *entities.Card) *entities.CardActivity {
	return m.ToCardActivityEntity(card)
}
