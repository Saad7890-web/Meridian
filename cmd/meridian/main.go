package main

import (
	"log"
	"net/http"
	"strconv"

	"github.com/Saad7890-web/meridian/internal/api"
	"github.com/Saad7890-web/meridian/internal/config"
)

func main() {
	cfg := config.Load()

	http.HandleFunc("/health", api.HealthHandler)

	go func() {
		log.Printf("HTTP server running on :%d\n", cfg.Port)
		err := http.ListenAndServe(":"+strconv.Itoa(cfg.Port), nil)
    	if err != nil {
        log.Fatalf("Server failed: %v", err)
    	}
	}()

	select {}
}