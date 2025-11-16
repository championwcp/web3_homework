package main

import (
	"fmt"
	"log"
	"blockchain-task-2/config"
	"blockchain-task-2/deploy"
	"blockchain-task-2/interact"
)

func main() {
	// 加载配置
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("配置加载失败:", err)
	}

	fmt.Println("=== 智能合约交互程序 ===")
	fmt.Println("1. 部署 Counter 合约")
	fmt.Println("2. 与 Counter 合约交互")
	fmt.Print("请选择操作: ")

	var choice int
	_, err = fmt.Scan(&choice)
	if err != nil {
		log.Fatal("输入错误:", err)
	}

	switch choice {
	case 1:
		deploy.DeployCounter(cfg)
	case 2:
		interact.InteractWithCounter(cfg)
	default:
		fmt.Println("无效选择")
	}
}