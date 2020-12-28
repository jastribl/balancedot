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

// GetAllExpensesOrdered fetches all SplitwiseExpenses ordered
func (m *SplitwiseExpenseRepo) GetAllExpensesOrdered() ([]*entities.SplitwiseExpense, error) {
	var expenses []*entities.SplitwiseExpense
	err := m.Order("date DESC").Find(&expenses).Error
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
