package models

import "time"

// User represents the Splitwise User object
type User struct {
	ID        int    `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Picture   struct {
		Small  string `json:"small"`
		Medium string `json:"medium"`
		Large  string `json:"large"`
	} `json:"picture"`
	CustomPicture      bool        `json:"custom_picture"`
	Email              string      `json:"email"`
	RegistrationStatus string      `json:"registration_status"`
	ForceRefreshAt     interface{} `json:"force_refresh_at"`
	Locale             string      `json:"locale"`
	CountryCode        string      `json:"country_code"`
	DateFormat         string      `json:"date_format"`
	DefaultCurrency    string      `json:"default_currency"`
	DefaultGroupID     int         `json:"default_group_id"`
	NotificationsRead  time.Time   `json:"notifications_read"`
	NotificationsCount int         `json:"notifications_count"`
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
