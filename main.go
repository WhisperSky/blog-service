package main

import (
	"blog-service/global"
	"blog-service/internal/model"
	"blog-service/internal/routers"
	"blog-service/pkg/logger"
	"blog-service/pkg/setting"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	val "github.com/go-playground/validator/v10"
	zhTranslations "github.com/go-playground/validator/v10/translations/zh"
	"gopkg.in/natefinch/lumberjack.v2"
	"log"
	"net/http"
	"time"
)

func init() {
	// 注册配置文件
	err := setupSetting()
	if err != nil {
		log.Fatalf("init.setupSetting err: %v", err)
	}

	// 注册数据库引擎
	err = setupDBEngine()
	if err != nil {
		log.Fatalf("init.setupDBEngine err: %v", err)
	}

	// 注册日志记录
	err = setupLogger()
	if err != nil {
		log.Fatalf("init.setupLogger err: %v", err)
	}

	// 注册国际化翻译处理
	err = setupTranslator()
	if err != nil {
		log.Fatalf("init.setupTranslator err: %v", err)
	}
}

// @title 博客系统
// @version 1.0
// @description Go 语言编程之旅：一起用 Go 做项目
// @termsOfService https://github.com/go-programming-tour-book
func main() {
	// 初级web监听demo
	//r := gin.Default()
	//r.GET("/ping", func(c *gin.Context) {
	//	c.JSON(200, gin.H{"message": "pong"})
	//})
	//r.Run()
	global.Logger.Infof("%s: go-programming-tour-book/%s", "eddycjy", "blog-service")

	// 设置运行方式
	gin.SetMode(global.ServerSetting.RunMode)

	// 设置路由
	router := routers.NewRouter()
	s := &http.Server{
		Addr:           ":" + global.ServerSetting.HttpPort,
		Handler:        router,
		ReadTimeout:    global.ServerSetting.ReadTimeout,
		WriteTimeout:   global.ServerSetting.WriteTimeout,
		MaxHeaderBytes: 1 << 20,
	}
	s.ListenAndServe()

}

// setupSetting 注册配置文件
func setupSetting() error {
	settingss, err := setting.NewSetting()
	if err != nil {
		return err
	}
	err = settingss.ReadSection("Server", &global.ServerSetting)
	if err != nil {
		return err
	}
	err = settingss.ReadSection("App", &global.AppSetting)
	if err != nil {
		return err
	}
	err = settingss.ReadSection("Database", &global.DatabaseSetting)
	if err != nil {
		return err
	}
	err = settingss.ReadSection("JWT", &global.JWTSetting)
	if err != nil {
		return err
	}

	global.ServerSetting.ReadTimeout *= time.Second
	global.ServerSetting.WriteTimeout *= time.Second
	global.JWTSetting.Expire *= time.Second
	return nil
}

// setupDBEngine 注册数据库引擎
func setupDBEngine() error {
	var err error
	global.DBEngine, err = model.NewDBEngine(global.DatabaseSetting)
	if err != nil {
		return err
	}

	return nil
}

// setupLogger 注册日志
func setupLogger() error {
	global.Logger = logger.NewLogger(&lumberjack.Logger{
		Filename: global.AppSetting.LogSavePath + "/" + global.AppSetting.LogFileName + global.AppSetting.LogFileExt,
		// 设置日志文件所允许的最大占用空间为 600MB
		MaxSize: 600,
		// 日志文件最大生存周期为 10 天
		MaxAge: 10,
		// 设置日志文件名的时间格式为本地时间
		LocalTime: true,
	}, "", log.LstdFlags).WithCaller(2)

	return nil
}

///*
//	以下init问题为解决如下问题：
//    https://golang2.eddycjy.com/posts/ch2/05-validator/
//	如果 RegisterDefaultTranslations 放在中间件里调用，如果多个用户同时访问的时候，后端很容易就挂掉，
//	出现 panic:concurrent map read and map write 并且无法 recovery，使用 ab 做压力测试可复现这个问题。
//*/
func setupTranslator() error {
	uni := ut.New(zh.New())
	global.Trans, _ = uni.GetTranslator("zh")
	v, ok := binding.Validator.Engine().(*val.Validate)
	if ok {
		err := zhTranslations.RegisterDefaultTranslations(v, global.Trans)
		if err != nil {
			return err
		}
	}

	return nil
}
