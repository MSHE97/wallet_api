package logger

import (
	"log"
	"os"
	"wallet/utils"

	"gopkg.in/natefinch/lumberjack.v2"
)

var File *log.Logger

func Init() {
	logfile, err := os.OpenFile(utils.Settings.ApiParams.LogFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	File = log.New(logfile, "", log.Ldate|log.Ltime)

	lumberLog := &lumberjack.Logger{
		Filename:  utils.Settings.ApiParams.LogFile,
		LocalTime: true,
	}
	File.SetOutput(lumberLog)
	File.Println("[LOGGER] initialized")
}
