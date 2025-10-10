package database

import (
	"fmt"
	"log"
	"time"

	"blog-system/config"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DB 全局数据库实例
var DB *gorm.DB

// InitDB 初始化数据库连接
func InitDB() (*gorm.DB, error) {
	cfg := config.GetConfig()
	
	// 构建 DSN (Data Source Name)
	dsn := cfg.Database.GetDSN()
	
	// GORM 配置
	gormConfig := &gorm.Config{
		// 在 debug 模式下显示详细的 SQL 日志
		Logger: logger.Default.LogMode(getLogLevel(cfg.Server.Mode)),
		// 禁用默认的事务
		SkipDefaultTransaction: false,
	}

	// 连接数据库
	var err error
	DB, err = gorm.Open(mysql.Open(dsn), gormConfig)
	if err != nil {
		return nil, fmt.Errorf("数据库连接失败: %v", err)
	}

	// 配置连接池
	sqlDB, err := DB.DB()
	if err != nil {
		return nil, fmt.Errorf("获取数据库实例失败: %v", err)
	}

	// 设置连接池参数
	dbConfig := cfg.Database
	sqlDB.SetMaxIdleConns(dbConfig.MaxIdleConns)    // 最大空闲连接数
	sqlDB.SetMaxOpenConns(dbConfig.MaxOpenConns)    // 最大打开连接数
	sqlDB.SetConnMaxLifetime(time.Hour)             // 连接最大存活时间

	log.Println("数据库连接成功!")
	log.Printf("连接池配置: 最大空闲连接=%d, 最大打开连接=%d", 
		dbConfig.MaxIdleConns, dbConfig.MaxOpenConns)

	return DB, nil
}

// getLogLevel 根据运行模式获取 GORM 日志级别
func getLogLevel(mode string) logger.LogLevel {
	switch mode {
	case "release":
		return logger.Silent
	case "test":
		return logger.Warn
	default: // debug
		return logger.Info
	}
}

// GetDB 获取数据库实例
func GetDB() *gorm.DB {
	return DB
}

// CloseDB 关闭数据库连接
func CloseDB() error {
	if DB != nil {
		sqlDB, err := DB.DB()
		if err != nil {
			return err
		}
		return sqlDB.Close()
	}
	return nil
}

// HealthCheck 数据库健康检查
func HealthCheck() error {
	if DB == nil {
		return fmt.Errorf("数据库连接未初始化")
	}
	
	sqlDB, err := DB.DB()
	if err != nil {
		return err
	}
	
	return sqlDB.Ping()
}