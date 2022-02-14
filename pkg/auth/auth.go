package auth

import (
	"crypto/sha256"
	"encoding/hex"

	"github.com/ethanmidgley/storage-bucket/pkg/config"
	"github.com/gin-gonic/gin"
)

// IsAuthenticated middleware to make sure all routes are protected
func IsAuthenticated() gin.HandlerFunc {
	return func(c *gin.Context) {
		keys := c.Request.Header["Api-Key"]

		if len(keys) == 0 {
			c.AbortWithStatusJSON(401, gin.H{"message": "authentication needed"})
			return
		}
		h := sha256.New()
		h.Write([]byte(keys[0]))

		if _, ok := config.Conf.KeysMap[hex.EncodeToString(h.Sum(nil))]; !ok {
			c.AbortWithStatusJSON(401, gin.H{"message": "invalid credentials"})
			return
		}
		c.Next()
	}
}
