package api

import (
	"net/http"
	"strconv"

	"gihub.com/jastribl/balancedot/entities"
	"gihub.com/jastribl/balancedot/helpers"
	"gihub.com/jastribl/balancedot/repos"
	"gihub.com/jastribl/balancedot/splitwise"
)

// SplitwiseLoginCheck checks to see if the current user has an authorized splitwise token
func (m *App) SplitwiseLoginCheck(w ResponseWriter, r *Request) WriterResponse {
	responseMap := map[string]string{
		"message": "Authentication Required",
	}
	splitwiseConfig := &m.config.Splitwise
	var responseCode int
	if splitwise.HasToken(splitwiseConfig) {
		responseCode = http.StatusOK
		responseMap["message"] = "success"
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

// SplitwiseOauthCallback takes in the code and state and configures the token
func (m *App) SplitwiseOauthCallback(w ResponseWriter, r *Request) WriterResponse {
	// todo: make work with user
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
	return m.genericGetAll(w, r, entities.SplitwiseExpense{}, &repos.GetAllOfOptions{
		Order: "date DESC",
	})
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

	expensesAdded := []*entities.SplitwiseExpense{}
	expensesUpdated := []*entities.SplitwiseExpense{}

	for _, expense := range *expenses {
		for _, user := range expense.Users {
			if user.UserID != currentUser.ID {
				continue
			}
			amountPaid, err := strconv.ParseFloat(user.PaidShare, 64)
			if err != nil {
				return w.SendUnexpectedError(err)
			}
			amountOwed, err := strconv.ParseFloat(user.OwedShare, 64)
			if err != nil {
				return w.SendUnexpectedError(err)
			}
			newSplitwiseExpense := &entities.SplitwiseExpense{
				SplitwiseID:        expense.ID,
				Description:        expense.Description,
				Details:            expense.Details,
				CurrencyCode:       expense.CurrencyCode,
				Amount:             amountOwed,
				AmountPaid:         amountPaid,
				Date:               expense.Date,
				SplitwiseCreatedAt: expense.CreatedAt,
				SplitwiseUpdatedAt: expense.UpdatedAt,
				SplitwiseDeletedAt: expense.DeletedAt,
				Category:           *expense.Category.Name,
			}

			idExists, err := helpers.RowExistsIncludingDeleted(
				m.db,
				&entities.SplitwiseExpense{},
				entities.SplitwiseExpense{
					SplitwiseID: expense.ID,
				},
			)
			if err != nil {
				return w.SendUnexpectedError(err)
			}
			if idExists {
				nothingChanged, err := helpers.RowExistsIncludingDeleted(
					m.db,
					&entities.SplitwiseExpense{},
					newSplitwiseExpense,
				)
				if err != nil {
					return w.SendUnexpectedError(err)
				}
				if nothingChanged {
					break
				}
				expensesUpdated = append(expensesUpdated, newSplitwiseExpense)
				err = m.db.
					Model(&entities.SplitwiseExpense{}).
					Where("splitwise_id = ?", newSplitwiseExpense.SplitwiseID).
					Updates(newSplitwiseExpense).
					Error
			} else {
				expensesAdded = append(expensesAdded, newSplitwiseExpense)
				err = m.db.Create(newSplitwiseExpense).Error
			}
			if err != nil {
				if helpers.IsUniqueConstraintError(err, "splitwise_expenses_splitwise_id_unique") {
					return w.SendError(
						"got duplicate even though this shouldn't be possible",
						http.StatusInternalServerError,
						newSplitwiseExpense,
					)
				}
				return w.SendUnexpectedError(err)
			}
			break
		}
	}

	return w.SendResponse(map[string][]*entities.SplitwiseExpense{
		"expenses_added":   expensesAdded,
		"expenses_updated": expensesUpdated,
	})
}
