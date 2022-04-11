package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/mux"
)

type Store struct {
	mu sync.Mutex
	m  map[string]string
}

var stores = &Store{m: make(map[string]string)}

func getStores(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stores.m)
}

func getStore(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for key, value := range stores.m {
		if key == params["key"] {
			json.NewEncoder(w).Encode(value)
			return
		}
	}
	json.NewEncoder(w).Encode(&Store{})
}

func updateStore(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for key := range stores.m {
		if key == params["key"] {
			stores.mu.Lock()
			stores.m[params["key"]] = params["value"]
			stores.mu.Unlock()
			json.NewEncoder(w).Encode(stores.m)
			return
		}
	}
	json.NewEncoder(w).Encode(&Store{})
}

func main() {
	r := mux.NewRouter()

	stores.m["1"] = "one"
	stores.m["2"] = "two"
	stores.m["3"] = "three"

	r.HandleFunc("/stores", getStores).Methods("GET")
	r.HandleFunc("/stores/{key}", getStore).Methods("GET")
	r.HandleFunc("/stores/{key}/{value}", updateStore).Methods("PUT")
	log.Fatal(http.ListenAndServe(":8000", r))
}
