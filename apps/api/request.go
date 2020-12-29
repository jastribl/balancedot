package api

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gocarina/gocsv"
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
func (m *Request) DecodeParams(params interface{}) error {
	decoder := json.NewDecoder(m.Request.Body)
	return decoder.Decode(&params)
}

// ReadMultipartCSV reads a multipart csv file to a given output variable
func (m *Request) ReadMultipartCSV(fileName string, out interface{}) error {
	m.ParseMultipartForm(10 << 20) // 10MB file size limit
	file, handler, err := m.FormFile(fileName)
	if err != nil {
		return err
	}
	defer file.Close()
	fmt.Printf("Uploaded File: %+v\n", handler.Filename)
	fmt.Printf("File Size: %+v\n", handler.Size)

	bufferedReader := bufio.NewReader(file)

	gocsv.SetCSVReader(func(in io.Reader) gocsv.CSVReader {
		r := csv.NewReader(in)
		// Ignore wrong number of items in a line
		r.FieldsPerRecord = -1
		return r
	})
	return gocsv.Unmarshal(bufferedReader, out)
}
