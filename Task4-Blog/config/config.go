package config

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/spf13/viper"
)

// Config 全局配置结构体
type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	JWT      JWTConfig      `mapstructure:"jwt"`
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Port         int    `mapstructure:"port"`
	Mode         string `mapstructure:"mode"`
	ReadTimeout  int    `mapstructure:"read_timeout"`
	WriteTimeout int    `mapstructure:"write_timeout"`
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	DBName   string `mapstructure:"dbname"`
	Charset  string `mapstructure:"charset"`
	// 连接池配置
	MaxIdleConns int `mapstructure:"max_idle_conns"`
	MaxOpenConns int `mapstructure:"max_open_conns"`
}

// JWTConfig JWT配置
type JWTConfig struct {
	Secret string `mapstructure:"secret"`
	Expire int    `mapstructure:"expire"` // 过期时间（小时）
	Issuer string `mapstructure:"issuer"`
}

// GlobalConfig 全局配置实例
var GlobalConfig *Config

// Init 初始化配置
func Init() error {
	// 设置配置文件路径
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")
	viper.AddConfigPath(".")

	// 设置环境变量前缀
	viper.SetEnvPrefix("BLOG")
	viper.AutomaticEnv()

	// 设置默认值
	setDefaults()

	// 读取配置文件
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Println("配置文件未找到，使用默认配置和环境变量")
		} else {
			return fmt.Errorf("读取配置文件失败: %v", err)
		}
	} else {
		log.Printf("加载配置文件: %s", viper.ConfigFileUsed())
	}

	// 解析配置到结构体
	if err := viper.Unmarshal(&GlobalConfig); err != nil {
		return fmt.Errorf("解析配置失败: %v", err)
	}

	// 从环境变量覆盖配置（如果存在）
	overrideFromEnv()

	log.Println("配置初始化完成")
	return nil
}

// setDefaults 设置默认配置
func setDefaults() {
	// 服务器配置默认值
	viper.SetDefault("server.port", 8080)
	viper.SetDefault("server.mode", "debug")
	viper.SetDefault("server.read_timeout", 30)
	viper.SetDefault("server.write_timeout", 30)

	// 数据库配置默认值
	viper.SetDefault("database.host", "localhost")
	viper.SetDefault("database.port", 3306)
	viper.SetDefault("database.user", "root")
	viper.SetDefault("database.password", "password")
	viper.SetDefault("database.dbname", "blog_system")
	viper.SetDefault("database.charset", "utf8mb4")
	viper.SetDefault("database.max_idle_conns", 10)
	viper.SetDefault("database.max_open_conns", 100)

	// JWT配置默认值
	viper.SetDefault("jwt.secret", "your-secret-key-change-in-production")
	viper.SetDefault("jwt.expire", 24) // 24小时
	viper.SetDefault("jwt.issuer", "blog-system")
}

// overrideFromEnv 从环境变量覆盖配置
func overrideFromEnv() {
	// 服务器配置环境变量
	if port := os.Getenv("BLOG_SERVER_PORT"); port != "" {
		if p, err := strconv.Atoi(port); err == nil {
			GlobalConfig.Server.Port = p
		}
	}
	if mode := os.Getenv("BLOG_SERVER_MODE"); mode != "" {
		GlobalConfig.Server.Mode = mode
	}

	// 数据库配置环境变量
	if host := os.Getenv("BLOG_DB_HOST"); host != "" {
		GlobalConfig.Database.Host = host
	}
	if port := os.Getenv("BLOG_DB_PORT"); port != "" {
		if p, err := strconv.Atoi(port); err == nil {
			GlobalConfig.Database.Port = p
		}
	}
	if user := os.Getenv("BLOG_DB_USER"); user != "" {
		GlobalConfig.Database.User = user
	}
	if password := os.Getenv("BLOG_DB_PASSWORD"); password != "" {
		GlobalConfig.Database.Password = password
	}
	if dbname := os.Getenv("BLOG_DB_NAME"); dbname != "" {
		GlobalConfig.Database.DBName = dbname
	}

	// JWT配置环境变量
	if secret := os.Getenv("BLOG_JWT_SECRET"); secret != "" {
		GlobalConfig.JWT.Secret = secret
	}
	if expire := os.Getenv("BLOG_JWT_EXPIRE"); expire != "" {
		if e, err := strconv.Atoi(expire); err == nil {
			GlobalConfig.JWT.Expire = e
		}
	}
}

// GetDSN 获取数据库连接字符串
func (d *DatabaseConfig) GetDSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local",
		d.User, d.Password, d.Host, d.Port, d.DBName, d.Charset)
}

// GetConfig 获取全局配置实例
func GetConfig() *Config {
	return GlobalConfig
}