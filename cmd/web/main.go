package main

import (
	"gin-learning/internal/logger"
	"gin-learning/internal/routes"
	"gin-learning/pkg/database"
)

func main() {
	logger.SetupLogger()
	database.Connect()
	r := routes.SetupRouter()
	r.Run(":9595")

}
