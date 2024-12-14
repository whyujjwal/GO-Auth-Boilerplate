package main

import (
	"auth/config"
	"auth/internal/api/routes"
	"auth/pkg/database"
	"auth/pkg/logger"
	"log"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Map cfg.Logger to logger.LoggerConfig
	loggerConfig := logger.LoggerConfig{
		LogLevel:    cfg.Logger.LogLevel,
		LogFilePath: cfg.Logger.LogFilePath,
		MaxSize:     cfg.Logger.MaxSize,
		MaxBackups:  cfg.Logger.MaxBackups,
		MaxAge:      cfg.Logger.MaxAge,
	}

	// Initialize logger
	l, err := logger.InitFileLogger(loggerConfig)
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}

	// Initialize database
	db, err := database.Initialize(cfg.Database)
	if err != nil {
		l.Fatal().Err(err).Msg("Failed to initialize database")
	}

	// Setup router
	router := routes.Setup(db, l)

	// Start server
	l.Info().Msgf("Starting server on %s", cfg.Server.Address)
	if err := router.Run(cfg.Server.Address); err != nil {
		l.Fatal().Err(err).Msg("Server failed to start")
	}
}
