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
	cfg := config.NewConfig()
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

	// Splitwise
	mainRouter.Handle(
		"/api/splitwise_login_check",
		api.Handler(apiApp.SplitwiseLoginCheck),
	)
	mainRouter.Handle(
		"/api/splitwise_oauth_callback",
		api.Handler(apiApp.SplitwiseOauthCallback),
	)
	mainRouter.Handle(
		"/api/splitwise_expenses/unlinked",
		api.Handler(apiApp.GetAllUnlinkedSplitwiseExpenses),
	).Methods("GET")
	mainRouter.Handle(
		"/api/splitwise_expenses/{splitwiseExpenseUUID}",
		api.Handler(apiApp.GetSplitwiseExpenseByUUID),
	).Methods("GET")
	mainRouter.Handle(
		"/api/splitwise_expenses/{splitwiseExpenseUUID}/for_linking",
		api.Handler(apiApp.GetSplitwiseExpenseByUUIDForLinking),
	).Methods("GET")
	mainRouter.Handle(
		"/api/splitwise_expenses/{splitwiseExpenseID}/raw",
		api.Handler(apiApp.GetRawSplitwiseExpense),
	).Methods("GET")
	mainRouter.Handle(
		"/api/splitwise_expenses",
		api.Handler(apiApp.GetAllSplitwiseExpenses),
	).Methods("GET")
	mainRouter.Handle(
		"/api/refresh_splitwise",
		api.Handler(apiApp.RefreshSplitwise),
	).Methods("POST")

	// Card Activities
	mainRouter.Handle(
		"/api/cards/{cardUUID}/activities",
		api.Handler(apiApp.UploadCardActivities),
	).Methods("POST")
	mainRouter.Handle(
		"/api/card_activities/{cardActivityUUID}",
		api.Handler(apiApp.GetCardActivityByUUID),
	).Methods("GET")
	mainRouter.Handle(
		"/api/card_activities/{cardActivityUUID}/link_splitwise/{splitwiseExpenseUUID}",
		api.Handler(apiApp.LinkCardActivityToSplitwiseExpense),
	).Methods("POST")
	mainRouter.Handle(
		"/api/card_activities/{cardActivityUUID}/unlink_splitwise/{splitwiseExpenseUUID}",
		api.Handler(apiApp.UnLinkCardActivityToSplitwiseExpense),
	).Methods("POST")

	// Cards
	mainRouter.Handle(
		"/api/cards/{cardUUID}",
		api.Handler(apiApp.GetCardByUUID),
	).Methods("GET")
	mainRouter.Handle(
		"/api/cards",
		api.Handler(apiApp.GetAllCards),
	).Methods("GET")
	mainRouter.Handle(
		"/api/card",
		api.Handler(apiApp.CreateNewCard),
	).Methods("POST")

	// Account Activities
	mainRouter.Handle(
		"/api/accounts/{accountUUID}/activities",
		api.Handler(apiApp.UploadAccountActivities),
	).Methods("POST")
	mainRouter.Handle(
		"/api/accounts/{accountUUID}/auto_link_with_card_activities",
		api.Handler(apiApp.AutoLinkAccountToCardActivities),
	).Methods("POST")
	mainRouter.Handle(
		"/api/account_activities/{accountActivityUUID}",
		api.Handler(apiApp.GetAccountActivityByUUID),
	).Methods("GET")
	mainRouter.Handle(
		"/api/account_activities/{accountActivityUUID}/link_splitwise/{splitwiseExpenseUUID}",
		api.Handler(apiApp.LinkAccountActivityToSplitwiseExpense),
	).Methods("POST")
	mainRouter.Handle(
		"/api/account_activities/{accountActivityUUID}/unlink_splitwise/{splitwiseExpenseUUID}",
		api.Handler(apiApp.UnLinkAccountActivityToSplitwiseExpense),
	).Methods("POST")

	// Chequing Accounts
	mainRouter.Handle(
		"/api/accounts/{accountUUID}",
		api.Handler(apiApp.GetAccountByUUID),
	).Methods("GET")
	mainRouter.Handle(
		"/api/accounts",
		api.Handler(apiApp.GetAllAccounts),
	).Methods("GET")
	mainRouter.Handle(
		"/api/account",
		api.Handler(apiApp.CreateNewAccount),
	).Methods("POST")

	mainRouter.PathPrefix("/static/").
		Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./client/public"))))
	mainRouter.PathPrefix("/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./client/public/index.html")
	})

	// Run the server
	log.Fatal(http.ListenAndServe("localhost:8080", mainRouter))
}
