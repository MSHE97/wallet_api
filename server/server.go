package server

import (
	"fmt"
	"net/http"
	"time"
	"wallet/handlers"
	"wallet/utils"
)

var Server *http.Server

func InitServer() *http.Server {
	routers := handlers.Init()
	endPoint := fmt.Sprintf(":%d", utils.Sets.ApiParams.PortRun)
	Server := &http.Server{
		Addr:           endPoint,
		Handler:        routers,
		ReadTimeout:    time.Duration(30) * time.Second,
		WriteTimeout:   time.Duration(30) * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	return Server
}
