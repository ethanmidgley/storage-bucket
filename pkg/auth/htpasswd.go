package auth

import (
	"encoding/base64"
	"strings"

	"github.com/ethanmidgley/storage-bucket/pkg/config"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// ValidateHTPasswd wil compare a username and password to a bcrypt htpasswd
func ValidateHTPasswd(username, password, htpassword string) bool {
	htuser := strings.Split(htpassword, ":")[0]
	hthash := strings.Split(htpassword, ":")[1]

	if htuser != username {
		return false
	}

	err := bcrypt.CompareHashAndPassword([]byte(hthash), []byte(password))
	if err != nil {
		return false
	}
	return true

}

// IsHTAuthenticated will use basic authentication to compare with a htpasswd set in the config
func IsHTAuthenticated() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.Request.Header["Authentication"]

		if len(authHeader) == 0 {
			c.AbortWithStatusJSON(401, gin.H{"message": "authentication needed"})
			return
		}

		splitHeader := strings.Split(authHeader[0], " ")
		if len(splitHeader) != 2 {
			c.AbortWithStatusJSON(401, gin.H{"message": "unsupported authentication type"})
			return
		}

		authType := strings.Split(authHeader[0], " ")[0]
		authData := strings.Split(authHeader[0], " ")[1]

		if authType != "Basic" {
			c.AbortWithStatusJSON(401, gin.H{"message": "unsupported authentication type"})
			return
		}

		userpass, err := base64.StdEncoding.DecodeString(authData)
		if err != nil {
			c.AbortWithStatusJSON(500, gin.H{"message": "internal server error when trying to authenticate"})
			return
		}

		username := strings.Split(string(userpass), ":")[0]
		password := strings.Split(string(userpass), ":")[1]

		if !ValidateHTPasswd(username, password, config.Conf.Yaml.ControlPlane.HTPasswd) {
			c.AbortWithStatusJSON(401, gin.H{"message": "invalid credentials"})
			return
		}

		c.Next()
	}
}
