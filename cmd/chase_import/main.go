package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"gihub.com/jastribl/balancedot/chase"
	"gihub.com/jastribl/balancedot/entities"
	"github.com/jinzhu/gorm"

	"gihub.com/jastribl/balancedot/helpers"
	"gihub.com/jastribl/balancedot/repos"
)

func main() {
	db, err := helpers.DbConnect()
	if err != nil {
		log.Panic(err)
	}

	card := getCardFromPrompt(db)
	fileName := getFilePathFromPrompt()

	log.Printf("Importing %s for card %s", fileName, card.LastFour)

	cardActivities, err := chase.GetCardActivitiesFromFile(fileName)
	if err != nil {
		log.Fatal(err)
	}

	for _, cardActivity := range cardActivities {
		newCardActivity := entities.CardActivity{
			CardUUID:        card.UUID,
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

func getFilePathFromPrompt() string {
	for {
		var fileName string
		fmt.Print("Input File: ")
		fmt.Scanln(&fileName)
		info, err := os.Stat(fileName)
		if !os.IsNotExist(err) && !info.IsDir() {
			return fileName
		}
		fmt.Printf("Invalid input file (%s). Please choose again.\n", fileName)
	}
}

func getCardFromPrompt(db *gorm.DB) *entities.Card {
	cardRepo := repos.NewCardRepo(db)
	cards, err := cardRepo.GetAllCards()
	if err != nil {
		log.Fatal(err)
	}

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
			var lastFour, cardDescription string
			for {
				fmt.Print("Please Enter Last Four: ")
				fmt.Scanln(&lastFour)
				// TODO: check for numbers only
				if len(lastFour) == 4 {
					break
				}
			}
			fmt.Print("Please Enter Bank Name: ")
			fmt.Scanln(&cardDescription)
			newCard := entities.Card{
				LastFour:    lastFour,
				Description: cardDescription,
			}
			err := db.Save(&newCard).Error
			if err != nil {
				log.Fatal(err)
			}
			cards, err = cardRepo.GetAllCards()
			if err != nil {
				log.Fatal(err)
			}
			continue
		}
		cardIndex, err := strconv.Atoi(cardSelection)
		if err == nil && cardIndex >= 0 && cardIndex < len(cards) {
			return cards[cardIndex]
		}
	}
}
