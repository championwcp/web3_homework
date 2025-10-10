package main

import (
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

// Employee 员工结构体
type Employee struct {
	ID         int    `db:"id" json:"id"`
	Name       string `db:"name" json:"name"`
	Department string `db:"department" json:"department"`
	Salary     int    `db:"salary" json:"salary"`
}

// EmployeeManager 员工管理器
type EmployeeManager struct {
	db *sqlx.DB
}

// NewEmployeeManager 创建员工管理器实例
func NewEmployeeManager() *EmployeeManager {
	dsn := "root:st123456@tcp(127.0.0.1:3306)/gorm?charset=utf8mb4&parseTime=True&loc=Local"

	db, err := sqlx.Connect("mysql", dsn)
	if err != nil {
		log.Fatalf("无法连接数据库: %v", err)
	}

	fmt.Println("✅ 数据库连接成功!")

	// 创建 employees 表（如果不存在）
	err = createEmployeesTable(db)
	if err != nil {
		log.Fatalf("创建表失败: %v", err)
	}

	fmt.Println("✅ 表结构准备完成!")

	return &EmployeeManager{db: db}
}

// createEmployeesTable 创建员工表
func createEmployeesTable(db *sqlx.DB) error {
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS employees (
		id INT AUTO_INCREMENT PRIMARY KEY,
		name VARCHAR(100) NOT NULL,
		department VARCHAR(50) NOT NULL,
		salary INT NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	)`

	_, err := db.Exec(createTableSQL)
	return err
}

// InsertSampleData 插入示例数据
func (em *EmployeeManager) InsertSampleData() error {
	// 先清空现有数据
	_, err := em.db.Exec("DELETE FROM employees")
	if err != nil {
		return fmt.Errorf("清空数据失败: %v", err)
	}

	// 插入示例数据
	employees := []Employee{
		{Name: "张三", Department: "技术部", Salary: 15000},
		{Name: "李四", Department: "技术部", Salary: 18000},
		{Name: "王五", Department: "销售部", Salary: 12000},
		{Name: "赵六", Department: "技术部", Salary: 22000},
		{Name: "孙七", Department: "人事部", Salary: 10000},
		{Name: "周八", Department: "技术部", Salary: 16000},
		{Name: "吴九", Department: "财务部", Salary: 13000},
		{Name: "郑十", Department: "技术部", Salary: 19000},
	}

	for _, emp := range employees {
		query := `INSERT INTO employees (name, department, salary) VALUES (?, ?, ?)`
		_, err := em.db.Exec(query, emp.Name, emp.Department, emp.Salary)
		if err != nil {
			return fmt.Errorf("插入数据失败: %v", err)
		}
	}

	fmt.Println("✅ 示例数据插入完成!")
	return nil
}

// GetEmployeesByDepartment 查询指定部门的所有员工
func (em *EmployeeManager) GetEmployeesByDepartment(department string) ([]Employee, error) {
	var employees []Employee

	// 使用 sqlx.Select 直接映射到结构体切片
	query := `SELECT id, name, department, salary FROM employees WHERE department = ?`
	err := em.db.Select(&employees, query, department)
	if err != nil {
		return nil, fmt.Errorf("查询失败: %v", err)
	}

	return employees, nil
}

// GetHighestPaidEmployee 查询工资最高的员工
func (em *EmployeeManager) GetHighestPaidEmployee() (*Employee, error) {
	var employee Employee

	// 使用 sqlx.Get 查询单条记录
	query := `SELECT id, name, department, salary FROM employees ORDER BY salary DESC LIMIT 1`
	err := em.db.Get(&employee, query)
	if err != nil {
		return nil, fmt.Errorf("查询失败: %v", err)
	}

	return &employee, nil
}

// GetAllEmployees 获取所有员工
func (em *EmployeeManager) GetAllEmployees() ([]Employee, error) {
	var employees []Employee

	query := `SELECT id, name, department, salary FROM employees ORDER BY id`
	err := em.db.Select(&employees, query)
	if err != nil {
		return nil, fmt.Errorf("查询失败: %v", err)
	}

	return employees, nil
}

// 演示 SQLx 查询操作
func demonstrateSQLxQueries(em *EmployeeManager) {
	log.Println("开始演示 SQLx 查询操作...")

	// 1. 插入示例数据
	log.Println("\n1. 插入示例数据...")
	err := em.InsertSampleData()
	if err != nil {
		log.Printf("插入数据失败: %v\n", err)
		return
	}

	// 显示所有员工
	log.Println("\n2. 所有员工信息:")
	allEmployees, err := em.GetAllEmployees()
	if err != nil {
		log.Printf("查询所有员工失败: %v\n", err)
	} else {
		fmt.Printf("\n所有员工数据 (%d 名):\n", len(allEmployees))
		for _, emp := range allEmployees {
			fmt.Printf("   ID: %d, 姓名: %s, 部门: %s, 工资: %d\n",
				emp.ID, emp.Name, emp.Department, emp.Salary)
		}
		fmt.Println()
	}

	// 3. 查询技术部所有员工（题目要求1）
	log.Println("\n3. 查询技术部所有员工...")
	techEmployees, err := em.GetEmployeesByDepartment("技术部")
	if err != nil {
		log.Printf("查询技术部员工失败: %v\n", err)
	} else {
		fmt.Printf("技术部员工 (%d 名):\n", len(techEmployees))
		for _, emp := range techEmployees {
			fmt.Printf("   ID: %d, 姓名: %s, 工资: %d\n",
				emp.ID, emp.Name, emp.Salary)
		}
		fmt.Println()
	}

	// 4. 查询工资最高的员工（题目要求2）
	log.Println("\n4. 查询工资最高的员工...")
	highestPaid, err := em.GetHighestPaidEmployee()
	if err != nil {
		log.Printf("查询最高工资员工失败: %v\n", err)
	} else {
		fmt.Printf("工资最高的员工:\n")
		fmt.Printf("   ID: %d, 姓名: %s, 部门: %s, 工资: %d\n",
			highestPaid.ID, highestPaid.Name, highestPaid.Department, highestPaid.Salary)
		fmt.Println()
	}

	// 5. 其他部门查询示例
	log.Println("\n5. 查询销售部员工...")
	salesEmployees, err := em.GetEmployeesByDepartment("销售部")
	if err != nil {
		log.Printf("查询销售部员工失败: %v\n", err)
	} else {
		fmt.Printf("销售部员工 (%d 名):\n", len(salesEmployees))
		for _, emp := range salesEmployees {
			fmt.Printf("   ID: %d, 姓名: %s, 工资: %d\n",
				emp.ID, emp.Name, emp.Salary)
		}
		fmt.Println()
	}

	log.Println("SQLx 查询操作演示完成!")
}

// 使用 Named Exec 和 Named Query 的额外示例
func demonstrateNamedQueries(em *EmployeeManager) {
	log.Println("\n命名查询示例...")

	// 使用 Named Exec 插入数据
	newEmployee := Employee{
		Name:       "钱十一",
		Department: "技术部",
		Salary:     21000,
	}

	query := `INSERT INTO employees (name, department, salary) VALUES (:name, :department, :salary)`
	_, err := em.db.NamedExec(query, newEmployee)
	if err != nil {
		log.Printf("命名插入失败: %v\n", err)
	} else {
		fmt.Println(" 使用 NamedExec 插入数据成功")
	}

	// 使用 Named Query 查询
	//var highSalaryEmployees []Employee
	query = `SELECT * FROM employees WHERE salary > :min_salary AND department = :dept`
	args := map[string]interface{}{
		"min_salary": 17000,
		"dept":       "技术部",
	}

	rows, err := em.db.NamedQuery(query, args)
	if err != nil {
		log.Printf("命名查询失败: %v\n", err)
		return
	}
	defer rows.Close()

	fmt.Println("高薪技术部员工:")
	for rows.Next() {
		var emp Employee
		err := rows.StructScan(&emp)
		if err != nil {
			log.Printf("扫描结构体失败: %v\n", err)
			continue
		}
		fmt.Printf("   ID: %d, 姓名: %s, 工资: %d\n", emp.ID, emp.Name, emp.Salary)
	}
}

func main() {
	// 创建员工管理器
	employeeManager := NewEmployeeManager()

	// 执行演示
	demonstrateSQLxQueries(employeeManager)

	// 额外的命名查询演示
	demonstrateNamedQueries(employeeManager)

	// 程序结束提示
	fmt.Println("\nSQLx 查询程序执行完毕！")
}
