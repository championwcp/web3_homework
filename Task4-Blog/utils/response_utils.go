package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response 统一响应结构
type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// SuccessResponse 成功响应
func SuccessResponse(c *gin.Context, statusCode int, message string, data interface{}) {
	c.JSON(statusCode, Response{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// ErrorResponse 错误响应
func ErrorResponse(c *gin.Context, statusCode int, message string, err string) {
	c.JSON(statusCode, Response{
		Success: false,
		Message: message,
		Error:   err,
	})
}

// ValidationErrorResponse 参数验证错误响应
func ValidationErrorResponse(c *gin.Context, message string, errors interface{}) {
	c.JSON(http.StatusBadRequest, Response{
		Success: false,
		Message: message,
		Error:   "参数验证失败",
		Data:    errors,
	})
}

// UnauthorizedResponse 未授权响应
func UnauthorizedResponse(c *gin.Context, message string) {
	c.JSON(http.StatusUnauthorized, Response{
		Success: false,
		Message: message,
		Error:   "未授权访问",
	})
}

// ForbiddenResponse 禁止访问响应
func ForbiddenResponse(c *gin.Context, message string) {
	c.JSON(http.StatusForbidden, Response{
		Success: false,
		Message: message,
		Error:   "权限不足",
	})
}

// NotFoundResponse 资源未找到响应
func NotFoundResponse(c *gin.Context, message string) {
	c.JSON(http.StatusNotFound, Response{
		Success: false,
		Message: message,
		Error:   "资源未找到",
	})
}

// InternalServerErrorResponse 服务器内部错误响应
func InternalServerErrorResponse(c *gin.Context, message string, err error) {
	c.JSON(http.StatusInternalServerError, Response{
		Success: false,
		Message: message,
		Error:   err.Error(),
	})
}

// PaginationResponse 分页响应
type PaginationResponse struct {
	Page      int   `json:"page"`
	PageSize  int   `json:"page_size"`
	Total     int64 `json:"total"`
	TotalPage int64 `json:"total_page"`
}

// SuccessResponseWithPagination 带分页的成功响应
func SuccessResponseWithPagination(c *gin.Context, message string, data interface{}, pagination PaginationResponse) {
	response := gin.H{
		"success":    true,
		"message":    message,
		"data":       data,
		"pagination": pagination,
	}
	c.JSON(http.StatusOK, response)
}