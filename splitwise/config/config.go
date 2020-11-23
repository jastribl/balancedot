package config

// Config is the struct that holds splitwise config info
type Config struct {
	ConsumerKey       string      `json:"consumer-key"`
	ConsumerSecret    string      `json:"consumer-secret"`
	OAuthCallbackURL  string      `json:"oauth-callback-url"`
	State             string      `json:"state"`
	TokenFileLocation string      `json:"token-file-location"`
	AuthURLChan       chan string `json:"-"`
	AuthCodeChan      chan string `json:"-"`
}
