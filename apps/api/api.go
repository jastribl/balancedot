package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/jinzhu/gorm"
)

// App contains API handler methods
type App struct {
	db *gorm.DB
}

// DecodeParams decodes the request body into the params structure
func (m *App) DecodeParams(r *http.Request, params interface{}) {
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&params)
	if err != nil {
		panic(err)
	}
}

// SaveEntity saves an entity and handles errors
func (m *App) SaveEntity(entity interface{}) error {
	return m.db.Save(entity).Error
}

// NewApp returns a new App
func NewApp(db *gorm.DB) (*App, error) {
	a := &App{
		db: db,
	}
	return a, nil
}

// Error is a generic error format for the api
type Error struct {
	Error   interface{}
	Message string
	Code    int
}

type errorResponse struct {
	Message string `json:"message"`
}

func (e *Error) renderError(w http.ResponseWriter, r *http.Request) {
	defer json.NewEncoder(log.Writer()).Encode(e)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(e.Code)
	json.NewEncoder(w).Encode(errorResponse{
		Message: e.Message,
	})
}

// Handler is a a generic handler for the api
type Handler func(http.ResponseWriter, *http.Request) *Error

func (fn Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if ree := recover(); ree != nil {
			unepxectedError := Error{
				Error:   ree,
				Message: "Unexpected Error",
				Code:    500,
			}
			unepxectedError.renderError(w, r)
		}
	}()

	if e := fn(w, r); e != nil {
		e.renderError(w, r)
	}
}
