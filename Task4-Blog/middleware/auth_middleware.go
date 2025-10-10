package middleware

import (
	"blog-system/services"
	"blog-system/utils"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

// AuthMiddleware JWT 认证中间件
type AuthMiddleware struct {
	authService *services.AuthService
}

// NewAuthMiddleware 创建认证中间件实例
func NewAuthMiddleware(authService *services.AuthService) *AuthMiddleware {
	return &AuthMiddleware{
		authService: authService,
	}
}

// AuthRequired 需要认证的中间件
func (am *AuthMiddleware) AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从请求头中获取 token
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			utils.UnauthorizedResponse(c, "缺少认证token")
			c.Abort()
			return
		}

		// 检查 token 格式 (Bearer token)
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			utils.UnauthorizedResponse(c, "token格式错误")
			c.Abort()
			return
		}

		tokenString := parts[1]

		// 验证 token
		token, err := am.authService.ValidateToken(tokenString)
		if err != nil || !token.Valid {
			utils.UnauthorizedResponse(c, "token无效或已过期")
			c.Abort()
			return
		}

		// 从 token 中提取用户信息
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			utils.UnauthorizedResponse(c, "token解析失败")
			c.Abort()
			return
		}

		// 将用户信息存储到上下文中
		userID := uint(claims["user_id"].(float64))
		username := claims["username"].(string)
		email := claims["email"].(string)
		role := claims["role"].(string)

		c.Set("userID", userID)
		c.Set("username", username)
		c.Set("email", email)
		c.Set("role", role)

		c.Next()
	}
}

// OptionalAuth 可选认证中间件（有token就验证，没有就跳过）
func (am *AuthMiddleware) OptionalAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Next()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.Next()
			return
		}

		tokenString := parts[1]

		token, err := am.authService.ValidateToken(tokenString)
		if err != nil || !token.Valid {
			c.Next()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.Next()
			return
		}

		userID := uint(claims["user_id"].(float64))
		username := claims["username"].(string)
		email := claims["email"].(string)
		role := claims["role"].(string)

		c.Set("userID", userID)
		c.Set("username", username)
		c.Set("email", email)
		c.Set("role", role)

		c.Next()
	}
}

// AdminRequired 需要管理员权限的中间件
func (am *AuthMiddleware) AdminRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 先进行认证
		am.AuthRequired()(c)
		if c.IsAborted() {
			return
		}

		// 检查用户角色
		role, exists := c.Get("role")
		if !exists || role != "admin" {
			utils.ForbiddenResponse(c, "需要管理员权限")
			c.Abort()
			return
		}

		c.Next()
	}
}

// GetUserFromContext 从上下文中获取用户信息
func GetUserFromContext(c *gin.Context) (userID uint, username, email, role string, exists bool) {
	userIDVal, exists := c.Get("userID")
	if !exists {
		return 0, "", "", "", false
	}

	usernameVal, _ := c.Get("username")
	emailVal, _ := c.Get("email")
	roleVal, _ := c.Get("role")

	return userIDVal.(uint), usernameVal.(string), emailVal.(string), roleVal.(string), true
}