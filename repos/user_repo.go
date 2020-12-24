package repos

import (
	"gihub.com/jastribl/balancedot/entities"
	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
)

// UserRepo is the repo for Users
type UserRepo struct {
	*gorm.DB
}

// NewUserRepo returns a new UserRepo using the given db
func NewUserRepo(db *gorm.DB) *UserRepo {
	return &UserRepo{
		DB: db,
	}
}

// GetUserByUUID todo
func (m *UserRepo) GetUserByUUID(uuid uuid.UUID) (*entities.User, error) {
	user := &entities.User{}
	err := m.Where("uuid = ?", uuid).Find(user).Error
	if err != nil {
		return nil, err
	}
	return user, nil
}
