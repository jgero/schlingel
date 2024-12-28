package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
)

func uploadFile(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	c.SaveUploadedFile(file, fmt.Sprintf("./%s", file.Filename))

	c.String(http.StatusOK, fmt.Sprintf("'%s' uploaded!", file.Filename))
}

func downloadFile(c *gin.Context) {
	filename := c.Param("filename")
	filepath := fmt.Sprintf("./%s", filename)

	if _, err := os.Stat(filepath); err != nil {
		c.String(http.StatusNotFound, fmt.Sprintf("'%s' does not exist.", filename))
		return
	}

	c.File(filepath)

	c.String(http.StatusOK, fmt.Sprintf("'%s' downloaded!", filename))
}
