package raft

import (
	"net"
	"os"
	"time"

	"github.com/hashicorp/raft"
	boltdb "github.com/hashicorp/raft-boltdb"
)

type Node struct {
	Raft          *raft.Raft
	LogStore      raft.LogStore
	StableStore   raft.StableStore
	SnapshotStore raft.SnapshotStore
}


func NewNode(nodeID, bindAddr string, fsm *FSM) (*Node, error) {
	config := raft.DefaultConfig()
	config.LocalID = raft.ServerID(nodeID)

	addr, err := net.ResolveTCPAddr("tcp", bindAddr)
	if err != nil {
		return nil, err
	}

	transport, err := raft.NewTCPTransport(bindAddr, addr, 3, 10*time.Second, os.Stdout)
	if err != nil {
		return nil, err
	}

	logStore, err := boltdb.NewBoltStore("raft-log.db")
	if err != nil {
		return nil, err
	}

	stableStore, err := boltdb.NewBoltStore("raft-stable.db")
	if err != nil {
		return nil, err
	}

	snapshots, err := raft.NewFileSnapshotStore(".", 1, os.Stdout)
	if err != nil {
		return nil, err
	}

	r, err := raft.NewRaft(config, fsm, logStore, stableStore, snapshots, transport)
	if err != nil {
		return nil, err
	}

	return &Node{
		Raft:          r,
		LogStore:      logStore,
		StableStore:   stableStore,
		SnapshotStore: snapshots,
	}, nil
}


func (n *Node) BootstrapIfNeeded(nodeID string, peers []string) {
	hasState, err := raft.HasExistingState(n.LogStore, n.StableStore, n.SnapshotStore)
	if err != nil {
		return
	}

	if hasState {
		return
	}

	var servers []raft.Server

	for _, peer := range peers {
		if peer == "" {
			continue
		}

		servers = append(servers, raft.Server{
			ID:      raft.ServerID(peer),
			Address: raft.ServerAddress(peer),
		})
	}

	n.Raft.BootstrapCluster(raft.Configuration{
		Servers: servers,
	})
}


func (n *Node) Apply(cmd []byte) error {
	f := n.Raft.Apply(cmd, 5*time.Second)
	return f.Error()
}