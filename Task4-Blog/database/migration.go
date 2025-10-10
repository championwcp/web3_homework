package database

import (
	"fmt"
	"log"

	"blog-system/models"
)

// Migrate 执行数据库迁移
func Migrate() error {
	if DB == nil {
		return fmt.Errorf("数据库连接未初始化")
	}

	log.Println("开始数据库迁移...")

	// 自动迁移所有模型
	err := DB.AutoMigrate(
		&models.User{},
		&models.Post{},
		&models.Comment{},
	)
	if err != nil {
		return fmt.Errorf("数据库迁移失败: %v", err)
	}

	log.Println("数据库迁移完成!")
	
	// 显示创建的表信息
	showTableInfo()

	return nil
}

// showTableInfo 显示表信息
func showTableInfo() {
	var tableInfo []struct {
		TableName string `gorm:"column:TABLE_NAME"`
	}
	
	DB.Raw("SELECT TABLE_NAME FROM INFORMATION_SCHEMA.TABLES WHERE TABLE_SCHEMA = DATABASE()").Scan(&tableInfo)
	
	log.Println("数据库表列表:")
	for _, table := range tableInfo {
		log.Printf("   - %s", table.TableName)
	}
}

// CreateTestData 创建测试数据
func CreateTestData() error {
	log.Println("创建测试数据...")
	
	// 这里可以添加一些初始测试数据
	
	log.Println("测试数据创建完成")
	return nil
}

// ResetDatabase 重置数据库（开发环境使用）
func ResetDatabase() error {
	if DB == nil {
		return fmt.Errorf("数据库连接未初始化")
	}

	log.Println("⚠️  重置数据库...")

	// 删除所有表（按依赖顺序）
	tables := []string{"comments", "posts", "users"}
	for _, table := range tables {
		if err := DB.Exec(fmt.Sprintf("DROP TABLE IF EXISTS %s", table)).Error; err != nil {
			return fmt.Errorf("删除表 %s 失败: %v", table, err)
		}
	}

	log.Println("数据库重置完成")
	return nil
}