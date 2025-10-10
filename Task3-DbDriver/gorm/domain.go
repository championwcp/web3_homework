package main

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

// User 用户模型
type User struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"size:100;not null" json:"name"`
	Email     string    `gorm:"size:100;uniqueIndex;not null" json:"email"`
	PostCount int       `gorm:"default:0" json:"post_count"` // 文章数量统计字段
	Posts     []Post    `gorm:"foreignKey:UserID" json:"posts,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Post 文章模型
type Post struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	Title         string    `gorm:"size:200;not null" json:"title"`
	Content       string    `gorm:"type:text;not null" json:"content"`
	CommentCount  int       `gorm:"default:0" json:"comment_count"`              // 评论数量
	CommentStatus string    `gorm:"size:20;default:'无评论'" json:"comment_status"` // 评论状态
	UserID        uint      `gorm:"not null" json:"user_id"`
	User          User      `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Comments      []Comment `gorm:"foreignKey:PostID" json:"comments,omitempty"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// Comment 评论模型
type Comment struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Content   string    `gorm:"type:text;not null" json:"content"`
	Author    string    `gorm:"size:100;not null" json:"author"`
	PostID    uint      `gorm:"not null" json:"post_id"`
	Post      Post      `gorm:"foreignKey:PostID" json:"post,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

// ============================== 钩子函数 ==============================

// BeforeCreate Post 创建前的钩子函数
func (p *Post) BeforeCreate(tx *gorm.DB) (err error) {
	fmt.Printf("正在创建文章: %s\n", p.Title)
	return nil
}

// AfterCreate Post 创建后的钩子函数 - 更新用户的文章数量
func (p *Post) AfterCreate(tx *gorm.DB) (err error) {
	// 更新用户的文章数量
	result := tx.Model(&User{}).Where("id = ?", p.UserID).
		Update("post_count", gorm.Expr("post_count + ?", 1))

	if result.Error != nil {
		return result.Error
	}

	fmt.Printf("已更新用户 %d 的文章数量\n", p.UserID)
	return nil
}

// BeforeDelete Comment 删除前的钩子函数
func (c *Comment) BeforeDelete(tx *gorm.DB) (err error) {
	fmt.Printf("正在删除评论: %s\n", c.Content)
	return nil
}

// AfterDelete Comment 删除后的钩子函数 - 检查并更新文章评论状态
func (c *Comment) AfterDelete(tx *gorm.DB) (err error) {
	// 查询文章的评论数量
	var commentCount int64
	if err := tx.Model(&Comment{}).Where("post_id = ?", c.PostID).Count(&commentCount).Error; err != nil {
		return err
	}

	// 更新文章的评论数量和状态
	updates := map[string]interface{}{
		"comment_count": commentCount,
	}

	if commentCount == 0 {
		updates["comment_status"] = "无评论"
		fmt.Printf("文章 %d 已无评论，更新状态为 '无评论'\n", c.PostID)
	} else {
		updates["comment_status"] = "有评论"
	}

	if err := tx.Model(&Post{}).Where("id = ?", c.PostID).Updates(updates).Error; err != nil {
		return err
	}

	fmt.Printf("已更新文章 %d 的评论数量为 %d\n", c.PostID, commentCount)
	return nil
}

// ============================== 业务逻辑 ==============================

// BlogManager 博客管理器
type BlogManager struct {
	db *gorm.DB
}

// NewBlogManager 创建博客管理器
func NewBlogManager(db *gorm.DB) *BlogManager {
	return &BlogManager{db: db}
}

// CreateSampleData 创建示例数据
func (bm *BlogManager) CreateSampleData() error {
	// 清空现有数据
	bm.db.Exec("DELETE FROM comments")
	bm.db.Exec("DELETE FROM posts")
	bm.db.Exec("DELETE FROM users")

	// 创建用户
	users := []User{
		{Name: "张三", Email: "zhangsan@example.com"},
		{Name: "李四", Email: "lisi@example.com"},
		{Name: "王五", Email: "wangwu@example.com"},
	}

	for i := range users {
		if err := bm.db.Create(&users[i]).Error; err != nil {
			return err
		}
	}
	fmt.Println("用户数据创建完成")

	// 创建文章
	posts := []Post{
		{Title: "Go语言入门指南", Content: "这是一篇关于Go语言的入门教程...", UserID: users[0].ID},
		{Title: "GORM使用技巧", Content: "本文将介绍GORM的高级用法...", UserID: users[0].ID},
		{Title: "数据库设计原则", Content: "良好的数据库设计是系统成功的关键...", UserID: users[1].ID},
		{Title: "Web开发最佳实践", Content: "现代Web开发的最佳实践...", UserID: users[2].ID},
	}

	for i := range posts {
		if err := bm.db.Create(&posts[i]).Error; err != nil {
			return err
		}
	}
	fmt.Println("✅ 文章数据创建完成")

	// 创建评论
	comments := []Comment{
		{Content: "写得很好！", Author: "读者A", PostID: posts[0].ID},
		{Content: "感谢分享", Author: "读者B", PostID: posts[0].ID},
		{Content: "很有用的教程", Author: "读者C", PostID: posts[0].ID},
		{Content: "期待更多内容", Author: "读者D", PostID: posts[1].ID},
		{Content: "赞同你的观点", Author: "读者E", PostID: posts[2].ID},
	}

	for i := range comments {
		if err := bm.db.Create(&comments[i]).Error; err != nil {
			return err
		}
	}
	fmt.Println("✅ 评论数据创建完成")

	return nil
}

// GetUserPostsWithComments 查询用户的所有文章及其评论（题目2要求1）
func (bm *BlogManager) GetUserPostsWithComments(userID uint) ([]Post, error) {
	var posts []Post

	// 预加载用户信息和评论信息
	err := bm.db.Preload("User").Preload("Comments").
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&posts).Error

	if err != nil {
		return nil, err
	}

	return posts, nil
}

// GetMostCommentedPost 查询评论数量最多的文章（题目2要求2）
func (bm *BlogManager) GetMostCommentedPost() (*Post, error) {
	var post Post

	// 使用子查询找到评论数量最多的文章
	subQuery := bm.db.Model(&Comment{}).
		Select("post_id, COUNT(*) as comment_count").
		Group("post_id").
		Order("comment_count DESC").
		Limit(1)

	err := bm.db.Preload("User").Preload("Comments").
		Joins("JOIN (?) AS c ON posts.id = c.post_id", subQuery).
		First(&post).Error

	if err != nil {
		// 如果没有评论，返回评论数量最多的文章（按comment_count字段）
		err = bm.db.Preload("User").Preload("Comments").
			Order("comment_count DESC").
			First(&post).Error
		if err != nil {
			return nil, err
		}
	}

	return &post, nil
}

// GetPostWithComments 获取文章及其评论
func (bm *BlogManager) GetPostWithComments(postID uint) (*Post, error) {
	var post Post
	err := bm.db.Preload("User").Preload("Comments").
		First(&post, postID).Error
	if err != nil {
		return nil, err
	}
	return &post, nil
}

// DeleteComment 删除评论（用于测试钩子函数）
func (bm *BlogManager) DeleteComment(commentID uint) error {
	var comment Comment
	if err := bm.db.First(&comment, commentID).Error; err != nil {
		return err
	}

	return bm.db.Delete(&comment).Error
}

// GetAllUsers 获取所有用户
func (bm *BlogManager) GetAllUsers() ([]User, error) {
	var users []User
	err := bm.db.Preload("Posts").Find(&users).Error
	return users, err
}

// GetAllPosts 获取所有文章
func (bm *BlogManager) GetAllPosts() ([]Post, error) {
	var posts []Post
	err := bm.db.Preload("User").Preload("Comments").Find(&posts).Error
	return posts, err
}
