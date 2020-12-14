package api

import (
	"encoding/json"
	"net/http"

	"gihub.com/jastribl/balancedot/entities"
	"gihub.com/jastribl/balancedot/repos"
)

// GetAllCards get all the Cards
func (m *App) GetAllCards(w http.ResponseWriter, r *http.Request) (e *Error) {
	cardRepo := repos.NewCardRepo(m.db)
	cards, err := cardRepo.GetAllCards()
	if err != nil {
		return &Error{
			Message: err.Error(), // todo: message
			Error:   err,
			Code:    500, // todo: better code
		}
	}
	err = json.NewEncoder(w).Encode(cards)
	if err != nil {
		return &Error{
			Message: err.Error(), // todo: better message
			Error:   err,
			Code:    500, // todo: better code
		}
	}
	return
}

type newCardParams struct {
	LastFour    string `json:"last_four"`
	Description string `json:"description"`
}

// CreateNewCard adds a new Card
func (m *App) CreateNewCard(w http.ResponseWriter, r *http.Request) (e *Error) {
	var p newCardParams
	m.DecodeParams(r, &p)
	card := entities.Card{
		LastFour:    p.LastFour,
		Description: p.Description,
	}
	err := m.SaveEntity(&card)
	if err != nil {
		return &Error{
			Message: err.Error(), // todo: better message
			Error:   err,
			Code:    500, // todo: better
		}
	}

	return m.GetAllCards(w, r)
}
