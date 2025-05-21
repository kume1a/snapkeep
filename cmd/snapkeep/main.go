package main

import (
	"log"
	"snapkeep/internal/backup"
	"snapkeep/internal/config"
)

func main() {
	cfg := config.Load()
	if err := backup.Run(cfg); err != nil {
		log.Fatalf("Backup failed: %v", err)
	}
}
