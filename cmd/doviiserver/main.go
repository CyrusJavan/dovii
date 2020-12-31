package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/CyrusJavan/dovii"
	"github.com/gorilla/mux"
)

type server struct {
	store  *dovii.KVStore
	router *mux.Router
}

func main() {
	store, err := dovii.NewKVStore(dovii.BitcaskEngine)
	if err != nil {
		log.Fatal(err)
	}

	s := &server{
		store:  store,
		router: mux.NewRouter(),
	}

	s.routes()

	log.Fatal(http.ListenAndServe("0.0.0.0:7070", s))
}

func (s *server) routes() {
	s.router.HandleFunc("/{key}", s.getHandler()).
		Methods("GET")

	s.router.HandleFunc("/{key}/{value}", s.setHandler()).
		Methods("POST")
}

func (s *server) setHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		key := vars["key"]
		value := vars["value"]
		err := (*s.store).Set(key, value)
		if err != nil {
			w.WriteHeader(http.StatusBadGateway)
			return
		}
	}
}

func (s *server) getHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		key := vars["key"]
		value, err := (*s.store).Get(key)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		b, err := json.Marshal(map[string]string{
			"value": value,
		})
		if err != nil {
			log.Fatal(err)
		}
		_, err = w.Write(b)
		if err != nil {
			log.Fatal(err)
		}
	}
}
func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}
