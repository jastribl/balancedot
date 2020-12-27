package main

import (
	"log"
	"net/http"
	"os"

	"gihub.com/jastribl/balancedot/apps/api"
	"gihub.com/jastribl/balancedot/config"
	"gihub.com/jastribl/balancedot/helpers"
	"github.com/gorilla/mux"
	"github.com/pkg/browser"
)

func main() {
	// Setup logging
	log.SetOutput(os.Stdout)
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// Setup config and splitwise auth
	var cfg *config.Config
	cfg = config.NewConfig()
	go func() {
		browser.OpenURL(<-cfg.Splitwise.AuthURLChan)
	}()

	// Setup DB
	db, err := helpers.DbConnect()
	if err != nil {
		log.Panic(err)
	}

	// Setup app
	apiApp, err := api.NewApp(db, cfg)
	if err != nil {
		log.Fatal(err)
	}

	// Setup routing
	mainRouter := mux.NewRouter()

	mainRouter.Handle("/api/splitwise_login_check", api.Handler(apiApp.SplitwiseLoginCheck))
	mainRouter.Handle("/api/splitwise_oauth_callback", api.Handler(apiApp.SplitiwseOatuhCallback))

	mainRouter.Handle("/api/cards/{cardUUID}/activities", api.Handler(apiApp.GetAllCardActivitiesForCard)).Methods("GET")
	mainRouter.Handle("/api/cards/{cardUUID}/activity", api.Handler(apiApp.UploadCardActivities)).Methods("POST")

	mainRouter.Handle("/api/splitwise_expenses", api.Handler(apiApp.GetAllSplitwiseExpenses)).Methods("GET")
	mainRouter.Handle("/api/refresh_splitwise", api.Handler(apiApp.RefreshSplitwise)).Methods("POST")

	mainRouter.Handle("/api/cards/{cardUUID}", api.Handler(apiApp.GetCardByUUID)).Methods("GET")
	mainRouter.Handle("/api/cards", api.Handler(apiApp.GetAllCards)).Methods("GET")

	mainRouter.Handle("/api/card", api.Handler(apiApp.CreateNewCard)).Methods("POST")

	mainRouter.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./client/public"))))
	mainRouter.PathPrefix("/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./client/public/index.html")
	})

	// Run the server
	log.Fatal(http.ListenAndServe("localhost:8080", mainRouter))
}
