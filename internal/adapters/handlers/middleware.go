package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"opc_ua_service/internal/middleware/logging"
	"time"
)

type RequestInfo struct {
	RemoteAddr string              `json:"remote_addr"`
	Method     string              `json:"method"`
	Path       string              `json:"path"`
	Headers    map[string][]string `json:"headers"`
}

func LoggingMiddleware(parentLogger *logging.Logger) gin.HandlerFunc {
	logger := parentLogger.WithPrefix("HTTP")

	return func(c *gin.Context) {
		if c.Request.Method == http.MethodOptions {
			c.Next()
			return
		}

		logger.Info("Request started",
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"remote_addr", c.Request.RemoteAddr,
		)

		start := time.Now()

		c.Next()

		status := c.Writer.Status()
		logger.Info("Request completed",
			"status", status,
			"latency", time.Since(start),
			"client_ip", c.ClientIP(),
		)
	}
}
