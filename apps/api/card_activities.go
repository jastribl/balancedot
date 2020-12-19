package api

import (
	"bufio"
	"fmt"
	"net/http"

	"gihub.com/jastribl/balancedot/entities"
	"gihub.com/jastribl/balancedot/helpers"

	"github.com/gocarina/gocsv"
	"github.com/jinzhu/gorm"

	"gihub.com/jastribl/balancedot/chase/models"
	"gihub.com/jastribl/balancedot/repos"
	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
)

// GetAllCardActivitiesForCard get all the Cards
func (m *App) GetAllCardActivitiesForCard(w ResponseWriter, r *http.Request) interface{} {
	params := mux.Vars(r)
	cardActivityRepo := repos.NewCardActivityRepo(m.db)
	cardActivities, err := cardActivityRepo.GetAllCardActivitiesForCard(params["cardUUID"])
	if err != nil {
		return err
	}

	return w.RenderJSON(cardActivities)
}

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

	var newCardActivities []*entities.CardActivity
	txError := helpers.NewTransaction(m.db, func(tx *gorm.DB) interface{} {
		for _, cardActivity := range cardActivities {
			search := entities.CardActivity{
				CardUUID:        cardUUID,
				TransactionDate: cardActivity.TransactionDate.Time,
				PostDate:        cardActivity.PostDate.Time,
				Description:     cardActivity.Description,
				Type:            cardActivity.Type,
				Amount:          cardActivity.Amount.ToFloat64(),
			}
			// note: use tx instead to check for duplicate inside the same transaction
			exists, err := helpers.RowExists(m.db, &entities.CardActivity{}, search)
			if err != nil {
				return err
			}
			if exists {
				return Error{
					Message: "Duplicate activity found",
					Error:   nil,
					Code:    409,
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
				return err
			}
		}
		return nil
	})

	if txError != nil {
		return txError
	}

	// todo: make the return type into the number of records inserted or something

	return w.RenderJSON(newCardActivities)
}
