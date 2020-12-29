package entities

import uuid "github.com/satori/go.uuid"

// Account represents a singel account from the user
type Account struct {
	UUID        uuid.UUID         `json:"uuid" gorm:"primary_key;"`
	LastFour    string            `json:"last_four"`
	Description string            `json:"description"`
	Activities  []AccountActivity `json:"activities"`
}
