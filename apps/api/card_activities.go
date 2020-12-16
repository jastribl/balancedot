package api

import (
	"gihub.com/jastribl/balancedot/repos"
	"github.com/gorilla/mux"
)

// GetAllCardActivitiesForCard get all the Cards
func (m *App) GetAllCardActivitiesForCard(w ResponseWriter, r Request) interface{} {
	params := mux.Vars(r)
	cardActivityRepo := repos.NewCardActivityRepo(m.db)
	cardActivities, err := cardActivityRepo.GetAllCardActivitiesForCard(params["cardUUID"])
	if err != nil {
		return err
	}

	return w.RenderJSON(cardActivities)
}

// type newCardParams struct {
// 	LastFour    string `json:"last_four"`
// 	Description string `json:"description"`
// }

// // CreateNewCard adds a new Card
// func (m *App) CreateNewCard(w ResponseWriter, r Request) interface{} {
// 	var p newCardParams
// 	m.DecodeParams(r, &p)
// 	card := entities.Card{
// 		LastFour:    p.LastFour,
// 		Description: p.Description,
// 	}
// 	err := m.SaveEntity(&card)
// 	if err != nil {
// 		if helpers.IsUniqueConstraintError(err, "cards_last_four_unique") {
// 			return &Error{
// 				Message: "Card already exists",
// 				Error:   err,
// 				Code:    409,
// 			}
// 		}
// 		return err
// 	}

// 	return w.RenderJSON(card)
// }
