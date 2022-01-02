package api

import (
	"database/sql"
	"log"
	"net/http"
	"regexp"
	"strings"
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
			Preload("CardActivites").
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
			newAccountActivity := accountActivity.ToAccountActivityEntity(&account)
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

func (m *App) AutoLinkAccountToCardActivities(w ResponseWriter, r *Request) WriterResponse {
	var account entities.Account
	err := repos.NewGenericRepo(
		m.db.Preload("Activities.CardActivites"),
	).GetByUUID(&account, r.GetParams()["accountUUID"])
	if err != nil {
		return w.SendUnexpectedError(err)
	}
	var cards []*entities.Card
	if m.db.Preload("Activities").Find(&cards).Error != nil {
		return w.SendUnexpectedError(err)
	}
	bofACards, chaseCards := []*entities.Card{}, []*entities.Card{}
	last4ToCard := map[string]*entities.Card{}
	for _, card := range cards {
		last4ToCard[card.LastFour] = card
		switch card.BankName {
		case entities.BofABankName:
			bofACards = append(bofACards, card)
		case entities.ChaseBankName:
			chaseCards = append(chaseCards, card)
		default:
			log.Fatal("Found card with unexpected bank name")
		}
	}

	if account.BankName == entities.ChaseBankName {
		matcher1 := regexp.MustCompile(`Payment to Chase card ending in [0-9]{4}`)
	ACCOUNT_ACTIVITY_LOOP1:
		for _, accountActivity := range account.Activities {
			if len(accountActivity.CardActivites) == 1 {
				continue
			}
			accountActivityDescription := accountActivity.Description
			if strings.Contains(accountActivityDescription, "Payment to Chase card ending in") {
				// Payment from Chase account to Chase card
				res := matcher1.FindString(accountActivityDescription)
				if card, ok := last4ToCard[res[len(res)-4:]]; ok {
					_, err := m.linkChaseAccountActivityToChaseCard(&accountActivity, card)
					if err != nil {
						return w.SendUnexpectedError(err)
					}
					continue ACCOUNT_ACTIVITY_LOOP1
				}
			} else if strings.Contains(accountActivityDescription, "CHASE CREDIT CRD EPAY ONUS") {
				// Payment from Chase account to Chase card (special case)
				for _, card := range chaseCards {
					linked, err := m.linkChaseAccountActivityToChaseCard(&accountActivity, card)
					if err != nil {
						return w.SendUnexpectedError(err)
					}
					if linked {
						continue ACCOUNT_ACTIVITY_LOOP1
					}
				}
			}
		}
	} else if account.BankName == entities.BofABankName {
		matcher2 := regexp.MustCompile(`Online Banking payment to CRD [0-9]{4}`)
	ACCOUNT_ACTIVITY_LOOP2:
		for _, accountActivity := range account.Activities {
			if len(accountActivity.CardActivites) == 1 {
				continue
			}
			accountActivityDescription := accountActivity.Description
			if strings.Contains(accountActivityDescription, "Online Banking payment to CRD") {
				// Payment from BofA account to Bofa Card
				res := matcher2.FindString(accountActivityDescription)
				if card, ok := last4ToCard[res[len(res)-4:]]; ok {
					_, err := m.linkBofAAccountActivityToBofaCard(&accountActivity, card)
					if err != nil {
						return w.SendUnexpectedError(err)
					}
					continue ACCOUNT_ACTIVITY_LOOP2
				}
			} else if strings.Contains(accountActivityDescription, "CHASE CREDIT CRD DES:EPAY") {
				// Payment from BofA account to Chase Card
				for _, card := range chaseCards {
					linked, err := m.linkBofAAccountActivityToChaseCard(&accountActivity, card)
					if err != nil {
						return w.SendUnexpectedError(err)
					}
					if linked {
						continue ACCOUNT_ACTIVITY_LOOP2
					}
				}
			}
		}
	}

	return w.SendSimpleMessage("success")
}

func (m *App) linkChaseAccountActivityToChaseCard(
	accountActivity *entities.AccountActivity,
	card *entities.Card,
) (bool, error) {
	return m.linkAccountActivityToCard("Payment Thank You", 2, accountActivity, card)
}

func (m *App) linkBofAAccountActivityToChaseCard(
	accountActivity *entities.AccountActivity,
	card *entities.Card,
) (bool, error) {
	return m.linkAccountActivityToCard("Payment Thank You", 4, accountActivity, card)
}

func (m *App) linkBofAAccountActivityToBofaCard(
	accountActivity *entities.AccountActivity,
	card *entities.Card,
) (bool, error) {
	return m.linkAccountActivityToCard("Online payment from CHK", 3, accountActivity, card)
}

func (m *App) linkAccountActivityToCard(
	descriptionContains string,
	maxDateSpread int,
	accountActivity *entities.AccountActivity,
	card *entities.Card,
) (bool, error) {
	dateRange := 0
START:
	for _, cardActivity := range card.Activities {
		if strings.Contains(cardActivity.Description, descriptionContains) {
			datDuration := time.Duration(int64(time.Hour * 24 * time.Duration(int64(dateRange))))
			if cardActivity.Amount == -accountActivity.Amount &&
				(cardActivity.PostDate.Add(datDuration) == accountActivity.PostingDate ||
					cardActivity.PostDate.Add(-datDuration) == accountActivity.PostingDate) {
				err := m.db.Exec(`
					INSERT INTO account_card_links (
						card_activity_uuid,
						account_activity_uuid
					)
					VALUES (?, ?)
					`,
					cardActivity.UUID,
					accountActivity.UUID,
				).Error
				if err != nil {
					return false, err
				}
				return true, nil
			}
		}
	}
	if dateRange < maxDateSpread {
		dateRange += 1
		goto START
	}
	return false, nil
}

func (m *App) getAllAccountActivitiesForSplitwiseExpense(
	splitwiseExpense *entities.SplitwiseExpense,
	daySpread int,
	amountSpread float64,
) ([]*entities.AccountActivity, error) {
	query := m.db
	numExistingLinks := len(splitwiseExpense.AccountActivities)
	if numExistingLinks > 0 {
		exisitngAccountActivityUUIDs := make([]string, numExistingLinks)
		for i, cardActivity := range splitwiseExpense.AccountActivities {
			exisitngAccountActivityUUIDs[i] = cardActivity.UUID.String()
		}
		query = query.Where(
			"uuid NOT IN @existing_links",
			sql.Named("existing_links", exisitngAccountActivityUUIDs),
		)
	}
	query = query.Where(
		`
			(
				(-amount BETWEEN @a1 AND @a2)
					OR
				(posting_date BETWEEN DATE(@d1) AND DATE(@d2))
			)
		`,
		sql.Named("a1", splitwiseExpense.AmountPaid-amountSpread),
		sql.Named("a2", splitwiseExpense.AmountPaid+amountSpread),
		sql.Named("d1", splitwiseExpense.Date.Add(-time.Hour*24*time.Duration(daySpread))),
		sql.Named("d2", splitwiseExpense.Date.Add(time.Hour*24*time.Duration(daySpread))),
	)

	var allAccountActivities []*entities.AccountActivity
	err := query.
		Preload("SplitwiseExpenses").
		Find(&allAccountActivities).Error
	if err != nil {
		return nil, err
	}

	return allAccountActivities, nil
}

// LinkAccountActivityToSplitwiseExpense links a account activity to a splitwise expense
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

// UnLinkAccountActivityToSplitwiseExpense links a account activity to a splitwise expense
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
