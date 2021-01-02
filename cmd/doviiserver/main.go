package main

import (
	"log"
	"net/http"
	"os"

	"github.com/CyrusJavan/dovii/internal/dovii"
	doviiraft "github.com/CyrusJavan/dovii/internal/raft"
	"github.com/gorilla/mux"
	"github.com/hashicorp/raft"
)

type server struct {
	store      *dovii.KVStore
	router     *mux.Router
	raftServer *raft.Raft
}

func main() {
	log.SetOutput(os.Stderr)

	store, err := dovii.NewKVStore(dovii.BitcaskEngine)
	if err != nil {
		log.Fatal(err)
	}

	r, err := doviiraft.NewRaftServer(store)
	if err != nil {
		log.Println("raft:", err)
	}

	s := &server{
		store:      store,
		router:     mux.NewRouter(),
		raftServer: r,
	}

	s.routes()

	err = http.ListenAndServe("0.0.0.0:7070", s)
	if err != nil {
		log.Println("listenandserve:", err)
	}
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}
