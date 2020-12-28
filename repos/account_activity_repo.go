package repos

import (
	"gihub.com/jastribl/balancedot/entities"
	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
)

// AccountActivityRepo is the repo for AccountActivities
type AccountActivityRepo struct {
	*gorm.DB
}

// GetAllAccountActivitiesForAccount fetches all AccountActivities for a given AccountUUID
func (m *AccountActivityRepo) GetAllAccountActivitiesForAccount(accountUUID string) ([]*entities.AccountActivity, error) {
	var accountActivities []*entities.AccountActivity
	err := m.Where("account_uuid = ?", accountUUID).Find(&accountActivities).Error
	if err != nil {
		return nil, err
	}
	return accountActivities, nil
}

// NewAccountActivityRepo returns a new AccountActivityRepo using the given db
func NewAccountActivityRepo(db *gorm.DB) *AccountActivityRepo {
	return &AccountActivityRepo{
		DB: db,
	}
}

// GetAccountActivity returns the AccountActivity for the given uuid
func (m *AccountActivityRepo) GetAccountActivity(uuid uuid.UUID) (*entities.AccountActivity, error) {
	accountActivity := &entities.AccountActivity{}
	err := m.Preload("Account").Find(accountActivity, &entities.AccountActivity{UUID: uuid}).Error
	if err != nil {
		return nil, err
	}
	return accountActivity, nil
}
