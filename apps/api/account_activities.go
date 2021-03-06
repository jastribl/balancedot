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

// GetAccountActivityByUUID gets a single Account Activity by UUID
func (m *App) GetAccountActivityByUUID(w ResponseWriter, r *Request) WriterResponse {
	return m.genericGetByUUID(
		w, r,
		m.db.
			Preload("Account").
			Preload("SplitwiseExpenses.CardActivities").
			Preload("SplitwiseExpenses.AccountActivities"),
		&entities.AccountActivity{},
		r.GetParams()["accountActivityUUID"],
	)
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

func (m *App) getAllAccountActivitiesForSplitwiseExpense(
	splitwiseExpense *entities.SplitwiseExpense,
) ([]*entities.AccountActivity, error) {
	exisitngAccountActivityUUIDs := make([]string, len(splitwiseExpense.AccountActivities))
	for i, accountActivity := range splitwiseExpense.AccountActivities {
		exisitngAccountActivityUUIDs[i] = accountActivity.UUID.String()
	}
	var allAccountActivities []*entities.AccountActivity
	err := m.db.Not(exisitngAccountActivityUUIDs).Where(
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
	).
		Preload("SplitwiseExpenses.CardActivities").
		Preload("SplitwiseExpenses.AccountActivities").
		Find(&allAccountActivities).Error
	if err != nil {
		return nil, err
	}

	return allAccountActivities, nil
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

// UnLinkAccountActivityToSplitwiseExpense links a account activitiy to a splitwise expense
func (m *App) UnLinkAccountActivityToSplitwiseExpense(w ResponseWriter, r *Request) WriterResponse {
	err := m.db.Exec(`
			DELETE FROM
			account_activity_links
			WHERE
				account_activity_uuid = ? AND
				splitwise_expense_uuid = ?
		`,
		r.GetParams()["accountActivityUUID"],
		r.GetParams()["splitwiseExpenseUUID"],
	).Error
	if err != nil {
		return w.SendUnexpectedError(err)
	}

	return w.SendSimpleMessage("success")
}
