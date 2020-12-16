package main

import (
	"log"
	"net/http"
	"os"

	"gihub.com/jastribl/balancedot/apps/api"
	"gihub.com/jastribl/balancedot/helpers"
	"github.com/gorilla/mux"
)

func main() {
	log.SetOutput(os.Stdout)
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	db, err := helpers.DbConnect()
	if err != nil {
		log.Panic(err)
	}

	mainRouter := mux.NewRouter()

	apiApp, err := api.NewApp(db)
	if err != nil {
		log.Fatal(err)
	}

	mainRouter.Handle("/api/cards/{cardUUID}/activities", api.Handler(apiApp.GetAllCardActivitiesForCard)).Methods("GET")

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
