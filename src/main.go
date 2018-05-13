package main

import (
	"github.com/yspk/frame/src/common/cache"
	"github.com/yspk/frame/src/common/constant"
	"github.com/yspk/frame/src/common/logger"
	"github.com/yspk/frame/src/common/middleware"
	"github.com/yspk/frame/src/config"
	"github.com/yspk/frame/src/controller/back"
	"github.com/yspk/frame/src/controller/front"
	"github.com/yspk/frame/src/model"
	"flag"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/yspk/frame/src/common/sensitive"
)

func main() {
	configPath := flag.String("conf", "config/config.json", "Config file path")
	flag.Parse()

	err := config.LoadConfig(*configPath)
	if err != nil {
		logger.Fatal("Config Failed!", err)
		return
	}

	logger.SetLevel(config.GetLoggerLevel())

	//TODO
	filter := sensitive.New()
	filter.LoadWordDict("/usr/local/gopath/src/github.com/yspk/frame/src/common/sensitive/dict/dict.txt")
	logger.Fatal(filter.Replace("静静是色魔",42))


	db, err := gorm.Open(config.GetDBName(), config.GetDBSource())
	if err != nil {
		logger.Fatal("Open db Failed!!!!", err)
		return
	}

	if err := cache.RedisTest(config.GetCacheRedisAddr(), config.GetCacheRedisPassword()); err != nil {
		logger.Fatal(err)
		return
	}

	cache.InitCache(config.GetCacheRedisAddr(), config.GetCacheRedisPassword())

	model.InitModel(db)

	r := gin.New()
	gin.SetMode(gin.DebugMode)

	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.OPTIONS("*f", func(c *gin.Context) {})

	dbMiddleware := middleware.GetDbPrepareHandler(config.GetDBName(), config.GetDBSource(), config.IsOrmLogEnabled(), constant.ContextDb)
	jwtMiddleware := middleware.JWTAuth()

	cms := r.Group("back")
	cms.Use(dbMiddleware, jwtMiddleware)
	{
		cms.GET("/admin/login", back.AdminLoginHandler)
	}

	ims := r.Group("front")
	ims.Use(dbMiddleware, jwtMiddleware)
	{
		ims.GET("/user/login", front.UserLoginHandler)
		ims.POST("/file/upload", front.FileUploadHandler)
	}

	r.NoRoute(front.FileServeHandler)
	r.Run(":8555")
}
