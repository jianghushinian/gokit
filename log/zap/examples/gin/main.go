package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	log "github.com/jianghushinian/gokit/log/zap"
)

func main() {
	// r := gin.Default()
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(LoggerMiddleware(), RecoveryMiddleware())

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	port := ":8000"
	log.Info(fmt.Sprintf("Listening and serving HTTP on %s", port))
	r.Run(port)
}
