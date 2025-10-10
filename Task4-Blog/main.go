package main

import (
	"fmt"
	"log"

	"blog-system/config"
	"blog-system/database"
	"blog-system/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	// 1. 初始化配置
	if err := config.Init(); err != nil {
		log.Fatalf("配置初始化失败: %v", err)
	}

	// 2. 初始化数据库连接
	_, err := database.InitDB()
	if err != nil {
		log.Fatalf("数据库初始化失败: %v", err)
	}
	defer func() {
		if err := database.CloseDB(); err != nil {
			log.Printf("关闭数据库连接失败: %v", err)
		} else {
			log.Println(" 数据库连接已关闭")
		}
	}()

	// 3. 执行数据库迁移
	if err := database.Migrate(); err != nil {
		log.Fatalf(" 数据库迁移失败: %v", err)
	}

	// 4. 设置 Gin 运行模式
	cfg := config.GetConfig()
	gin.SetMode(cfg.Server.Mode)

	// 5. 初始化 Gin
	r := gin.Default()

	// 6. 设置路由
	routes.SetupRoutes(r)

	// 7. 启动服务器
	serverConfig := cfg.Server
	log.Printf("服务器启动在 :%d 端口 [%s 模式]", serverConfig.Port, serverConfig.Mode)
	
	if err := r.Run(fmt.Sprintf(":%d", serverConfig.Port)); err != nil {
		log.Fatalf("服务器启动失败: %v", err)
	}
}