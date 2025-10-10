package controllers

import (
	"net/http"
	"strconv"

	"blog-system/services"
	"blog-system/utils"
	"blog-system/models"

	"github.com/gin-gonic/gin"
)

// UserController 用户控制器
type UserController struct {
	userService *services.UserService
}

// NewUserController 创建用户控制器实例
func NewUserController(userService *services.UserService) *UserController {
	return &UserController{
		userService: userService,
	}
}

// GetUserByID 根据ID获取用户信息
func (uc *UserController) GetUserByID(c *gin.Context) {
	userIDStr := c.Param("id")
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "参数错误", "用户ID格式不正确")
		return
	}

	user, err := uc.userService.GetUserByID(uint(userID))
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "用户不存在", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "获取用户信息成功", user.ToResponse())
}

// GetUserPosts 获取用户的文章列表
func (uc *UserController) GetUserPosts(c *gin.Context) {
	userIDStr := c.Param("id")
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "参数错误", "用户ID格式不正确")
		return
	}

	// 获取分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	posts, total, err := uc.userService.GetUserPosts(uint(userID), page, pageSize)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "获取文章列表失败", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "获取文章列表成功", gin.H{
		"posts": posts,
		"pagination": gin.H{
			"page":       page,
			"page_size":  pageSize,
			"total":      total,
			"total_page": (total + int64(pageSize) - 1) / int64(pageSize),
		},
	})
}

// GetUsers 获取用户列表（管理员功能）
func (uc *UserController) GetUsers(c *gin.Context) {
	// 获取分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	users, total, err := uc.userService.GetUsers(page, pageSize)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "获取用户列表失败", err.Error())
		return
	}

	// 转换为响应格式
	var userResponses []models.UserResponse
	for _, user := range users {
		userResponses = append(userResponses, user.ToResponse())
	}

	utils.SuccessResponse(c, http.StatusOK, "获取用户列表成功", gin.H{
		"users": userResponses,
		"pagination": gin.H{
			"page":       page,
			"page_size":  pageSize,
			"total":      total,
			"total_page": (total + int64(pageSize) - 1) / int64(pageSize),
		},
	})
}