package models

// CardActivity holds line items from the Chase Card Activity Report file
type CardActivity struct {
	TransactionDate ChaseDate   `csv:"Transaction Date"`
	PostDate        ChaseDate   `csv:"Post Date"`
	Description     string      `csv:"Description"`
	Category        string      `csv:"Category"`
	Type            string      `csv:"Type"`
	Amount          MoneyAmount `csv:"Amount"`
	Memo            string      `csv:"Memo"`
}
