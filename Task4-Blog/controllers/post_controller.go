package controllers

import (
	"net/http"
	"strconv"

	"blog-system/models"
	"blog-system/services"
	"blog-system/utils"

	"github.com/gin-gonic/gin"
)

// PostController 文章控制器
type PostController struct {
	postService *services.PostService
}

// NewPostController 创建文章控制器实例
func NewPostController(postService *services.PostService) *PostController {
	return &PostController{
		postService: postService,
	}
}

// CreatePostRequest 创建文章请求结构
type CreatePostRequest struct {
	Title   string            `json:"title" binding:"required,min=1,max=200"`
	Content string            `json:"content" binding:"required,min=1"`
	Summary string            `json:"summary,omitempty"`
	Slug    string            `json:"slug,omitempty"`
	Status  models.PostStatus `json:"status,omitempty"`
	IsPublic bool            `json:"is_public,omitempty"`
}

// UpdatePostRequest 更新文章请求结构
type UpdatePostRequest struct {
	Title   string            `json:"title,omitempty"`
	Content string            `json:"content,omitempty"`
	Summary string            `json:"summary,omitempty"`
	Slug    string            `json:"slug,omitempty"`
	Status  models.PostStatus `json:"status,omitempty"`
	IsPublic *bool           `json:"is_public,omitempty"`
}

// CreatePost 创建文章
func (pc *PostController) CreatePost(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "未授权", "用户信息不存在")
		return
	}

	var req CreatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "请求参数错误", err.Error())
		return
	}

	post, err := pc.postService.CreatePost(userID.(uint), req.Title, req.Content, req.Summary, req.Slug, req.Status, req.IsPublic)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "创建文章失败", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "创建文章成功", post.ToResponse())
}

// GetPostByID 根据ID获取文章
func (pc *PostController) GetPostByID(c *gin.Context) {
	postIDStr := c.Param("id")
	postID, err := strconv.ParseUint(postIDStr, 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "参数错误", "文章ID格式不正确")
		return
	}

	post, err := pc.postService.GetPostByID(uint(postID))
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "文章不存在", err.Error())
		return
	}

	// 增加阅读次数
	pc.postService.IncrementViewCount(uint(postID))

	utils.SuccessResponse(c, http.StatusOK, "获取文章成功", post.ToResponse())
}

// GetPosts 获取文章列表
func (pc *PostController) GetPosts(c *gin.Context) {
	// 获取分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	status := c.DefaultQuery("status", "published") // 默认只获取已发布的文章

	posts, total, err := pc.postService.GetPosts(page, pageSize, models.PostStatus(status))
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

// UpdatePost 更新文章
func (pc *PostController) UpdatePost(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "未授权", "用户信息不存在")
		return
	}

	postIDStr := c.Param("id")
	postID, err := strconv.ParseUint(postIDStr, 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "参数错误", "文章ID格式不正确")
		return
	}

	var req UpdatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "请求参数错误", err.Error())
		return
	}

	// 检查用户是否有权限修改这篇文章
	if !pc.postService.IsPostOwner(uint(postID), userID.(uint)) {
		utils.ErrorResponse(c, http.StatusForbidden, "权限不足", "只能修改自己的文章")
		return
	}

	post, err := pc.postService.UpdatePost(uint(postID), req.Title, req.Content, req.Summary, req.Slug, req.Status, req.IsPublic)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "更新文章失败", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "更新文章成功", post.ToResponse())
}

// DeletePost 删除文章
func (pc *PostController) DeletePost(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "未授权", "用户信息不存在")
		return
	}

	postIDStr := c.Param("id")
	postID, err := strconv.ParseUint(postIDStr, 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "参数错误", "文章ID格式不正确")
		return
	}

	// 检查用户是否有权限删除这篇文章
	if !pc.postService.IsPostOwner(uint(postID), userID.(uint)) {
		utils.ErrorResponse(c, http.StatusForbidden, "权限不足", "只能删除自己的文章")
		return
	}

	if err := pc.postService.DeletePost(uint(postID)); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "删除文章失败", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "删除文章成功", nil)
}

// GetUserPosts 获取当前用户的文章
func (pc *PostController) GetUserPosts(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "未授权", "用户信息不存在")
		return
	}

	// 获取分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	status := c.DefaultQuery("status", "")

	var postStatus models.PostStatus
	if status != "" {
		postStatus = models.PostStatus(status)
	}

	posts, total, err := pc.postService.GetUserPosts(userID.(uint), page, pageSize, postStatus)
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