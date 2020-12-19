package api

import (
	"bufio"
	"fmt"
	"net/http"

	"gihub.com/jastribl/balancedot/entities"
	"gihub.com/jastribl/balancedot/helpers"

	"github.com/gocarina/gocsv"

	"gihub.com/jastribl/balancedot/chase/models"
	"gihub.com/jastribl/balancedot/repos"
	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
)

// GetAllCardActivitiesForCard get all the Cards
func (m *App) GetAllCardActivitiesForCard(w ResponseWriter, r *http.Request) interface{} {
	params := mux.Vars(r)
	// todo: verify cardUUID is a param, prob assert
	cardActivityRepo := repos.NewCardActivityRepo(m.db)
	cardActivities, err := cardActivityRepo.GetAllCardActivitiesForCard(params["cardUUID"])
	if err != nil {
		return err
	}

	return w.RenderJSON(cardActivities)
}

// type cardActivitiesParams struct {
// 	LastFour    string `json:"last_four"`
// 	Description string `json:"description"`
// }

// UploadCardActivities uploads new CardActivities
func (m *App) UploadCardActivities(w ResponseWriter, r *http.Request) interface{} {
	params := mux.Vars(r)
	cardUUID, err := uuid.FromString(params["cardUUID"])
	if err != nil {
		return err
	}

	r.ParseMultipartForm(10 << 20) // 10MB file size limit
	file, handler, err := r.FormFile("file")
	if err != nil {
		return err
	}
	defer file.Close()
	fmt.Printf("Uploaded File: %+v\n", handler.Filename)
	fmt.Printf("File Size: %+v\n", handler.Size)

	bufferedReader := bufio.NewReader(file)

	cardActivities := []*models.CardActivity{}
	err = gocsv.Unmarshal(bufferedReader, &cardActivities)
	if err != nil {
		return err
	}

	// todo: make helper function for this that i like the retuns value of instead of db.Transaction
	tx := m.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	var newCardActivities []*entities.CardActivity
	for _, cardActivity := range cardActivities {
		search := entities.CardActivity{
			CardUUID:        cardUUID,
			TransactionDate: cardActivity.TransactionDate.Time,
			PostDate:        cardActivity.PostDate.Time,
			Description:     cardActivity.Description,
			Type:            cardActivity.Type,
			Amount:          cardActivity.Amount.ToFloat64(),
		}
		// todo: use tx instead to check for duplicate inside the same transaction
		exists, err := helpers.RowExists(m.db, &entities.CardActivity{}, search)
		if err != nil {
			tx.Rollback()
			return err
		}
		if exists {
			tx.Rollback()
			return Error{
				Message: "Duplicate activity found", // todo: better messaging
				Error:   nil,
				Code:    409, // todo: is this the right code?
			}
		}
		newCardActivity := &entities.CardActivity{
			CardUUID:        cardUUID,
			TransactionDate: cardActivity.TransactionDate.Time,
			PostDate:        cardActivity.PostDate.Time,
			Description:     cardActivity.Description,
			Category:        cardActivity.Category,
			Type:            cardActivity.Type,
			Amount:          cardActivity.Amount.ToFloat64(),
		}
		newCardActivities = append(newCardActivities, newCardActivity)
		err = tx.Save(newCardActivity).Error
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	err = tx.Commit().Error
	if err != nil {
		return err
	}

	// todo: make the return type into the number of records inserted or something

	return w.RenderJSON(newCardActivities)
}

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
