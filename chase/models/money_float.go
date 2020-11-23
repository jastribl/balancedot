package models

import (
	"fmt"
	"strconv"
)

// MoneyAmount represents the money object
// including marshalling and unmarshalling logic.
type MoneyAmount struct {
	float64
}

// MarshalCSV converts the representation to a string
func (amount *MoneyAmount) MarshalCSV() (string, error) {
	return fmt.Sprintf("%.2f", amount.float64), nil
}

// UnmarshalCSV converts the string to the representation
func (amount *MoneyAmount) UnmarshalCSV(csv string) (err error) {
	amount.float64, err = strconv.ParseFloat(csv, 64)
	return err
}
