package main

import (
	"io"
	"log"
	"net/http"
	"strconv"

	"gihub.com/jastribl/balancedot/entities"
	"gihub.com/jastribl/balancedot/helpers"
	"gihub.com/jastribl/balancedot/splitwise/models"

	"gihub.com/jastribl/balancedot/config"
	"gihub.com/jastribl/balancedot/splitwise"
	"github.com/pkg/browser"
)

var cfg *config.Config

func oauthCallback(w http.ResponseWriter, r *http.Request) {
	if r.URL.Query().Get("state") != cfg.Splitwise.State {
		io.WriteString(w, "Invalid state in callback")
		return
	}

	cfg.Splitwise.AuthCodeChan <- r.URL.Query().Get("code")
	io.WriteString(w, "Processing Token...\nPlease close this tab and return to your application.")
}

func main() {
	cfg = config.NewConfig()
	db, err := helpers.DbConnect()
	if err != nil {
		log.Panic(err)
	}

	go func() {
		browser.OpenURL(<-cfg.Splitwise.AuthURLChan)
	}()

	m := http.NewServeMux()
	m.HandleFunc("/oauth_callback", oauthCallback)
	s := http.Server{Addr: cfg.ServerURL, Handler: m}
	defer s.Close()

	// Run the server
	go func() {
		log.Fatal(s.ListenAndServe())
	}()

	splitwiseClient, err := splitwise.NewClientForCLI(&cfg.Splitwise)
	if err != nil {
		log.Fatal(err)
	}

	currentUser, err := splitwiseClient.GetCurrentUser()
	if err != nil {
		log.Fatal(err)
	}

	expenses, err := splitwiseClient.GetExpenses()
	if err != nil {
		log.Fatal(err)
	}

	// jsonEncoder := json.NewEncoder(log.Writer())
	// jsonEncoder.Encode(expenses)

	enterExpenseForUser := func(expense *models.Expense, user *models.ExpenseUser) {
		amountPaid, err := strconv.ParseFloat(user.PaidShare, 64)
		if err != nil {
			log.Fatal(err)
		}
		amountOwed, err := strconv.ParseFloat(user.OwedShare, 64)
		if err != nil {
			log.Fatal(err)
		}
		newSplitwiseExpnese := entities.SplitwiseExpense{
			SplitwiseID:  expense.ID,
			Description:  expense.Description,
			Details:      expense.Details,
			CurrencyCode: expense.CurrencyCode,
			Amount:       amountOwed,
			AmountPaid:   amountPaid,
			Date:         expense.Date,
			CreatedAt:    expense.CreatedAt,
			UpdatedAt:    expense.UpdatedAt,
			DeletedAt:    expense.DeletedAt,
			Category:     *expense.Category.Name,
		}
		// TODO: update entry on duplicate somehow
		err = db.Create(&newSplitwiseExpnese).Error
		if err != nil {
			if helpers.IsUniqueConstraintError(err, "splitwise_expenses_splitwise_id_key") {
				log.Printf("Attempted to add duplicate splitwise expense entry: %#v\n", newSplitwiseExpnese.SplitwiseID)
			}
		}
	}

	for _, expense := range *expenses {
		for _, user := range expense.Users {
			if user.UserID == currentUser.ID {
				enterExpenseForUser(expense, user)
				break
			}
		}
	}
}
