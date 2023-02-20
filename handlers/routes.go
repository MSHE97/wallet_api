package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Init() *gin.Engine {
	router := gin.Default()

	auth := router.Group("")
	auth.Use(ApiAuth())
	{
		auth.POST("/ping", ping)
		auth.POST("/account/check", check4Exist)
		auth.POST("/cash_in", cashIn)
		auth.POST("/totals", totals)
		auth.POST("/balance", checkBalance)
	}
	router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{"error": "Page not found"})
	})
	return router
}

func ping(c *gin.Context) {
	c.JSON(200, "pong")
}
