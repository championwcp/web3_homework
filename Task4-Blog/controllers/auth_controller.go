package controllers

import (
	"net/http"
	"time"

	"blog-system/models"
	"blog-system/services"
	"blog-system/utils"

	"github.com/gin-gonic/gin"
)

// AuthController 认证控制器
type AuthController struct {
	authService *services.AuthService
	userService *services.UserService
}

// NewAuthController 创建认证控制器实例
func NewAuthController(authService *services.AuthService, userService *services.UserService) *AuthController {
	return &AuthController{
		authService: authService,
		userService: userService,
	}
}

// RegisterRequest 注册请求结构
type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	Bio      string `json:"bio,omitempty"`
}

// LoginRequest 登录请求结构
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// AuthResponse 认证响应结构
type AuthResponse struct {
	Token string             `json:"token"`
	User  models.UserResponse `json:"user"`
}

// Register 用户注册
func (ac *AuthController) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "请求参数错误", err.Error())
		return
	}

	// 检查用户名是否已存在
	if ac.userService.UsernameExists(req.Username) {
		utils.ErrorResponse(c, http.StatusBadRequest, "注册失败", "用户名已存在")
		return
	}

	// 检查邮箱是否已存在
	if ac.userService.EmailExists(req.Email) {
		utils.ErrorResponse(c, http.StatusBadRequest, "注册失败", "邮箱已存在")
		return
	}

	// 创建用户
	user, err := ac.authService.Register(req.Username, req.Email, req.Password, req.Bio)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "注册失败", err.Error())
		return
	}

	// 生成 JWT token
	token, err := ac.authService.GenerateToken(user)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "生成token失败", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "注册成功", AuthResponse{
		Token: token,
		User:  user.ToResponse(),
	})
}

// Login 用户登录
func (ac *AuthController) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "请求参数错误", err.Error())
		return
	}

	// 验证用户凭证
	user, err := ac.authService.Login(req.Username, req.Password)
	if err != nil {
		utils.ErrorResponse(c, http.StatusUnauthorized, "登录失败", "用户名或密码错误")
		return
	}

	// 更新最后登录时间
	now := time.Now()
	ac.userService.UpdateLastLogin(user.ID, &now)

	// 生成 JWT token
	token, err := ac.authService.GenerateToken(user)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "生成token失败", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "登录成功", AuthResponse{
		Token: token,
		User:  user.ToResponse(),
	})
}

// GetProfile 获取当前用户信息
func (ac *AuthController) GetProfile(c *gin.Context) {
	// 从上下文中获取用户ID（由中间件设置）
	userID, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "未授权", "用户信息不存在")
		return
	}

	user, err := ac.userService.GetUserByID(userID.(uint))
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "用户不存在", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "获取用户信息成功", user.ToResponse())
}

// UpdateProfile 更新用户信息
func (ac *AuthController) UpdateProfile(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "未授权", "用户信息不存在")
		return
	}

	var updateData struct {
		Bio    string `json:"bio,omitempty"`
		Avatar string `json:"avatar,omitempty"`
	}

	if err := c.ShouldBindJSON(&updateData); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "请求参数错误", err.Error())
		return
	}

	user, err := ac.userService.UpdateUser(userID.(uint), updateData.Bio, updateData.Avatar)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "更新用户信息失败", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "更新用户信息成功", user.ToResponse())
}