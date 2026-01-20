package main

import (
	"Tendabox/internal/logger"
	//"Tendabox/internal/routes"
	"Tendabox/pkg/database"
)

func main() {
	logger.SetupLogger()
	database.Connect()
	//r := routes.SetupRouter()
	//r.Run()

}
