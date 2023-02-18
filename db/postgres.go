package db

import (
	"fmt"
	"log"
	"strconv"
	"wallet/logger"
	"wallet/utils"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var database *gorm.DB

func initDB() *gorm.DB {
	params := utils.Sets.PostgresParams

	connString := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable",
		params.Server, strconv.Itoa(params.Port),
		params.User, params.Schema,
		params.Password)

	db, err := gorm.Open("postgres", connString)
	if err != nil {
		log.Printf("\n[DB] Coudn't connect to postgresql database . !! %v !!", connString)
		log.Panic(err)
	}
	//logger.File.Println("PostgreSQL connection string: ", connString)
	db.LogMode(true)
	db.SingularTable(true)
	return db
}

// Creates connection to database
func Connect() {
	database = initDB()
}

func GetConn() *gorm.DB {
	return database
}

func Close() {
	err := database.Close()
	if err != nil {
		logger.File.Println("[DB] Error on closing local DB: ", err)
	}
}
