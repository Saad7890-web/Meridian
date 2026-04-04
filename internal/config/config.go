package config

import (
	"os"
	"strconv"
	"strings"
)

type Config struct {
	NodeID string
	Port   int
	Peers []string
	RaftPort int
}

func Load() *Config {
	port, _ := strconv.Atoi(getEnv("PORT", "8080"))
	raftPort, _ := strconv.Atoi(getEnv("RAFT_PORT", "9000"))

	peers := strings.Split(getEnv("PEERS", ""), ",")

	return &Config{
		NodeID:   getEnv("NODE_ID", "node-1"),
		Port:     port,
		RaftPort: raftPort,
		Peers:    peers,
	}
}

func getEnv(key, fallback string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return fallback
}