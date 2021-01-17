package api

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"

	"gihub.com/jastribl/balancedot/chase/models"
	"gihub.com/jastribl/balancedot/entities"
	"gihub.com/jastribl/balancedot/helpers"
	"gihub.com/jastribl/balancedot/lib/regex"
	"gihub.com/jastribl/balancedot/lib/textscanner"
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
	r.ParseMultipartForm(10 << 20) // 10MB file size limit
	parsed := []*models.ChaseCardActivity{}
	for fileName := range r.MultipartForm.File {
		batch := []*models.ChaseCardActivity{}
		if err := r.ReadMultipartCSV(fileName, &batch); err != nil {
			return nil, err
		}
		parsed = append(parsed, batch...)
	}

	cardActivities := make([]models.CardActivity, len(parsed))
	for i := range parsed {
		cardActivities[i] = parsed[i]
	}
	return cardActivities, nil
}

// todo: move this into the banking folder once created
func getBofAActivitiesForFilename(r *Request, fileName string) ([]*models.BofACardActivity, error) {
	parsed := []*models.BofACardActivity{}
	tempFile, err := r.ReceiveFileToTemp(fileName)
	if err != nil {
		return nil, err
	}
	defer func() {
		tempFile.Close()
		os.Remove(tempFile.Name())
	}()
	body, err := exec.Command(
		"pdftotext",
		"-raw",
		"-q",
		"-nopgbrk",
		"-enc",
		"UTF-8",
		"-eol",
		"unix",
		tempFile.Name(),
		"-", // output to stdout
	).Output()
	scanner := textscanner.NewScanner(body)
	var startingYear int
	foundYear, err := scanner.EatToLineContainsWithCallback("JUSTIN A STRIBLING ! Account #", func(s string) error {
		regexResult, err := regex.NewResult(s, "", []regex.ResultInput{
			{Pattern: `JUSTIN A STRIBLING ! Account # \d{4} \d{4} \d{4} \d{4} ! `},
			{Pattern: `\w+`, Key: "start_month"},
			{Pattern: ` `},
			{Pattern: `\d{1,2}`, Key: "start_day"},
			{Pattern: ` - `},
			{Pattern: `\w+`, Key: "end_month"},
			{Pattern: ` `},
			{Pattern: `\d{1,2}`, Key: "end_day"},
			{Pattern: `, `},
			{Pattern: `\d{4}`, Key: "year"},
		})
		if err != nil {
			return err
		}

		startMonth := regexResult.GetString("start_month")
		startDay := regexResult.GetString("start_day")
		endMonth := regexResult.GetString("end_month")
		endDay := regexResult.GetString("end_day")
		year, err := regexResult.GetInt("year")
		if err != nil {
			return err
		}
		startTime, err := time.Parse("January 2 2006", fmt.Sprintf("%s %s %d", startMonth, startDay, year))
		if err != nil {
			return err
		}
		endTime, err := time.Parse("January 2 2006", fmt.Sprintf("%s %s %d", endMonth, endDay, year))
		if err != nil {
			return err
		}
		if startTime.After(endTime) {
			// this wraps, so need to bump start month back
			startTime, err = time.Parse("January 2 2006", fmt.Sprintf("%s %s %d", startMonth, startDay, year-1))
			if err != nil {
				return err
			}
		}
		startingYear = startTime.Year()
		return nil
	})
	if err != nil {
		return nil, err
	}
	if !foundYear {
		return nil, errors.New("Unable to find year in file")
	}

	if scanner.EatToLine("Payments and Other Credits") {
		var lastDate *time.Time
		amounts := 0
		err = scanner.ProcessToAndEatLine("TOTAL PAYMENTS AND OTHER CREDITS FOR THIS PERIOD", func(s string) error {
			regexResult, err := regex.NewResult(s, " ", []regex.ResultInput{
				{Pattern: `\d{2}/\d{2}`, Key: "transaction_date"},
				{Pattern: `\d{2}/\d{2}`, Key: "posting_date"},
				{Pattern: `(?s).*`, Key: "description"},
				{Pattern: `\d{4}`, Key: "reference_number"},
				{Pattern: `\d{4}`, Key: "account_number"},
				{Pattern: `\d+,?\d*?\.\d{2}`, Key: "amount"},
			})
			if err != nil {
				return err
			}

			transactionDateTime, err := time.Parse(
				"2006/01/02",
				fmt.Sprintf("%d/%s", startingYear, regexResult.GetString("transaction_date")),
			)
			if err != nil {
				return err
			}
			postingDateTime, err := time.Parse(
				"2006/01/02",
				fmt.Sprintf("%d/%s", startingYear, regexResult.GetString("posting_date")),
			)
			if err != nil {
				return err
			}

			// Some error handling for statements that go over months
			if lastDate == nil {
				lastDate = &postingDateTime
			} else if lastDate.After(postingDateTime) {
				return errors.New("This situation isn't handled")
			}
			lastDate = &postingDateTime

			moneyAmountInt, err := regexResult.GetMoneyAmountAsInt("amount", true)
			if err != nil {
				return err
			}
			amounts += moneyAmountInt
			newEntry := models.BofACardActivity{
				TransactionDate: models.BofADate{Time: transactionDateTime},
				PostingDate:     models.BofADate{Time: postingDateTime},
				Description:     regexResult.GetString("description"),
				ReferenceNumber: regexResult.GetString("reference_number"),
				AccountNumber:   regexResult.GetString("account_number"),
				Amount:          models.MoneyAmountFromFloat64(float64(moneyAmountInt) / 100.0),
			}
			parsed = append(parsed, &newEntry)
			return nil
		}, func(s string) error {
			regexResult, err := regex.NewResult(s, "", []regex.ResultInput{
				{Pattern: `TOTAL PAYMENTS AND OTHER CREDITS FOR THIS PERIOD 6\$`},
				{Pattern: `\d+,?\d*?\.\d{2}`, Key: "amount"},
			})
			if err != nil {
				return err
			}

			totalAmountInt, err := regexResult.GetMoneyAmountAsInt("amount", false)
			if err != nil {
				return err
			}
			if totalAmountInt != amounts {
				return fmt.Errorf(
					"Amounts don't equal total for TOTAL PAYMENTS AND OTHER CREDITS FOR THIS PERIOD (%d, %d)",
					totalAmountInt,
					amounts,
				)
			}
			return nil
		})
		if err != nil {
			return nil, err
		}
	}

	if scanner.EatToLine("Purchases and Adjustments") {
		var lastDate *time.Time
		prevLine := ""
		amounts := 0
		err = scanner.ProcessToAndEatLine("TOTAL PURCHASES AND ADJUSTMENTS FOR THIS PERIOD", func(s string) error {
			triedWithPrevLine := false

		start:
			regexResult, err := regex.NewResult(s, " ", []regex.ResultInput{
				{Pattern: `\d{2}/\d{2}`, Key: "transaction_date"},
				{Pattern: `\d{2}/\d{2}`, Key: "posting_date"},
				{Pattern: `(?s).*`, Key: "description"},
				{Pattern: `\d{4}`, Key: "reference_number"},
				{Pattern: `\d{4}`, Key: "account_number"},
				{Pattern: `\d+,?\d*?\.\d{2}`, Key: "amount"},
			})
			if err != nil { // Line doesn't match (wrong number of matches)
				if prevLine == "" {
					// Line doesn't work, maybe try using it with the next line
					prevLine = s
					return nil
				}
				if triedWithPrevLine {
					// Already tried with the previous line
					// Set previous line to the combination we already have in s
					prevLine = s
					return nil
				}
				// Add newline and space (regex handles newlines with ?s and space is removed later, but needed for parsing)
				s = prevLine + "\n " + s
				triedWithPrevLine = true
				goto start
			} else if prevLine != "" {
				if triedWithPrevLine {
					// We successfully used the previous line and can clear it now
					prevLine = ""
				} else {
					if strings.HasPrefix(prevLine, "continued on next page...") && strings.HasSuffix(prevLine, "Purchases and Adjustments") {
						// This should be okay - clear this
						prevLine = ""
					} else {
						return fmt.Errorf("Discarding line that might be important: %s", prevLine)
					}
				}
			}
			// We found something, clear our the prevLine
			if prevLine != "" {
				panic("prevLine wasn't empty. Shouldn't be continuing")
			}
			transactionDateTime, err := time.Parse(
				"2006/01/02",
				fmt.Sprintf("%d/%s", startingYear, regexResult.GetString("transaction_date")),
			)
			if err != nil {
				return err
			}
			postingDateTime, err := time.Parse(
				"2006/01/02",
				fmt.Sprintf("%d/%s", startingYear, regexResult.GetString("posting_date")),
			)
			if err != nil {
				return err
			}

			// Some error handling for statements that go over months
			if lastDate == nil {
				lastDate = &postingDateTime
			} else if lastDate.After(postingDateTime) {
				return errors.New("This situation isn't handled")
			}
			lastDate = &postingDateTime

			moneyAmountInt, err := regexResult.GetMoneyAmountAsInt("amount", false)
			if err != nil {
				return err
			}
			amounts += moneyAmountInt
			newEntry := models.BofACardActivity{
				TransactionDate: models.BofADate{Time: transactionDateTime},
				PostingDate:     models.BofADate{Time: postingDateTime},
				Description:     strings.ReplaceAll(regexResult.GetString("description"), "\n ", "\n"),
				ReferenceNumber: regexResult.GetString("reference_number"),
				AccountNumber:   regexResult.GetString("account_number"),
				Amount:          models.MoneyAmountFromFloat64(-float64(moneyAmountInt) / 100.0), // negative since these are payments
			}

			parsed = append(parsed, &newEntry)
			return nil
		}, func(s string) error {
			regexResult, err := regex.NewResult(s, "", []regex.ResultInput{
				{Pattern: `TOTAL PURCHASES AND ADJUSTMENTS FOR THIS PERIOD \$`, Key: "description"},
				{Pattern: `\d+,?\d*?\.\d{2}`, Key: "amount"},
			})
			if err != nil {
				return err
			}

			totalAmountInt, err := regexResult.GetMoneyAmountAsInt("amount", false)
			if err != nil {
				return err
			}
			if totalAmountInt != amounts {
				return fmt.Errorf(
					"Amounts don't equal total for TOTAL PURCHASES AND ADJUSTMENTS FOR THIS PERIOD (%d, %d)",
					totalAmountInt,
					amounts,
				)
			}
			return nil
		})
		if err != nil {
			return nil, err
		}
	}

	return parsed, nil
}

func readBofACardActivities(w ResponseWriter, r *Request) ([]models.CardActivity, error) {
	r.ParseMultipartForm(10 << 20) // 10MB file size limit
	parsed := []*models.BofACardActivity{}
	for fileName := range r.MultipartForm.File {
		batch, err := getBofAActivitiesForFilename(r, fileName)
		if err != nil {
			return nil, err
		}
		parsed = append(parsed, batch...)
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
	case entities.ChaseBankName:
		readFunction = readChaseCardActivities
	case entities.BofABankName:
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
