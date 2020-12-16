package api

import (
	"gihub.com/jastribl/balancedot/entities"
	"gihub.com/jastribl/balancedot/helpers"
	"gihub.com/jastribl/balancedot/repos"
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
