package api

import (
	"fmt"
	"net/http"
	"strconv"

	"gihub.com/jastribl/balancedot/entities"
	"gihub.com/jastribl/balancedot/helpers"
	"gihub.com/jastribl/balancedot/repos"
	"gihub.com/jastribl/balancedot/splitwise"
	"gihub.com/jastribl/balancedot/splitwise/models"
)

// SplitwiseLoginCheck todo
func (m *App) SplitwiseLoginCheck(w ResponseWriter, r *Request) WriterResponse {
	responseMap := map[string]string{
		"message": "Authentication Response",
	}
	splitwiseConfig := &m.config.Splitwise
	var responseCode int
	if splitwise.HasToken(splitwiseConfig) {
		responseCode = http.StatusOK
	} else {
		responseCode = http.StatusUnauthorized
		responseMap["redirect_url"] = splitwise.GetAuthPortalURL(splitwiseConfig)
	}
	return w.SendResponseWithCode(responseMap, responseCode)

}

type oauthCallbackParams struct {
	Code  string `json:"code"`
	State string `json:"state"`
}

// SplitiwseOatuhCallback todo
func (m *App) SplitiwseOatuhCallback(w ResponseWriter, r *Request) WriterResponse {
	var p oauthCallbackParams
	err := r.DecodeParams(&p)
	if err != nil {
		return w.SendUnexpectedError(err)
	}
	splitwiseConfig := m.config.Splitwise

	if p.State != splitwiseConfig.State {
		return w.SendError("Invalid state in callback", http.StatusInternalServerError)
	}
	token, err := splitwise.GetTokenFromCode(&splitwiseConfig, p.Code)
	if err != nil {
		return w.SendUnexpectedError(err)
	}

	splitwise.SaveToken(&splitwiseConfig, token)

	return w.SendSimpleMessage("success")
}

// GetAllSplitwiseExpenses get all the SplitwiseExpenses
func (m *App) GetAllSplitwiseExpenses(w ResponseWriter, r *Request) WriterResponse {
	splitwiseExpenseRepo := repos.NewSplitwiseExpenseRepo(m.db)
	splitwiseExpenses, err := splitwiseExpenseRepo.GetAllExpensesOrdered()
	if err != nil {
		return w.SendUnexpectedError(err)
	}

	return w.SendResponse(splitwiseExpenses)
}

// RefreshSplitwise refreshes the data from the Splitwise API
func (m *App) RefreshSplitwise(w ResponseWriter, r *Request) WriterResponse {
	splitwiseClient, err := splitwise.NewClientForUser(&m.config.Splitwise)
	if err != nil {
		return w.SendUnexpectedError(err)
	}

	currentUser, err := splitwiseClient.GetCurrentUser()
	if err != nil {
		return w.SendUnexpectedError(err)
	}

	expenses, err := splitwiseClient.GetExpenses()
	if err != nil {
		return w.SendUnexpectedError(err)
	}

	enterExpenseForUser := func(expense *models.Expense, user *models.ExpenseUser) error {
		amountPaid, err := strconv.ParseFloat(user.PaidShare, 64)
		if err != nil {
			return err
		}
		amountOwed, err := strconv.ParseFloat(user.OwedShare, 64)
		if err != nil {
			return err
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
		return m.db.Create(&newSplitwiseExpnese).Error
	}

	for _, expense := range *expenses {
		for _, user := range expense.Users {
			if user.UserID == currentUser.ID {
				err := enterExpenseForUser(expense, user)
				if err != nil {
					if helpers.IsUniqueConstraintError(err, "splitwise_expenses_splitwise_id_key") {
						break
						return w.SendError(
							fmt.Sprintf("Attempted to add duplicate splitwise expense entry: %#v\n", expense.ID),
							http.StatusConflict,
						)
					}
					return w.SendUnexpectedError(err)
				}
				break
			}
		}
	}

	return m.GetAllSplitwiseExpenses(w, r)
}
