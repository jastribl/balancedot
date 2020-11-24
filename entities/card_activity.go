package entities

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

// CardActivity represents a singel card activity record from the user
type CardActivity struct {
	UUID            uuid.UUID `json:"uuid"`
	CardUUID        uuid.UUID `json:"card_uuid"`
	Card            *Card     `json:"card"  gorm:"ForeignKey:CardUUID"`
	TransactionDate time.Time `json:"transaction_date"`
	PostDate        time.Time `json:"post_date"`
	Description     string    `json:"description"`
	Category        string    `json:"category"`
	Type            string    `json:"type"`
	Amount          float64   `json:"amount"`
}

// BeforeCreate is the before trigger for CardActivity
func (u *CardActivity) BeforeCreate() error {
	if u.UUID == uuid.Nil {
		u.UUID = uuid.NewV4()
	}

	return nil
}
