package handlers

import (
	"errors"
	"wallet/logger"
	"wallet/models"

	"github.com/gin-gonic/gin"
)

var ErrAccInactive = errors.New("account inactive")

func check4Exist(c *gin.Context) {
	var (
		req     models.CheckRequest
		resp    models.Response
		account models.Accounts
	)

	if err := c.ShouldBindJSON(&req); err != nil {
		logger.File.Println("[CHECK] request handling ", err)
		c.JSON(resp.BadRequest(err))
		return
	}

	if err := account.GetByPhone(req.Phone); err != nil || account.ID == 0 {
		c.JSON(resp.NotFound(err))
		return
	}

	if err := models.CheckUser(account); err != nil {
		c.JSON(resp.Inactive(err))
		return
	}

	c.JSON(resp.Found(account.ID, account.UserUuid))
}
