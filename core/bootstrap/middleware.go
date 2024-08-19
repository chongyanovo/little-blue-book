package bootstrap

import (
	"context"
	"github.com/ChongYanOvO/little-blue-book/pkg/ginx/middleware/logger"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func NewMiddleware(l *zap.Logger) []gin.HandlerFunc {
	return []gin.HandlerFunc{
		NewLoggerMiddleware(l),
	}
}

func NewLoggerMiddleware(l *zap.Logger) gin.HandlerFunc {
	return logger.
		NewBuilder(func(ctx context.Context, log *logger.AccessLog) {
			l.Info("Http请求", zap.Any("日志", log))
		}).
		AllowRequestBody(true).
		AllowResponseBody(true).
		Build()
}
