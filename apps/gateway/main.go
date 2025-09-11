package main

import (
	"log"
	"os"

	"github.com/tienhung-ho/smart-document/common/config"
	"github.com/tienhung-ho/smart-document/common/errors"
	"github.com/tienhung-ho/smart-document/common/logging"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig("./etc", "gateway")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize logger
	logger, err := logging.InitLogger(&cfg.Logging)
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	defer logger.Sync()

	// Log startup
	logging.Info("Starting Gateway Service")
	logging.Infof("Environment: %s", cfg.Environment)
	logging.Infof("Server will start on %s:%d", cfg.Server.Host, cfg.Server.Port)

	// Example of structured logging
	logging.WithFields(map[string]any{
		"service": "gateway",
		"version": "1.0.0",
		"port":    cfg.Server.Port,
	}).Info("Service configuration loaded")

	// Example of error handling
	if err := simulateOperation(); err != nil {
		if appErr, ok := errors.GetAppError(err); ok {
			logging.Errorf("Application error occurred: [%d] %s", appErr.Code, appErr.Message)
			if appErr.Context != nil {
				logging.WithFields(appErr.Context).Error("Error context")
			}
		} else {
			logging.Errorf("Unexpected error: %v", err)
		}
		os.Exit(1)
	}

	logging.Info("Gateway service started successfully")

	// In a real application, this would start the HTTP server
	// For demonstration, we'll just exit
	logging.Info("Demo completed, shutting down...")
}

func simulateOperation() error {
	// Simulate a validation error
	return errors.Validation("Invalid request format").
		WithContext("field", "email").
		WithContext("value", "invalid-email").
		WithDetails("Email format is not valid")
}
