package api

import (
	"gihub.com/jastribl/balancedot/entities"
	"gihub.com/jastribl/balancedot/helpers"
	"gihub.com/jastribl/balancedot/repos"
	"github.com/gorilla/mux"
)

// GetAllCards get all the Cards
func (m *App) GetAllCards(w ResponseWriter, r Request) interface{} {
	cardRepo := repos.NewCardRepo(m.db)
	cards, err := cardRepo.GetAllCards()
	if err != nil {
		return err
	}

	return w.RenderJSON(cards)
}

type newCardParams struct {
	LastFour    string `json:"last_four"`
	Description string `json:"description"`
}

// CreateNewCard adds a new Card
func (m *App) CreateNewCard(w ResponseWriter, r Request) interface{} {
	var p newCardParams
	m.DecodeParams(r, &p)
	card := entities.Card{
		LastFour:    p.LastFour,
		Description: p.Description,
	}
	err := m.SaveEntity(&card)
	if err != nil {
		if helpers.IsUniqueConstraintError(err, "cards_last_four_unique") {
			return &Error{
				Message: "Card already exists",
				Error:   err,
				Code:    409,
			}
		}
		return err
	}

	return w.RenderJSON(card)
}

// GetCardByUUID gets a single Card by UUID
func (m *App) GetCardByUUID(w ResponseWriter, r Request) interface{} {
	params := mux.Vars(r)
	cardRepo := repos.NewCardRepo(m.db)
	card, err := cardRepo.GetCardByUUID(params["cardUUID"])
	if err != nil {
		return err
	}

	return w.RenderJSON(card)
}
