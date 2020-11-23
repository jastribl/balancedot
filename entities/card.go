package entities

// Card represents a singel card from the user
type Card struct {
	ID           uint           `json"id" gorm:"AUTO_INCREMENT"`
	LastFour     string         `json:"last_four"`
	BankName     string         `json:"bank_name"`
	CardActivity []CardActivity `json:"card_activities"`
}
