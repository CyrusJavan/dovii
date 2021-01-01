package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/CyrusJavan/dovii"
	"github.com/gorilla/mux"
	"github.com/hashicorp/raft"
	"github.com/hashicorp/raft-boltdb"
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

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		err := http.ListenAndServe("0.0.0.0:7070", s)
		if err != nil {
			log.Println("listenandserve:", err)
		}
		wg.Done()
	}()

	go func() {
		err := doRaft()
		if err != nil {
			log.Println("raft:", err)
		}
		wg.Done()
	}()

	wg.Wait()
}

func doRaft() error {
	conf := raft.DefaultConfig()
	bs, err := raftboltdb.NewBoltStore("/tmp/raftboltdbstore")
	if err != nil {
		return fmt.Errorf("init boltstore: %w", err)
	}

	ss := &raft.DiscardSnapshotStore{}

	t, err := raft.NewTCPTransport("localhost:7777", nil, 3, time.Second, os.Stdout)
	if err != nil {
		return fmt.Errorf("new tcp transport: %w", err)
	}

	servConf := raft.Configuration{Servers: []raft.Server{
		{
			Suffrage: raft.Voter,
			ID:       "node1",
			Address:  "0.0.0.0:7777",
		},
	}}

	err = raft.BootstrapCluster(conf, bs, bs, ss, t, servConf)
	if err != nil {
		return fmt.Errorf("bootstrapping cluster: %w", err)
	}

	return nil
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
