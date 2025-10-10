package routes

import (
	"blog-system/controllers"
	"blog-system/middleware"
	"blog-system/services"

	"github.com/gin-gonic/gin"
)

// SetupRoutes 设置路由
func SetupRoutes(r *gin.Engine) {
	// 初始化服务层
	authService := services.NewAuthService()
	userService := services.NewUserService()
	postService := services.NewPostService()
	commentService := services.NewCommentService()

	// 初始化控制器
	authController := controllers.NewAuthController(authService, userService)
	userController := controllers.NewUserController(userService)
	postController := controllers.NewPostController(postService)
	commentController := controllers.NewCommentController(commentService)

	// 初始化中间件
	authMiddleware := middleware.NewAuthMiddleware(authService)

	// 全局中间件
	setupGlobalMiddleware(r)

	// API 路由组
	api := r.Group("/api/v1")
	{
		// 公开路由 - 不需要认证
		public := api.Group("")
		{
			setupPublicRoutes(public, authController, userController, postController, commentController)
		}

		// 受保护路由 - 需要认证
		protected := api.Group("")
		protected.Use(authMiddleware.AuthRequired())
		{
			setupProtectedRoutes(protected, authController, userController, postController, commentController)
		}

		// 管理员路由 - 需要管理员权限
		admin := api.Group("/admin")
		admin.Use(authMiddleware.AdminRequired())
		{
			setupAdminRoutes(admin, userController, postController, commentController)
		}
	}

	// 健康检查路由（不在 API 分组内）
	setupHealthRoutes(r)
}

// setupGlobalMiddleware 设置全局中间件
func setupGlobalMiddleware(r *gin.Engine) {
	// 跨域中间件
	r.Use(middleware.CORS())
	// 日志中间件
	r.Use(middleware.Logger())
	// 恢复中间件
	r.Use(middleware.Recovery())
	// 安全中间件
	r.Use(middleware.SecurityHeaders())
	// 限流中间件（每秒10个请求，突发20个）
	// r.Use(middleware.GlobalRateLimit(10, 20))
}

// setupPublicRoutes 设置公开路由
func setupPublicRoutes(public *gin.RouterGroup, authController *controllers.AuthController, userController *controllers.UserController, postController *controllers.PostController, commentController *controllers.CommentController) {
	// 认证相关
	auth := public.Group("/auth")
	{
		auth.POST("/register", authController.Register)
		auth.POST("/login", authController.Login)
	}

	// 用户相关
	users := public.Group("/users")
	{
		users.GET("/:id", userController.GetUserByID)
		users.GET("/:id/posts", userController.GetUserPosts)
	}

	// 文章相关
	posts := public.Group("/posts")
	{
		posts.GET("", postController.GetPosts)
		posts.GET("/:id", postController.GetPostByID)
	}

	// 评论相关
	comments := public.Group("/comments")
	{
		comments.GET("/posts/:postId", commentController.GetPostComments)
		comments.GET("/:id", commentController.GetCommentByID)
	}
}

// setupProtectedRoutes 设置受保护路由（需要登录）
func setupProtectedRoutes(protected *gin.RouterGroup, authController *controllers.AuthController, userController *controllers.UserController, postController *controllers.PostController, commentController *controllers.CommentController) {
	// 用户相关
	users := protected.Group("/users")
	{
		users.GET("/profile", authController.GetProfile)
		users.PUT("/profile", authController.UpdateProfile)
		users.GET("/my/posts", postController.GetUserPosts)
	}

	// 文章相关
	posts := protected.Group("/posts")
	{
		posts.POST("", postController.CreatePost)
		posts.PUT("/:id", postController.UpdatePost)
		posts.DELETE("/:id", postController.DeletePost)
	}

	// 评论相关
	comments := protected.Group("/comments")
	{
		comments.POST("", commentController.CreateComment)
		comments.DELETE("/:id", commentController.DeleteComment)
	}
}

// setupAdminRoutes 设置管理员路由
func setupAdminRoutes(admin *gin.RouterGroup, userController *controllers.UserController, postController *controllers.PostController, commentController *controllers.CommentController) {
	// 用户管理
	users := admin.Group("/users")
	{
		users.GET("", userController.GetUsers)
		// 可以添加更多管理员功能：用户封禁、角色修改等
	}

	// // 文章管理
	// posts := admin.Group("/posts")
	// {
	// 	// 可以添加文章审核、推荐等功能
	// }

	// // 评论管理
	// comments := admin.Group("/comments")
	// {
	// 	// 可以添加评论审核、删除等功能
	// }
}

// setupHealthRoutes 设置健康检查路由
func setupHealthRoutes(r *gin.Engine) {
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"message": "Blog system is running",
			"version": "1.0.0",
		})
	})

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"name":        "Personal Blog System API",
			"version":     "1.0.0",
			"description": "A blog system built with Gin and GORM",
			"docs":        "/api/v1",
		})
	})

	// 数据库信息
	r.GET("/db-info", func(c *gin.Context) {
		// 这里可以添加数据库连接检查
		c.JSON(200, gin.H{
			"database": "connected",
			"status":   "healthy",
		})
	})
}

// SetupSwaggerRoutes 设置 Swagger 文档路由（可选）
func SetupSwaggerRoutes(r *gin.Engine) {
	// r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}

// SetupTestRoutes 设置测试路由（开发环境使用）
func SetupTestRoutes(r *gin.Engine) {
	test := r.Group("/test")
	{
		test.GET("/ping", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"message": "pong",
			})
		})

		test.GET("/error", func(c *gin.Context) {
			panic("测试异常恢复")
		})

		test.GET("/slow", func(c *gin.Context) {
			// 模拟慢请求
			// time.Sleep(2 * time.Second)
			c.JSON(200, gin.H{
				"message": "slow response",
			})
		})
	}
}