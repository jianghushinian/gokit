package main

import (
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	log "github.com/jianghushinian/gokit/log/zap"
)

// LoggerMiddleware ref: https://github.com/gin-gonic/gin/blob/v1.9.0/logger.go#L182
func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Start timer
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		// Process request
		c.Next()

		// Stop timer
		stop := time.Now()
		latency := stop.Sub(start)
		if latency > time.Minute {
			latency = latency.Truncate(time.Second)
		}

		log.Info("GIN request",
			log.String("start", start.Format(time.RFC3339)),
			log.Int("status", c.Writer.Status()),
			log.String("latency", fmt.Sprintf("%s", latency)),
			log.String("method", c.Request.Method),
			log.String("path", path),
			log.String("query", raw),
			log.String("clientIP", c.ClientIP()),
			log.String("userAgent", c.Request.UserAgent()),
			log.String("error", c.Errors.ByType(gin.ErrorTypePrivate).String()),
		)
	}
}

// RecoveryMiddleware ref: https://github.com/gin-gonic/gin/blob/v1.9.0/recovery.go#L33
func RecoveryMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// Check for a broken connection, as it is not really a
				// condition that warrants a panic stack trace.
				var brokenPipe bool
				if ne, ok := err.(*net.OpError); ok {
					var se *os.SyscallError
					if errors.As(ne, &se) {
						seStr := strings.ToLower(se.Error())
						if strings.Contains(seStr, "broken pipe") ||
							strings.Contains(seStr, "connection reset by peer") {
							brokenPipe = true
						}
					}
				}
				httpRequest, _ := httputil.DumpRequest(c.Request, false)
				headers := strings.Split(string(httpRequest), "\r\n")
				for idx, header := range headers {
					current := strings.Split(header, ":")
					if current[0] == "Authorization" {
						headers[idx] = current[0] + ": *"
					}
				}
				headersToStr := strings.Join(headers, "\r\n")
				if brokenPipe {
					log.Error(c.Request.URL.String(),
						log.Any("err", err),
						log.String("headers", headersToStr),
						log.Stack("stack"),
					)
					// If the connection is dead, we can't write a status to it.
					c.Error(err.(error)) // nolint: errcheck
					c.Abort()
				} else {
					log.Error(c.Request.URL.String(),
						log.Any("err", err),
						log.String("headers", headersToStr),
						log.Stack("stack"),
						log.String("panicRecoveredTime", time.Now().Format(time.RFC3339)),
					)
					c.AbortWithStatus(http.StatusInternalServerError)
				}
			}
		}()
		c.Next()
	}
}

func init() {
	// custom logger
	// logger := log.New(os.Stderr, log.InfoLevel, log.AddCallerSkip(2))
	// log.ReplaceDefault(logger)
}
