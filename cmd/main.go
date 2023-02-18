package main

import (
	"log"
	"wallet/db"
	"wallet/logger"
	"wallet/utils"
)

func main() {
	log.Println("[MAIN] Work has started")
	utils.ReadConfigs()
	logger.Init()

	log.Print("[MAIN] Datebase connecting...")
	db.Connect()
}
