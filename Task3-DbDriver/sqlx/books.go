package main

import (
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

// Book 书籍结构体 - 确保类型安全映射
type Book struct {
	ID          int       `db:"id" json:"id"`
	Title       string    `db:"title" json:"title"`
	Author      string    `db:"author" json:"author"`
	Price       float64   `db:"price" json:"price"`
	PublishDate time.Time `db:"publish_date" json:"publish_date"`
	ISBN        string    `db:"isbn" json:"isbn"`
	InStock     bool      `db:"in_stock" json:"in_stock"`
}

// BookManager 书籍管理器
type BookManager struct {
	db *sqlx.DB
}

// NewBookManager 创建书籍管理器实例
func NewBookManager() *BookManager {
	dsn := "root:st123456@tcp(127.0.0.1:3306)/gorm?charset=utf8mb4&parseTime=True&loc=Local"

	db, err := sqlx.Connect("mysql", dsn)
	if err != nil {
		log.Fatalf("无法连接数据库: %v", err)
	}

	fmt.Println("数据库连接成功!")

	// 创建 books 表（如果不存在）
	err = createBooksTable(db)
	if err != nil {
		log.Fatalf("创建表失败: %v", err)
	}

	fmt.Println("表结构准备完成!")

	return &BookManager{db: db}
}

// createBooksTable 创建书籍表
func createBooksTable(db *sqlx.DB) error {
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS books (
		id INT AUTO_INCREMENT PRIMARY KEY,
		title VARCHAR(200) NOT NULL,
		author VARCHAR(100) NOT NULL,
		price DECIMAL(10, 2) NOT NULL,
		publish_date DATE,
		isbn VARCHAR(20),
		in_stock BOOLEAN DEFAULT true,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		INDEX idx_price (price),
		INDEX idx_author (author)
	)`

	_, err := db.Exec(createTableSQL)
	return err
}

// InsertSampleData 插入示例数据
func (bm *BookManager) InsertSampleData() error {
	// 先清空现有数据
	_, err := bm.db.Exec("DELETE FROM books")
	if err != nil {
		return fmt.Errorf("清空数据失败: %v", err)
	}

	// 插入示例数据
	books := []Book{
		{
			Title:       "Go语言编程",
			Author:      "张三",
			Price:       65.50,
			PublishDate: time.Date(2022, 1, 15, 0, 0, 0, 0, time.UTC),
			ISBN:        "978-7-121-12345-6",
			InStock:     true,
		},
		{
			Title:       "数据库系统概念",
			Author:      "李四",
			Price:       89.00,
			PublishDate: time.Date(2021, 8, 20, 0, 0, 0, 0, time.UTC),
			ISBN:        "978-7-115-23456-7",
			InStock:     true,
		},
		{
			Title:       "算法导论",
			Author:      "王五",
			Price:       128.00,
			PublishDate: time.Date(2020, 5, 10, 0, 0, 0, 0, time.UTC),
			ISBN:        "978-7-111-34567-8",
			InStock:     false,
		},
		{
			Title:       "深入理解计算机系统",
			Author:      "赵六",
			Price:       99.80,
			PublishDate: time.Date(2019, 12, 5, 0, 0, 0, 0, time.UTC),
			ISBN:        "978-7-121-45678-9",
			InStock:     true,
		},
		{
			Title:       "设计模式",
			Author:      "陈七",
			Price:       45.00,
			PublishDate: time.Date(2023, 3, 25, 0, 0, 0, 0, time.UTC),
			ISBN:        "978-7-302-56789-0",
			InStock:     true,
		},
		{
			Title:       "计算机网络",
			Author:      "刘八",
			Price:       75.50,
			PublishDate: time.Date(2021, 11, 30, 0, 0, 0, 0, time.UTC),
			ISBN:        "978-7-04-67890-1",
			InStock:     true,
		},
		{
			Title:       "软件工程",
			Author:      "孙九",
			Price:       38.00,
			PublishDate: time.Date(2020, 7, 18, 0, 0, 0, 0, time.UTC),
			ISBN:        "978-7-113-78901-2",
			InStock:     false,
		},
		{
			Title:       "人工智能",
			Author:      "周十",
			Price:       156.00,
			PublishDate: time.Date(2022, 9, 12, 0, 0, 0, 0, time.UTC),
			ISBN:        "978-7-121-89012-3",
			InStock:     true,
		},
	}

	for _, book := range books {
		query := `INSERT INTO books (title, author, price, publish_date, isbn, in_stock) 
		          VALUES (:title, :author, :price, :publish_date, :isbn, :in_stock)`
		_, err := bm.db.NamedExec(query, book)
		if err != nil {
			return fmt.Errorf("插入数据失败: %v", err)
		}
	}

	fmt.Println("✅ 示例数据插入完成!")
	return nil
}

// GetExpensiveBooks 查询价格大于指定值的书籍（题目要求）
func (bm *BookManager) GetExpensiveBooks(minPrice float64) ([]Book, error) {
	var books []Book

	// 使用 sqlx.Select 进行类型安全映射
	query := `SELECT id, title, author, price, publish_date, isbn, in_stock 
	          FROM books 
	          WHERE price > ? 
	          ORDER BY price DESC`

	err := bm.db.Select(&books, query, minPrice)
	if err != nil {
		return nil, fmt.Errorf("查询高价书籍失败: %v", err)
	}

	return books, nil
}

// GetAllBooks 获取所有书籍
func (bm *BookManager) GetAllBooks() ([]Book, error) {
	var books []Book

	query := `SELECT id, title, author, price, publish_date, isbn, in_stock 
	          FROM books 
	          ORDER BY id`
	err := bm.db.Select(&books, query)
	if err != nil {
		return nil, fmt.Errorf("查询所有书籍失败: %v", err)
	}

	return books, nil
}

// GetBooksByAuthor 按作者查询书籍
func (bm *BookManager) GetBooksByAuthor(author string) ([]Book, error) {
	var books []Book

	query := `SELECT id, title, author, price, publish_date, isbn, in_stock 
	          FROM books 
	          WHERE author = ? 
	          ORDER BY price DESC`
	err := bm.db.Select(&books, query, author)
	if err != nil {
		return nil, fmt.Errorf("按作者查询失败: %v", err)
	}

	return books, nil
}

// GetAvailableExpensiveBooks 查询有库存的高价书籍（复杂查询示例）
func (bm *BookManager) GetAvailableExpensiveBooks(minPrice float64) ([]Book, error) {
	var books []Book

	// 复杂的多条件查询
	query := `SELECT id, title, author, price, publish_date, isbn, in_stock 
	          FROM books 
	          WHERE price > ? AND in_stock = true
	          ORDER BY price DESC, publish_date DESC`

	err := bm.db.Select(&books, query, minPrice)
	if err != nil {
		return nil, fmt.Errorf("查询有库存高价书籍失败: %v", err)
	}

	return books, nil
}

// GetBooksPriceRange 查询价格区间的书籍
func (bm *BookManager) GetBooksPriceRange(minPrice, maxPrice float64) ([]Book, error) {
	var books []Book

	query := `SELECT id, title, author, price, publish_date, isbn, in_stock 
	          FROM books 
	          WHERE price BETWEEN ? AND ?
	          ORDER BY price ASC`

	err := bm.db.Select(&books, query, minPrice, maxPrice)
	if err != nil {
		return nil, fmt.Errorf("查询价格区间书籍失败: %v", err)
	}

	return books, nil
}

// GetBooksStatistics 获取书籍统计信息
func (bm *BookManager) GetBooksStatistics() (map[string]interface{}, error) {
	var stats struct {
		TotalBooks      int     `db:"total_books"`
		AveragePrice    float64 `db:"avg_price"`
		MaxPrice        float64 `db:"max_price"`
		MinPrice        float64 `db:"min_price"`
		BooksInStock    int     `db:"books_in_stock"`
		BooksOutOfStock int     `db:"books_out_of_stock"`
	}

	query := `SELECT 
		COUNT(*) as total_books,
		AVG(price) as avg_price,
		MAX(price) as max_price,
		MIN(price) as min_price,
		SUM(CASE WHEN in_stock = true THEN 1 ELSE 0 END) as books_in_stock,
		SUM(CASE WHEN in_stock = false THEN 1 ELSE 0 END) as books_out_of_stock
	FROM books`

	err := bm.db.Get(&stats, query)
	if err != nil {
		return nil, fmt.Errorf("获取统计信息失败: %v", err)
	}

	result := map[string]interface{}{
		"total_books":        stats.TotalBooks,
		"average_price":      fmt.Sprintf("%.2f", stats.AveragePrice),
		"max_price":          stats.MaxPrice,
		"min_price":          stats.MinPrice,
		"books_in_stock":     stats.BooksInStock,
		"books_out_of_stock": stats.BooksOutOfStock,
	}

	return result, nil
}

// 演示类型安全映射和复杂查询
func demonstrateTypeSafeMapping(bm *BookManager) {
	log.Println("开始演示类型安全映射和复杂查询...")

	// 1. 插入示例数据
	log.Println("\n1. 插入示例数据...")
	err := bm.InsertSampleData()
	if err != nil {
		log.Printf("插入数据失败: %v\n", err)
		return
	}

	// 2. 显示所有书籍
	log.Println("\n2. 所有书籍信息:")
	allBooks, err := bm.GetAllBooks()
	if err != nil {
		log.Printf("查询所有书籍失败: %v\n", err)
	} else {
		fmt.Printf("\n📚 所有书籍数据 (%d 本):\n", len(allBooks))
		for _, book := range allBooks {
			stockStatus := "有货"
			if !book.InStock {
				stockStatus = "缺货"
			}
			fmt.Printf("   ID: %d, 书名: 《%s》, 作者: %s, 价格: ¥%.2f, 出版日期: %s, 库存: %s\n",
				book.ID, book.Title, book.Author, book.Price,
				book.PublishDate.Format("2006-01-02"), stockStatus)
		}
		fmt.Println()
	}

	// 3. 题目要求：查询价格大于50元的书籍
	log.Println("\n3. 查询价格大于50元的书籍（题目要求）...")
	expensiveBooks, err := bm.GetExpensiveBooks(50.0)
	if err != nil {
		log.Printf("查询高价书籍失败: %v\n", err)
	} else {
		fmt.Printf("价格大于50元的书籍 (%d 本):\n", len(expensiveBooks))
		for _, book := range expensiveBooks {
			fmt.Printf("   《%s》- %s | 价格: ¥%.2f | ISBN: %s\n",
				book.Title, book.Author, book.Price, book.ISBN)
		}
		fmt.Println()
	}

	// 4. 复杂查询：有库存的高价书籍
	log.Println("\n4. 查询有库存且价格大于80元的书籍...")
	availableExpensiveBooks, err := bm.GetAvailableExpensiveBooks(80.0)
	if err != nil {
		log.Printf("查询有库存高价书籍失败: %v\n", err)
	} else {
		fmt.Printf("有库存且价格大于80元的书籍 (%d 本):\n", len(availableExpensiveBooks))
		for _, book := range availableExpensiveBooks {
			fmt.Printf("   《%s》- %s | 价格: ¥%.2f | 出版日期: %s\n",
				book.Title, book.Author, book.Price, book.PublishDate.Format("2006-01-02"))
		}
		fmt.Println()
	}

	// 5. 价格区间查询
	log.Println("\n5. 查询价格在40-100元之间的书籍...")
	midRangeBooks, err := bm.GetBooksPriceRange(40.0, 100.0)
	if err != nil {
		log.Printf("查询价格区间书籍失败: %v\n", err)
	} else {
		fmt.Printf("价格在40-100元之间的书籍 (%d 本):\n", len(midRangeBooks))
		for _, book := range midRangeBooks {
			fmt.Printf("   《%s》- %s | 价格: ¥%.2f\n",
				book.Title, book.Author, book.Price)
		}
		fmt.Println()
	}

	// 6. 统计信息
	log.Println("\n6. 书籍统计信息...")
	stats, err := bm.GetBooksStatistics()
	if err != nil {
		log.Printf("获取统计信息失败: %v\n", err)
	} else {
		fmt.Println("书籍统计信息:")
		for key, value := range stats {
			fmt.Printf("   %s: %v\n", key, value)
		}
		fmt.Println()
	}

	log.Println("类型安全映射演示完成!")
}

func main() {
	// 创建书籍管理器
	bookManager := NewBookManager()

	// 执行演示
	demonstrateTypeSafeMapping(bookManager)

	// 程序结束提示
	fmt.Println("\n✨ 类型安全映射程序执行完毕！")
}
