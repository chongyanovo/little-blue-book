package middleware

import (
	"encoding/gob"
	"github.com/ChongYanOvO/little-blue-book/internal/web"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"log"
	"net/http"
	"strings"
	"time"
)

// LoginBuilder JWT 登录校验
type LoginBuilder struct {
	paths []string
}

func NewLoginBuilder() *LoginBuilder {
	return &LoginBuilder{}
}

func (l *LoginBuilder) IgnorePaths(path string) *LoginBuilder {
	l.paths = append(l.paths, path)
	return l
}

func (l *LoginBuilder) Build() gin.HandlerFunc {
	gob.Register(time.Time{})
	return func(ctx *gin.Context) {
		for _, path := range l.paths {
			if ctx.Request.URL.Path == path {
				return
			}
		}

		// 不需要登录校验的
		if ctx.Request.URL.Path == "/users/login" || ctx.Request.URL.Path == "/users/signup" {
			return
		}

		// 我现在使用 JWT 来校验
		tokenHeader := ctx.GetHeader("x-jwt-token")
		if tokenHeader == "" {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		segs := strings.SplitN(tokenHeader, " ", 2)
		if len(segs) != 2 {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		tokenStr := segs[1]
		userClaims := &web.UserClaims{}
		token, err := jwt.ParseWithClaims(tokenStr, userClaims, func(token *jwt.Token) (interface{}, error) {
			return []byte("secret"), nil
		})
		if err != nil {
			// 没登陆 Bearer
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		// err 为 nil token 不为 nil
		if token == nil || !token.Valid || userClaims.UserId == 0 {
			// 没登陆
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		now := time.Now()
		if userClaims.ExpiresAt.Sub(now) < time.Minute {
			userClaims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(time.Minute))
			tokenStr, err = token.SignedString([]byte("secret"))
			if err != nil {
				// 记录日志
				log.Println("jwt 续约失败", err)
			}
			ctx.Header("x-jwt-token", tokenStr)
		}
		ctx.Set("userClaims", userClaims)
	}
}
