package repos

import (
	"gihub.com/jastribl/balancedot/entities"
	"github.com/jinzhu/gorm"
)

// SplitwiseExpenseRepo is the repo for SplitwiseExpenses
type SplitwiseExpenseRepo struct {
	*gorm.DB
}

// NewSplitwiseExpenseRepo returns a new SplitwiseExpenseRepo using the given db
func NewSplitwiseExpenseRepo(db *gorm.DB) *SplitwiseExpenseRepo {
	return &SplitwiseExpenseRepo{
		DB: db,
	}
}

// GetAllExpenses fetches all SplitwiseExpenses
func (m *SplitwiseExpenseRepo) GetAllExpenses() ([]*entities.SplitwiseExpense, error) {
	var expenses []*entities.SplitwiseExpense
	err := m.Find(&expenses).Error
	if err != nil {
		return nil, err
	}
	return expenses, nil
}

// GetAllExpensesWithCardActivities fetches all SplitwiseExpenses with the linked card expenses
func (m *SplitwiseExpenseRepo) GetAllExpensesWithCardActivities() ([]*entities.SplitwiseExpense, error) {
	var expenses []*entities.SplitwiseExpense
	err := m.Preload("CardActivities").Where("deleted_at is NULL").Find(&expenses).Error
	if err != nil {
		return nil, err
	}
	return expenses, nil
}
