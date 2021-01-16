package api

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"

	"gihub.com/jastribl/balancedot/chase/models"
	"gihub.com/jastribl/balancedot/entities"
	"gihub.com/jastribl/balancedot/helpers"
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
		matchers := []string{
			`(?P<pre>JUSTIN A STRIBLING ! Account # \d{4} \d{4} \d{4} \d{4} ! )`,
			`(?P<start_month>\w+)`, // 2
			`(?P<trash1> )`,
			`(?P<start_day>\d{1,2})`,
			`(?P<trash2> - )`,
			`(?P<end_month>\w+)`,
			`(?P<trash3> )`,
			`(?P<end_day>\d{1,2})`,
			`(?P<trash4>, )`,
			`(?P<year>\d{4})`,
		}
		r := regexp.MustCompile(strings.Join(matchers, ""))
		results := r.FindStringSubmatch(s)
		if len(results) != len(matchers)+1 || len(r.SubexpNames()) != len(matchers)+1 {
			return fmt.Errorf("Wrong number of entries when looking for year on the line: '%s'. Expected (%d, %d), but got (%d, %d)",
				s,
				len(results),
				len(matchers)+1,
				len(r.SubexpNames()),
				len(matchers)+1,
			)
		}

		startMonth := results[2]
		startDay := results[4]
		endMonth := results[6]
		endDay := results[8]
		year, err := strconv.Atoi(results[10])
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
			matchers := []string{
				`(?P<transaction_date>\d{2}/\d{2})`,
				`(?P<posting_date>\d{2}/\d{2})`,
				`(?P<description>(?s).*)`,
				`(?P<reference_number>\d{4})`,
				`(?P<account_number>\d{4})`,
				`(?P<amount>\d+,?\d*?\.\d{2})`,
			}
			r := regexp.MustCompile(strings.Join(matchers, " "))
			results := r.FindStringSubmatch(s)
			if len(results) != len(matchers)+1 || len(r.SubexpNames()) != len(matchers)+1 {
				return fmt.Errorf("Wrong number of entries when looking for payments on the line: '%s'. Expected (%d, %d), but got (%d, %d)",
					s,
					len(results),
					len(matchers)+1,
					len(r.SubexpNames()),
					len(matchers)+1,
				)
			}
			transactionDateTime, err := time.Parse("2006/01/02", fmt.Sprintf("%d/%s", startingYear, results[1]))
			if err != nil {
				return err
			}
			postingDateTime, err := time.Parse("2006/01/02", fmt.Sprintf("%d/%s", startingYear, results[2]))
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

			// 1: to remove first '6' since it's there instead of negative
			// ",", "" to remove commas for thousands
			// ".", "" to remove decimal to parse and store as int
			moneyAmountInt, err := strconv.Atoi(strings.ReplaceAll(strings.ReplaceAll(results[6][1:], ",", ""), ".", ""))
			if err != nil {
				return err
			}
			amounts += moneyAmountInt
			newEntry := models.BofACardActivity{
				TransactionDate: models.BofADate{Time: transactionDateTime},
				PostingDate:     models.BofADate{Time: postingDateTime},
				Description:     results[3],
				ReferenceNumber: results[4],
				AccountNumber:   results[5],
				Amount:          models.MoneyAmountFromFloat64(float64(moneyAmountInt) / 100.0),
			}
			parsed = append(parsed, &newEntry)
			return nil
		}, func(s string) error {
			matchers := []string{
				`(?P<description>TOTAL PAYMENTS AND OTHER CREDITS FOR THIS PERIOD 6\$)`,
				`(?P<amount>\d+,?\d*?\.\d{2})`,
			}
			r := regexp.MustCompile(strings.Join(matchers, ""))
			results := r.FindStringSubmatch(s)
			if len(results) != len(matchers)+1 || len(r.SubexpNames()) != len(matchers)+1 {
				return fmt.Errorf("Wrong number of entries when looking for payments on the line: '%s'. Expected (%d, %d), but got (%d, %d)",
					s,
					len(results),
					len(matchers)+1,
					len(r.SubexpNames()),
					len(matchers)+1,
				)
			}
			// ",", "" to remove commas for thousands
			// ".", "" to remove decimal to parse and store as int
			totalAmountInt, err := strconv.Atoi(strings.ReplaceAll(strings.ReplaceAll(results[2], ",", ""), ".", ""))
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
			matchers := []string{
				`(?P<transaction_date>\d{2}/\d{2})`,
				`(?P<posting_date>\d{2}/\d{2})`,
				`(?P<description>(?s).*)`,
				`(?P<reference_number>\d{4})`,
				`(?P<account_number>\d{4})`,
				`(?P<amount>\d+,?\d*?\.\d{2})`,
			}
			r := regexp.MustCompile(strings.Join(matchers, " "))

		start:
			results := r.FindStringSubmatch(s)
			if len(results) != len(matchers)+1 || len(r.SubexpNames()) != len(matchers)+1 {
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
			transactionDateTime, err := time.Parse("2006/01/02", fmt.Sprintf("%d/%s", startingYear, results[1]))
			if err != nil {
				return err
			}
			postingDateTime, err := time.Parse("2006/01/02", fmt.Sprintf("%d/%s", startingYear, results[2]))
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
			// ",", "" to remove commas for thousands
			// ".", "" to remove decimal to parse and store as int
			moneyAmountInt, err := strconv.Atoi(strings.ReplaceAll(strings.ReplaceAll(results[6], ",", ""), ".", ""))
			if err != nil {
				return err
			}
			amounts += moneyAmountInt
			newEntry := models.BofACardActivity{
				TransactionDate: models.BofADate{Time: transactionDateTime},
				PostingDate:     models.BofADate{Time: postingDateTime},
				Description:     strings.ReplaceAll(results[3], "\n ", "\n"),
				ReferenceNumber: results[4],
				AccountNumber:   results[5],
				Amount:          models.MoneyAmountFromFloat64(-float64(moneyAmountInt) / 100.0), // negative since these are payments
			}

			parsed = append(parsed, &newEntry)
			return nil
		}, func(s string) error {
			matchers := []string{
				`(?P<description>TOTAL PURCHASES AND ADJUSTMENTS FOR THIS PERIOD \$)`,
				`(?P<amount>\d+,?\d*?\.\d{2})`,
			}
			r := regexp.MustCompile(strings.Join(matchers, ""))
			results := r.FindStringSubmatch(s)
			if len(results) != len(matchers)+1 || len(r.SubexpNames()) != len(matchers)+1 {
				return fmt.Errorf("Wrong number of entries when looking for payments on the line: '%s'. Expected (%d, %d), but got (%d, %d)",
					s,
					len(results),
					len(matchers)+1,
					len(r.SubexpNames()),
					len(matchers)+1,
				)
			}
			// ",", "" to remove commas for thousands
			// ".", "" to remove decimal to parse and store as int
			totalAmountInt, err := strconv.Atoi(strings.ReplaceAll(strings.ReplaceAll(results[2], ",", ""), ".", ""))
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
