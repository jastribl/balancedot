package entities

import uuid "github.com/satori/go.uuid"

// Card represents a singel card from the user
type Card struct {
	UUID         uuid.UUID      `json:"uuid"`
	LastFour     string         `json:"last_four"`
	BankName     string         `json:"bank_name"`
	CardActivity []CardActivity `json:"card_activities"`
}

// BeforeCreate is the before trigger for Card
func (u *Card) BeforeCreate() error {
	if u.UUID == uuid.Nil {
		u.UUID = uuid.NewV4()
	}

	return nil
}
