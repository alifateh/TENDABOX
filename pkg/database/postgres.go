package database

import (
	"log/slog"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
	dsn := "host=localhost user=root password=Xper!@DB1885 dbname=mydatabase port=5432 sslmode=disable"
	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		//Logger: logger.Default.LogMode(logger.Warn),
	})

	if err != nil {
		slog.Warn("Error in Connecting", "Fatal_Error", err)
	}
	slog.Info("DB Connectio is Established")
}
