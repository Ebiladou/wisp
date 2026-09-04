package middleware

import (
	"context"
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type contextKey string

const requestIDContextKey contextKey = "request_id"

const RequestIDKey = "request_id"

func LoggingMiddleware(logger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// with a reverse proxy, a request id might be provided in the request, so we get it from the header context, if not available, we just generate a new one here. this is only after push to prod and the infra is set up.

		requestID := uuid.NewString()

		c.Set(RequestIDKey, requestID)

		ctx := context.WithValue(
			c.Request.Context(),
			requestIDContextKey,
			requestID,
		)

		c.Request = c.Request.WithContext(ctx)

		c.Header("X-Request-ID", requestID)
		c.Next()

		logger.Info(
			"http request",
			"request_id", requestID,
			"method", c.Request.Method,
			"path", c.FullPath(),
			"status", c.Writer.Status(),
			"duration_ms", time.Since(start).Milliseconds(),
			"client_ip", c.ClientIP(),
		)
	}
}

func RequestID(ctx context.Context) string {
	requestID, ok := ctx.Value(requestIDContextKey).(string)
	if !ok {
		return ""
	}

	return requestID
}
