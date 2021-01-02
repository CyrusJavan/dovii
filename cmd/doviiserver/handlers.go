package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	doviiraft "github.com/CyrusJavan/dovii/internal/raft"
	"github.com/gorilla/mux"
)

func (s *server) setHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		key := vars["key"]
		value := vars["value"]
		log.Println("received POST for key=", key)
		lm := doviiraft.LogMessage{
			Key:   key,
			Value: value,
		}
		lmData, err := json.Marshal(lm)
		if err != nil {
			w.Write([]byte("marshal failed" + err.Error()))
			w.WriteHeader(http.StatusBadGateway)
			return
		}
		future := s.raftServer.Apply(lmData, time.Second)
		err = future.Error()
		if err != nil {
			w.Write([]byte(err.Error()))
			w.WriteHeader(http.StatusBadGateway)
			return
		}
	}
}

func (s *server) getHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		key := vars["key"]
		log.Println("received GET for key=", key)
		value, err := (*s.store).Get(key)
		if err != nil {
			w.Write([]byte(err.Error()))
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
