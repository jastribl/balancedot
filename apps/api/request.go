package api

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"reflect"
	"strings"

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

func getCSVHeaderString(out interface{}) string {
	parts := []string{}
	t := reflect.TypeOf(out).Elem().Elem().Elem()
	for i := 0; i < t.NumField(); i++ {
		parts = append(parts, t.Field(i).Tag.Get("csv"))
	}
	return strings.Join(parts, ",")
}

// ReadMultipartCSV reads a multipart csv file to a given output variable
func (m *Request) ReadMultipartCSV(fileName string, out interface{}) error {
	headerString := getCSVHeaderString(out)
	m.ParseMultipartForm(10 << 20) // 10MB file size limit
	file, _, err := m.FormFile(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	// Move to the start of the header
	bufferedReader := bufio.NewReader(file)
	bytesRead := int64(0)
	for {
		line, _, err := bufferedReader.ReadLine()
		if err != nil {
			if err == io.EOF {
				return errors.New("Found EOF before header")
			}
			return err
		}
		if string(line) == headerString {
			file.Seek(bytesRead, io.SeekStart)
			bufferedReader = bufio.NewReader(file)
			break
		}
		bytesRead += int64(len(line)) + 2 // 2 for new line chars
	}

	gocsv.SetCSVReader(func(in io.Reader) gocsv.CSVReader {
		r := csv.NewReader(in)
		// Ignore wrong number of items in a line
		r.FieldsPerRecord = -1
		return r
	})
	return gocsv.Unmarshal(bufferedReader, out)
}

// ReceiveFileToTemp receives a file to a temp file location
func (m *Request) ReceiveFileToTemp(fileName string) (*os.File, error) {
	m.ParseMultipartForm(10 << 20) // 10MB file size limit
	inputFile, header, err := m.FormFile(fileName)
	if err != nil {
		return nil, err
	}
	defer inputFile.Close()

	tmpFolder := "tmp"
	if _, err = os.Stat(tmpFolder); err != nil {
		if os.IsNotExist(err) {
			os.Mkdir(tmpFolder, 0700)
		} else {
			return nil, err
		}
	}
	outputFile, err := ioutil.TempFile(tmpFolder, header.Filename)
	if err != nil {
		return nil, err
	}

	if _, err = io.Copy(outputFile, inputFile); err != nil {
		return nil, err
	}

	return outputFile, nil
}
