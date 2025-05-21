package main

import (
	"log"
	"snapkeep/internal/backup"
	"snapkeep/internal/config"
)

func main() {
	config.LoadEnv()

	// db, err := db.InitializeDB()
	// if err != nil {
	// 	logger.Fatal("Failed to initialize database: ", err)
	//   return
	// }

	cfg := config.Load()
	if err := backup.Run(cfg); err != nil {
		log.Fatalf("Backup failed: %v", err)
	}
}
