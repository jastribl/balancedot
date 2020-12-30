package entities

import uuid "github.com/satori/go.uuid"

// Card represents a singel card from the user
type Card struct {
	UUID        uuid.UUID      `json:"uuid" gorm:"primary_key;"`
	LastFour    string         `json:"last_four"`
	Description string         `json:"description"`
	Activities  []CardActivity `json:"activities"`
	BankName    BankNames      `json:"bank_name"`
}
