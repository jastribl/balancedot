package models

import "time"

// RepeatType described the type of repeat on the Expense
type RepeatType string

const (
	// NeverRepeat repeats never
	NeverRepeat RepeatType = "never"
	// WeeklyRepeat repeats ever week
	WeeklyRepeat RepeatType = "weekly"
	// FornightlyRepeat repeats ever 2 weeks
	FornightlyRepeat RepeatType = "fortnightly"
	// MonthlyRepeat repeats ever month
	MonthlyRepeat RepeatType = "monthly"
	// YearlyRepeat repeats ever year
	YearlyRepeat RepeatType = "yearly"
)

// Expense represents the Splitwise Expense object
type Expense struct {
	ID                     int          `json:"id"`
	GroupID                *int         `json:"group_id"`
	Description            string       `json:"description"`
	Repeats                bool         `json:"repeats"`
	RepeatInterval         RepeatType   `json:"repeat_interval"`
	EmailReminder          bool         `json:"email_reminder"`
	EmailReminderInAdvance int          `json:"email_reminder_in_advance"`
	NextRepeat             *interface{} `json:"next_repeat"`
	Details                string       `json:"details"`
	CommentsCount          int          `json:"comments_count"`
	Payment                bool         `json:"payment"`
	CreationMethod         *string      `json:"creation_method"`
	TransactionMethod      *string      `json:"transaction_method"`
	TransactionConfirmed   bool         `json:"transaction_confirmed"`
	TransactionID          *interface{} `json:"transaction_id"`
	Cost                   string       `json:"cost"`
	CurrencyCode           string       `json:"currency_code"`
	Repayments             *[]struct {
		From   int    `json:"from"`
		To     int    `json:"to"`
		Amount string `json:"amount"`
	} `json:"repayments"`
	Date      time.Time  `json:"date"`
	CreatedAt time.Time  `json:"created_at"`
	CreatedBy User       `json:"created_by"`
	UpdatedAt *time.Time `json:"updated_at"`
	UpdatedBy *User      `json:"updated_by"`
	DeletedAt *time.Time `json:"deleted_at"`
	DeletedBy *User      `json:"deleted_by"`
	Category  *struct {
		ID   *int    `json:"id"`
		Name *string `json:"name"`
	} `json:"category"`
	Receipt *struct {
		Large    *interface{} `json:"large"`
		Original *interface{} `json:"original"`
	} `json:"receipt"`
	Users []struct {
		User       User   `json:"user"`
		UserID     int    `json:"user_id"`
		PaidShare  string `json:"paid_share"`
		OwedShare  string `json:"owed_share"`
		NetBalance string `json:"net_balance"`
	} `json:"users"`
}
