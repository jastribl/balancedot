package api

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"net/http"

	"gihub.com/jastribl/balancedot/chase/models"
	"gihub.com/jastribl/balancedot/entities"
	"gihub.com/jastribl/balancedot/helpers"
	"gihub.com/jastribl/balancedot/repos"
	"github.com/gocarina/gocsv"
	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
)

// GetAllAccountActivitiesForAccount gets all the Account Activities
func (m *App) GetAllAccountActivitiesForAccount(w ResponseWriter, r *Request) WriterResponse {
	params := r.GetParams()
	accountActivityRepo := repos.NewAccountActivityRepo(m.db)
	accountActivities, err := accountActivityRepo.GetAllAccountActivitiesForAccount(params["accountUUID"])
	if err != nil {
		return w.SendUnexpectedError(err)
	}

	return w.SendResponse(accountActivities)
}

// UploadAccountActivities uploads new AccountActivities
func (m *App) UploadAccountActivities(w ResponseWriter, r *Request) WriterResponse {
	params := r.GetParams()
	accountUUID, err := uuid.FromString(params["accountUUID"])
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

	accountActivities := []*models.AccountActivity{}

	gocsv.SetCSVReader(func(in io.Reader) gocsv.CSVReader {
		r := csv.NewReader(in)
		// Ignore wrong number of items in a line (chase account output seems wrong)
		r.FieldsPerRecord = -1
		return r
	})
	err = gocsv.Unmarshal(bufferedReader, &accountActivities)
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
