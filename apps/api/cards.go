package api

import (
	"net/http"

	"gihub.com/jastribl/balancedot/entities"
	"gihub.com/jastribl/balancedot/helpers"
	"gihub.com/jastribl/balancedot/repos"
)

// GetAllCards get all the Cards
func (m *App) GetAllCards(w ResponseWriter, r *Request) WriterResponse {
	return m.genericGetAll(w, r, entities.Card{}, nil)
}

type newCardParams struct {
	LastFour    string `json:"last_four"`
	Description string `json:"description"`
}

// CreateNewCard adds a new Card
func (m *App) CreateNewCard(w ResponseWriter, r *Request) WriterResponse {
	var p newCardParams
	err := r.DecodeParams(&p)
	if err != nil {
		return w.SendUnexpectedError(err)
	}
	card := entities.Card{
		LastFour:    p.LastFour,
		Description: p.Description,
	}
	err = m.db.Create(&card).Error
	if err != nil {
		if helpers.IsUniqueConstraintError(err, "cards_last_four_unique") {
			return w.SendError("Card already exists", http.StatusConflict, err)
		}
		return w.SendUnexpectedError(err)
	}

	return w.SendResponse(card)
}

// GetCardByUUID gets a single Card by UUID
func (m *App) GetCardByUUID(w ResponseWriter, r *Request) WriterResponse {
	var card entities.Card
	err := repos.NewGenericRepo(m.db).GetByUUID(&card, r.GetParams()["cardUUID"])
	if err != nil {
		return w.SendUnexpectedError(err)
	}

	return w.SendResponse(card)
}
