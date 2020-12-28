package api

import (
	"net/http"

	"gihub.com/jastribl/balancedot/entities"
	"gihub.com/jastribl/balancedot/helpers"
	"gihub.com/jastribl/balancedot/repos"
)

// GetAllAccounts get all the Accounts
func (m *App) GetAllAccounts(w ResponseWriter, r *Request) WriterResponse {
	accountRepo := repos.NewAccountRepo(m.db)
	accounts, err := accountRepo.GetAllAccounts()
	if err != nil {
		return w.SendUnexpectedError(err)
	}

	return w.SendResponse(accounts)
}

type newAccountParams struct {
	// todo: these
	LastFour    string `json:"last_four"`
	Description string `json:"description"`
}

// CreateNewAccount adds a new Card
func (m *App) CreateNewAccount(w ResponseWriter, r *Request) WriterResponse {
	var p newAccountParams
	err := r.DecodeParams(&p)
	if err != nil {
		return w.SendUnexpectedError(err)
	}
	account := entities.Account{
		LastFour:    p.LastFour,
		Description: p.Description,
	}
	err = m.db.Create(&account).Error
	if err != nil {
		// todo: this
		if helpers.IsUniqueConstraintError(err, "accounts_last_four_unique") {
			return w.SendError("Account already exists", http.StatusConflict, err)
		}
		return w.SendUnexpectedError(err)
	}

	return w.SendResponse(account)
}

// GetAccountByUUID gets a single Account by UUID
func (m *App) GetAccountByUUID(w ResponseWriter, r *Request) WriterResponse {
	params := r.GetParams()
	accountRepo := repos.NewAccountRepo(m.db)
	account, err := accountRepo.GetAccountByUUID(params["accountUUID"])
	if err != nil {
		return w.SendUnexpectedError(err)
	}

	return w.SendResponse(account)
}
