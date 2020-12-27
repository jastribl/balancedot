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

// TransactionResponse is a response for a transaction callback
type TransactionResponse struct {
	Rollback bool
	Handle   func() error
}

// TransactionResponseFromUnexpectedError returns an unexpected error response
func TransactionResponseFromUnexpectedError(err error) TransactionResponse {
	return TransactionResponse{
		Rollback: true,
		Handle: func() error {
			return err
		},
	}
}

// TransactionFailedHandle returns a normal failed transation with proper handling
func TransactionFailedHandle(fn func() error) TransactionResponse {
	return TransactionResponse{
		Rollback: true,
		Handle:   fn,
	}
}

// TransactionSuccess returns a successful transaction result
func TransactionSuccess() TransactionResponse {
	return TransactionResponse{
		Rollback: false,
		Handle: func() error {
			return nil
		},
	}
}

// TransactionAction is the action that wishes to be taken from the transaction
type TransactionAction int

const (
	// TransactionActionCommit means the transaction should be committed
	TransactionActionCommit TransactionAction = iota
	// TransactionActionRollback means the transaction should be rolled back
	TransactionActionRollback
)

// NewTransaction wraps some transaction logic with rollback and commit
func NewTransaction(db *gorm.DB, fn func(tx *gorm.DB) TransactionAction) bool {
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()
	action := fn(tx)
	if action == TransactionActionCommit {
		tx.Commit()
		return true
	}

	tx.Rollback()
	return false
}
