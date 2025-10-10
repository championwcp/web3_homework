package models

import (
	"time"
)

// Tag 标签模型
type Tag struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"size:50;uniqueIndex;not null" json:"name"`
	Slug      string    `gorm:"size:100;uniqueIndex;not null" json:"slug"` // URL 友好标识
	Color     string    `gorm:"size:7" json:"color"` // 标签颜色，如 #FF0000
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// 多对多关系
	Posts []Post `gorm:"many2many:post_tags;" json:"posts,omitempty"`
}

// TableName 指定表名
func (Tag) TableName() string {
	return "tags"
}

// PostTag 文章标签关联表
type PostTag struct {
	PostID    uint      `gorm:"primaryKey" json:"post_id"`
	TagID     uint      `gorm:"primaryKey" json:"tag_id"`
	CreatedAt time.Time `json:"created_at"`
}

// TableName 指定表名
func (PostTag) TableName() string {
	return "post_tags"
}