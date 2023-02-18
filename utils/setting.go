package utils

import (
	"encoding/json"
	"log"
	"os"
	"wallet/models"
)

var Settings models.Settings

// ReadConfigs to init api settings
func ReadConfigs() {

	doc, err := os.ReadFile("files/conf.json")
	if err != nil {
		log.Println("Fail in conf file reading. ", err)
		panic(err)
	}

	err = json.Unmarshal(doc, &Settings)
	if err != nil {
		log.Println(err)
		panic(err.Error())
	}
}
