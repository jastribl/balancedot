package models

// AccountActivity holds line items from the Chase Chequing Account Activity Report file
type AccountActivity struct {
	Details           string      `csv:"Details"`
	PostingDate       ChaseDate   `csv:"Posting Date"`
	Description       string      `csv:"Description"`
	Amount            MoneyAmount `csv:"Amount"`
	Type              string      `csv:"Type"`
	Balance           MoneyAmount `csv:"Balance"`
	CheckOrSlipNumber string      `csv:"Check or Slip #"`
}
