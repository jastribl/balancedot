package api

import (
	"net/http"

	"gihub.com/jastribl/balancedot/repos"
	"gihub.com/jastribl/balancedot/splitwise"
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
	splitwiseExpenses, err := splitwiseExpenseRepo.GetAllExpenses()
	if err != nil {
		return w.SendUnexpectedError(err)
	}

	return w.SendResponse(splitwiseExpenses)
}

// RefreshSplitwise refreshes the data from the Splitwise API
func (m *App) RefreshSplitwise(w ResponseWriter, r *Request) WriterResponse {
	splitwiseClient, err := splitwise.NewClient(&m.config.Splitwise)
	if err != nil {
		if e, ok := err.(splitwise.ClientSetupError); ok {
			return w.SendResponseWithCode(map[string]interface{}{
				"message":      "Splitwise Redirect Required",
				"redirect_url": e.RedirectURL,
			}, http.StatusInternalServerError)
		}
		return w.SendUnexpectedError(err)
	}

	currentUser, err := splitwiseClient.GetCurrentUser()
	if err != nil {
		return w.SendUnexpectedError(err)
	}

	return w.SendResponse(currentUser)
}
