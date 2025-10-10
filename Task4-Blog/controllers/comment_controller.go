package controllers

import (
	"net/http"
	"strconv"

	"blog-system/services"
	"blog-system/utils"

	"github.com/gin-gonic/gin"
)

// CommentController 评论控制器
type CommentController struct {
	commentService *services.CommentService
}

// NewCommentController 创建评论控制器实例
func NewCommentController(commentService *services.CommentService) *CommentController {
	return &CommentController{
		commentService: commentService,
	}
}

// CreateCommentRequest 创建评论请求结构
type CreateCommentRequest struct {
	Content  string `json:"content" binding:"required,min=1,max=1000"`
	PostID   uint   `json:"post_id" binding:"required"`
	ParentID *uint  `json:"parent_id,omitempty"`
}

// CreateComment 创建评论
func (cc *CommentController) CreateComment(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "未授权", "用户信息不存在")
		return
	}

	var req CreateCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "请求参数错误", err.Error())
		return
	}

	comment, err := cc.commentService.CreateComment(userID.(uint), req.PostID, req.Content, req.ParentID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "创建评论失败", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "创建评论成功", comment.ToResponse())
}

// GetPostComments 获取文章评论列表
func (cc *CommentController) GetPostComments(c *gin.Context) {
	postIDStr := c.Param("postId")
	postID, err := strconv.ParseUint(postIDStr, 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "参数错误", "文章ID格式不正确")
		return
	}

	// 获取分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	comments, total, err := cc.commentService.GetPostComments(uint(postID), page, pageSize)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "获取评论列表失败", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "获取评论列表成功", gin.H{
		"comments": comments,
		"pagination": gin.H{
			"page":       page,
			"page_size":  pageSize,
			"total":      total,
			"total_page": (total + int64(pageSize) - 1) / int64(pageSize),
		},
	})
}

// DeleteComment 删除评论
func (cc *CommentController) DeleteComment(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "未授权", "用户信息不存在")
		return
	}

	commentIDStr := c.Param("id")
	commentID, err := strconv.ParseUint(commentIDStr, 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "参数错误", "评论ID格式不正确")
		return
	}

	// 检查用户是否有权限删除这条评论
	if !cc.commentService.IsCommentOwner(uint(commentID), userID.(uint)) {
		utils.ErrorResponse(c, http.StatusForbidden, "权限不足", "只能删除自己的评论")
		return
	}

	if err := cc.commentService.DeleteComment(uint(commentID)); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "删除评论失败", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "删除评论成功", nil)
}

// GetCommentByID 根据ID获取评论
func (cc *CommentController) GetCommentByID(c *gin.Context) {
	commentIDStr := c.Param("id")
	commentID, err := strconv.ParseUint(commentIDStr, 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "参数错误", "评论ID格式不正确")
		return
	}

	comment, err := cc.commentService.GetCommentByID(uint(commentID))
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "评论不存在", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "获取评论成功", comment.ToResponse())
}