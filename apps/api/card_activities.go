package api

import (
	"net/http"

	"gihub.com/jastribl/balancedot/chase/models"
	"gihub.com/jastribl/balancedot/entities"
	"gihub.com/jastribl/balancedot/helpers"
	"gihub.com/jastribl/balancedot/repos"
	"gorm.io/gorm"
)

// GetAllCardActivitiesForCard gets all the Card Activities
func (m *App) GetAllCardActivitiesForCard(w ResponseWriter, r *Request) WriterResponse {
	var card entities.Card
	err := repos.NewGenericRepo(m.db.Preload("Activities")).
		GetByUUID(&card, r.GetParams()["cardUUID"])
	if err != nil {
		return w.SendUnexpectedError(err)
	}
	return w.SendResponse(card.Activities)
}

func readChaseCardActivities(w ResponseWriter, r *Request) ([]models.CardActivity, error) {
	parsed := []*models.ChaseCardActivity{}
	if err := r.ReadMultipartCSV("file", &parsed); err != nil {
		return nil, err
	}
	cardActivities := make([]models.CardActivity, len(parsed))
	for i := range parsed {
		cardActivities[i] = parsed[i]
	}
	return cardActivities, nil
}

func readBofACardActivities(w ResponseWriter, r *Request) ([]models.CardActivity, error) {
	parsed := []*models.BofACardActivity{}
	if err := r.ReadMultipartCSV("file", &parsed); err != nil {
		return nil, err
	}

	cardActivities := make([]models.CardActivity, len(parsed))
	for i := range parsed {
		cardActivities[i] = parsed[i]
	}
	return cardActivities, nil
}

// UploadCardActivities uploads new CardActivities
func (m *App) UploadCardActivities(w ResponseWriter, r *Request) WriterResponse {
	var card entities.Card
	err := repos.NewGenericRepo(m.db).GetByUUID(&card, r.GetParams()["cardUUID"])
	if err != nil {
		return w.SendUnexpectedError(err)
	}

	var readFunction func(w ResponseWriter, r *Request) ([]models.CardActivity, error)
	switch card.BankName {
	case entities.Chase:
		readFunction = readChaseCardActivities
	case entities.BofA:
		readFunction = readBofACardActivities
	default:
		return w.SendError("Hit unexpected bank name", http.StatusInternalServerError)
	}

	cardActivities, err := readFunction(w, r)
	if err != nil {
		return w.SendUnexpectedError(err)
	}

	newCardActivities := make([]*entities.CardActivity, len(cardActivities))
	success := helpers.NewTransaction(m.db, func(tx *gorm.DB) helpers.TransactionAction {
		for i, cardActivity := range cardActivities {
			search := cardActivity.ToCardActivitiyUniqueMatcher(&card)
			exists, err := helpers.RowExists(tx, &entities.CardActivity{}, search)
			if err != nil {
				w.SendUnexpectedError(err)
				return helpers.TransactionActionRollback
			}
			if exists {
				w.SendError("Duplicate activity found", http.StatusConflict)
				return helpers.TransactionActionRollback
			}
			newCardActivities[i] = cardActivity.ToCardActivitiyEntity(&card)
		}
		return helpers.TransactionActionCommit
	})

	// todo: make the return type into the number of records inserted or something
	if success {
		err = m.db.Create(&newCardActivities).Error
		if err != nil {
			return w.SendUnexpectedError(err)
		}
		return w.SendResponse(newCardActivities)
	}

	return WriterResponseSuccess
}
