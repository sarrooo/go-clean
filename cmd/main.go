package main

import (
	"fmt"
	"os"

	"github.com/sarrooo/go-clean/internal/controllers"
	"github.com/sarrooo/go-clean/internal/database"
	"github.com/sarrooo/go-clean/internal/logger"
	"github.com/sarrooo/go-clean/internal/repositories"
	"github.com/sarrooo/go-clean/internal/services"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func main() {
	initViper()

	// Initialize logger
	logger, err := logger.New()
	if err != nil {
		panic(fmt.Errorf("error logger: %w", err))
	}
	defer func(Logger *zap.Logger) {
		// Sync logger before exit
		if err := Logger.Sync(); err != nil {
			logger.Fatal("Error syncing logger", zap.Error(err))
		}
	}(logger)

	// Initialize database
	gormClient, err := database.NewGormClient()
	if err != nil {
		logger.Fatal("Error initializing database", zap.Error(err))
	}

	// Initialize repositories
	globalRepository := repositories.NewGlobalRepository(gormClient)

	// Initialize services
	service := services.New(logger, globalRepository)

	// Initialize handlers
	routing := controllers.NewRouter(logger, service)
	err = routing.Run(":" + viper.GetString("PORT"))
	if err != nil {
		logger.Fatal("Error running router", zap.Error(err))
	}
}

func initViper() {
	// Load environment variables
	viper.SetConfigFile(".env")
	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("error config file: %w", err))
	}
	if viper.GetString("HTTP_PROXY") != "" {
		os.Setenv("HTTP_PROXY", viper.GetString("HTTP_PROXY"))
	}
	if viper.GetString("HTTPS_PROXY") != "" {
		os.Setenv("HTTPS_PROXY", viper.GetString("HTTPS_PROXY"))
	}
}
