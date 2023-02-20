package handlers

import (
	"io"
	"time"
	"wallet/models"

	"github.com/gin-gonic/gin"
)

func totals(c *gin.Context) {
	var (
		req  models.TotalsRequest
		resp models.Response
		dflt bool
	)

	_, err := io.ReadAll(c.Request.Body)
	if err == io.EOF {
		dflt = true
	}

	if dflt {
		req.From = models.TruncateToMounth(time.Now()).Format(time.RFC3339)
		req.To = models.TruncateToDay(time.Now()).Format(time.RFC3339)
	} else {
		if err := c.ShouldBind(&req); err != nil {
			c.JSON(resp.BadRequest(err))
			return
		}
	}
	resp.Payload, err = req.GetTotals()
	if err != nil {
		c.JSON(resp.InternalError(err))
		return
	}
	c.JSON(resp.Totals(req))
}
