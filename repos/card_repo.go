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

// GetAllCards fetches all Cards
func (m *CardRepo) GetAllCards() ([]*entities.Card, error) {
	var cards []*entities.Card
	err := m.Find(&cards).Error
	if err != nil {
		return nil, err
	}
	return cards, nil
}

// GetCard returns the Card for the given lastFour
func (m *CardRepo) GetCard(lastFour string) (*entities.Card, error) {
	card := &entities.Card{}
	err := m.Preload("CardActivity").Find(card, &entities.Card{LastFour: lastFour}).Error
	if err != nil {
		return nil, err
	}
	return card, nil
}
