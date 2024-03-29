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

// GetAllUnlinkedSplitwiseExpenses gets all the SplitwiseExpenses that aren't already linked
func (m *App) GetAllUnlinkedSplitwiseExpenses(w ResponseWriter, r *Request) WriterResponse {
	return m.genricRawFindAll(
		w, r,
		m.db.Preload("CardActivities").Preload("AccountActivities"),
		entities.SplitwiseExpense{},
		`
			SELECT e.* 
			FROM splitwise_expenses e 
				LEFT JOIN expense_links el
						ON e.uuid = el.splitwise_expense_uuid
				LEFT JOIN account_activity_links al
						ON e.uuid = al.splitwise_expense_uuid
			WHERE e.splitwise_deleted_at IS NULL 
				AND e.amount_paid > 0 
				AND e.creation_method NOT IN ('venmo', 'payment', 'debt_consolidation')
				AND el.splitwise_expense_uuid IS NULL
				AND al.splitwise_expense_uuid IS NULL
				AND e.date > '2019-08-25'::date -- todo: remove this
		`,
	)
}

// GetSplitwiseExpenseByUUID gets a single SplitwiseExpense by UUID
func (m *App) GetSplitwiseExpenseByUUID(w ResponseWriter, r *Request) WriterResponse {
	return m.genericGetByUUID(
		w, r,
		m.db.
			Preload("CardActivities.SplitwiseExpenses").
			Preload("AccountActivities.SplitwiseExpenses"),
		&entities.SplitwiseExpense{},
		r.GetParams()["splitwiseExpenseUUID"],
	)
}

type splitwiseLinkResponse struct {
	*entities.SplitwiseExpense
	CardActivityLinks    []*entities.CardActivity    `json:"card_activity_links"`
	AccountActivityLinks []*entities.AccountActivity `json:"account_activity_links"`
}

// GetSplitwiseExpenseByUUIDForLinking gets a single SplitwiseExpense by UUID along with all linking info
func (m *App) GetSplitwiseExpenseByUUIDForLinking(w ResponseWriter, r *Request) WriterResponse {
	splitwiseExpenseUUID := r.GetParams()["splitwiseExpenseUUID"]

	splitwiseExpense := &entities.SplitwiseExpense{}
	err := repos.NewGenericRepo(m.db.
		Preload("CardActivities.SplitwiseExpenses").
		Preload("AccountActivities.SplitwiseExpenses"),
	).GetByUUID(splitwiseExpense, splitwiseExpenseUUID)
	if err != nil {
		return w.SendUnexpectedError(err)
	}

	daySpread := r.GetQueryIntDefault("day_spread", 3)
	amountSpreadCents := r.GetQueryIntDefault("amount_spread", 3)

	if daySpread < 0 {
		daySpread = 999999999999999999
	}
	var amountSpread float64
	if amountSpreadCents < 0 {
		amountSpread = 999999999999999999.99
	} else {
		amountSpread = float64(amountSpreadCents) / 100
	}

	allCardActivityLinks, err := m.getAllCardActivitiesForSplitwiseExpense(
		splitwiseExpense,
		daySpread,
		amountSpread,
	)
	if err != nil {
		return w.SendUnexpectedError(err)
	}

	allAccountActivityLinks, err := m.getAllAccountActivitiesForSplitwiseExpense(
		splitwiseExpense,
		daySpread,
		amountSpread,
	)
	if err != nil {
		return w.SendUnexpectedError(err)
	}

	return w.SendResponse(&splitwiseLinkResponse{
		SplitwiseExpense:     splitwiseExpense,
		CardActivityLinks:    allCardActivityLinks,
		AccountActivityLinks: allAccountActivityLinks,
	})
}

// GetRawSplitwiseExpense gets the raw splitwise expense from the API
func (m *App) GetRawSplitwiseExpense(w ResponseWriter, r *Request) WriterResponse {
	splitwiseClient, err := splitwise.NewClientForUser(&m.config.Splitwise)
	if err != nil {
		return w.SendUnexpectedError(err)
	}

	res, err := splitwiseClient.GetRawExpense(r.GetParams()["splitwiseExpenseID"])
	if err != nil {
		return w.SendUnexpectedError(err)
	}

	return w.SendResponse(res)
}

// GetAllSplitwiseExpenses gets all the SplitwiseExpenses
func (m *App) GetAllSplitwiseExpenses(w ResponseWriter, r *Request) WriterResponse {
	return m.genericGetAll(
		w, r,
		m.db.Preload("CardActivities").Preload("AccountActivities"),
		entities.SplitwiseExpense{},
		&repos.GetAllOfOptions{
			Where: "splitwise_deleted_at IS NULL", // Don't load deleted expense
		},
	)
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
				CreationMethod:     expense.CreationMethod,
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
