package regex

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// Result is the result structure for a parsed regex
type Result struct {
	inputMap map[string]int
	results  []string
}

// ResultInput is the input for a Result
type ResultInput struct {
	Pattern string
	Key     string
}

// NewResult gets a new regex result
func NewResult(s string, joiner string, input []ResultInput) (*Result, error) {
	trashCounter := 0
	inputMap := map[string]int{}
	matchers := []string{}
	for i, d := range input {
		if d.Key == "" {
			d.Key = fmt.Sprintf("trash%d", trashCounter)
			trashCounter++
		} else {
			inputMap[d.Key] = i + 1
		}
		matchers = append(matchers, fmt.Sprintf(`(?P<%s>%s)`, d.Key, d.Pattern))
	}
	r := regexp.MustCompile(strings.Join(matchers, joiner))
	results := r.FindStringSubmatch(s)
	if len(results) != len(matchers)+1 || len(r.SubexpNames()) != len(matchers)+1 {
		return nil, fmt.Errorf("Wrong number of entries on the line: '%s'. Expected (%d, %d), but got (%d, %d)",
			s,
			len(results),
			len(matchers)+1,
			len(r.SubexpNames()),
			len(matchers)+1,
		)
	}
	return &Result{
		inputMap: inputMap,
		results:  results,
	}, nil
}

// GetString gets a string by key
func (r *Result) GetString(k string) string {
	if key, ok := r.inputMap[k]; ok {
		return r.results[key]
	}
	panic(fmt.Sprintf("Unable to access string result using: '%s'", k))
}

// GetMoneyAmountAsInt gets a money amount as an int by skey
func (r *Result) GetMoneyAmountAsInt(k string, removeLeadingSix bool) (int, error) {
	s := r.GetString(k)
	if removeLeadingSix {
		// 1: to remove first '6' since it's there instead of negative
		s = s[1:]
	}
	// ",", "" to remove commas for thousands
	// ".", "" to remove decimal to parse and store as int
	return strconv.Atoi(strings.ReplaceAll(strings.ReplaceAll(s, ",", ""), ".", ""))
}

// GetInt gets an int by key
func (r *Result) GetInt(k string) (int, error) {
	return strconv.Atoi(r.GetString(k))
}
