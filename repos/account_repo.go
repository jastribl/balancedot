package repos

import (
	"gihub.com/jastribl/balancedot/entities"
	"github.com/jinzhu/gorm"
)

// AccountRepo is the repo for Accounts
type AccountRepo struct {
	*gorm.DB
}

// NewAccountRepo returns a new AccountRepo using the given db
func NewAccountRepo(db *gorm.DB) *AccountRepo {
	return &AccountRepo{
		DB: db,
	}
}

// GetAllAccounts fetches all Accounts
func (m *AccountRepo) GetAllAccounts() ([]*entities.Account, error) {
	var accounts []*entities.Account
	err := m.Find(&accounts).Error
	if err != nil {
		return nil, err
	}
	return accounts, nil
}

// GetAccountByUUID fetches a single Account by UUID
func (m *AccountRepo) GetAccountByUUID(uuid string) (*entities.Account, error) {
	account := &entities.Account{}
	err := m.Where("uuid = ?", uuid).Find(account).Error
	if err != nil {
		return nil, err
	}
	return account, nil
}

// // GetAccount returns the Account for the given lastFour
// func (m *AccountRepo) GetCard(lastFour string) (*entities.Card, error) {
// 	card := &entities.Card{}
// 	err := m.Preload("CardActivity").Find(card, &entities.Card{LastFour: lastFour}).Error
// 	if err != nil {
// 		return nil, err
// 	}
// 	return card, nil
// }
