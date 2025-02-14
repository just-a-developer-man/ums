package middleware

import (
	"context"
	"log/slog"
	"time"
	"ums/internal/logger"

	"github.com/gin-gonic/gin"
)

func AdminAuth() gin.HandlerFunc {
	slog.Debug("Registering admin authentication middleware")
	return func(c *gin.Context) {
		slog.DebugContext(c.Request.Context(), "Perfoming admin authentication")
		c.Next()
	}
}

func Logging() gin.HandlerFunc {
	slog.Debug("Registering logging middleware")
	return func(c *gin.Context) {
		ctx := enrichContextWithRequestInfo(c)
		c.Request = c.Request.Clone(ctx)

		startTime := time.Now()
		c.Next()

		duration := time.Since(startTime)
		ctx = enrichContextWithAfterRequestInfo(c.Request.Context(), c.Writer.Status(), duration)
		c.Request = c.Request.Clone(ctx)

		slog.InfoContext(ctx, "Processed request")
	}
}

// enrichContextWithRequestInfo adds request data to the context.
func enrichContextWithRequestInfo(c *gin.Context) context.Context {
	ctx := c.Request.Context()
	ctx = logger.WithCallerIP(ctx, c.ClientIP())
	ctx = logger.WithEndpoint(ctx, c.FullPath())
	ctx = logger.WithMethod(ctx, c.Request.Method)
	return ctx
}

// enrichContextWithAfterRequestInfo adds request processing duration and response status to the context.
func enrichContextWithAfterRequestInfo(ctx context.Context, status int, duration time.Duration) context.Context {
	ctx = logger.WithStatusCode(ctx, status)
	return logger.WithDuration(ctx, duration)
}
