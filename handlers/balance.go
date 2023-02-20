package handlers

import (
	"wallet/logger"
	"wallet/models"

	"github.com/gin-gonic/gin"
)

func checkBalance(c *gin.Context) {
	var (
		req  models.CheckRequest
		resp models.Response
		user models.Users
		acc  models.Accounts
	)
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.File.Println("	[WARN] binding error ", err)
		c.JSON(resp.BadRequest(err))
		return
	}

	if err := acc.GetByPhone(req.Phone); err != nil {
		logger.File.Printf("	ACCOUNT: %v not found. %v ", req.Phone, err)
		c.JSON(resp.NotFound(models.ErrAccNotFound))
		return
	}
	if err := user.GetByUUID(acc.UserUuid); err != nil {
		logger.File.Printf("	ACCOUNT: %v user not found. %v ", req.Phone, err)
		c.JSON(resp.NotFound(models.ErrInactiveUser))
		return
	}
	xUser := c.MustGet("userId").(int)
	if xUser != user.ID {
		c.JSON(resp.NotAllowed(models.ErrNotAllow))
		return
	}
	c.JSON(resp.Balance(acc.Balance))
}
