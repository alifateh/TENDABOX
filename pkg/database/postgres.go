package database

import (
	"fmt"
	"log/slog"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
	//dbHost := os.Getenv("DB_HOST")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbPort := os.Getenv("DB_PORT")
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", "postgres", dbUser, dbPassword, dbName, dbPort)
	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		//Logger: logger.Default.LogMode(logger.Warn),
		//Logger: logger.Default.LogMode(logger.Silent),
	})

	if err != nil {
		slog.Warn("Error in Connecting", "Fatal_Error", err)
		os.Exit(1)
	} else {
		slog.Info("DB Connectio is Established")
	}
}
