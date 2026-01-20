package main

import (
	"gin-learning/internal/logger"
	"gin-learning/internal/models"
	"gin-learning/pkg/database"
	"log/slog"
	"time"
)

func main() {
	logger.SetupLogger()
	database.Connect()

	err := database.DB.AutoMigrate(&models.User{})
	if err != nil {
		slog.Warn("Failed to migrate database: ", "Error", err)
	} else {
		slog.Info("Migration of Users DB is Done")
		seeder()
	}
}

func seeder() {
	var New_user = models.User{
		FirstName: "Ali",
		LastName:  "Fateh",
		Age:       time.Date(1986, 6, 6, 0, 0, 0, 0, time.UTC),
		Email:     "ali@fateh.com",
		Password:  "123456",
		Level:     "Admin",
	}
	result := database.DB.Create(&New_user)
	if result.Error != nil {
		slog.Warn("Can NOT Seed", "Error", result.Error)
	} else {
		slog.Info("First User is seeded")
	}
}
