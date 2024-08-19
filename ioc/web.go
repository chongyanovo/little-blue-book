package ioc

import (
	"github.com/ChongYanOvO/little-blue-book/internal/web"
	"github.com/ChongYanOvO/little-blue-book/internal/web/middleware"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"strings"
	"time"
)

func InitGin(middlewares []gin.HandlerFunc, uh *web.UserHandler) *gin.Engine {
	server := gin.Default()
	server.Use(middlewares...)
	uh.RegisterRoutes(server)
	return server
}

func InitMiddlewares() []gin.HandlerFunc {
	return []gin.HandlerFunc{
		InitCorsMiddleware(),
		InitLoginMiddleware(),
	}
}

func InitCorsMiddleware() gin.HandlerFunc {
	return cors.New(cors.Config{
		AllowAllOrigins: false,
		AllowOrigins:    []string{"http://localhost:3000"},
		// 在使用 JWT 的时候，因为我们使用了 Authorizaition 的头部，所以需要加上
		AllowHeaders: []string{"Content-Type", "Authorization"},
		// 为了 JWT
		ExposeHeaders:    []string{"x-jwt-token", "Authorization"},
		AllowMethods:     []string{"POST", "GET", "PUT"},
		AllowCredentials: true,
		// 你不加这个 前端是拿不到的
		AllowOriginFunc: func(origin string) bool {
			if strings.HasPrefix(origin, "http://localhost") {
				return true
			}
			return strings.Contains(origin, "abc")
		},
		MaxAge: 12 * time.Hour,
	})
}

func InitLoginMiddleware() gin.HandlerFunc {
	return middleware.NewLoginBuilder().
		IgnorePaths("/users/signup").
		IgnorePaths("/users/login").
		IgnorePaths("/users/login/code").Build()
}
