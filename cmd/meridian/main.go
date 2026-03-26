package main

import (
	"log"
	"net/http"
	"strconv"

	"github.com/Saad7890-web/meridian/internal/api"
	"github.com/Saad7890-web/meridian/internal/config"
	"github.com/Saad7890-web/meridian/internal/raft"
	"github.com/Saad7890-web/meridian/internal/storage"
)

func main() {
	cfg := config.Load()

	wal, err := storage.NewWAL("data.log")
	if err != nil {
		log.Fatalf("failed to init WAL: %v", err)
	}

	store := storage.NewStore(wal)

	if err := store.Recover(); err != nil {
		log.Fatalf("recovery failed: %v", err)
	}

	log.Println("Storage initialized & recovered")


	fsm := raft.NewFSM(store)

	// 4. Raft Node
	raftNode, err := raft.NewNode(
		cfg.NodeID,
		"0.0.0.0:"+strconv.Itoa(cfg.Port+1000),
		fsm,
	)
	if err != nil {
		log.Fatalf("raft init failed: %v", err)
	}

	// 5. KV Service (IMPORTANT: pass raftNode)
	

	http.HandleFunc("/health", api.HealthHandler)

	// kvHandler := api.NewKVHandler(store)

	// http.HandleFunc("/put", kvHandler.Put)
	// http.HandleFunc("/get", kvHandler.Get)
	// http.HandleFunc("/delete", kvHandler.Delete)

	
	kvService := api.NewKVService(store, raftNode)

	go api.StartGRPCServer(cfg.Port, kvService)

	// go func() {
	// 	addr := ":" + strconv.Itoa(cfg.Port)
	// 	log.Printf("HTTP server running on %s\n", addr)

	// 	if err := http.ListenAndServe(addr, nil); err != nil {
	// 		log.Fatalf("server failed: %v", err)
	// 	}
	// }()

	select {}
}