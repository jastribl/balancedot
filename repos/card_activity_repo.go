package repos

import (
	"gihub.com/jastribl/balancedot/entities"
	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
)

// CardActivityRepo is the repo for Cards
type CardActivityRepo struct {
	*gorm.DB
}

// NewCardActivityRepo returns a new CardActivityRepo using the given db
func NewCardActivityRepo(db *gorm.DB) *CardActivityRepo {
	return &CardActivityRepo{
		DB: db,
	}
}

// GetCardActivity returns the CarActivity for the given uuid
func (m *CardActivityRepo) GetCardActivity(uuid uuid.UUID) (*entities.CardActivity, error) {
	cardActivity := &entities.CardActivity{}
	err := m.Preload("Card").Find(cardActivity, &entities.CardActivity{UUID: uuid}).Error
	if err != nil {
		return nil, err
	}
	return cardActivity, nil
}
