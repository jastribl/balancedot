package entities

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

// SplitwiseExpense represents a singel splitwise expense
type SplitwiseExpense struct {
	UUID               uuid.UUID       `json:"uuid" gorm:"primary_key;"`
	SplitwiseID        int             `json:"splitwise_id"`
	Description        string          `json:"description"`
	Details            string          `json:"details"`
	CurrencyCode       string          `json:"currency_code"`
	Amount             float64         `json:"amount"`
	AmountPaid         float64         `json:"amount_paid"`
	Date               time.Time       `json:"date"`
	SplitwiseCreatedAt time.Time       `json:"splitwise_created_at"`
	SplitwiseUpdatedAt *time.Time      `json:"splitwise_updated_at"`
	SplitwiseDeletedAt *time.Time      `json:"splitwise_deleted_at"`
	Category           string          `json:"category"`
	CardActivities     []*CardActivity `json:"card_activities" gorm:"many2many:expense_links;"`
}
