package ginx

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

func WrapperBody[T any](l *zap.Logger, fn func(ctx *gin.Context, req T) (Result, error)) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req T
		if err := ctx.Bind(&req); err != nil {
			return
		}
		res, err := fn(ctx, req)
		if err != nil {
			l.Error("处理业务逻辑错误", zap.Error(err),
				zap.String("path", ctx.Request.URL.String()),
				zap.String("router", ctx.FullPath()),
			)
		}
		ctx.JSON(http.StatusOK, res)
	}
}
func WrapperBodyWitJwt[T any](l *zap.Logger, fn func(ctx *gin.Context, req T, uc *UserClaims) (Result, error)) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req T
		if err := ctx.Bind(&req); err != nil {
			return
		}
		var uc *UserClaims
		uc, err := ExtractJwtClaims(ctx)
		if err != nil {
			l.Error("获取UserClaims错误", zap.Error(err))
			ctx.AbortWithStatus(http.StatusUnauthorized)
		}
		res, err := fn(ctx, req, uc)
		if err != nil {
			l.Error("处理业务逻辑错误", zap.Error(err),
				zap.String("path", ctx.Request.URL.String()),
				zap.String("router", ctx.FullPath()),
			)
			ctx.AbortWithStatus(http.StatusUnauthorized)
		}
		if err != nil {
			l.Error("获取UserClaims错误", zap.Error(err))
		}
		ctx.JSON(http.StatusOK, res)
	}
}
