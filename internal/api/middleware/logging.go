package middleware

import (
	"crypto/rand"
	"fmt"
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/techcollectiveassam/flood-response-api/internal/pkg/logging"
)

func LoggingMiddleware(logger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		requestID := generateRequestID()

		c.Set("request_id", requestID)
		c.Header("X-Request-Id", requestID)

		loggerWithID := logger.With("request_id", requestID)
		c.Request = c.Request.WithContext(logging.WithContext(c.Request.Context(), loggerWithID))

		c.Next()

		attrs := []any{
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"status", c.Writer.Status(),
			"latency_ms", time.Since(start).Milliseconds(),
			"client_ip", c.ClientIP(),
		}

		switch {
		case c.Writer.Status() >= 500:
			loggerWithID.Error("request", attrs...)
		case c.Writer.Status() >= 400:
			loggerWithID.Warn("request", attrs...)
		default:
			loggerWithID.Info("request", attrs...)
		}
	}
}

func generateRequestID() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
}
