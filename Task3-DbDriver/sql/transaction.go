package main

import (
	"fmt"
	"log"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// Account 账户模型
type Account struct {
	ID      int     `gorm:"primaryKey;autoIncrement" json:"id"`
	Balance float64 `gorm:"not null;default:0" json:"balance"`
}

// Transaction 交易记录模型
type Transaction struct {
	ID            int       `gorm:"primaryKey;autoIncrement" json:"id"`
	FromAccountID int       `gorm:"not null" json:"from_account_id"`
	ToAccountID   int       `gorm:"not null" json:"to_account_id"`
	Amount        float64   `gorm:"not null" json:"amount"`
	CreatedAt     time.Time `gorm:"autoCreateTime" json:"created_at"`
}

// TransferManager
type TransferManager struct {
	db *gorm.DB
}

// NewTransferManager 创建转账管理器实例
func NewTransferManager() *TransferManager {
	dsn := "root:st123456@tcp(127.0.0.1:3306)/gorm?charset=utf8mb4&parseTime=True&loc=Local"

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("无法连接数据库: %v", err)
	}

	fmt.Println("数据库连接成功!")

	// 自动迁移表结构
	err = db.AutoMigrate(&Account{}, &Transaction{})
	if err != nil {
		log.Fatalf("自动迁移失败: %v", err)
	}

	fmt.Println("✅ 表结构迁移完成!")

	return &TransferManager{db: db}
}

// Transfer 转账事务
func (tm *TransferManager) Transfer(fromAccountID, toAccountID int, amount float64) error {
	// 开始事务
	tx := tm.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	// 延迟函数处理事务提交或回滚
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			log.Printf("事务回滚（发生异常）: %v", r)
		}
	}()

	fmt.Printf("\n🔄 开始转账: 从账户 %d 向账户 %d 转账 %.2f 元\n", fromAccountID, toAccountID, amount)

	// 1. 检查转出账户是否存在并锁定账户记录
	var fromAccount Account
	result := tx.Set("gorm:query_option", "FOR UPDATE").First(&fromAccount, fromAccountID)
	if result.Error != nil {
		tx.Rollback()
		return fmt.Errorf("转出账户 %d 不存在", fromAccountID)
	}

	// 2. 检查余额是否足够
	if fromAccount.Balance < amount {
		tx.Rollback()
		return fmt.Errorf("转账失败: 账户 %d 余额不足 (当前余额: %.2f, 需要: %.2f)",
			fromAccountID, fromAccount.Balance, amount)
	}

	// 3. 检查转入账户是否存在
	var toAccount Account
	result = tx.First(&toAccount, toAccountID)
	if result.Error != nil {
		tx.Rollback()
		return fmt.Errorf("转入账户 %d 不存在", toAccountID)
	}

	// 4. 从转出账户扣除金额
	result = tx.Model(&fromAccount).Update("balance", gorm.Expr("balance - ?", amount))
	if result.Error != nil {
		tx.Rollback()
		return fmt.Errorf("扣除转出账户余额失败: %v", result.Error)
	}

	// 5. 向转入账户增加金额
	result = tx.Model(&toAccount).Update("balance", gorm.Expr("balance + ?", amount))
	if result.Error != nil {
		tx.Rollback()
		return fmt.Errorf("增加转入账户余额失败: %v", result.Error)
	}

	// 6. 记录交易信息
	transaction := Transaction{
		FromAccountID: fromAccountID,
		ToAccountID:   toAccountID,
		Amount:        amount,
		CreatedAt:     time.Now(),
	}
	result = tx.Create(&transaction)
	if result.Error != nil {
		tx.Rollback()
		return fmt.Errorf("记录交易失败: %v", result.Error)
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("提交事务失败: %v", err)
	}

	fmt.Printf("✅ 转账成功! 交易ID: %d\n", transaction.ID)
	fmt.Printf("   账户 %d 新余额: %.2f\n", fromAccountID, fromAccount.Balance-amount)
	fmt.Printf("   账户 %d 新余额: %.2f\n", toAccountID, toAccount.Balance+amount)

	return nil
}

// CreateAccount 创建账户
func (tm *TransferManager) CreateAccount(initialBalance float64) (*Account, error) {
	account := Account{
		Balance: initialBalance,
	}
	result := tm.db.Create(&account)
	if result.Error != nil {
		return nil, result.Error
	}
	fmt.Printf("✅ 创建账户成功: ID=%d, 初始余额=%.2f\n", account.ID, account.Balance)
	return &account, nil
}

