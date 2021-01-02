package raft

import (
	"encoding/json"
	"io"
	"log"

	"github.com/CyrusJavan/dovii/internal/dovii"
	"github.com/hashicorp/raft"
)

type LogMessage struct {
	Key   string
	Value string
}

type BitcaskFSM struct {
	Store *dovii.KVStore
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
	err = fsm.Store.Set(lm.Key, lm.Value)
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
		err = fsm.Store.Set(lm.Key, lm.Value)
		if err != nil {
			log.Println("restore setting value to store:", err)
			return err
		}
	}
	return nil
}
