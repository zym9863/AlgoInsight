package main

import (
	"log"
	"net/http"

	"gin/config"
	"gin/handlers"

	"github.com/gin-gonic/gin"
)

func main() {
	// 加载配置
	cfg := config.Load()

	// 设置Gin模式
	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// 创建Gin路由器
	router := gin.Default()

	// 设置中间件
	setupMiddleware(router, cfg)

	// 设置路由
	setupRoutes(router)

	// 启动服务器
	port := cfg.Port
	if port == "" {
		port = "8080"
	}

	log.Printf("服务器启动在端口 %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatal("服务器启动失败:", err)
	}
}

// setupMiddleware 设置中间件
func setupMiddleware(router *gin.Engine, cfg *config.Config) {
	// CORS中间件
	router.Use(func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")

		// 在开发环境允许所有来源，生产环境限制来源
		if cfg.Environment == "development" {
			c.Header("Access-Control-Allow-Origin", "*")
		} else {
			// 生产环境可以配置允许的域名
			allowedOrigins := []string{"http://localhost:5173", "http://localhost:4173"}
			for _, allowedOrigin := range allowedOrigins {
				if origin == allowedOrigin {
					c.Header("Access-Control-Allow-Origin", origin)
					break
				}
			}
		}

		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		c.Header("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	})

	// 日志中间件
	router.Use(gin.Logger())

	// 恢复中间件
	router.Use(gin.Recovery())

	// 静态文件服务（用于生产环境）
	if cfg.Environment == "production" {
		router.Static("/assets", "./web/dist/assets")
		router.StaticFile("/", "./web/dist/index.html")
		router.StaticFile("/favicon.ico", "./web/dist/favicon.ico")

		// 处理SPA路由
		router.NoRoute(func(c *gin.Context) {
			c.File("./web/dist/index.html")
		})
	}
}

// setupRoutes 设置API路由
func setupRoutes(router *gin.Engine) {
	// API路由组
	api := router.Group("/api")
	{
		// 健康检查
		api.GET("/health", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"status":  "ok",
				"message": "算法洞察平台API服务正常运行",
			})
		})

		// 算法相关路由
		algorithms := api.Group("/algorithms")
		{
			algorithms.GET("", handlers.GetAlgorithms)
			algorithms.GET("/category/:category", handlers.GetAlgorithmsByCategory)
			algorithms.GET("/info/:id", handlers.GetAlgorithmInfo)
		}

		// 可视化相关路由
		visualize := api.Group("/visualize")
		{
			visualize.POST("/execute", handlers.ExecuteVisualization)
			visualize.GET("/step/:sessionId/:stepId", handlers.GetVisualizationStep)
			visualize.POST("/reset", handlers.ResetVisualization)
		}

		// 性能测试相关路由
		benchmark := api.Group("/benchmark")
		{
			benchmark.POST("/run", handlers.RunBenchmark)
			benchmark.GET("/results/:testId", handlers.GetBenchmarkResults)
			benchmark.POST("/compare", handlers.CompareBenchmarks)
		}

		// 数据生成相关路由
		data := api.Group("/data")
		{
			data.POST("/generate", handlers.GenerateData)
			data.GET("/presets", handlers.GetDataPresets)
		}
	}
}
