package entities

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

// CardActivity represents a singel card activity record from the user
type CardActivity struct {
	ID              uint      `json"id" gorm:"AUTO_INCREMENT"`
	UUID            uuid.UUID `json:"uuid"`
	CardID          int       `json:"card_id"`
	Card            Card      `json:"card"  gorm:"ForeignKey:CardID"`
	TransactionDate time.Time `json:"transaction_date"`
	PostDate        time.Time `json:"post_date"`
	Description     string    `json:"description"`
	Category        string    `json:"category"`
	Type            string    `json:"type"`
	Amount          float64   `json:"amount"` // todo: use models.MoneyAmount
}

// BeforeCreate is the before trigger for CardActivity
func (u *CardActivity) BeforeCreate() error {
	if u.UUID == uuid.Nil {
		u.UUID = uuid.NewV4()
	}

	return nil
}
