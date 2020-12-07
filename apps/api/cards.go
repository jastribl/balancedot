package api

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"gihub.com/jastribl/balancedot/entities"
	"gihub.com/jastribl/balancedot/repos"
)

// GetAllCards get all the Cards
func (m *App) GetAllCards(w http.ResponseWriter, r *http.Request) {
	cardRepo := repos.NewCardRepo(m.db)
	cards, err := cardRepo.GetAllCards()
	if err != nil {
		log.Fatal(err)
	}
	if err := json.NewEncoder(w).Encode(cards); err != nil {
		log.Panic(err)
	}
}

type newCardParams struct {
	LastFour    string `json:"last_four"`
	Description string `json:"description"`
}

// CreateNewCard adds a new Card
func (m *App) CreateNewCard(w http.ResponseWriter, r *http.Request) {
	time.Sleep(time.Second * 5)
	decoder := json.NewDecoder(r.Body)
	var p newCardParams
	err := decoder.Decode(&p)
	if err != nil {
		panic(err)
	}
	err = m.db.Save(&entities.Card{
		LastFour:    p.LastFour,
		Description: p.Description,
	}).Error
	if err != nil {
		log.Fatal(err)
	}
	m.GetAllCards(w, r)

	// var newCardParam newCardParams
	// body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	// if err != nil {
	// 	panic(err)
	// }
	// if err := r.Body.Close(); err != nil {
	// 	panic(err)
	// }
	// if err := json.Unmarshal(body, &newCardParam); err != nil {
	// 	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	// 	w.WriteHeader(422) // unprocessable entity
	// 	if err := json.NewEncoder(w).Encode(err); err != nil {
	// 		panic(err)
	// 	}
	// }

	// log.Printf("%#v\n", newCardParam)

	// t := RepoCreateTodo(newCardParam)
	// w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	// w.WriteHeader(http.StatusCreated)
	// if err := json.NewEncoder(w).Encode(t); err != nil {
	//     panic(err)
	// }
}
