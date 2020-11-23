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
	getCurrentUserURL = "https://secure.splitwise.com/api/v3.0/get_current_user"
)

// Client holds all things for Splitwise requests
type Client struct {
	httpClient *http.Client
}

func getTokenFromWeb(oauthConfig *oauth2.Config, cfg *config.Config) (*oauth2.Token, error) {
	// Get the url and pass it through the channel to be followed
	cfg.AuthURLChan <- oauthConfig.AuthCodeURL(cfg.State, oauth2.AccessTypeOffline)

	// Get the code passsed back and exchange client token info
	return oauthConfig.Exchange(context.TODO(), <-cfg.AuthCodeChan)
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

// Saves a token to a file path.
func saveToken(path string, token *oauth2.Token) {
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

// NewClient gets a new client
func NewClient(cfg *config.Config) (*Client, error) {
	oauthConfig := &oauth2.Config{
		ClientID:     cfg.ConsumerKey,
		ClientSecret: cfg.ConsumerSecret,
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://secure.splitwise.com/oauth/authorize",
			TokenURL: "https://secure.splitwise.com/oauth/token",
		},
		RedirectURL: cfg.OAuthCallbackURL,
	}

	tokFile := cfg.TokenFileLocation
	tok, err := tokenFromFile(tokFile)
	if err != nil {
		tok, err = getTokenFromWeb(oauthConfig, cfg)
		if err != nil {
			return nil, err
		}
		saveToken(tokFile, tok)
	}

	return &Client{
		oauthConfig.Client(context.Background(), tok),
	}, nil
}
