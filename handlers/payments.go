package handlers

import (
	"wallet/db"
	"wallet/logger"
	"wallet/models"

	"github.com/gin-gonic/gin"
)

func cashIn(c *gin.Context) {
	var (
		req  models.PaymentRequest
		resp models.Response
	)
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.File.Println("[CASH-IN] request handling ", err)
		c.JSON(resp.BadRequest(err))
		return
	}

	if req.Amount <= 0 {
		c.JSON(resp.BadRequest(models.ErrWrongAmount))
		return
	}

	// check sender account
	xUserId := c.MustGet("userId").(int)
	if err := req.PreCheckSender(xUserId); err != nil {
		c.JSON(resp.BadRequest(models.ErrUserNotFound))
		return
	}

	// check receiver account
	if err := req.PreCheckReceiver(); err != nil {
		switch err {
		case models.ErrOutOfLimit:
			c.JSON(resp.NotFound(err))
			return
		default:
			c.JSON(resp.NotAllowed(models.ErrWrongReceiveAcc))
			return
		}
	}

	tx := db.GetConn().Begin()
	defer tx.Commit()
	// create payment

	transaction := models.Payments{
		FromUser:  xUserId,
		ToAccount: req.ToAccount,
		Amount:    req.Amount,
		Type:      req.Type,
		State:     models.PaymentStatusInSaved,
	}
	if err := transaction.Create(tx); err != nil || transaction.ID == 0 {
		c.JSON(resp.TransactionError(err))
		tx.Rollback()
		return
	}
	if err := transaction.RunProcessing(tx); err != nil {
		c.JSON(resp.TransactionError(err))
		tx.Rollback()
		return
	}

}
