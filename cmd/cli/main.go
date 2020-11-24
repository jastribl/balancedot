package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"gihub.com/jastribl/balancedot/config"
	"gihub.com/jastribl/balancedot/helpers"
	"gihub.com/jastribl/balancedot/repos"
	"gihub.com/jastribl/balancedot/splitwise"
	"github.com/pkg/browser"
	uuid "github.com/satori/go.uuid"
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
	}

	jsonEncoder := json.NewEncoder(log.Writer())

	jsonEncoder.Encode(user)

	expenses, err := splitwiseClient.GetExpenses()
	if err != nil {
		log.Fatal(err)
	}
	jsonEncoder.Encode(expenses)

	db, err := helpers.DbConnect()
	if err != nil {
		log.Panic(err)
	}

	cardRepo := repos.NewCardRepo(db)
	card, err := cardRepo.GetCard("2427")
	if err != nil {
		log.Panic(err)
	}
	jsonEncoder.Encode(card)

	cardActivityRepo := repos.NewCardActivityRepo(db)
	uuid, err := uuid.FromString("d2300ddd-e048-4aad-93e2-4e9d07850714")
	if err != nil {
		log.Panic(err)
	}
	cardActivity, err := cardActivityRepo.GetCardActivity(uuid)
	if err != nil {
		log.Panic(err)
	}
	jsonEncoder.Encode(cardActivity)
}
