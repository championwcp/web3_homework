package models

import (
	"time"

	"gorm.io/gorm"
)

// PostStatus 文章状态类型
type PostStatus string

const (
	PostStatusDraft     PostStatus = "draft"     // 草稿
	PostStatusPublished PostStatus = "published" // 已发布
	PostStatusArchived  PostStatus = "archived"  // 已归档
)

// Post 文章模型
type Post struct {
	ID          uint       `gorm:"primaryKey" json:"id"`
	Title       string     `gorm:"size:200;not null" json:"title"`
	Content     string     `gorm:"type:longtext;not null" json:"content"`
	Summary     string     `gorm:"type:text" json:"summary"`                    // 文章摘要
	Slug        string     `gorm:"size:255;uniqueIndex" json:"slug"`           // URL 友好标识
	Status      PostStatus `gorm:"size:20;default:'draft'" json:"status"`      // 文章状态
	IsPublic    bool       `gorm:"default:true" json:"is_public"`              // 是否公开
	ViewCount   int        `gorm:"default:0" json:"view_count"`                // 阅读次数
	LikeCount   int        `gorm:"default:0" json:"like_count"`                // 点赞数
	CommentCount int       `gorm:"default:0" json:"comment_count"`             // 评论数
	PublishedAt *time.Time `json:"published_at,omitempty"`                     // 发布时间
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"` // 软删除

	// 外键关系
	UserID uint `gorm:"not null;index" json:"user_id"`
	User   User `gorm:"foreignKey:UserID" json:"user,omitempty"`

	// 关联关系
	Comments []Comment `gorm:"foreignKey:PostID" json:"comments,omitempty"`
	Tags     []Tag     `gorm:"many2many:post_tags;" json:"tags,omitempty"`
}

// TableName 指定表名
func (Post) TableName() string {
	return "posts"
}

// BeforeCreate 创建前的钩子函数
func (p *Post) BeforeCreate(tx *gorm.DB) error {
	// 自动生成摘要
	if p.Summary == "" && len(p.Content) > 150 {
		p.Summary = p.Content[:150] + "..."
	} else if p.Summary == "" {
		p.Summary = p.Content
	}
	return nil
}

// BeforeUpdate 更新前的钩子函数
func (p *Post) BeforeUpdate(tx *gorm.DB) error {
	// 如果文章状态变为已发布，设置发布时间
	if p.Status == PostStatusPublished && p.PublishedAt == nil {
		now := time.Now()
		p.PublishedAt = &now
	}
	return nil
}

// AfterCreate 创建后的钩子函数 - 更新用户的文章数量
func (p *Post) AfterCreate(tx *gorm.DB) error {
	// 更新用户的文章数量
	return tx.Model(&User{}).Where("id = ?", p.UserID).
		Update("post_count", gorm.Expr("post_count + ?", 1)).Error
}

// AfterDelete 删除后的钩子函数 - 更新用户的文章数量
func (p *Post) AfterDelete(tx *gorm.DB) error {
	// 更新用户的文章数量
	return tx.Model(&User{}).Where("id = ?", p.UserID).
		Update("post_count", gorm.Expr("post_count - ?", 1)).Error
}

// PostResponse 文章响应结构
type PostResponse struct {
	ID          uint       `json:"id"`
	Title       string     `json:"title"`
	Content     string     `json:"content"`
	Summary     string     `json:"summary"`
	Slug        string     `json:"slug"`
	Status      PostStatus `json:"status"`
	IsPublic    bool       `json:"is_public"`
	ViewCount   int        `json:"view_count"`
	LikeCount   int        `json:"like_count"`
	CommentCount int       `json:"comment_count"`
	PublishedAt *time.Time `json:"published_at,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	User        UserResponse `json:"user"`
}

// ToResponse 转换为响应结构体
func (p *Post) ToResponse() PostResponse {
	return PostResponse{
		ID:          p.ID,
		Title:       p.Title,
		Content:     p.Content,
		Summary:     p.Summary,
		Slug:        p.Slug,
		Status:      p.Status,
		IsPublic:    p.IsPublic,
		ViewCount:   p.ViewCount,
		LikeCount:   p.LikeCount,
		CommentCount: p.CommentCount,
		PublishedAt: p.PublishedAt,
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
		User:        p.User.ToResponse(),
	}
}