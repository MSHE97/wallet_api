package handlers

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"io/ioutil"
	"net/http"
	"strconv"
	"wallet/logger"
	"wallet/models"
	"wallet/utils"

	"github.com/gin-gonic/gin"
)

func ApiAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		xUserId := c.Params.ByName("X-UserId")
		xDigest := c.Params.ByName("X-Digest")
		if len(xUserId) < 1 {
			logger.File.Println("[AUTH] request Abortion. empty X-UserId ")
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "empty X-UserId "})
			return
		}
		user := models.Users{}
		userId, err := strconv.Atoi(xUserId)
		if err != nil {
			logger.File.Println("[AUTH] request Abortion. wrong X-UserId ")
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "wrong X-UserId "})
			return
		}

		if err := user.GetByID(int64(userId)); err != nil {
			logger.File.Println("[AUTH] request Abortion. X-UserId not found ")
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "wrong X-UserId "})
			return
		}

		if len(xDigest) < 1 {
			logger.File.Println("[AUTH] request Abortion. empty X-Digest ")
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "empty X-Digest "})
			return
		}

		body, err := ioutil.ReadAll(c.Request.Body)
		if err != nil {
			logger.File.Println("[AUTH] request Abortion. ", err)
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err})
			return
		}

		if GetSignature(string(body), utils.Sets.ApiParams.Secret) != xDigest {
			logger.File.Println("[AUTH] request Abortion. WRONG X-Digest ")
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "wrong X-Digest "})
			return
		}

		c.Next()
	}
}

func GetSignature(input string, key string) string {
	key_for_sign := []byte(key)
	h := hmac.New(sha1.New, key_for_sign)
	h.Write([]byte(input))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}
