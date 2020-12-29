package api

import (
	"encoding/json"
	"log"
	"net/http"
	"runtime/debug"
)

// WriterResponse is the response of the writer to be used mostly for internal return types and tracking
type WriterResponse int

const (
	// WriterResponseSuccess is for success
	WriterResponseSuccess WriterResponse = iota

	// WriterResponseError is for error
	WriterResponseError
)

// ResponseWriter is a normal http ResponseWriter but with extras
type ResponseWriter interface {
	SendResponse(i interface{}, extras ...interface{}) WriterResponse
	SendResponseWithCode(i interface{}, code int, extras ...interface{}) WriterResponse
	SendSimpleMessage(message string, extras ...interface{}) WriterResponse
	SendSimpleMessageWithCode(message string, code int, extras ...interface{}) WriterResponse
	SendError(message string, code int, extras ...interface{}) WriterResponse
	SendUnexpectedError(err interface{}, extras ...interface{}) WriterResponse
}

type writer struct {
	http.ResponseWriter
}

// NewResponseWriter constructs and returns a new ResponseWriter instance given the http.ResponseWriter
func NewResponseWriter(w http.ResponseWriter) ResponseWriter {
	return &writer{
		ResponseWriter: w,
	}
}

// SendResponseWithCode sends arbitrary data as JSON with the given code
func (w *writer) SendResponseWithCode(i interface{}, code int, extras ...interface{}) WriterResponse {
	w.WriteHeader(code)
	err := json.NewEncoder(w).Encode(i)
	w.printExtras(extras...)
	if err != nil {
		http.Error(w, "Serialization error", http.StatusInternalServerError)
	}
	return WriterResponseSuccess
}

// SendResponse sends arbitrary data as JSON
func (w *writer) SendResponse(i interface{}, extras ...interface{}) WriterResponse {
	return w.SendResponseWithCode(i, http.StatusOK, extras...)
}

// SendSimpleMessage sends a simple "message" response with http.StatusOK
func (w *writer) SendSimpleMessage(message string, extras ...interface{}) WriterResponse {
	return w.SendSimpleMessageWithCode(message, http.StatusOK, extras...)
}

// SendSimpleMessageWithCode sends a simple "message" response with a given code
func (w *writer) SendSimpleMessageWithCode(message string, code int, extras ...interface{}) WriterResponse {
	return w.SendResponseWithCode(&map[string]string{"message": message}, code, extras...)
}

// SendError sends an error with message and code, and logs any extras for debugging
func (w *writer) SendError(message string, code int, extras ...interface{}) WriterResponse {
	log.Printf("Http error occured, sending back: '%s', %d", message, code)
	return w.SendResponseWithCode(&map[string]string{"message": message}, code, extras...)
}

// SendUnexpectedError sends a response when an unexpected error is found along with extras
func (w *writer) SendUnexpectedError(err interface{}, extras ...interface{}) WriterResponse {
	return w.SendError(
		"Unexpected Error",
		http.StatusInternalServerError,
		append([]interface{}{err, debug.Stack()}, extras...)...,
	)
}

func (w *writer) printExtras(extras ...interface{}) {
	go func() {
		for i, extra := range extras {
			log.Printf("Extra %d: %+v", i, extra)
		}
	}()
}
