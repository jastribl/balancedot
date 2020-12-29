package api

import (
	"net/http"

	"gihub.com/jastribl/balancedot/entities"
	"gihub.com/jastribl/balancedot/helpers"
)

// GetAllCards get all the Cards
func (m *App) GetAllCards(w ResponseWriter, r *Request) WriterResponse {
	return m.genericGetAll(w, r, m.db, entities.Card{}, nil)
}

type newCardParams struct {
	LastFour    string             `json:"last_four"`
	Description string             `json:"description"`
	BankName    entities.BankNames `json:"bank_name"`
}

// CreateNewCard adds a new Card
func (m *App) CreateNewCard(w ResponseWriter, r *Request) WriterResponse {
	var p newCardParams
	err := r.DecodeParams(&p)
	if err != nil {
		return w.SendUnexpectedError(err)
	}
	err = m.db.Create(&entities.Card{
		LastFour:    p.LastFour,
		Description: p.Description,
		BankName:    p.BankName,
	}).Error
	if err != nil {
		if helpers.IsUniqueConstraintError(err, "cards_last_four_unique") {
			return w.SendError("Card already exists", http.StatusConflict, err)
		}
		return w.SendUnexpectedError(err)
	}

	return w.SendSimpleMessageWithCode("success", http.StatusCreated)
}

// GetCardByUUID gets a single Card by UUID
func (m *App) GetCardByUUID(w ResponseWriter, r *Request) WriterResponse {
	return m.genericGetByUUID(w, r, m.db, &entities.Card{}, r.GetParams()["cardUUID"])
}
