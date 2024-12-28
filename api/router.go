package api

import (
	"bytes"
	"context"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

type HtmxPartialRenderer = func() (string, func(context.Context, io.Writer, map[string]string) error)

func BuildRouter(pages []HtmxPartialRenderer) *gin.Engine {
	router := gin.Default()
	// Set a lower memory limit for multipart forms (default is 32 MiB)
	router.MaxMultipartMemory = 8 << 20 // 8 MiB
	router.POST("/files/upload", uploadFile)
	router.GET("/files/:filename", downloadFile)

	for _, p := range pages {
		path, r := p()
		router.GET(path, func(c *gin.Context) {
			var buf bytes.Buffer
			err := r(c.Request.Context(), &buf, paramsToMap(c.Params))
			if err != nil {
				panic(err)
			}
			c.Data(http.StatusOK, "text/html; charset=utf-8", buf.Bytes())
		})
	}
	return router
}

func paramsToMap(params gin.Params) map[string]string {
	m := make(map[string]string)
	for _, p := range params {
		m[p.Key] = p.Value
	}
	return m
}

