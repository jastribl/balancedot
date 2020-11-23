package models

// ChequingActivity holds line items from the Chase Chequing Activity Report file
type ChequingActivity struct {
	Details     string      `csv:"Details"`
	PostingDate ChaseDate   `csv:"Posting Date"`
	Description string      `csv:"Description"`
	Amount      MoneyAmount `csv:"Amount"`
	Type        string      `csv:"Type"`
}
