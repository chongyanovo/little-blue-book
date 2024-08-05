package main

import (
	"github.com/ChongYanOvO/little-blue-book/config"
	"github.com/ChongYanOvO/little-blue-book/internal/repository"
	"github.com/ChongYanOvO/little-blue-book/internal/repository/cache"
	"github.com/ChongYanOvO/little-blue-book/internal/repository/dao"
	"github.com/ChongYanOvO/little-blue-book/internal/service"
	"github.com/ChongYanOvO/little-blue-book/internal/service/sms/memory"
	"github.com/ChongYanOvO/little-blue-book/internal/web"
	"github.com/ChongYanOvO/little-blue-book/internal/web/middleware"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"strings"
	"time"
)

func main() {
	db := initDB()
	rdb := initRDB()
	server := initWebServer()

	u := initUser(db, rdb)
	u.RegisterRoutes(server)

	server.Run(":8088")
}

func initWebServer() *gin.Engine {
	server := gin.Default()

	server.Use(cors.New(cors.Config{
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
	}))

	server.Use(middleware.NewLoginBuilder().
		IgnorePaths("/users/signup").
		IgnorePaths("/users/login").
		IgnorePaths("/users/login/code").Build())
	return server
}

func initUser(db *gorm.DB, rdb *redis.Client) *web.UserHandler {
	ud := dao.NewUserDao(db)
	uc := cache.NewUserCache(rdb, 30*time.Minute)
	cc := cache.NewCodeCache(rdb)
	repo := repository.NewUserRepository(ud, uc)
	codeRepo := repository.NewCodeRepository(cc)
	svc := service.NewUserService(repo)
	smsSvc := memory.NewService()
	codeSvc := service.NewCodeService(codeRepo, smsSvc)
	u := web.NewUserHandler(svc, codeSvc)
	return u
}

func initDB() *gorm.DB {
	db, err := gorm.Open(mysql.Open(config.Config.DB.DSN))
	if err != nil {
		panic(err)
	}
	err = dao.InitTable(db)
	if err != nil {
		panic(err)
	}
	return db.Debug()
}

func initRDB() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr: config.Config.Redis.Addr,
	})
}
