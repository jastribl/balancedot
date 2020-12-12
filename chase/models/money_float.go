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

// ToFloat64 returns the raw float value from the MoneyAmount
func (m *MoneyAmount) ToFloat64() float64 {
	return m.float64
}

// MarshalCSV converts the representation to a string
func (m *MoneyAmount) MarshalCSV() (string, error) {
	return fmt.Sprintf("%.2f", m.float64), nil
}

// UnmarshalCSV converts the string to the representation
func (m *MoneyAmount) UnmarshalCSV(csv string) (err error) {
	m.float64, err = strconv.ParseFloat(csv, 64)
	return err
}