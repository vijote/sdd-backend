package main

import (
	"log"
	"os"

	"github.com/vijote/sdd-backend/internal/config"
	"github.com/vijote/sdd-backend/internal/database"
	"github.com/vijote/sdd-backend/internal/server"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	command := "serve"
	if len(os.Args) > 1 {
		command = os.Args[1]
	}

	switch command {
	case "migrate":
		// One-shot migration runner. The infra repo wraps this in an init
		// container or Job; the serve process never mutates the schema.
		db, err := database.Open(cfg)
		if err != nil {
			log.Fatalf("failed to connect to database: %v", err)
		}
		if err := database.Migrate(db); err != nil {
			log.Fatalf("migration failed: %v", err)
		}
	case "serve":
		db, err := database.Open(cfg)
		if err != nil {
			log.Fatalf("failed to connect to database: %v", err)
		}
		if err := server.New(cfg, db).Start(); err != nil {
			log.Fatalf("server failed: %v", err)
		}
	default:
		log.Fatalf("unknown command %q (want \"serve\" or \"migrate\")", command)
	}
}
