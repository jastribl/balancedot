package splitwise

import (
	"encoding/json"
	"io/ioutil"

	"gihub.com/jastribl/balancedot/splitwise/models"
)

const getCurrentUserURL = "https://secure.splitwise.com/api/v3.0/get_current_user"

// GetCurrentUser fetches and returns the currently logged in user
func (c *Client) GetCurrentUser() (*models.User, error) {
	type userReponse struct {
		User models.User `json:"user"`
	}

	resp, err := c.httpClient.Get(getCurrentUserURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	response := userReponse{}
	err = json.Unmarshal(b, &response)
	if err != nil {
		return nil, err
	}

	return &response.User, nil
}
