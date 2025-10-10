package models

import (
	"time"

	"gorm.io/gorm"
)

// User 用户模型
type User struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Username  string         `gorm:"size:50;uniqueIndex;not null" json:"username"`
	Email     string         `gorm:"size:100;uniqueIndex;not null" json:"email"`
	Password  string         `gorm:"size:255;not null" json:"-"` // 不序列化到 JSON
	Bio       string         `gorm:"type:text" json:"bio"`       // 个人简介
	Avatar    string         `gorm:"size:255" json:"avatar"`     // 头像 URL
	Role      string         `gorm:"size:20;default:'user'" json:"role"` // user, admin
	IsActive  bool           `gorm:"default:true" json:"is_active"`
	LastLogin *time.Time     `json:"last_login,omitempty"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"` // 软删除

	// 关联关系
	Posts    []Post    `gorm:"foreignKey:UserID" json:"posts,omitempty"`
	Comments []Comment `gorm:"foreignKey:UserID" json:"comments,omitempty"`
}

// TableName 指定表名
func (User) TableName() string {
	return "users"
}

// BeforeCreate 创建前的钩子函数
func (u *User) BeforeCreate(tx *gorm.DB) error {
	// 可以在这里设置默认值或验证
	if u.Role == "" {
		u.Role = "user"
	}
	return nil
}

// UserResponse 用户响应结构（不包含敏感信息）
type UserResponse struct {
	ID        uint      `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Bio       string    `json:"bio"`
	Avatar    string    `json:"avatar"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
}

// ToResponse 转换为响应结构体
func (u *User) ToResponse() UserResponse {
	return UserResponse{
		ID:        u.ID,
		Username:  u.Username,
		Email:     u.Email,
		Bio:       u.Bio,
		Avatar:    u.Avatar,
		Role:      u.Role,
		CreatedAt: u.CreatedAt,
	}
}