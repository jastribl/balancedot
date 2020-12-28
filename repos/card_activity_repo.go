package repos

import (
	"gihub.com/jastribl/balancedot/entities"
	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
)

// CardActivityRepo is the repo for CardActivities
type CardActivityRepo struct {
	*gorm.DB
}

// GetAllCardActivitiesForCard fetches all CardActivities for a given CardUUID
func (m *CardActivityRepo) GetAllCardActivitiesForCard(cardUUID string) ([]*entities.CardActivity, error) {
	var cardActivities []*entities.CardActivity
	err := m.Where("card_uuid = ?", cardUUID).Find(&cardActivities).Error
	if err != nil {
		return nil, err
	}
	return cardActivities, nil
}

// NewCardActivityRepo returns a new CardActivityRepo using the given db
func NewCardActivityRepo(db *gorm.DB) *CardActivityRepo {
	return &CardActivityRepo{
		DB: db,
	}
}

// GetCardActivity returns the CardActivity for the given uuid
func (m *CardActivityRepo) GetCardActivity(uuid uuid.UUID) (*entities.CardActivity, error) {
	cardActivity := &entities.CardActivity{}
	err := m.Preload("Card").Find(cardActivity, &entities.CardActivity{UUID: uuid}).Error
	if err != nil {
		return nil, err
	}
	return cardActivity, nil
}
