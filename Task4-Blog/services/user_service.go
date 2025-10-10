package services

import (
	"blog-system/database"
	"blog-system/models"
	"blog-system/utils"
	"time"

	"gorm.io/gorm"
)

// UserService 用户服务
type UserService struct {
	db *gorm.DB
}

// NewUserService 创建用户服务实例
func NewUserService() *UserService {
	return &UserService{
		db: database.GetDB(),
	}
}

// GetUserByID 根据ID获取用户
func (us *UserService) GetUserByID(userID uint) (*models.User, error) {
	var user models.User
	if err := us.db.First(&user, userID).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// GetUserByUsername 根据用户名获取用户
func (us *UserService) GetUserByUsername(username string) (*models.User, error) {
	var user models.User
	if err := us.db.Where("username = ?", username).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// GetUserByEmail 根据邮箱获取用户
func (us *UserService) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	if err := us.db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// UsernameExists 检查用户名是否存在
func (us *UserService) UsernameExists(username string) bool {
	var count int64
	us.db.Model(&models.User{}).Where("username = ?", username).Count(&count)
	return count > 0
}

// EmailExists 检查邮箱是否存在
func (us *UserService) EmailExists(email string) bool {
	var count int64
	us.db.Model(&models.User{}).Where("email = ?", email).Count(&count)
	return count > 0
}

// UpdateUser 更新用户信息
func (us *UserService) UpdateUser(userID uint, bio, avatar string) (*models.User, error) {
	var user models.User
	if err := us.db.First(&user, userID).Error; err != nil {
		return nil, err
	}

	updates := make(map[string]interface{})
	if bio != "" {
		updates["bio"] = bio
	}
	if avatar != "" {
		updates["avatar"] = avatar
	}

	if len(updates) > 0 {
		if err := us.db.Model(&user).Updates(updates).Error; err != nil {
			return nil, err
		}
	}

	return &user, nil
}

// UpdateLastLogin 更新最后登录时间
func (us *UserService) UpdateLastLogin(userID uint, loginTime *time.Time) error {
	return us.db.Model(&models.User{}).Where("id = ?", userID).Update("last_login", loginTime).Error
}

// GetUserPosts 获取用户的文章列表
func (us *UserService) GetUserPosts(userID uint, page, pageSize int) ([]models.PostResponse, int64, error) {
	var posts []models.Post
	var total int64

	// 计算偏移量
	offset := (page - 1) * pageSize

	// 获取总数
	if err := us.db.Model(&models.Post{}).Where("user_id = ?", userID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 获取文章列表
	if err := us.db.Preload("User").
		Where("user_id = ?", userID).
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

// GetUsers 获取用户列表
func (us *UserService) GetUsers(page, pageSize int) ([]models.User, int64, error) {
	var users []models.User
	var total int64

	// 计算偏移量
	offset := (page - 1) * pageSize

	// 获取总数
	if err := us.db.Model(&models.User{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 获取用户列表
	if err := us.db.
		Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

// ChangePassword 修改密码
func (us *UserService) ChangePassword(userID uint, oldPassword, newPassword string) error {
	var user models.User
	if err := us.db.First(&user, userID).Error; err != nil {
		return err
	}

	// 验证旧密码
	if !utils.CheckPasswordHash(oldPassword, user.Password) {
		return gorm.ErrInvalidData
	}

	// 加密新密码
	hashedPassword, err := utils.HashPassword(newPassword)
	if err != nil {
		return err
	}

	// 更新密码
	return us.db.Model(&user).Update("password", hashedPassword).Error
}