package deploy

import (
	"blockchain-task-2/bindings"
	"blockchain-task-2/config"
	"context"
	"crypto/ecdsa"
	"fmt"
	"log"
	"math/big"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	//"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

// 私钥解析辅助函数
func parsePrivateKey(privateKeyStr string) (*ecdsa.PrivateKey, error) {
	// 去掉 0x 前缀（如果存在）
	if len(privateKeyStr) > 2 && privateKeyStr[:2] == "0x" {
		privateKeyStr = privateKeyStr[2:]
	}
	
	// 验证私钥长度
	if len(privateKeyStr) != 64 {
		return nil, fmt.Errorf("私钥长度不正确，应该是64个字符，当前长度: %d", len(privateKeyStr))
	}
	
	return crypto.HexToECDSA(privateKeyStr)
}

func DeployCounter(cfg *config.Config) {
	fmt.Println("开始部署 Counter 合约...")

	// 连接到 Sepolia 网络
	client, err := ethclient.Dial(cfg.AlchemyURL)
	if err != nil {
		log.Fatal("连接失败:", err)
	}
	defer client.Close()

	// 准备私钥（处理 0x 前缀）
	privateKey, err := parsePrivateKey(cfg.PrivateKey)
	if err != nil {
		log.Fatal("私钥格式错误:", err)
	}

	// 直接从私钥获取地址
	fromAddress := crypto.PubkeyToAddress(privateKey.PublicKey)

	// 获取账户余额
	balance, err := client.BalanceAt(context.Background(), fromAddress, nil)
	if err != nil {
		log.Fatal("获取余额失败:", err)
	}
	fmt.Printf("账户余额: %s wei\n", balance.String())

	// 检查余额是否足够
	minBalance := new(big.Int).Mul(big.NewInt(1000000), big.NewInt(2000000000)) // 1M gas * 2 gwei
	if balance.Cmp(minBalance) < 0 {
		log.Fatalf("余额可能不足！当前: %s wei, 建议至少: %s wei", balance.String(), minBalance.String())
	}

	// 获取 nonce
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Fatal("获取 nonce 失败:", err)
	}

	// 获取 gas 价格
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatal("获取 gas 价格失败:", err)
	}

	// 增加 gas 价格以确保快速确认
	gasPrice = gasPrice.Mul(gasPrice, big.NewInt(12))
	gasPrice = gasPrice.Div(gasPrice, big.NewInt(10)) // 增加 20%

	// 获取链 ID
	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		log.Fatal("获取网络ID失败:", err)
	}

	// 创建认证对象
	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	if err != nil {
		log.Fatal("创建交易认证失败:", err)
	}
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)      // 部署合约不需要发送 ETH
	auth.GasLimit = uint64(1000000) // 大幅增加 Gas Limit
	auth.GasPrice = gasPrice

	fmt.Printf("\n=== 部署参数 ===\n")
	fmt.Printf("部署账户: %s\n", fromAddress.Hex())
	fmt.Printf("Nonce: %d\n", nonce)
	fmt.Printf("Gas 价格: %s wei\n", gasPrice.String())
	fmt.Printf("Gas 限制: %d\n", auth.GasLimit)
	
	// 计算预估费用
	estimatedCost := new(big.Int).Mul(gasPrice, big.NewInt(int64(auth.GasLimit)))
	fmt.Printf("预估部署费用: %s wei\n", estimatedCost.String())

	// 部署合约
	fmt.Println("\n正在部署合约，请等待...")
	address, tx, _, err := bindings.DeployCounter(auth, client)
	if err != nil {
		log.Fatal("合约部署失败:", err)
	}

	fmt.Printf("\n=== 合约部署成功 ===\n")
	fmt.Printf("合约地址: %s\n", address.Hex())
	fmt.Printf("交易哈希: %s\n", tx.Hash().Hex())
	fmt.Printf("Gas 使用量: %d\n", tx.Gas())
	fmt.Printf("可以在 https://sepolia.etherscan.io/tx/%s 查看交易状态\n", tx.Hash().Hex())
	fmt.Printf("可以在 https://sepolia.etherscan.io/address/%s 查看合约\n", address.Hex())

	// 保存合约地址
	saveContractAddress(address.Hex())
}

func saveContractAddress(address string) {
	fmt.Printf("请记下合约地址用于后续交互: %s\n", address)
}