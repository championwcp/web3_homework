package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	AlchemyURL       string
	PrivateKey       string
	RecipientAddress string
}

func LoadConfig() (*Config, error) {
	// 加载 .env 文件
	err := godotenv.Load()
	if err != nil {
		log.Println("警告: 无法加载 .env 文件，使用系统环境变量")
	}

	url := os.Getenv("ALCHEMY_SEPOLIA_URL")
	if url == "" {
		return nil, fmt.Errorf("请设置 ALCHEMY_SEPOLIA_URL 环境变量")
	}

	privateKey := os.Getenv("PRIVATE_KEY")
	if privateKey == "" {
		return nil, fmt.Errorf("请设置 PRIVATE_KEY 环境变量")
	}

	return &Config{
		AlchemyURL:       url,
		PrivateKey:       privateKey,
		RecipientAddress: os.Getenv("RECIPIENT_ADDRESS"),
	}, nil
}