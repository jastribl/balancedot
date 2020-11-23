package repos

import (
	"gihub.com/jastribl/balancedot/entities"
	"github.com/jinzhu/gorm"
)

// CardRepo is the repo for Cards
type CardRepo struct {
	*gorm.DB
}

// NewCardRepo returns a new CardRepo using the given db
func NewCardRepo(db *gorm.DB) *CardRepo {
	return &CardRepo{
		DB: db,
	}
}

// GetCard returns the Card for the given lastFour
func (m *CardRepo) GetCard(lastFour string) (*entities.Card, error) {
	card := &entities.Card{}
	err := m.Preload("CardActivity").Find(card, "last_four = ?", lastFour).Error
	if err != nil {
		return nil, err
	}
	return card, nil
}
