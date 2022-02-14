package controllers

import (
	"github.com/ethanmidgley/storage-bucket/pkg/auth"
	"github.com/ethanmidgley/storage-bucket/pkg/config"
	"github.com/gin-gonic/gin"
)

// CreateKeys is the controller to automatically generate keys providing keys have not already been generated
func CreateKeys(c *gin.Context) {
	if len(config.Conf.Yaml.ControlPlane.Keys) == config.Conf.Yaml.ControlPlane.NumberOfKeys {
		c.AbortWithStatusJSON(401, gin.H{"message": "keys already generated"})
		return
	}

	keys, keyshashed := auth.GenerateKeys()
	config.Conf.Yaml.ControlPlane.Keys = keyshashed
	config.Conf.Update()

	c.JSON(200, gin.H{"message": "keys generated successfully", "keys": keys})
	return
}
