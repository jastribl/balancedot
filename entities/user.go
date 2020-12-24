package entities

import uuid "github.com/satori/go.uuid"

// User represents a singel user
type User struct {
	UUID     uuid.UUID `json:"uuid" gorm:"primary_key;"`
	Username string    `json:"username"`
	Password string    `json:"-"`
}
