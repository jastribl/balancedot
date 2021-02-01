package api

import (
	"net/http"

	"gihub.com/jastribl/balancedot/entities"
	"gihub.com/jastribl/balancedot/helpers"
)

// GetAllAccounts get all the Accounts
func (m *App) GetAllAccounts(w ResponseWriter, r *Request) WriterResponse {
	return m.genericGetAll(w, r, m.db, entities.Account{}, nil)
}

type newAccountParams struct {
	LastFour    string             `json:"last_four"`
	Description string             `json:"description"`
	BankName    entities.BankNames `json:"bank_name"`
}

// CreateNewAccount adds a new Card
func (m *App) CreateNewAccount(w ResponseWriter, r *Request) WriterResponse {
	var p newAccountParams
	err := r.DecodeParams(&p)
	if err != nil {
		return w.SendUnexpectedError(err)
	}
	err = m.db.Create(&entities.Account{
		LastFour:    p.LastFour,
		Description: p.Description,
		BankName:    p.BankName,
	}).Error
	if err != nil {
		if helpers.IsUniqueConstraintError(err, "accounts_last_four_unique") {
			return w.SendError("Account already exists", http.StatusConflict, err)
		}
		return w.SendUnexpectedError(err)
	}

	return w.SendSimpleMessageWithCode("success", http.StatusCreated)
}

// GetAccountByUUID gets a single Account by UUID
func (m *App) GetAccountByUUID(w ResponseWriter, r *Request) WriterResponse {
	return m.genericGetByUUID(w, r, m.db, &entities.Account{}, r.GetParams()["accountUUID"])
}
