package interact

import (
	"blockchain-task-2/bindings"
	"blockchain-task-2/config"
	"context"
	"crypto/ecdsa"
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
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

func InteractWithCounter(cfg *config.Config) {
	fmt.Println("=== Counter 合约交互 ===")

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

	// 用户输入合约地址
	var contractAddressStr string
	fmt.Print("请输入 Counter 合约地址: ")
	_, err = fmt.Scan(&contractAddressStr)
	if err != nil {
		log.Fatal("输入错误:", err)
	}

	contractAddress := common.HexToAddress(contractAddressStr)

	// 创建合约实例
	counter, err := bindings.NewCounter(contractAddress, client)
	if err != nil {
		log.Fatal("创建合约实例失败:", err)
	}

	fmt.Printf("成功连接到合约: %s\n", contractAddress.Hex())
	fmt.Printf("当前账户: %s\n", fromAddress.Hex())

	// 显示交互菜单
	for {
		fmt.Println("\n请选择操作:")
		fmt.Println("1. 获取当前计数")
		fmt.Println("2. 增加计数")
		fmt.Println("3. 增加指定数值")
		fmt.Println("4. 重置计数 (仅合约所有者)")
		fmt.Println("5. 获取合约所有者")
		fmt.Println("0. 退出")

		var choice int
		fmt.Scan(&choice)

		switch choice {
		case 1:
			getCount(counter)
		case 2:
			increaseCount(counter, auth, client, fromAddress)
		case 3:
			increaseBy(counter, auth, client, fromAddress)
		case 4:
			resetCount(counter, auth, client, fromAddress)
		case 5:
			getOwner(counter)
		case 0:
			fmt.Println("退出交互")
			return
		default:
			fmt.Println("无效选择")
		}
	}
}

func getCount(counter *bindings.Counter) {
	count, err := counter.GetCount(nil)
	if err != nil {
		fmt.Printf("获取计数失败: %v\n", err)
		return
	}
	fmt.Printf("当前计数: %s\n", count.String())
}

func increaseCount(counter *bindings.Counter, auth *bind.TransactOpts, client *ethclient.Client, fromAddress common.Address) {
	// 更新 nonce
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		fmt.Printf("获取 nonce 失败: %v\n", err)
		return
	}
	auth.Nonce = big.NewInt(int64(nonce))

	// 获取 gas 价格
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		fmt.Printf("获取 gas 价格失败: %v\n", err)
		return
	}
	auth.GasPrice = gasPrice
	auth.GasLimit = uint64(100000)

	fmt.Printf("发送增加计数交易...\n")
	tx, err := counter.Increase(auth)
	if err != nil {
		fmt.Printf("增加计数失败: %v\n", err)
		return
	}

	fmt.Printf("交易已发送，哈希: %s\n", tx.Hash().Hex())
	fmt.Printf("可以在 https://sepolia.etherscan.io/tx/%s 查看状态\n", tx.Hash().Hex())
}

func increaseBy(counter *bindings.Counter, auth *bind.TransactOpts, client *ethclient.Client, fromAddress common.Address) {
	var number int64
	fmt.Print("请输入要增加的数值: ")
	_, err := fmt.Scan(&number)
	if err != nil {
		fmt.Printf("输入错误: %v\n", err)
		return
	}

	// 更新 nonce
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		fmt.Printf("获取 nonce 失败: %v\n", err)
		return
	}
	auth.Nonce = big.NewInt(int64(nonce))

	// 获取 gas 价格
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		fmt.Printf("获取 gas 价格失败: %v\n", err)
		return
	}
	auth.GasPrice = gasPrice
	auth.GasLimit = uint64(100000)

	fmt.Printf("发送增加 %d 计数交易...\n", number)
	tx, err := counter.IncreaseBy(auth, big.NewInt(number))
	if err != nil {
		fmt.Printf("增加计数失败: %v\n", err)
		return
	}

	fmt.Printf("交易已发送，哈希: %s\n", tx.Hash().Hex())
	fmt.Printf("可以在 https://sepolia.etherscan.io/tx/%s 查看状态\n", tx.Hash().Hex())
}

func resetCount(counter *bindings.Counter, auth *bind.TransactOpts, client *ethclient.Client, fromAddress common.Address) {
	// 更新 nonce
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		fmt.Printf("获取 nonce 失败: %v\n", err)
		return
	}
	auth.Nonce = big.NewInt(int64(nonce))

	// 获取 gas 价格
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		fmt.Printf("获取 gas 价格失败: %v\n", err)
		return
	}
	auth.GasPrice = gasPrice
	auth.GasLimit = uint64(100000)

	fmt.Printf("发送重置计数交易...\n")
	tx, err := counter.Reset(auth)
	if err != nil {
		fmt.Printf("重置计数失败: %v\n", err)
		return
	}

	fmt.Printf("重置交易已发送，哈希: %s\n", tx.Hash().Hex())
	fmt.Printf("可以在 https://sepolia.etherscan.io/tx/%s 查看状态\n", tx.Hash().Hex())
}

func getOwner(counter *bindings.Counter) {
	owner, err := counter.Owner(nil)
	if err != nil {
		fmt.Printf("获取所有者失败: %v\n", err)
		return
	}
	fmt.Printf("合约所有者: %s\n", owner.Hex())
}