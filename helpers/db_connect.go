package helpers

import (
	"fmt"
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres" // imported to allow postgres connections
)

// DbConnect opens a connection to the database and returns the db object
func DbConnect() (*gorm.DB, error) {
	t := fmt.Sprintf(
		"host=%s user=%s dbname=%s sslmode=disable password=%s",
		os.Getenv("DB_URL"),
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_DB"),
		os.Getenv("POSTGRES_PASSWORD"),
	)
	return gorm.Open("postgres", t)
}
