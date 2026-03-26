package raft

import (
	"net"
	"os"
	"time"

	"github.com/hashicorp/raft"
	boltdb "github.com/hashicorp/raft-boltdb"
)

type Node struct {
	Raft *raft.Raft
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

	
	r.BootstrapCluster(raft.Configuration{
		Servers: []raft.Server{
			{
				ID:      config.LocalID,
				Address: transport.LocalAddr(),
			},
		},
	})

	return &Node{Raft: r}, nil
}

func (n *Node) Apply(cmd []byte) error {
	f := n.Raft.Apply(cmd, 5*time.Second)
	return f.Error()
}