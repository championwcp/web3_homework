package utils

import (
	"blog-system/config"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// CreateToken 创建 JWT token
func CreateToken(userID uint, username, email, role string) (string, error) {
	cfg := config.GetConfig()
	
	// 创建 token 声明
	claims := jwt.MapClaims{
		"user_id":  userID,
		"username": username,
		"email":    email,
		"role":     role,
		"exp":      time.Now().Add(time.Hour * time.Duration(cfg.JWT.Expire)).Unix(),
		"iat":      time.Now().Unix(),
		"iss":      cfg.JWT.Issuer,
	}

	// 创建 token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 签名 token
	return token.SignedString([]byte(cfg.JWT.Secret))
}

// ParseToken 解析 JWT token
func ParseToken(tokenString string) (*jwt.Token, error) {
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

// ValidateToken 验证 token 是否有效
func ValidateToken(tokenString string) bool {
	token, err := ParseToken(tokenString)
	if err != nil {
		return false
	}
	return token.Valid
}

// ExtractUserIDFromToken 从 token 中提取用户ID
func ExtractUserIDFromToken(tokenString string) (uint, error) {
	token, err := ParseToken(tokenString)
	if err != nil {
		return 0, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userID := uint(claims["user_id"].(float64))
		return userID, nil
	}

	return 0, jwt.ErrInvalidKey
}

// ExtractUserInfoFromToken 从 token 中提取用户信息
func ExtractUserInfoFromToken(tokenString string) (uint, string, string, string, error) {
	token, err := ParseToken(tokenString)
	if err != nil {
		return 0, "", "", "", err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userID := uint(claims["user_id"].(float64))
		username := claims["username"].(string)
		email := claims["email"].(string)
		role := claims["role"].(string)
		
		return userID, username, email, role, nil
	}

	return 0, "", "", "", jwt.ErrInvalidKey
}