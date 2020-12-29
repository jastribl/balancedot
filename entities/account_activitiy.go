package entities

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

// AccountActivity represents a singel account activity record from the user
type AccountActivity struct {
	UUID              uuid.UUID           `json:"uuid" gorm:"primary_key; default:uuid_generate_v4();"`
	AccountUUID       uuid.UUID           `json:"account_uuid"`
	Account           *Account            `json:"account" gorm:"foreignKey:AccountUUID"`
	Details           string              `json:"details"`
	PostingDate       time.Time           `json:"posting_date"`
	Description       string              `json:"description"`
	Amount            float64             `json:"amount"`
	Type              string              `json:"type"`
	SplitwiseExpenses []*SplitwiseExpense `json:"splitwise_expenses" gorm:"many2many:account_activity_links;"`
}
