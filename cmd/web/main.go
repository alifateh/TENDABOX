package main

import (
	"Tendabox/internal/logger"
	"Tendabox/internal/routes"
	"Tendabox/pkg/database"
	"log/slog"
)

func main() {
	logger.SetupLogger()
	database.Connect()
	r := routes.SetupRouter()
	r.Run()
	slog.Info("just for fun")
}
