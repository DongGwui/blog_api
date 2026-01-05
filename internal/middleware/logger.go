package middleware

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/ydonggwui/blog-api/internal/pkg/logger"
)

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		// Generate request ID
		requestID := uuid.New().String()
		c.Set("request_id", requestID)

		// Add request ID to context
		ctx := context.WithValue(c.Request.Context(), logger.RequestIDKey, requestID)
		c.Request = c.Request.WithContext(ctx)

		c.Next()

		latency := time.Since(start)
		status := c.Writer.Status()

		if query != "" {
			path = path + "?" + query
		}

		// Log request
		logger.Info(ctx, "HTTP Request",
			"status", status,
			"method", c.Request.Method,
			"path", path,
			"client_ip", c.ClientIP(),
			"latency_ms", latency.Milliseconds(),
		)
	}
}
