package api

import (
	"bufio"
	"fmt"
	"net/http"

	"gihub.com/jastribl/balancedot/chase/models"
	"gihub.com/jastribl/balancedot/entities"
	"gihub.com/jastribl/balancedot/helpers"
	"gihub.com/jastribl/balancedot/repos"
	"github.com/gocarina/gocsv"
	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
)

// GetAllCardActivitiesForCard gets all the Card Activities
func (m *App) GetAllCardActivitiesForCard(w ResponseWriter, r *Request) WriterResponse {
	params := r.GetParams()
	cardActivityRepo := repos.NewCardActivityRepo(m.db)
	cardActivities, err := cardActivityRepo.GetAllCardActivitiesForCard(params["cardUUID"])
	if err != nil {
		return w.SendUnexpectedError(err)
	}

	return w.SendResponse(cardActivities)
}

// UploadCardActivities uploads new CardActivities
func (m *App) UploadCardActivities(w ResponseWriter, r *Request) WriterResponse {
	params := r.GetParams()
	cardUUID, err := uuid.FromString(params["cardUUID"])
	if err != nil {
		return w.SendUnexpectedError(err)
	}

	r.ParseMultipartForm(10 << 20) // 10MB file size limit
	file, handler, err := r.FormFile("file")
	if err != nil {
		return w.SendUnexpectedError(err)
	}
	defer file.Close()
	fmt.Printf("Uploaded File: %+v\n", handler.Filename)
	fmt.Printf("File Size: %+v\n", handler.Size)

	bufferedReader := bufio.NewReader(file)

	cardActivities := []*models.CardActivity{}
	err = gocsv.Unmarshal(bufferedReader, &cardActivities)
	if err != nil {
		return w.SendUnexpectedError(err)
	}

	var newCardActivities []*entities.CardActivity
	success := helpers.NewTransaction(m.db, func(tx *gorm.DB) helpers.TransactionAction {
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
				w.SendUnexpectedError(err)
				return helpers.TransactionActionRollback
			}
			if exists {
				w.SendError("Duplicate activity found", http.StatusConflict)
				return helpers.TransactionActionRollback
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
			err = tx.Create(newCardActivity).Error
			if err != nil {
				w.SendUnexpectedError(err)
				return helpers.TransactionActionRollback
			}
		}
		return helpers.TransactionActionCommit
	})

	// todo: make the return type into the number of records inserted or something
	if success {
		return w.SendResponse(newCardActivities)
	}

	return WriterResponseSuccess
}
