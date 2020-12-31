package models

import "gihub.com/jastribl/balancedot/entities"

// CardActivity is the interface for card activities
type CardActivity interface {
	ToCardActivitiyEntity(card *entities.Card) *entities.CardActivity
	ToCardActivitiyUniqueMatcher(card *entities.Card) *entities.CardActivity
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

// ToCardActivitiyEntity converts to an CardActivity entity
func (m *ChaseCardActivity) ToCardActivitiyEntity(card *entities.Card) *entities.CardActivity {
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

// ToCardActivitiyUniqueMatcher converts to an CardActivity entity matcher
func (m *ChaseCardActivity) ToCardActivitiyUniqueMatcher(card *entities.Card) *entities.CardActivity {
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
	TransactionDate BofADate    `csv:"Transaction Date"`
	PostDate        BofADate    `csv:"Post Date"`
	Description     string      `csv:"Description"`
	Category        string      `csv:"Category"`
	Type            string      `csv:"Type"`
	Amount          MoneyAmount `csv:"Amount"`
	Memo            string      `csv:"Memo"`
}

// ToCardActivitiyEntity converts to an CardActivity entity
func (m *BofACardActivity) ToCardActivitiyEntity(card *entities.Card) *entities.CardActivity {
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

// ToCardActivitiyUniqueMatcher  converts to an CardActivity entity macher
func (m *BofACardActivity) ToCardActivitiyUniqueMatcher(card *entities.Card) *entities.CardActivity {
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
