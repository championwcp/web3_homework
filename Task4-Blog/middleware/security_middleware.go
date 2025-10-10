package middleware

import (
	"github.com/gin-gonic/gin"
)

// SecurityHeaders 安全头中间件
func SecurityHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 防止 XSS 攻击
		c.Writer.Header().Set("X-XSS-Protection", "1; mode=block")
		// 防止点击劫持
		c.Writer.Header().Set("X-Frame-Options", "DENY")
		// 防止 MIME 类型嗅探
		c.Writer.Header().Set("X-Content-Type-Options", "nosniff")
		// HSTS 强制 HTTPS
		c.Writer.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		// 内容安全策略
		c.Writer.Header().Set("Content-Security-Policy", "default-src 'self'")
		// 引用策略
		c.Writer.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")

		c.Next()
	}
}

// NoCache 禁止缓存中间件
func NoCache() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Cache-Control", "no-cache, no-store, max-age=0, must-revalidate, value")
		c.Writer.Header().Set("Expires", "Thu, 01 Jan 1970 00:00:00 GMT")
		c.Writer.Header().Set("Last-Modified", "Thu, 01 Jan 1970 00:00:00 GMT")
		c.Next()
	}
}

// Secure 安全中间件
func Secure() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 设置安全头
		SecurityHeaders()(c)
		
		// 对于 API 请求，禁止缓存
		if c.Request.URL.Path == "/api/" || c.Request.URL.Path[:4] == "/api" {
			NoCache()(c)
		}
		
		c.Next()
	}
}