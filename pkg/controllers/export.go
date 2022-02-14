package controllers

import (
	"fmt"

	"github.com/ethanmidgley/storage-bucket/pkg/config"
	"github.com/ethanmidgley/storage-bucket/pkg/storage"
	"github.com/gin-gonic/gin"
)

// Export controller will trigger an export of the bucket to the server root directory
func Export(c *gin.Context) {

	if !storage.CheckBucket() {
		c.JSON(404, gin.H{
			"message": "no bucket detected to export",
		})
		return
	}

	err := storage.Export()
	if err != nil {
		c.JSON(500, gin.H{"message": "failed to export bucket", "error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": fmt.Sprintf("exported bucket to %s.%s", config.Conf.Yaml.Bucket.Location, config.Conf.Yaml.Bucket.Export.Compression)})
}
