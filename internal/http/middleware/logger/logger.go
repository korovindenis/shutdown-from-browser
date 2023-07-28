package logger

import (
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func LoggingMiddleware(log *zap.Logger) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		startTime := time.Now()

		// Processing request
		ctx.Next()

		endTime := time.Now()

		loggerWithContext := log.With(
			zap.Any("HTTP REQUEST", struct {
				REQUEST_ID string
				METHOD     string
				URI        string
				STATUS     int
				LATENCY    time.Duration
				CLIENT_IP  string
			}{
				ctx.GetString("requestId"),
				ctx.Request.Method,
				ctx.Request.RequestURI,
				ctx.Writer.Status(),
				endTime.Sub(startTime),
				ctx.ClientIP(),
			}),
		)

		loggerWithContext.Info("Logging with context")

		ctx.Next()
	}
}
