package api

import (
	"encoding/json"
	"log"
	"net/http"
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
	SendResponse(i interface{}) WriterResponse
	SendResponseWithCode(i interface{}, code int) WriterResponse
	SendSimpleMessage(message string) WriterResponse
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
func (w *writer) SendResponseWithCode(i interface{}, code int) WriterResponse {
	w.WriteHeader(code)
	err := json.NewEncoder(w).Encode(i)
	if err != nil {
		http.Error(w, "Serialization error", http.StatusInternalServerError)
	}
	return WriterResponseSuccess
}

// SendResponse sends arbitrary data as JSON
func (w *writer) SendResponse(i interface{}) WriterResponse {
	return w.SendResponseWithCode(i, http.StatusOK)
}

// SendSimpleMessage sends a simple "message" response with http.StatusOK
func (w *writer) SendSimpleMessage(message string) WriterResponse {
	return w.SendResponseWithCode(&map[string]string{"message": message}, http.StatusOK)
}

// SendError sends an error with message and code, and logs any extras for debugging
func (w *writer) SendError(message string, code int, extras ...interface{}) WriterResponse {
	log.Printf("Http error occured, sending back: '%s', %d", message, code)
	for i, extra := range extras {
		log.Printf("Extra %d: %+v", i, extra)
	}
	return w.SendResponseWithCode(&map[string]string{"message": message}, code)
}

// SendUnexpectedError sends a response when an unexpected error is found along with extras
func (w *writer) SendUnexpectedError(err interface{}, extras ...interface{}) WriterResponse {
	return w.SendError("Unexpected Error", http.StatusInternalServerError, err, extras)
}
