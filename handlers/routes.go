package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Init() *gin.Engine {
	router := gin.Default()
	router.GET("/ping", ping)
	auth := router.Group("/v1")
	auth.Use(ApiAuth())
	{
		router.POST("/ping")
	}
	router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{"error": "Page not found"})
	})
	return router
}

func ping(c *gin.Context) {
	c.JSON(200, "pong")
}