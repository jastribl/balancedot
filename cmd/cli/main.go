package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"gihub.com/jastribl/balancedot/chase"
	"gihub.com/jastribl/balancedot/config"
	"gihub.com/jastribl/balancedot/splitwise"
	"github.com/pkg/browser"
)

var cfg *config.Config

func oauthCallback(w http.ResponseWriter, r *http.Request) {
	if r.URL.Query().Get("state") != cfg.Splitwise.State {
		io.WriteString(w, "Invalid state in callback")
		return
	}

	cfg.Splitwise.AuthCodeChan <- r.URL.Query().Get("code")
	io.WriteString(w, "Processing Token...\nPlease close this tab and return to your application.")
}

func main() {
	cfg = config.NewConfig()
	chaseCardActivities, err := chase.GetCardActivitiesFromFile("in.csv")
	if err != nil {
		log.Fatal(err)
	}

	err = chase.PrintCardActivitiesToFile(chaseCardActivities, "out.csv")
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		browser.OpenURL(<-cfg.Splitwise.AuthURLChan)
	}()

	m := http.NewServeMux()
	m.HandleFunc("/oauth_callback", oauthCallback)
	s := http.Server{Addr: cfg.ServerURL, Handler: m}
	defer s.Close()

	// Run the server
	go func() {
		log.Fatal(s.ListenAndServe())
	}()

	splitwiseClient, err := splitwise.NewClient(&cfg.Splitwise)
	if err != nil {
		log.Fatal(err)
		return
	}

	user, err := splitwiseClient.GetCurrentUser()
	if err != nil {
		log.Fatal(err)
		return
	}

	jsonEncoder := json.NewEncoder(log.Writer())
	jsonEncoder.Encode(user)
}
