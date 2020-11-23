package config

import (
	"encoding/json"
	"fmt"
	"os"

	"gihub.com/jastribl/balancedot/splitwise/config"
)

// Config is the struct that holds application config info
type Config struct {
	Splitwise config.Config `json:"splitwise"`
	ServerURL string        `json:"server-url"`
}

var configCache *Config

// NewConfig gets a new Config
func NewConfig() *Config {
	if configCache == nil {
		configCache = new(Config)
		configFile, err := os.Open("config/config.json")
		defer configFile.Close()
		if err != nil {
			fmt.Println(err.Error())
		}
		jsonParser := json.NewDecoder(configFile)
		jsonParser.Decode(configCache)

		// init channels
		configCache.Splitwise.AuthURLChan = make(chan string)
		configCache.Splitwise.AuthCodeChan = make(chan string)
	}
	return configCache
}
