package api

import (
	"net/http"
	"time"

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
	case entities.ChaseBankName:
		readFunction = readChaseAccountActivities
	case entities.BofABankName:
		readFunction = readBofAAccountActivities
	default:
		return w.SendError("Hit unexpected bank name", http.StatusInternalServerError)
	}

	accountActivities, err := readFunction(w, r)
	if err != nil {
		return w.SendUnexpectedError(err)
	}

	newAccountActivities := make([]*entities.AccountActivity, len(accountActivities))
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
			newAccountActivities[i] = newAccountActivity
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

// GetAllAccountActivitiesForSplitwiseExpenseUUID gets all account activities that might be related to a given splitwise expense
func (m *App) GetAllAccountActivitiesForSplitwiseExpenseUUID(w ResponseWriter, r *Request) WriterResponse {
	var splitwiseExpense entities.SplitwiseExpense
	err := repos.NewGenericRepo(m.db).GetByUUID(&splitwiseExpense, r.GetParams()["splitwiseExpenseUUID"])
	if err != nil {
		return w.SendUnexpectedError(err)
	}

	var allAccountActivities []*entities.AccountActivity
	err = m.db.Where(
		`
			(amount = ?) OR
			(-amount >= (?) AND -amount <= (?)) OR
			(-amount >= (?) AND posting_date BETWEEN (?) AND (?))
		`,
		-splitwiseExpense.AmountPaid,
		splitwiseExpense.AmountPaid-0.03,
		splitwiseExpense.AmountPaid+0.03,
		splitwiseExpense.AmountPaid-1,
		splitwiseExpense.Date.Add(-time.Hour*72),
		splitwiseExpense.Date.Add(time.Hour*72),
	).Find(&allAccountActivities).Error
	if err != nil {
		return w.SendUnexpectedError(err)
	}

	return w.SendResponse(allAccountActivities)
}

// LinkAccountActivityToSplitwiseExpense links a account activitiy to a splitwise expense
func (m *App) LinkAccountActivityToSplitwiseExpense(w ResponseWriter, r *Request) WriterResponse {
	err := m.db.Exec(`
			INSERT INTO account_activity_links (
				account_activity_uuid,
				splitwise_expense_uuid
			)
			VALUES (?, ?)
		`,
		r.GetParams()["accountActivityUUID"],
		r.GetParams()["splitwiseExpenseUUID"],
	).Error
	if err != nil {
		return w.SendUnexpectedError(err)
	}

	return w.SendSimpleMessage("success")
}
