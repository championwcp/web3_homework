package models

import (
	"time"

	"gorm.io/gorm"
)

// Comment 评论模型
type Comment struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Content   string    `gorm:"type:text;not null" json:"content"`
	IsApproved bool     `gorm:"default:true" json:"is_approved"` // 评论是否审核通过
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"` // 软删除

	// 外键关系
	UserID uint `gorm:"not null;index" json:"user_id"`
	User   User `gorm:"foreignKey:UserID" json:"user,omitempty"`

	PostID uint `gorm:"not null;index" json:"post_id"`
	Post   Post `gorm:"foreignKey:PostID" json:"post,omitempty"`

	// 父评论ID（支持回复功能）
	ParentID *uint    `gorm:"index" json:"parent_id,omitempty"`
	Replies  []Comment `gorm:"foreignKey:ParentID" json:"replies,omitempty"`
}

// TableName 指定表名
func (Comment) TableName() string {
	return "comments"
}

// BeforeCreate 创建前的钩子函数
func (c *Comment) BeforeCreate(tx *gorm.DB) error {
	// 可以在这里设置默认值
	return nil
}

// AfterCreate 创建后的钩子函数 - 更新文章的评论数量
func (c *Comment) AfterCreate(tx *gorm.DB) error {
	if c.IsApproved {
		return tx.Model(&Post{}).Where("id = ?", c.PostID).
			Update("comment_count", gorm.Expr("comment_count + ?", 1)).Error
	}
	return nil
}

// BeforeDelete 删除前的钩子函数
func (c *Comment) BeforeDelete(tx *gorm.DB) error {
	return nil
}

// AfterDelete 删除后的钩子函数 - 检查并更新文章评论状态
func (c *Comment) AfterDelete(tx *gorm.DB) error {
	if c.IsApproved {
		// 更新文章的评论数量
		err := tx.Model(&Post{}).Where("id = ?", c.PostID).
			Update("comment_count", gorm.Expr("comment_count - ?", 1)).Error
		if err != nil {
			return err
		}

		// 检查文章的评论数量
		var commentCount int64
		if err := tx.Model(&Comment{}).Where("post_id = ? AND is_approved = ?", c.PostID, true).Count(&commentCount).Error; err != nil {
			return err
		}

		// 可以在这里添加更新文章状态等其它逻辑
		if commentCount == 0 {
			// 文章没有评论了，可以更新状态
			// 例如: tx.Model(&Post{}).Where("id = ?", c.PostID).Update("has_comments", false)
		}
	}
	return nil
}

// CommentResponse 评论响应结构
type CommentResponse struct {
	ID         uint      `json:"id"`
	Content    string    `json:"content"`
	IsApproved bool      `json:"is_approved"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	User       UserResponse `json:"user"`
	PostID     uint      `json:"post_id"`
	ParentID   *uint     `json:"parent_id,omitempty"`
	Replies    []CommentResponse `json:"replies,omitempty"`
}

// ToResponse 转换为响应结构体
func (c *Comment) ToResponse() CommentResponse {
	// 递归转换回复
	var replies []CommentResponse
	for _, reply := range c.Replies {
		replies = append(replies, reply.ToResponse())
	}

	return CommentResponse{
		ID:         c.ID,
		Content:    c.Content,
		IsApproved: c.IsApproved,
		CreatedAt:  c.CreatedAt,
		UpdatedAt:  c.UpdatedAt,
		User:       c.User.ToResponse(),
		PostID:     c.PostID,
		ParentID:   c.ParentID,
		Replies:    replies,
	}
}