package models

import "time"

// ChaseDate represents the date used by Chase in their activity logs
// including marshalling and unmarshalling logic.
type ChaseDate struct {
	time.Time
}

// MarshalCSV converts the internal date as CSV string
func (date *ChaseDate) MarshalCSV() (string, error) {
	return date.Time.Format("01/02/2006"), nil
}

// UnmarshalCSV converts the CSV string as internal date
func (date *ChaseDate) UnmarshalCSV(csv string) (err error) {
	date.Time, err = time.Parse("01/02/2006", csv)
	return err
}
