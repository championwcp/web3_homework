package main

import (
	"fmt"
	"log"
	"blockchain-task/config"
	"blockchain-task/query"
	"blockchain-task/transaction"
)

func main() {
	// 加载配置
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("配置加载失败:", err)
	}

	fmt.Println("=== 区块链交互程序 ===")
	fmt.Println("1. 查询最新区块信息")
	fmt.Println("2. 发送交易")
	fmt.Println("请选择操作:")

	var choice int
	_, err = fmt.Scan(&choice)
	if err != nil {
		log.Fatal("输入错误:", err)
	}

	switch choice {
	case 1:
		query.QueryBlock(cfg)
	case 2:
		transaction.SendTransaction(cfg)
	default:
		fmt.Println("无效选择")
	}
}