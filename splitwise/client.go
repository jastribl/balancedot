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

// Client holds all things for Splitwise requests
type Client struct {
	httpClient *http.Client
}

func getTokenFromWeb(cfg *config.Config) (*oauth2.Token, error) {
	// Get the url and pass it through the channel to be followed
	cfg.AuthURLChan <- GetAuthPortalURL(cfg)

	// Get the code passsed back and exchange client token info
	return GetTokenFromCode(cfg, <-cfg.AuthCodeChan)
}

// GetAuthPortalURL todo
func GetAuthPortalURL(cfg *config.Config) string {
	return getAuthConfig(cfg).AuthCodeURL(cfg.State, oauth2.AccessTypeOffline)
}

// GetTokenFromCode todo
func GetTokenFromCode(cfg *config.Config, code string) (*oauth2.Token, error) {
	return getAuthConfig(cfg).Exchange(context.TODO(), code)
}

// HasToken todo
func HasToken(cfg *config.Config) bool {
	f, err := os.Open(cfg.TokenFileLocation)
	if err != nil {
		return false
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return err == nil
}

// Retrieves a token from a local file.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
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

// ClientSetupError is a custom client setup error
type ClientSetupError struct {
	RedirectURL string
	Err         error
}

// Error returns the error
func (e ClientSetupError) Error() string {
	return e.Err.Error()
}

// NewClientForCLI gets a new client for a CLI using local token
func NewClientForCLI(cfg *config.Config) (*Client, error) {
	tokFile := cfg.TokenFileLocation
	tok, err := tokenFromFile(tokFile)
	if err != nil {
		tok, err = getTokenFromWeb(cfg)
		if err != nil {
			return nil, err
		}
		err = SaveToken(cfg, tok)
		if err != nil {
			return nil, err
		}
	}

	return &Client{
		getAuthConfig(cfg).Client(context.Background(), tok),
	}, nil
}

// NewClientForUser gets a new client for a user using the user token
func NewClientForUser(cfg *config.Config) (*Client, error) {
	// todo: make this work with users
	tokFile := cfg.TokenFileLocation
	tok, err := tokenFromFile(tokFile)
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
