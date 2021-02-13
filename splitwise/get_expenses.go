package splitwise

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"gihub.com/jastribl/balancedot/splitwise/models"
)

// GetExpenses fetches all expenses for the currently logged in user
func (c *Client) GetExpenses() (*[]*models.Expense, error) {
	var allExpenses []*models.Expense
	var offset = 0
	for {
		type expensesResponse struct {
			Expenses []*models.Expense `json:"expenses"`
		}

		params := map[string]string{
			"offset": fmt.Sprintf("%v", offset),
			"limit":  fmt.Sprintf("%v", 100),
		}

		resp, err := c.getURLWithParams(getExpensesURL, params)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		response := expensesResponse{}
		b, err := ioutil.ReadAll(resp.Body)
		if err := json.Unmarshal(b, &response); err != nil {
			return nil, err
		}
		newExpenses := response.Expenses
		numNewExpenses := len(newExpenses)
		allExpenses = append(allExpenses, newExpenses...)
		if numNewExpenses == 0 {
			break
		}
		offset += numNewExpenses
	}
	return &allExpenses, nil
}
