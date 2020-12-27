package api

import (
	"gihub.com/jastribl/balancedot/config"
	"github.com/jinzhu/gorm"
)

// App contains API handler methods
type App struct {
	db     *gorm.DB
	config *config.Config
}

// SaveEntity saves an entity and handles errors
func (m *App) SaveEntity(entity interface{}) error {
	return m.db.Save(entity).Error
}

// NewApp returns a new App
func NewApp(db *gorm.DB, config *config.Config) (*App, error) {
	return &App{
		db:     db,
		config: config,
	}, nil
}
