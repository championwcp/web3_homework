package services

import (
	"blog-system/database"
	"blog-system/models"

	"gorm.io/gorm"
)

// CommentService 评论服务
type CommentService struct {
	db *gorm.DB
}

// NewCommentService 创建评论服务实例
func NewCommentService() *CommentService {
	return &CommentService{
		db: database.GetDB(),
	}
}

// CreateComment 创建评论
func (cs *CommentService) CreateComment(userID, postID uint, content string, parentID *uint) (*models.Comment, error) {
	// 检查文章是否存在
	var post models.Post
	if err := cs.db.First(&post, postID).Error; err != nil {
		return nil, err
	}

	comment := &models.Comment{
		Content:    content,
		UserID:     userID,
		PostID:     postID,
		ParentID:   parentID,
		IsApproved: true, // 默认审核通过
	}

	if err := cs.db.Create(comment).Error; err != nil {
		return nil, err
	}

	// 预加载关联数据
	if err := cs.db.Preload("User").Preload("Post").First(comment, comment.ID).Error; err != nil {
		return nil, err
	}

	return comment, nil
}

// GetCommentByID 根据ID获取评论
func (cs *CommentService) GetCommentByID(commentID uint) (*models.Comment, error) {
	var comment models.Comment
	if err := cs.db.Preload("User").Preload("Post").
		First(&comment, commentID).Error; err != nil {
		return nil, err
	}
	return &comment, nil
}

// GetPostComments 获取文章的评论列表
func (cs *CommentService) GetPostComments(postID uint, page, pageSize int) ([]models.CommentResponse, int64, error) {
	var comments []models.Comment
	var total int64

	// 计算偏移量
	offset := (page - 1) * pageSize

	// 获取总数（只计算顶级评论）
	if err := cs.db.Model(&models.Comment{}).
		Where("post_id = ? AND parent_id IS NULL AND is_approved = ?", postID, true).
		Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 获取评论列表（只获取顶级评论，预加载回复）
	if err := cs.db.Preload("User").Preload("Replies").Preload("Replies.User").
		Where("post_id = ? AND parent_id IS NULL AND is_approved = ?", postID, true).
		Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&comments).Error; err != nil {
		return nil, 0, err
	}

	// 转换为响应格式
	var commentResponses []models.CommentResponse
	for _, comment := range comments {
		commentResponses = append(commentResponses, comment.ToResponse())
	}

	return commentResponses, total, nil
}

// DeleteComment 删除评论
func (cs *CommentService) DeleteComment(commentID uint) error {
	return cs.db.Delete(&models.Comment{}, commentID).Error
}

// IsCommentOwner 检查用户是否是评论的作者
func (cs *CommentService) IsCommentOwner(commentID, userID uint) bool {
	var comment models.Comment
	if err := cs.db.Select("user_id").First(&comment, commentID).Error; err != nil {
		return false
	}
	return comment.UserID == userID
}

// GetCommentReplies 获取评论的回复
func (cs *CommentService) GetCommentReplies(commentID uint) ([]models.Comment, error) {
	var replies []models.Comment
	if err := cs.db.Preload("User").
		Where("parent_id = ? AND is_approved = ?", commentID, true).
		Order("created_at ASC").
		Find(&replies).Error; err != nil {
		return nil, err
	}
	return replies, nil
}

// ToggleCommentApproval 切换评论审核状态（管理员功能）
func (cs *CommentService) ToggleCommentApproval(commentID uint) (*models.Comment, error) {
	var comment models.Comment
	if err := cs.db.First(&comment, commentID).Error; err != nil {
		return nil, err
	}

	// 切换审核状态
	comment.IsApproved = !comment.IsApproved
	if err := cs.db.Save(&comment).Error; err != nil {
		return nil, err
	}

	return &comment, nil
}