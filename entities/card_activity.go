package entities

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

// CardActivity represents a singel card activity record from the user
type CardActivity struct {
	UUID              uuid.UUID           `json:"uuid" gorm:"primary_key; default:uuid_generate_v4();"`
	CardUUID          uuid.UUID           `json:"card_uuid"`
	Card              *Card               `json:"card" gorm:"foreignKey:CardUUID"`
	TransactionDate   time.Time           `json:"transaction_date"`
	PostDate          time.Time           `json:"post_date"`
	Description       string              `json:"description"`
	Category          string              `json:"category"`
	Type              string              `json:"type"`
	Amount            float64             `json:"amount"`
	SplitwiseExpenses []*SplitwiseExpense `json:"splitwise_expenses" gorm:"many2many:expense_links;"`
}