// GetAccount 获取账户信息
func (tm *TransferManager) GetAccount(accountID int) (*Account, error) {
	var account Account
	result := tm.db.First(&account, accountID)
	if result.Error != nil {
		return nil, result.Error
	}
	return &account, nil
}

// GetAllAccounts 获取所有账户
func (tm *TransferManager) GetAllAccounts() ([]Account, error) {
	var accounts []Account
	result := tm.db.Find(&accounts)
	if result.Error != nil {
		return nil, result.Error
	}

	fmt.Println("\n所有账户信息:")
	for _, account := range accounts {
		fmt.Printf("   账户 ID: %d, 余额: %.2f\n", account.ID, account.Balance)
	}
	fmt.Println()

	return accounts, nil
}

// GetTransactionHistory 获取交易记录
func (tm *TransferManager) GetTransactionHistory() ([]Transaction, error) {
	var transactions []Transaction
	result := tm.db.Order("created_at DESC").Find(&transactions)
	if result.Error != nil {
		return nil, result.Error
	}

	fmt.Println("\n交易记录:")
	if len(transactions) == 0 {
		fmt.Println("   暂无交易记录")
	} else {
		for _, tx := range transactions {
			fmt.Printf("   交易ID: %d, 从账户 %d 到账户 %d, 金额: %.2f, 时间: %s\n",
				tx.ID, tx.FromAccountID, tx.ToAccountID, tx.Amount, tx.CreatedAt.Format("2006-01-02 15:04:05"))
		}
	}
	fmt.Println()

	return transactions, nil
}

// ClearAllTransferData 清空转账数据
func (tm *TransferManager) ClearAllTransferData() error {
	tx := tm.db.Begin()

	if err := tx.Exec("DELETE FROM transactions").Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Exec("DELETE FROM accounts").Error; err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	fmt.Println("🗑️  已清空所有账户和交易数据")
	return nil
}

// 演示转账事务操作
func demonstrateTransferOperations(tm *TransferManager) {
	log.Println("开始演示转账事务操作...")

	// 清空现有数据，确保演示环境干净
	tm.ClearAllTransferData()
	time.Sleep(1 * time.Second)

	// 1. 创建测试账户
	log.Println("\n1. 创建测试账户...")
	accountA, err := tm.CreateAccount(500.0) // 账户A初始余额500元
	if err != nil {
		log.Printf("创建账户A失败: %v\n", err)
		return
	}

	accountB, err := tm.CreateAccount(300.0) // 账户B初始余额300元
	if err != nil {
		log.Printf("创建账户B失败: %v\n", err)
		return
	}

	accountC, err := tm.CreateAccount(50.0) // 账户C初始余额50元
	if err != nil {
		log.Printf("创建账户C失败: %v\n", err)
		return
	}

	// 显示初始账户状态
	tm.GetAllAccounts()

	// 2. 正常转账：账户A向账户B转账100元
	log.Println("\n2. 正常转账测试...")
	err = tm.Transfer(accountA.ID, accountB.ID, 100.0)
	if err != nil {
		log.Printf("转账失败: %v\n", err)
	}

	// 显示转账后的账户状态
	tm.GetAllAccounts()
	tm.GetTransactionHistory()

	// 3. 余额不足测试：账户C向账户A转账100元（应该失败）
	log.Println("\n3. 余额不足测试...")
	err = tm.Transfer(accountC.ID, accountA.ID, 100.0)
	if err != nil {
		log.Printf("预期中的转账失败: %v\n", err)
	} else {
		log.Println("余额不足测试未按预期失败")
	}

	// 显示测试后的账户状态（余额应该不变）
	tm.GetAllAccounts()

	// 4. 成功的小额转账：账户C向账户B转账30元
	log.Println("\n4. 小额转账测试...")
	err = tm.Transfer(accountC.ID, accountB.ID, 30.0)
	if err != nil {
		log.Printf("小额转账失败: %v\n", err)
	}

	// 显示最终账户状态和交易记录
	log.Println("最终状态:")
	tm.GetAllAccounts()
	tm.GetTransactionHistory()

	log.Println("所有转账事务操作演示完成!")
}

func main() {
	// 创建转账管理器
	transferManager := NewTransferManager()

	// 执行演示
	demonstrateTransferOperations(transferManager)

	// 程序结束提示
	fmt.Println("转账事务程序执行完毕！")
}
