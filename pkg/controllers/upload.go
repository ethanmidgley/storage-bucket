package controllers

import (
	"fmt"

	"github.com/ethanmidgley/storage-bucket/pkg/config"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Upload this allows the uploading of files to the bucket
func Upload(c *gin.Context) {
	form, _ := c.MultipartForm()
	files := form.File["Uploads[]"]

	var fileKeys []string

	for _, file := range files {
		name := uuid.New().String()
		fileKeys = append(fileKeys, name)
		c.SaveUploadedFile(file, fmt.Sprintf("%s/%s", config.Conf.Yaml.Bucket.Location, name))
	}
	c.JSON(200, gin.H{"message": "uploaded all files", "file-keys": fileKeys})
}
