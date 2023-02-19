package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"wallet/db"
	"wallet/logger"
	"wallet/models"
	"wallet/server"
	"wallet/utils"
)

func main() {
	log.Println("[MAIN] Work has started")
	utils.ReadConfigs()
	logger.Init()

	log.Println("[MAIN] Datebase connecting...")
	db.Connect()
	logger.File.Println("[MAIN] Automigration...")
	err := autoMigrateDB()
	if err != nil {
		logger.File.Println("	[AUTOMIGRATE] ", err)
	}
	defer Stop()
	logger.File.Println("	[GIN] server init...")
	server := server.InitServer()
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.File.Fatalln("listen: ", err)
		}
	}()
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.File.Println("[MAIN] Shutdown Server ...")
	log.Println("[MAIN] Shutdown Server ...")
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	if err := server.Shutdown(ctx); err != nil {
		logger.File.Fatalln("Server Shutdown: ", err)
	}
	// catching ctx.Done(). timeout of 5 seconds.
	<-ctx.Done()
	cancel()
}

func Stop() {
	logger.File.Println("[MAIN] Work has stopped!")
	db.Close()
	logger.File.Print("_____________________________________________________________________________________\n")
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
