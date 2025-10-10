package services

import (
	"blog-system/config"
	"blog-system/database"
	"blog-system/models"
	"blog-system/utils"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"gorm.io/gorm"
)

// AuthService 认证服务
type AuthService struct {
	db *gorm.DB
}

// NewAuthService 创建认证服务实例
func NewAuthService() *AuthService {
	return &AuthService{
		db: database.GetDB(),
	}
}

// Register 用户注册
func (as *AuthService) Register(username, email, password, bio string) (*models.User, error) {
	// 加密密码
	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		Username: username,
		Email:    email,
		Password: hashedPassword,
		Bio:      bio,
		Role:     "user",
		IsActive: true,
	}

	// 创建用户
	if err := as.db.Create(user).Error; err != nil {
		return nil, err
	}

	return user, nil
}

// Login 用户登录
func (as *AuthService) Login(username, password string) (*models.User, error) {
	var user models.User
	
	// 根据用户名或邮箱查找用户
	if err := as.db.Where("username = ? OR email = ?", username, username).First(&user).Error; err != nil {
		return nil, err
	}

	// 检查用户是否激活
	if !user.IsActive {
		return nil, gorm.ErrRecordNotFound
	}

	// 验证密码
	if !utils.CheckPasswordHash(password, user.Password) {
		return nil, gorm.ErrRecordNotFound
	}

	return &user, nil
}

// GenerateToken 生成 JWT token
func (as *AuthService) GenerateToken(user *models.User) (string, error) {
	cfg := config.GetConfig()
	
	// 创建 token 声明
	claims := jwt.MapClaims{
		"user_id":  user.ID,
		"username": user.Username,
		"email":    user.Email,
		"role":     user.Role,
		"exp":      time.Now().Add(time.Hour * time.Duration(cfg.JWT.Expire)).Unix(),
		"iat":      time.Now().Unix(),
		"iss":      cfg.JWT.Issuer,
	}

	// 创建 token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 签名 token
	return token.SignedString([]byte(cfg.JWT.Secret))
}

// ValidateToken 验证 JWT token
func (as *AuthService) ValidateToken(tokenString string) (*jwt.Token, error) {
	cfg := config.GetConfig()
	
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// 验证签名方法
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(cfg.JWT.Secret), nil
	})

	if err != nil {
		return nil, err
	}

	return token, nil
}

// GetUserFromToken 从 token 中获取用户信息
func (as *AuthService) GetUserFromToken(tokenString string) (*models.User, error) {
	token, err := as.ValidateToken(tokenString)
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userID := uint(claims["user_id"].(float64))
		
		var user models.User
		if err := as.db.First(&user, userID).Error; err != nil {
			return nil, err
		}

		return &user, nil
	}

	return nil, jwt.ErrInvalidKey
}