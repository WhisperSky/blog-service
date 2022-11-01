package routers

import (
	_ "blog-service/docs"
	"blog-service/global"
	"blog-service/internal/middleware"
	"blog-service/internal/routers/api"
	v1 "blog-service/internal/routers/api/v1"
	"blog-service/pkg/limiter"
	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
	"time"
)

var methodLimiters = limiter.NewMethodLimiter().AddBuckets(limiter.LimiterBucketRule{
	Key:          "/auth",
	FillInterval: time.Second,
	Capacity:     10,
	Quantum:      10,
})

func NewRouter() *gin.Engine {
	r := gin.New()
	if global.ServerSetting.RunMode == "debug" {
		// 使用gin默认日志
		r.Use(gin.Logger())
		// 使用gin默认异常捕获处理
		r.Use(gin.Recovery())
	} else {
		// 访问日志记录
		r.Use(middleware.AccessLog())
		// 异常捕获处理（发送邮件）
		r.Use(middleware.Recovery())
	}
	// 限流器
	r.Use(middleware.RateLimiter(methodLimiters))
	// 统一超时控制
	r.Use(middleware.ContextTimeout(global.AppSetting.DefaultContextTimeout))
	// 注册中间件——国际化翻译处理
	// 已将注册事件添加到main.go init文件中
	//r.Use(middleware.Translations())

	// Swagger API文档
	url := ginSwagger.URL("http://127.0.0.1:8000/swagger/doc.json")
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))

	// 文件上传接口
	upload := api.NewUpload()
	r.POST("/upload/file", upload.UploadFile)
	// 设置文件服务去提供静态资源的访问，才能实现让外部请求本项目 HTTP Server 时同时提供静态资源的访问，实际上在 gin 中实现 File Server
	// 如下可访问：http://127.0.0.1:8000/static/里面的所有资源
	//r.StaticFS("/static", http.Dir(global.AppSetting.UploadSavePath))
	// 文件服务路由这样写，可以防止服务列出目录：
	// 如下可访问：只可访问返回的资源
	r.Static("/static", global.AppSetting.UploadSavePath)

	// 获取 Token的API
	r.POST("/auth", api.GetAuth)

	// 业务接口 API
	article := v1.NewArticle()
	tag := v1.NewTag()

	apiv1 := r.Group("/api/v1")
	// 在完成了 JWT 的中间件编写后，我们需要将其接入到应用流程中，但是需要注意的是，并非所有的接口都需要用到 JWT 中间件，
	// 因此我们需要利用 gin 中的分组路由的概念，只针对 apiv1 的路由分组进行 JWT 中间件的引用，
	// 也就是只有 apiv1 路由分组里的路由方法会受此中间件的约束
	apiv1.Use(middleware.JWT())
	{
		// 创建标签
		apiv1.POST("/tags", tag.Create)
		// 删除指定标签
		apiv1.DELETE("/tags/:id", tag.Delete)
		// 更新指定标签
		apiv1.PUT("/tags/:id", tag.Update)
		apiv1.PATCH("/tags/:id/state", tag.Update)
		// 获取指定标签
		apiv1.GET("/tags/:id", tag.Get)
		// 获取标签列表
		apiv1.GET("/tags", tag.List)

		// 创建文章
		apiv1.POST("/articles", article.Create)
		// 删除指定文章
		apiv1.DELETE("/articles/:id", article.Delete)
		// 更新指定文章
		apiv1.PUT("/articles/:id", article.Update)
		apiv1.PATCH("/articles/:id/state", article.Update)
		// 获取指定文章
		apiv1.GET("/articles/:id", article.Get)
		// 获取文章列表
		apiv1.GET("/articles", article.List)
	}

	return r
}
