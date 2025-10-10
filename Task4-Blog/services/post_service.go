package services

import (
	"blog-system/database"
	"blog-system/models"

	"gorm.io/gorm"
)

// PostService 文章服务
type PostService struct {
	db *gorm.DB
}

// NewPostService 创建文章服务实例
func NewPostService() *PostService {
	return &PostService{
		db: database.GetDB(),
	}
}

// CreatePost 创建文章
func (ps *PostService) CreatePost(userID uint, title, content, summary, slug string, status models.PostStatus, isPublic bool) (*models.Post, error) {
	post := &models.Post{
		Title:    title,
		Content:  content,
		Summary:  summary,
		Slug:     slug,
		Status:   status,
		IsPublic: isPublic,
		UserID:   userID,
	}

	if err := ps.db.Create(post).Error; err != nil {
		return nil, err
	}

	// 预加载用户信息
	if err := ps.db.Preload("User").First(post, post.ID).Error; err != nil {
		return nil, err
	}

	return post, nil
}

// GetPostByID 根据ID获取文章
func (ps *PostService) GetPostByID(postID uint) (*models.Post, error) {
	var post models.Post
	if err := ps.db.Preload("User").Preload("Comments").Preload("Comments.User").
		First(&post, postID).Error; err != nil {
		return nil, err
	}
	return &post, nil
}

// GetPosts 获取文章列表
func (ps *PostService) GetPosts(page, pageSize int, status models.PostStatus) ([]models.PostResponse, int64, error) {
	var posts []models.Post
	var total int64

	// 计算偏移量
	offset := (page - 1) * pageSize

	// 构建查询条件
	query := ps.db.Model(&models.Post{})
	if status != "" {
		query = query.Where("status = ?", status)
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 获取文章列表
	if err := query.Preload("User").
		Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&posts).Error; err != nil {
		return nil, 0, err
	}

	// 转换为响应格式
	var postResponses []models.PostResponse
	for _, post := range posts {
		postResponses = append(postResponses, post.ToResponse())
	}

	return postResponses, total, nil
}

// UpdatePost 更新文章
func (ps *PostService) UpdatePost(postID uint, title, content, summary, slug string, status models.PostStatus, isPublic *bool) (*models.Post, error) {
	var post models.Post
	if err := ps.db.First(&post, postID).Error; err != nil {
		return nil, err
	}

	updates := make(map[string]interface{})
	if title != "" {
		updates["title"] = title
	}
	if content != "" {
		updates["content"] = content
	}
	if summary != "" {
		updates["summary"] = summary
	}
	if slug != "" {
		updates["slug"] = slug
	}
	if status != "" {
		updates["status"] = status
	}
	if isPublic != nil {
		updates["is_public"] = *isPublic
	}

	if len(updates) > 0 {
		if err := ps.db.Model(&post).Updates(updates).Error; err != nil {
			return nil, err
		}
	}

	// 重新加载关联数据
	if err := ps.db.Preload("User").First(&post, postID).Error; err != nil {
		return nil, err
	}

	return &post, nil
}

// DeletePost 删除文章
func (ps *PostService) DeletePost(postID uint) error {
	return ps.db.Delete(&models.Post{}, postID).Error
}

// IsPostOwner 检查用户是否是文章的作者
func (ps *PostService) IsPostOwner(postID, userID uint) bool {
	var post models.Post
	if err := ps.db.Select("user_id").First(&post, postID).Error; err != nil {
		return false
	}
	return post.UserID == userID
}

// IncrementViewCount 增加文章阅读次数
func (ps *PostService) IncrementViewCount(postID uint) {
	ps.db.Model(&models.Post{}).Where("id = ?", postID).
		Update("view_count", gorm.Expr("view_count + ?", 1))
}

// GetUserPosts 获取用户的文章列表
func (ps *PostService) GetUserPosts(userID uint, page, pageSize int, status models.PostStatus) ([]models.PostResponse, int64, error) {
	var posts []models.Post
	var total int64

	// 计算偏移量
	offset := (page - 1) * pageSize

	// 构建查询条件
	query := ps.db.Model(&models.Post{}).Where("user_id = ?", userID)
	if status != "" {
		query = query.Where("status = ?", status)
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 获取文章列表
	if err := query.Preload("User").
		Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&posts).Error; err != nil {
		return nil, 0, err
	}

	// 转换为响应格式
	var postResponses []models.PostResponse
	for _, post := range posts {
		postResponses = append(postResponses, post.ToResponse())
	}

	return postResponses, total, nil
}