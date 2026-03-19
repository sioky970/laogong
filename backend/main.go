package main

import (
	"log"
	"os"
	"time"

	"laogong-visa/internal/handler"
	"laogong-visa/internal/repository"
	"laogong-visa/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// 初始化数据库
	_, err := repository.InitDB()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// 自动迁移
	if err := repository.AutoMigrate(); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	// 创建默认管理员
	if err := repository.CreateDefaultAdmin(); err != nil {
		log.Printf("Warning: Failed to create default admin: %v", err)
	}

	log.Println("Database connected and migrated successfully")

	// 设置 Gin 模式
	if os.Getenv("GIN_MODE") == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	// 创建路由
	r := gin.Default()

	// CORS 配置
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization"}
	r.Use(cors.New(config))

	// 静态文件服务（上传的图片）
	r.Static("/uploads", "/app/uploads")

	// 创建处理器
	authHandler := handler.NewAuthHandler()
	recordHandler := handler.NewRecordHandler()

	// API 路由组
	api := r.Group("/api")
	{
		// 公开路由
		api.POST("/auth/login", authHandler.Login)
		
		// 前台公开路由
		api.GET("/records", recordHandler.GetAll)
		api.GET("/records/:id", recordHandler.GetByID)
		api.GET("/records/:id/image", recordHandler.GetRecordImage)
		api.GET("/records/:id/qrcode", recordHandler.GetQRCode)

		// 需要认证的路由
		authorized := api.Group("/")
		authorized.Use(middleware.JWTAuth())
		{
			// 认证相关
			authorized.GET("/auth/profile", authHandler.GetProfile)
			
			// 管理后台路由
			authorized.GET("/admin/records", recordHandler.GetAllForAdmin)
			authorized.POST("/admin/records", recordHandler.Create)
			authorized.PUT("/admin/records/:id", recordHandler.Update)
			authorized.DELETE("/admin/records/:id", recordHandler.Delete)
			authorized.POST("/admin/upload", recordHandler.UploadImage)
		}
	}

	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
			"time":   time.Now().Format(time.RFC3339),
		})
	})

	// 启动服务器
	port := os.Getenv("PORT")
	if port == "" {
		port = "8083"
	}

	log.Printf("Server starting on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
