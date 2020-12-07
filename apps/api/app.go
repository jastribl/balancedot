package api

import "github.com/jinzhu/gorm"

// App contains API handler methods
type App struct {
	db *gorm.DB
}

// NewApp returns a new App
func NewApp(db *gorm.DB) (*App, error) {
	a := &App{
		db: db,
	}
	return a, nil
}
