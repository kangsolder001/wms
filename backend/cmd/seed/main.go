package main

import (
	"log"

	"wms/config"
	"wms/infrastructure/database"
	loggerInfra "wms/infrastructure/logger"
)

func main() {
	cfg, err := config.Load("config.json")
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	appLogger := loggerInfra.New(cfg.Logger.Level, cfg.Logger.Output)

	db, err := database.New(cfg.Database)
	if err != nil {
		appLogger.Fatal("failed to connect database", "error", err)
	}
	defer db.Close()

	if err := database.Migrate(db); err != nil {
		appLogger.Fatal("failed to run migrations", "error", err)
	}

	database.Seed(db, appLogger)
}
