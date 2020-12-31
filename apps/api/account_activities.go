package api

import (
	"net/http"

	"gihub.com/jastribl/balancedot/chase/models"
	"gihub.com/jastribl/balancedot/entities"
	"gihub.com/jastribl/balancedot/helpers"
	"gihub.com/jastribl/balancedot/repos"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

// GetAllAccountActivitiesForAccount gets all the Account Activities
func (m *App) GetAllAccountActivitiesForAccount(w ResponseWriter, r *Request) WriterResponse {
	var account entities.Account
	err := repos.NewGenericRepo(m.db.Preload("Activities")).
		GetByUUID(&account, r.GetParams()["accountUUID"])
	if err != nil {
		return w.SendUnexpectedError(err)
	}

	return w.SendResponse(account.Activities)
}

// UploadAccountActivities uploads new AccountActivities
func (m *App) UploadAccountActivities(w ResponseWriter, r *Request) WriterResponse {
	accountUUID, err := uuid.FromString(r.GetParams()["accountUUID"])
	if err != nil {
		return w.SendError("Invalid accountUUID provided", http.StatusUnprocessableEntity, err)
	}

	accountActivities := []*models.AccountActivity{}
	err = r.ReadMultipartCSV("file", &accountActivities)
	if err != nil {
		return w.SendUnexpectedError(err)
	}

	var newAccountActivities []*entities.AccountActivity
	success := helpers.NewTransaction(m.db, func(tx *gorm.DB) helpers.TransactionAction {
		for _, accountActivity := range accountActivities {
			newAccountActivity := &entities.AccountActivity{
				AccountUUID: accountUUID,
				Details:     accountActivity.Details,
				PostingDate: accountActivity.PostingDate.Time,
				Description: accountActivity.Description,
				Amount:      accountActivity.Amount.ToFloat64(),
				Type:        accountActivity.Type,
			}
			// note: use tx instead to check for duplicate inside the same transaction
			exists, err := helpers.RowExists(m.db, &entities.AccountActivity{}, newAccountActivity)
			if err != nil {
				w.SendUnexpectedError(err)
				return helpers.TransactionActionRollback
			}
			if exists {
				w.SendError("Duplicate activity found", http.StatusConflict)
				return helpers.TransactionActionRollback
			}
			newAccountActivities = append(newAccountActivities, newAccountActivity)
			err = tx.Create(newAccountActivity).Error
			if err != nil {
				w.SendUnexpectedError(err)
				return helpers.TransactionActionRollback
			}
		}
		return helpers.TransactionActionCommit
	})

	// todo: make the return type into the number of records inserted or something
	if success {
		return w.SendResponse(newAccountActivities)
	}

	return WriterResponseSuccess
}
