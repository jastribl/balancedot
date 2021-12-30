package splitwise

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"gihub.com/jastribl/balancedot/splitwise/config"
	"golang.org/x/oauth2"
)

const (
	getExpensesURL = "https://secure.splitwise.com/api/v3.0/get_expenses"
	getExpenseURL  = "https://secure.splitwise.com/api/v3.0/get_expense/%s"
)

// Client holds all things for Splitwise requests
type Client struct {
	httpClient *http.Client
}

// GetAuthPortalURL fetches the url required to redirect for Splitwise OAuth Authentication
func GetAuthPortalURL(cfg *config.Config) string {
	return getAuthConfig(cfg).AuthCodeURL(cfg.State, oauth2.AccessTypeOffline)
}

// GetTokenFromCode takes the code and exchanes it for the token
func GetTokenFromCode(cfg *config.Config, code string) (*oauth2.Token, error) {
	return getAuthConfig(cfg).Exchange(context.TODO(), code)
}

// HasToken returns if the user has a token
func HasToken(cfg *config.Config) bool {
	// todo: make work with user
	_, err := tokenFromFile(cfg)
	return err == nil
}

// Retrieves a token from a local file.
func tokenFromFile(cfg *config.Config) (*oauth2.Token, error) {
	f, err := os.Open(cfg.TokenFileLocation)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

// SaveToken saves a token given a config
func SaveToken(cfg *config.Config, token *oauth2.Token) error {
	f, err := os.OpenFile(cfg.TokenFileLocation, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	return json.NewEncoder(f).Encode(token)
}

// NewClientForUser gets a new client for a user using the user token
func NewClientForUser(cfg *config.Config) (*Client, error) {
	// todo: make this work with users
	tok, err := tokenFromFile(cfg)
	if err != nil {
		return nil, err
	}

	return &Client{
		getAuthConfig(cfg).Client(context.Background(), tok),
	}, nil
}

func getAuthConfig(cfg *config.Config) *oauth2.Config {
	return &oauth2.Config{
		ClientID:     cfg.ConsumerKey,
		ClientSecret: cfg.ConsumerSecret,
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://secure.splitwise.com/oauth/authorize",
			TokenURL: "https://secure.splitwise.com/oauth/token",
		},
		RedirectURL: cfg.OAuthCallbackURL,
	}
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

func (c *Client) getURLWithParams(url string, params map[string]string) (*http.Response, error) {
	url, err := buildURLWithParams(getExpensesURL, params)
	if err != nil {
		return nil, err
	}
	return c.httpClient.Get(url)
}

func (c *Client) getURLWithoutParams(url string) (*http.Response, error) {
	url, err := buildURLWithParams(url, map[string]string{})
	if err != nil {
		return nil, err
	}
	return c.httpClient.Get(url)
}
