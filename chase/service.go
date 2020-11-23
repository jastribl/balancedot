package chase

import (
	"os"

	"gihub.com/jastribl/balancedot/chase/models"
	"github.com/gocarina/gocsv"
)

// GetCardActivitiesFromFile gets chase activites
func GetCardActivitiesFromFile(inputFile string) ([]*models.CardActivity, error) {
	clientsFile, err := os.OpenFile(
		inputFile,
		os.O_RDWR|os.O_CREATE,
		os.ModePerm,
	)
	if err != nil {
		return nil, err
	}
	defer clientsFile.Close()

	chaseActivities := []*models.CardActivity{}
	err = gocsv.UnmarshalFile(clientsFile, &chaseActivities)
	return chaseActivities, err
}

// PrintCardActivitiesToFile gets chase activites
func PrintCardActivitiesToFile(chaseActivities []*models.CardActivity, outputFile string) error {
	clientsFile, err := os.OpenFile(
		outputFile,
		os.O_RDWR|os.O_CREATE,
		os.ModePerm,
	)
	if err != nil {
		return err
	}

	return gocsv.MarshalFile(&chaseActivities, clientsFile)
}
