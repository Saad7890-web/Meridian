package raft

import (
	"encoding/json"
	"io"

	"github.com/Saad7890-web/meridian/internal/storage"
	"github.com/hashicorp/raft"
)

type FSM struct {
	store *storage.Store
}

func NewFSM(store *storage.Store) *FSM {
	return &FSM{store: store}
}


func (f *FSM) Apply(log *raft.Log) interface{} {
	var cmd map[string]string

	if err := json.Unmarshal(log.Data, &cmd); err != nil {
		return nil
	}

	switch cmd["op"] {
	case "put":
		f.store.Put(cmd["key"], cmd["value"])
	case "delete":
		f.store.Delete(cmd["key"])
	}

	return nil
}

func (f *FSM) Snapshot() (raft.FSMSnapshot, error) {
	return &noopSnapshot{}, nil
}

func (f *FSM) Restore(rc io.ReadCloser) error {
	return nil
}

type noopSnapshot struct{}

func (n *noopSnapshot) Persist(sink raft.SnapshotSink) error {
	sink.Cancel()
	return nil
}

func (n *noopSnapshot) Release() {}