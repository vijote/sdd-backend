package main

import (
	"log"

	"github.com/vijote/sdd-backend/internal/config"
	"github.com/vijote/sdd-backend/internal/server"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	if err := server.New(cfg).Start(); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
