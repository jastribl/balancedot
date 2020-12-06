package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func helloServer(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello %s!", r.URL.Path[1:])
}

type item struct {
	ID int    `json:"id"`
	A  string `json:"a"`
	B  string `json:"b"`
}

func helloJSON(w http.ResponseWriter, r *http.Request) {
	items := []item{
		item{
			ID: 1,
			A:  "Something",
			B:  "Something else",
		},
		item{
			ID: 2,
			A:  "Another things",
			B:  "Yet another thing",
		},
	}
	if err := json.NewEncoder(w).Encode(items); err != nil {
		log.Panic(err)
	}
}

func main() {
	fs := http.FileServer(http.Dir("./client/public"))
	http.Handle("/", fs)
	http.HandleFunc("/api", helloServer)
	http.HandleFunc("/api/json", helloJSON)

	log.Println("Listening on :5000...")
	err := http.ListenAndServe("127.0.0.1:5000", nil)
	if err != nil {
		log.Fatal(err)
	}
}
