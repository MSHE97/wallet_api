package main

import (
	"log"
	"wallet/db"
	"wallet/logger"
	"wallet/models"
	"wallet/utils"
)

func main() {
	log.Println("[MAIN] Work has started")
	utils.ReadConfigs()
	logger.Init()

	log.Print("[MAIN] Datebase connecting...")
	db.Connect()
	logger.File.Println("[MAIN] Automigration...")
	err := autoMigrateDB()
	if err != nil {
		logger.File.Println("	[AUTOMIGRATE] ", err)
	}
}

func autoMigrateDB() error {
	err := db.GetConn().AutoMigrate(
		models.Accounts{},
		models.Payments{},
		models.Users{},
		models.Sessions{},
	).Error
	return err
}
