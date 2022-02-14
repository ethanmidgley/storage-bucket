package controllers

import (
	"fmt"

	"github.com/ethanmidgley/storage-bucket/pkg/storage"
	"github.com/gin-gonic/gin"
)

// DeleteRequest is the stucture which should be passed to the delete controller in the form of an json object
type DeleteRequest struct {
	Key string `json:"key"`
}

// Delete controller gets key from request and then performs a series of checks before deleting
func Delete(c *gin.Context) {
	var req DeleteRequest

	err := c.BindJSON(&req)
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{"message": "a file key is needed to delete"})
		return
	}

	if req.Key == "" {
		c.AbortWithStatusJSON(400, gin.H{"message": "a file key is needed to delete"})
		return
	}

	if !storage.FileExist(req.Key) {
		c.AbortWithStatusJSON(404, gin.H{"message": "no file found with the key provided"})
		return
	}

	err = storage.Delete(req.Key)
	if err != nil {
		c.AbortWithStatusJSON(500, gin.H{"message": "an error occured whilst trying to delete the file"})
		fmt.Println(err.Error())
		return
	}

	c.JSON(200, gin.H{"message": "file deleted from storage"})

}
