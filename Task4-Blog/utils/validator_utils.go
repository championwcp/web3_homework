package utils

import (
	"regexp"
	"strings"
	"unicode"

	"github.com/go-playground/validator/v10"
)

// CustomValidator 自定义验证器
type CustomValidator struct {
	validator *validator.Validate
}

// NewCustomValidator 创建自定义验证器
func NewCustomValidator() *CustomValidator {
	v := validator.New()
	
	// 注册自定义验证规则
	registerCustomValidations(v)
	
	return &CustomValidator{validator: v}
}

// Validate 验证结构体
func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

// registerCustomValidations 注册自定义验证规则
func registerCustomValidations(v *validator.Validate) {
	// 注册用户名验证
	v.RegisterValidation("username", func(fl validator.FieldLevel) bool {
		username := fl.Field().String()
		
		// 用户名长度 3-50 个字符
		if len(username) < 3 || len(username) > 50 {
			return false
		}
		
		// 只能包含字母、数字、下划线
		matched, _ := regexp.MatchString("^[a-zA-Z0-9_]+$", username)
		return matched
	})

	// 注册密码强度验证
	v.RegisterValidation("password", func(fl validator.FieldLevel) bool {
		password := fl.Field().String()
		
		// 密码长度至少6位
		if len(password) < 6 {
			return false
		}
		
		// 检查是否包含数字、字母
		var hasLetter, hasNumber bool
		for _, char := range password {
			switch {
			case unicode.IsLetter(char):
				hasLetter = true
			case unicode.IsNumber(char):
				hasNumber = true
			}
		}
		
		return hasLetter && hasNumber
	})

	// 注册 slug 验证（URL友好标识）
	v.RegisterValidation("slug", func(fl validator.FieldLevel) bool {
		slug := fl.Field().String()
		
		// slug 格式：小写字母、数字、连字符
		matched, _ := regexp.MatchString("^[a-z0-9-]+$", slug)
		return matched
	})
}

// ValidationErrorsToMap 将验证错误转换为 map
func ValidationErrorsToMap(err error) map[string]string {
	errorMap := make(map[string]string)
	
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, fieldError := range validationErrors {
			fieldName := strings.ToLower(fieldError.Field())
			
			switch fieldError.Tag() {
			case "required":
				errorMap[fieldName] = "该字段为必填项"
			case "email":
				errorMap[fieldName] = "邮箱格式不正确"
			case "min":
				errorMap[fieldName] = "长度太短"
			case "max":
				errorMap[fieldName] = "长度太长"
			case "username":
				errorMap[fieldName] = "用户名只能包含字母、数字、下划线，长度3-50位"
			case "password":
				errorMap[fieldName] = "密码必须包含字母和数字，长度至少6位"
			case "slug":
				errorMap[fieldName] = "Slug只能包含小写字母、数字和连字符"
			default:
				errorMap[fieldName] = "字段验证失败"
			}
		}
	}
	
	return errorMap
}