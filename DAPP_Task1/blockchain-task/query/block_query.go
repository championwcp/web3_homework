package query

import (
	"blockchain-task/config"
	"context"
	"fmt"
	"log"

	"github.com/ethereum/go-ethereum/ethclient"
)

func QueryBlock(cfg *config.Config) {
	fmt.Println("开始查询区块信息...")

	// 连接到 Sepolia 网络
	client, err := ethclient.Dial(cfg.AlchemyURL)
	if err != nil {
		log.Fatal("连接失败:", err)
	}
	defer client.Close()

	// 查询最新区块
	header, err := client.HeaderByNumber(context.Background(), nil)
	if err != nil {
		log.Fatal("查询区块头失败:", err)
	}

	// 获取完整区块信息
	block, err := client.BlockByNumber(context.Background(), header.Number)
	if err != nil {
		log.Fatal("查询区块失败:", err)
	}

	// 输出区块信息
	fmt.Printf("\n=== 区块信息 ===\n")
	fmt.Printf("区块号: %s\n", block.Number().String())
	fmt.Printf("区块哈希: %s\n", block.Hash().Hex())                                                                       
	fmt.Printf("时间戳: %d\n", block.Time())
	fmt.Printf("交易数量: %d\n", len(block.Transactions()))
	fmt.Printf("难度: %s\n", block.Difficulty().String())
fmt.Printf("Gas 限制: %d\n", block.GasLimit())
fmt.Printf("Gas 使用量: %d\n", block.GasUsed())
}