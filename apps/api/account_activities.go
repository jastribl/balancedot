package api

import (
	"net/http"

	"gihub.com/jastribl/balancedot/chase/models"
	"gihub.com/jastribl/balancedot/entities"
	"gihub.com/jastribl/balancedot/helpers"
	"gihub.com/jastribl/balancedot/repos"
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

func readChaseAccountActivities(w ResponseWriter, r *Request) ([]models.AccountActivity, error) {
	parsed := []*models.ChaseAccountActivity{}
	if err := r.ReadMultipartCSV("file", &parsed); err != nil {
		return nil, err
	}
	accountActivities := make([]models.AccountActivity, len(parsed))
	for i := range parsed {
		accountActivities[i] = parsed[i]
	}
	return accountActivities, nil
}

func readBofAAccountActivities(w ResponseWriter, r *Request) ([]models.AccountActivity, error) {
	parsed := []*models.BofAAccountActivity{}
	if err := r.ReadMultipartCSV("file", &parsed); err != nil {
		return nil, err
	}

	accountActivities := make([]models.AccountActivity, len(parsed))
	for i := range parsed {
		accountActivities[i] = parsed[i]
	}
	return accountActivities, nil
}

// UploadAccountActivities uploads new AccountActivities
func (m *App) UploadAccountActivities(w ResponseWriter, r *Request) WriterResponse {
	var account entities.Account
	err := repos.NewGenericRepo(m.db).GetByUUID(&account, r.GetParams()["accountUUID"])
	if err != nil {
		return w.SendUnexpectedError(err)
	}

	var readFunction func(w ResponseWriter, r *Request) ([]models.AccountActivity, error)
	switch account.BankName {
	case entities.Chase:
		readFunction = readChaseAccountActivities
	case entities.BofA:
		readFunction = readBofAAccountActivities
	default:
		return w.SendError("Hit unexpected bank name", http.StatusInternalServerError)
	}

	accountActivities, err := readFunction(w, r)
	if err != nil {
		return w.SendUnexpectedError(err)
	}

	newAccountActivities := make([]entities.AccountActivity, len(accountActivities))
	success := helpers.NewTransaction(m.db, func(tx *gorm.DB) helpers.TransactionAction {
		for i, accountActivity := range accountActivities {
			newAccountActivity := accountActivity.ToAccountActivitiyEntity(&account)
			exists, err := helpers.RowExists(tx, &entities.AccountActivity{}, newAccountActivity)
			if err != nil {
				w.SendUnexpectedError(err)
				return helpers.TransactionActionRollback
			}
			if exists {
				w.SendError("Duplicate activity found", http.StatusConflict)
				return helpers.TransactionActionRollback
			}
			newAccountActivities[i] = *newAccountActivity
		}
		return helpers.TransactionActionCommit
	})

	// todo: make the return type into the number of records inserted or something
	if success {
		err = m.db.Create(&newAccountActivities).Error
		if err != nil {
			return w.SendUnexpectedError(err)
		}
		return w.SendResponse(newAccountActivities)
	}

	return WriterResponseSuccess
}
