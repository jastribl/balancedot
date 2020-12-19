package helpers

import (
	"fmt"
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres" // imported to allow postgres connections
	"github.com/lib/pq"
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

// IsUniqueConstraintError returns whether an error is a unique constraint error
func IsUniqueConstraintError(err error, constraintName string) bool {
	if pqErr, ok := err.(*pq.Error); ok {
		return pqErr.Code == "23505" && pqErr.Constraint == constraintName
	}
	return false
}

// RowExists checks if a given search row exists for a given model
func RowExists(db *gorm.DB, model interface{}, search interface{}) (bool, error) {
	foundRows, err := db.Model(model).Where(search).Select("*").Rows()
	if err != nil {
		return false, err
	}
	if foundRows.Next() {
		return true, nil
	}

	return false, nil
}

// NewTransaction wraps some transaction logic with rollback and commit
func NewTransaction(db *gorm.DB, fn func(tx *gorm.DB) interface{}) interface{} {
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()
	err := fn(tx)
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}
