package config

import (
	"os"
	"strconv"
)

type Config struct {
	NodeID string
	Port   int
}

func Load() *Config {
	port, _ := strconv.Atoi(getEnv("PORT", "8080"))

	return &Config{
		NodeID: getEnv("NODE_ID", "node-1"),
		Port:   port,
	}
}

func getEnv(key, fallback string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return fallback
}