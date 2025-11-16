package transaction

import (
	"blockchain-task/config"
	"context"
	"crypto/ecdsa"
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
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

func SendTransaction(cfg *config.Config) {
	fmt.Println("准备发送交易...")

	// 连接到 Sepolia 网络
	client, err := ethclient.Dial(cfg.AlchemyURL)
	if err != nil {
		log.Fatal("连接失败:", err)
	}
	defer client.Close()

	// 1. 准备私钥（处理 0x 前缀）
	privateKey, err := parsePrivateKey(cfg.PrivateKey)
	if err != nil {
		log.Fatal("私钥格式错误:", err)
	}

	// 2. 从私钥获取公钥和地址
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("无法生成公钥")
	}
	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

	// 3. 获取账户余额（检查是否有足够 ETH）
	balance, err := client.BalanceAt(context.Background(), fromAddress, nil)
	if err != nil {
		log.Fatal("获取余额失败:", err)
	}
	fmt.Printf("账户余额: %s wei\n", balance.String())

	// 4. 获取账户当前 nonce
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Fatal("获取 nonce 失败:", err)
	}

	// 5. 设置转账金额（0.001 ETH）
	value := big.NewInt(1000000000000000) // 0.001 ETH in wei

	// 6. 设置 gas 限制
	gasLimit := uint64(21000) // 标准转账的 gas 限制

	// 7. 获取当前 gas 价格
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatal("获取 gas 价格失败:", err)
	}

	// 8. 计算所需总费用
	totalCost := new(big.Int).Mul(gasPrice, big.NewInt(int64(gasLimit)))
	totalCost.Add(totalCost, value)

	// 检查余额是否足够
	if balance.Cmp(totalCost) < 0 {
		log.Fatalf("余额不足！需要: %s wei, 当前余额: %s wei", totalCost.String(), balance.String())
	}

	// 9. 设置接收方地址
	var toAddress common.Address
	if cfg.RecipientAddress != "" {
		// 使用配置的接收方地址
		toAddress = common.HexToAddress(cfg.RecipientAddress)
	} else {
		// 如果没有设置接收方，默认发给自己
		toAddress = fromAddress
	}

	fmt.Printf("\n=== 交易详情 ===\n")
	fmt.Printf("发送方: %s\n", fromAddress.Hex())
	fmt.Printf("接收方: %s\n", toAddress.Hex())
	fmt.Printf("金额: 0.001 ETH (%s wei)\n", value.String())
	fmt.Printf("Nonce: %d\n", nonce)
	fmt.Printf("Gas 价格: %s wei\n", gasPrice.String())
	fmt.Printf("Gas 限制: %d\n", gasLimit)
	fmt.Printf("预计手续费: %s wei\n", new(big.Int).Mul(gasPrice, big.NewInt(int64(gasLimit))).String())

	// 10. 创建交易对象
	tx := types.NewTransaction(nonce, toAddress, value, gasLimit, gasPrice, nil)

	// 11. 获取链ID
	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		log.Fatal("获取网络ID失败:", err)
	}

	// 12. 签名交易
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		log.Fatal("交易签名失败:", err)
	}

	// 13. 发送交易
	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		log.Fatal("发送交易失败:", err)
	}

	// 14. 输出交易哈希
	fmt.Printf("\n=== 交易发送成功 ===\n")
	fmt.Printf("交易哈希: %s\n", signedTx.Hash().Hex())
	fmt.Printf("可以在 https://sepolia.etherscan.io/tx/%s 查看交易状态\n", signedTx.Hash().Hex())
}