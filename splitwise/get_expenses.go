package splitwise

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"gihub.com/jastribl/balancedot/splitwise/models"
)

const getExpensesURL = "https://secure.splitwise.com/api/v3.0/get_expenses"

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

		url, err := buildURLWithParams(getExpensesURL, params)
		if err != nil {
			return nil, err
		}
		resp, err := c.httpClient.Get(url)
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

func buildURLWithParams(url string, params map[string]string) (string, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}

	// If we have any parameters, add them here.
	if len(params) > 0 {
		query := req.URL.Query()
		for k, v := range params {
			query.Add(k, v)
		}
		req.URL.RawQuery = query.Encode()
	}
	return req.URL.String(), nil
}
