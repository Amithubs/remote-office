package main

import (
	"github.com/sirupsen/logrus"
	"github.com/tejashwikalptaru/remote-office/database"
	"github.com/tejashwikalptaru/remote-office/server"
)

//todo: load variables like db settings from env file https://github.com/joho/godotenv
func main() {
	// connect with database
	if err := database.ConnectAndMigrate("localhost", "5432", "postgres", "postgres", "", database.SSLModeDisable); err != nil {
		logrus.Panicf("Failed to initialize and migrate database with error: %+v", err)
	}

	// create server instance
	srv := server.SetupRoutes()
	if err := srv.Run(":8080"); err != nil {
		logrus.Panicf("Failed to run server with error: %+v", err)
	}
}
