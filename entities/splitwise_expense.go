package entities

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

// SplitwiseExpense represents a singel splitwise expense
type SplitwiseExpense struct {
	UUID               uuid.UUID          `json:"uuid" gorm:"primary_key; default:uuid_generate_v4();"`
	SplitwiseID        int                `json:"splitwise_id"`
	Description        string             `json:"description"`
	Details            string             `json:"details"`
	CurrencyCode       string             `json:"currency_code"`
	Amount             float64            `json:"amount"`
	AmountPaid         float64            `json:"amount_paid"`
	Date               time.Time          `json:"date"`
	SplitwiseCreatedAt time.Time          `json:"splitwise_created_at"`
	SplitwiseUpdatedAt *time.Time         `json:"splitwise_updated_at"`
	SplitwiseDeletedAt *time.Time         `json:"splitwise_deleted_at"`
	Category           string             `json:"category"`
	CreationMethod     *string            `json:"creation_method"`
	CardActivities     []*CardActivity    `json:"card_activities" gorm:"many2many:expense_links;"`
	AccountActivities  []*AccountActivity `json:"account_activities" gorm:"many2many:account_activity_links;"`
}
