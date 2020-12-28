package main

import (
	"encoding/json"
	"fmt"
	"log"

	"gihub.com/jastribl/balancedot/entities"
	"gihub.com/jastribl/balancedot/helpers"
)

func main() {
	db, err := helpers.DbConnect()
	if err != nil {
		log.Panic(err)

	}
	// db.LogMode(true)

	var expenses []*entities.SplitwiseExpense
	err = db.Preload("CardActivities").Where("deleted_at is NULL AND amount_paid > 0").Find(&expenses).Error
	if err != nil {
		log.Fatal(err)
	}

	// jsonEncoder := json.NewEncoder(log.Writer())
	for _, expense := range expenses {
		if len(expense.CardActivities) > 0 {
			continue
		}

		var cardActivities []*entities.CardActivity
		err := db.Where("amount = ?", -expense.AmountPaid).Find(&cardActivities).Error
		if err != nil {
			log.Fatal(err)
		}
		if len(cardActivities) == 0 {
			continue
		}
		fmt.Printf("Expense: %s\n\n", prettyJSON(expense))
		// log.Print("Found the following card activities")
		// jsonEncoder.Encode(cardActivities)
	Input:
		for {
			for i, activity := range cardActivities {
				fmt.Printf("%d)\n%s\n\n", i, prettyJSON(activity))
			}
			var activityIndex int
			for {
				fmt.Printf("Select index: ")
				fmt.Scanln(&activityIndex)
				fmt.Printf("You selected index %d\n", activityIndex)
				if activityIndex < len(cardActivities) || activityIndex == -1 {
					break
				}
			}
			if activityIndex == -1 {
				break Input
			}
			expense.CardActivities = append(expense.CardActivities, cardActivities[activityIndex])
			err = db.Create(&expense).Error
			if err != nil {
				log.Fatal(err)
			}
			break Input
		}
	}
}

func prettyJSON(thing interface{}) string {
	prettyJSON, err := json.MarshalIndent(thing, "", "    ")
	if err != nil {
		log.Fatal(err)
	}
	return string(prettyJSON)
}
