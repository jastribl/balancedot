package models

import "time"

// RegistractionStatus represents the registraction status of the User
type RegistractionStatus string

const (
	// ConfirmedRegistractionStatus is the confirmed status
	ConfirmedRegistractionStatus RegistractionStatus = "confirmed"
	// DummyRegistractionStatus is the dummy status
	DummyRegistractionStatus RegistractionStatus = "dummy"
	// InvitedRegistractionStatus is the invited status
	InvitedRegistractionStatus RegistractionStatus = "invited"
)

// User represents the Splitwise User object
type User struct {
	ID        int    `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Picture   struct {
		Small  *string `json:"small"`
		Medium *string `json:"medium"`
		Large  *string `json:"large"`
	} `json:"picture"`
	CustomPicture      bool                `json:"custom_picture"`
	Email              string              `json:"email"`
	RegistrationStatus RegistractionStatus `json:"registration_status"`
	Locale             string              `json:"locale"`
	DateFormat         string              `json:"date_format"`
	DefaultCurrency    string              `json:"default_currency"`
	DefaultGroupID     int                 `json:"default_group_id"`
	NotificationsRead  time.Time           `json:"notifications_read"`
	NotificationsCount int                 `json:"notifications_count"`
	Notifications      struct {
		AddedAsFriend  bool `json:"added_as_friend"`
		AddedToGroup   bool `json:"added_to_group"`
		ExpenseAdded   bool `json:"expense_added"`
		ExpenseUpdated bool `json:"expense_updated"`
		Bills          bool `json:"bills"`
		Payments       bool `json:"payments"`
		MonthlySummary bool `json:"monthly_summary"`
		Announcements  bool `json:"announcements"`
	} `json:"notifications"`
}
