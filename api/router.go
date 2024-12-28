package api

import (
	"github.com/gin-gonic/gin"
)

func BuildRouter() *gin.Engine {
	router := gin.Default()
	// Set a lower memory limit for multipart forms (default is 32 MiB)
	router.MaxMultipartMemory = 8 << 20 // 8 MiB
	router.POST("/files/upload", uploadFile)
	router.GET("/files/:filename", downloadFile)
	return router
}
