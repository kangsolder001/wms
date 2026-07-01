package middleware

import (
	"time"

	"wms/pkg/logger"

	"github.com/gin-gonic/gin"
)

type LoggingMiddleware struct {
	log logger.Logger
}

func NewLoggingMiddleware(log logger.Logger) *LoggingMiddleware {
	return &LoggingMiddleware{log: log}
}

func (m *LoggingMiddleware) Handle() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		m.log.Info("request",
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"status", c.Writer.Status(),
			"duration", time.Since(start).String(),
			"remote_addr", c.ClientIP(),
		)
	}
}
