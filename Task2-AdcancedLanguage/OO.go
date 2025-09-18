package main

import (
	"fmt"
	"math"
)

// ==================== 第一题：接口实现 ====================

// Shape 接口定义
type Shape interface {
	Area() float64
	Perimeter() float64
}

// Rectangle 结构体
type Rectangle struct {
	Width  float64
	Height float64
}

// Rectangle 实现 Shape 接口的 Area 方法
func (r Rectangle) Area() float64 {
	return r.Width * r.Height
}

// Rectangle 实现 Shape 接口的 Perimeter 方法
func (r Rectangle) Perimeter() float64 {
	return 2 * (r.Width + r.Height)
}

// Rectangle 的字符串表示
func (r Rectangle) String() string {
	return fmt.Sprintf("矩形(宽: %.2f, 高: %.2f)", r.Width, r.Height)
}

// Circle 结构体
type Circle struct {
	Radius float64
}

// Circle 实现 Shape 接口的 Area 方法
func (c Circle) Area() float64 {
	return math.Pi * c.Radius * c.Radius
}

// Circle 实现 Shape 接口的 Perimeter 方法
func (c Circle) Perimeter() float64 {
	return 2 * math.Pi * c.Radius
}

// Circle 的字符串表示
func (c Circle) String() string {
	return fmt.Sprintf("圆形(半径: %.2f)", c.Radius)
}

// 打印形状信息的通用函数
func printShapeInfo(s Shape) {
	fmt.Printf("%s - 面积: %.2f, 周长: %.2f\n", s, s.Area(), s.Perimeter())
}

// ==================== 第二题：组合使用 ====================

// Person 结构体
type Person struct {
	Name string
	Age  int
}

// Person 的方法
func (p Person) Introduce() string {
	return fmt.Sprintf("我叫%s，今年%d岁", p.Name, p.Age)
}

// Employee 结构体，组合了 Person
type Employee struct {
	Person      // 匿名嵌入，实现组合
	EmployeeID  string
	Department  string
	Salary      float64
}

// Employee 的 PrintInfo 方法
func (e Employee) PrintInfo() {
	fmt.Println("=== 员工信息 ===")
	fmt.Printf("员工ID: %s\n", e.EmployeeID)
	fmt.Printf("姓名: %s\n", e.Name)    // 直接访问嵌入结构的字段
	fmt.Printf("年龄: %d\n", e.Age)      // 直接访问嵌入结构的字段
	fmt.Printf("部门: %s\n", e.Department)
	fmt.Printf("薪资: %.2f\n", e.Salary)
	fmt.Printf("介绍: %s\n", e.Introduce()) // 直接调用嵌入结构的方法
	fmt.Println()
}

// Employee 的详细信息方法
func (e Employee) GetDetailedInfo() string {
	return fmt.Sprintf("员工%s(%s)在%s部门工作，月薪%.2f元", 
		e.Name, e.EmployeeID, e.Department, e.Salary)
}

// 创建员工的工厂函数
func NewEmployee(name string, age int, id, department string, salary float64) Employee {
	return Employee{
		Person:     Person{Name: name, Age: age},
		EmployeeID: id,
		Department: department,
		Salary:     salary,
	}
}

func main() {
	fmt.Println("=============== 第一题：接口实现 ===============")
	
	// 创建 Rectangle 实例
	rect := Rectangle{Width: 5.0, Height: 3.0}
	printShapeInfo(rect)
	
	// 创建 Circle 实例
	circle := Circle{Radius: 2.5}
	printShapeInfo(circle)
	
	// 使用接口类型的切片
	fmt.Println("\n使用接口切片处理多种形状:")
	shapes := []Shape{
		Rectangle{Width: 4.0, Height: 6.0},
		Circle{Radius: 3.0},
		Rectangle{Width: 2.5, Height: 1.5},
	}
	
	for i, shape := range shapes {
		fmt.Printf("形状%d: ", i+1)
		printShapeInfo(shape)
	}
	
	// 类型断言示例
	fmt.Println("\n类型断言示例:")
	for i, shape := range shapes {
		switch s := shape.(type) {
		case Rectangle:
			fmt.Printf("形状%d是矩形，对角线长度: %.2f\n", i+1, math.Sqrt(s.Width*s.Width+s.Height*s.Height))
		case Circle:
			fmt.Printf("形状%d是圆形，直径: %.2f\n", i+1, 2*s.Radius)
		}
	}

	fmt.Println("\n=============== 第二题：组合使用 ===============")
	
	// 创建 Employee 实例
	emp1 := Employee{
		Person: Person{
			Name: "张存才",
			Age:  28,
		},
		EmployeeID: "E1001",
		Department: "技术部",
		Salary:     15000.50,
	}
	
	emp1.PrintInfo()
	
	// 使用工厂函数创建员工
	emp2 := NewEmployee("李仙海", 32, "E1002", "市场部", 12000.75)
	emp2.PrintInfo()
	
	// 演示组合的特性
	fmt.Println("演示组合的特性:")
	fmt.Println("可以直接访问嵌入结构的字段:")
	fmt.Printf("员工姓名: %s, 年龄: %d\n", emp2.Name, emp2.Age)
	
	fmt.Println("可以直接调用嵌入结构的方法:")
	fmt.Println(emp2.Introduce())
	
	fmt.Println("调用Employee特有的方法:")
	fmt.Println(emp2.GetDetailedInfo())
	
	// 创建多个员工
	fmt.Println("\n员工列表:")
	employees := []Employee{
		NewEmployee("王体", 25, "E1003", "财务部", 8000.00),
		NewEmployee("赵楼", 35, "E1004", "人事部", 9500.00),
		NewEmployee("钱藏", 29, "E1005", "研发部", 18000.00),
	}
	
	for _, emp := range employees {
		fmt.Printf("- %s (%s部门)\n", emp.Name, emp.Department)
	}
}