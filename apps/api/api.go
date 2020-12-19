package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/jinzhu/gorm"
)

// UnexpectedError is the string for an unexpected error
const UnexpectedError = "Unexpected Error"

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

// FromErrUnexpected returns an unexpected 500 error from an err
func FromErrUnexpected(err *error) *Error {
	return &Error{
		Message: UnexpectedError,
		Error:   err,
		Code:    500,
	}
}

type errorResponse struct {
	Message string `json:"message"`
}

// ResponseWriter is a normal http ResponseWriter but with extras
type ResponseWriter struct {
	http.ResponseWriter
}

// RenderJSON renders antying as json as a 200 response
func (w *ResponseWriter) RenderJSON(toRender interface{}) interface{} {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(200)
	return json.NewEncoder(w).Encode(toRender)
}

// RenderError renders an error as json with the given code
func (w *ResponseWriter) RenderError(e *Error) {
	defer json.NewEncoder(log.Writer()).Encode(e)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(e.Code)
	json.NewEncoder(w).Encode(errorResponse{
		Message: e.Message,
	})
}

// Handler is a a generic handler for the api
type Handler func(ResponseWriter, *http.Request) interface{}

func (fn Handler) ServeHTTP(wIn http.ResponseWriter, r *http.Request) {
	w := ResponseWriter{wIn}
	defer func() {
		if ree := recover(); ree != nil {
			log.Printf("recovering from fatal: %+v", ree)
			w.RenderError(&Error{
				Message: UnexpectedError,
				Error:   ree,
				Code:    500,
			})
		}
	}()

	if e := fn(w, r); e != nil {
		var toRender *Error
		if err, ok := e.(error); ok {
			toRender = FromErrUnexpected(&err)
		} else if err, ok := e.(*error); ok {
			toRender = FromErrUnexpected(err)
		} else if err, ok := e.(Error); ok {
			toRender = &err
		} else if err, ok := e.(*Error); ok {
			toRender = err
		} else {
			log.Panic(err)
			toRender = &Error{
				Message: UnexpectedError,
				Error:   e,
				Code:    500,
			}
		}
		w.RenderError(toRender)
	}
}
