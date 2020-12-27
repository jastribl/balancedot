package api

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

// Request is a normal http Request but with extras
type Request struct {
	*http.Request
}

// NewRequest constructs and returns a new Request instance given the http.Request
func NewRequest(r *http.Request) *Request {
	return &Request{
		Request: r,
	}
}

// GetParams get the params from the mux request
func (m *Request) GetParams() map[string]string {
	return mux.Vars(m.Request)
}

// DecodeParams decodes the request body into the params structure
func (m *Request) DecodeParams(params interface{}) {
	decoder := json.NewDecoder(m.Request.Body)
	err := decoder.Decode(&params)
	if err != nil {
		panic(err)
	}
}
