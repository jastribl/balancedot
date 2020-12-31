package entities

// BankNames is the enum type for bank_names
type BankNames string

// todo: rename to ChaseBankName and BofABankName
const (
	// Chase is chase
	Chase BankNames = "chase"

	// BofA if BofA
	BofA BankNames = "bofa"
)
