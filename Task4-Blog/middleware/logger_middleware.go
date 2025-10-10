package middleware

import (
	"bytes"
	"io"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Logger 日志中间件
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 开始时间
		start := time.Now()
		
		// 请求路径
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery
		
		// 记录请求体（对于非GET请求）
		var requestBody []byte
		if c.Request.Method != "GET" {
			requestBody, _ = io.ReadAll(c.Request.Body)
			c.Request.Body = io.NopCloser(bytes.NewBuffer(requestBody))
		}

		// 处理请求
		c.Next()

		// 结束时间
		cost := time.Since(start)
		
		// 状态码
		statusCode := c.Writer.Status()
		
		// 客户端IP
		clientIP := c.ClientIP()
		
		// 请求方法
		method := c.Request.Method
		
		// 用户代理
		userAgent := c.Request.UserAgent()

		// 错误信息
		errors := c.Errors.ByType(gin.ErrorTypePrivate).String()

		// 记录日志
		log.Printf("| %3d | %13v | %15s | %-7s %s | %s | %s | %s",
			statusCode,
			cost,
			clientIP,
			method,
			path,
			query,
			userAgent,
			errors,
		)

		// 如果是错误响应，记录更多信息
		if statusCode >= 400 {
			log.Printf("Error Request Body: %s", string(requestBody))
		}
	}
}

// ZapLogger 使用 Zap 的日志中间件
func ZapLogger(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		// 处理请求
		c.Next()

		end := time.Now()
		latency := end.Sub(start)

		// 记录日志
		logger.Info("HTTP Request",
			zap.Int("status", c.Writer.Status()),
			zap.String("method", c.Request.Method),
			zap.String("path", path),
			zap.String("query", query),
			zap.String("ip", c.ClientIP()),
			zap.String("user-agent", c.Request.UserAgent()),
			zap.Duration("latency", latency),
			zap.String("errors", c.Errors.ByType(gin.ErrorTypePrivate).String()),
		)
	}
}

// Recovery 恢复中间件（防止 panic 导致服务崩溃）
func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// 记录 panic 信息
				log.Printf("Panic recovered: %v", err)
				
				// 返回 500 错误
				c.JSON(500, gin.H{
					"success": false,
					"message": "服务器内部错误",
					"error":   "Internal Server Error",
				})
				
				c.Abort()
			}
		}()
		
		c.Next()
	}
}