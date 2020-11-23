package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"gihub.com/jastribl/balancedot/chase"
	"gihub.com/jastribl/balancedot/entities"

	"gihub.com/jastribl/balancedot/helpers"
	"gihub.com/jastribl/balancedot/repos"
)

func main() {
	db, err := helpers.DbConnect()
	if err != nil {
		log.Panic(err)
	}

	cardRepo := repos.NewCardRepo(db)

	// cfg = config.NewConfig()
	cards, err := cardRepo.GetAllCards()
	if err != nil {
		log.Fatal(err)
	}

	var cardIndex uint
	for {
		fmt.Println("Please choose your card: ")
		for i, card := range cards {
			fmt.Printf("\t%d) %s\n", i, card.LastFour)
		}
		fmt.Printf("\tc) Create New Card\n")
		fmt.Print("\nCard index: ")
		var cardSelection string
		fmt.Scanln(&cardSelection)
		if cardSelection == "c" {
			var lastFour, bankName string
			for {
				fmt.Print("Please Enter Last Four: ")
				fmt.Scanln(&lastFour)
				// TODO: check for numbers only
				if len(lastFour) == 4 {
					break
				}
			}
			fmt.Print("Please Enter Bank Name: ")
			fmt.Scanln(&bankName)
			newCard := entities.Card{
				LastFour: lastFour,
				BankName: bankName,
			}
			err := db.Save(&newCard).Error
			if err != nil {
				log.Fatal(err)
			}
			cardIndex = newCard.ID
		}
		cardIndex, err := strconv.Atoi(cardSelection)
		if err == nil && cardIndex >= 0 && cardIndex < len(cards) {
			break
		}
	}

	card := cards[cardIndex]

	var fileName string
	for {
		fmt.Print("Input File: ")
		fmt.Scanln(&fileName)
		info, err := os.Stat(fileName)
		if !os.IsNotExist(err) && !info.IsDir() {
			break
		}
		fmt.Printf("Invalid input file (%s). Please choose again.\n", fileName)
	}

	log.Printf("Importing %s for card %s", fileName, card.LastFour)

	cardActivities, err := chase.GetCardActivitiesFromFile(fileName)
	if err != nil {
		log.Fatal(err)
	}

	for _, cardActivity := range cardActivities {
		newCardActivity := entities.CardActivity{
			CardID:          card.ID,
			TransactionDate: cardActivity.TransactionDate.Time,
			PostDate:        cardActivity.PostDate.Time,
			Description:     cardActivity.Description,
			Category:        cardActivity.Category,
			Type:            cardActivity.Type,
			Amount:          cardActivity.Amount.ToFloat64(),
		}
		// TODO: check for duplicates somehow
		err = db.Save(&newCardActivity).Error
		if err != nil {
			log.Fatal(err)
		}
	}
}
