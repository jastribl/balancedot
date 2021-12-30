package splitwise

import (
	"fmt"
	"io/ioutil"
)

// GetRawExpense fetches the raw expense for a given splitwise expense ID
func (c *Client) GetRawExpense(id string) (string, error) {
	resp, err := c.getURLWithoutParams(fmt.Sprintf(getExpenseURL, id))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
