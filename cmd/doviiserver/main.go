package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/CyrusJavan/dovii"
	"github.com/gorilla/mux"
	"github.com/hashicorp/raft"
	"github.com/hashicorp/raft-boltdb"
	"github.com/kelseyhightower/envconfig"
)

type server struct {
	store      *dovii.KVStore
	router     *mux.Router
	raftServer *raft.Raft
}

type ServerConfig struct {
	ID string
}

func main() {
	var sc ServerConfig
	err := envconfig.Process("dovii", &sc)
	if err != nil {
		log.Fatal(err)
	}

	store, err := dovii.NewKVStore(dovii.BitcaskEngine)
	if err != nil {
		log.Fatal(err)
	}

	r, err := newRaftServer(&sc, store)
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

func newRaftServer(sc *ServerConfig, store *dovii.KVStore) (*raft.Raft, error) {
	log.Println("unique node ID:", sc.ID)
	conf := raft.DefaultConfig()
	conf.LocalID = raft.ServerID(sc.ID)
	bs, err := raftboltdb.NewBoltStore(fmt.Sprintf("/tmp/raftboltdbstore-%s", sc.ID))
	if err != nil {
		return nil, fmt.Errorf("init boltstore: %w", err)
	}

	ss := raft.DiscardSnapshotStore{}

	adAddr := net.TCPAddr{
		IP:   []byte("10.0.0.2"),
		Port: 7777,
	}
	log.Println("advertise address:", string(adAddr.IP)+":"+strconv.Itoa(adAddr.Port))
	t, err := raft.NewTCPTransport("0.0.0.0:7777", &adAddr, 3, time.Second, os.Stdout)
	if err != nil {
		return nil, fmt.Errorf("new tcp transport: %w", err)
	}

	fsm := BitcaskFSM{}
	raftServer, err := raft.NewRaft(conf, &fsm, bs, bs, &ss, t)
	if err != nil {
		return nil, fmt.Errorf("starting raft rpc server: %w", err)
	}

	return raftServer, nil
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

type LogMessage struct {
	Key   string
	Value string
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

type BitcaskFSM struct {
	store *dovii.KVStore
}

func (fsm *BitcaskFSM) Apply(raftLog *raft.Log) interface{} {
	if raftLog.Type != raft.LogCommand {
		log.Println("raftLog.Type is not LogCommand")
		return nil
	}
	var lm LogMessage
	err := json.Unmarshal(raftLog.Data, &lm)
	if err != nil {
		log.Println("decode raftLog.Data:", err)
		return err
	}
	err = fsm.store.Set(lm.Key, lm.Value)
	if err != nil {
		log.Println("setting value to store:", err)
		return err
	}
	return nil
}

func (fsm *BitcaskFSM) Snapshot() (raft.FSMSnapshot, error) {
	return newSnapshotNoop()
}

func (fsm *BitcaskFSM) Restore(rc io.ReadCloser) error {
	decoder := json.NewDecoder(rc)
	for decoder.More() {
		var lm LogMessage
		err := decoder.Decode(&lm)
		if err != nil {
			log.Println("restore decode message:", err)
			return err
		}
		err = fsm.store.Set(lm.Key, lm.Value)
		if err != nil {
			log.Println("restore setting value to store:", err)
			return err
		}
	}
	return nil
}

type snapshotNoop struct{}

func (s snapshotNoop) Persist(_ raft.SnapshotSink) error { return nil }

func (s snapshotNoop) Release() {}

func newSnapshotNoop() (raft.FSMSnapshot, error) {
	return &snapshotNoop{}, nil
}
