package main

import (
	"fmt"
	"log"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// Student 学生模型
type Student struct {
	ID    int    `gorm:"primaryKey;autoIncrement" json:"id"`
	Name  string `gorm:"type:varchar(50);not null" json:"name"`
	Age   int    `gorm:"not null" json:"age"`
	Grade string `gorm:"type:varchar(20);not null" json:"grade"`
}

// TableName 表名
func (Student) TableName() string {
	return "students"
}

type StudentManager struct {
	db *gorm.DB
}

// NewStudentManager 创建学生管理器实例
func NewStudentManager() *StudentManager {
	//
	dsn := "root:pasword@tcp(127.0.0.1:3306)/gorm?charset=utf8mb4&parseTime=True&loc=Local"

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("无法连接数据库: %v", err)
	}

	fmt.Println("数据库连接成功!")

	// 自动迁移表结构
	err = db.AutoMigrate(&Student{})
	if err != nil {
		log.Fatalf("自动迁移失败: %v", err)
	}

	fmt.Println("表结构迁移完成!")

	return &StudentManager{db: db}
}

// CreateStudent 插入新学生
func (sm *StudentManager) CreateStudent(name string, age int, grade string) error {
	student := Student{
		Name:  name,
		Age:   age,
		Grade: grade,
	}

	result := sm.db.Create(&student)
	if result.Error != nil {
		return result.Error
	}

	fmt.Printf("插入成功: ID=%d, 姓名=%s, 年龄=%d, 年级=%s\n",
		student.ID, student.Name, student.Age, student.Grade)
	return nil
}

// GetStudentsByAge 查询年龄大于指定值的学生
func (sm *StudentManager) GetStudentsByAge(minAge int) ([]Student, error) {
	var students []Student
	result := sm.db.Where("age > ?", minAge).Find(&students)

	if result.Error != nil {
		return nil, result.Error
	}

	fmt.Printf("查询到 %d 名年龄大于 %d 的学生:\n", len(students), minAge)
	for _, student := range students {
		fmt.Printf("   ID: %d, 姓名: %s, 年龄: %d, 年级: %s\n",
			student.ID, student.Name, student.Age, student.Grade)
	}

	return students, nil
}

// UpdateStudentGrade 更新学生年级
func (sm *StudentManager) UpdateStudentGrade(name string, newGrade string) error {
	result := sm.db.Model(&Student{}).
		Where("name = ?", name).
		Update("grade", newGrade)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("未找到姓名为 %s 的学生", name)
	}

	fmt.Printf("更新成功: %s 的年级已更新为 %s\n", name, newGrade)
	return nil
}

// DeleteStudentsByAge 删除年龄小于指定值的学生
func (sm *StudentManager) DeleteStudentsByAge(maxAge int) error {
	result := sm.db.Where("age < ?", maxAge).Delete(&Student{})

	if result.Error != nil {
		return result.Error
	}

	fmt.Printf("删除成功: 删除了 %d 名年龄小于 %d 的学生\n", result.RowsAffected, maxAge)
	return nil
}

// GetAllStudents 获取所有学生
func (sm *StudentManager) GetAllStudents() ([]Student, error) {
	var students []Student
	result := sm.db.Find(&students)

	if result.Error != nil {
		return nil, result.Error
	}

	fmt.Printf("\n当前所有学生数据 (%d 名):\n", len(students))
	for _, student := range students {
		fmt.Printf("   ID: %d, 姓名: %s, 年龄: %d, 年级: %s\n",
			student.ID, student.Name, student.Age, student.Grade)
	}
	fmt.Println()

	return students, nil
}

// ClearAllStudents 清空所有学生数据（测试）
func (sm *StudentManager) ClearAllStudents() error {
	result := sm.db.Where("1 = 1").Delete(&Student{})
	if result.Error != nil {
		return result.Error
	}
	fmt.Println("已清空所有学生数据")
	return nil
}

// 演示所有 CRUD 操作
func demonstrateCRUDOperations(sm *StudentManager) {
	log.Println("开始执行学生管理系统 CRUD 操作...")

	// 清空现有数据
	sm.ClearAllStudents()
	time.Sleep(1 * time.Second)

	// 显示初始状态
	sm.GetAllStudents()
	log.Println("\n1. 插入新学生记录...")
	err := sm.CreateStudent("张三", 20, "三年级")
	if err != nil {
		log.Printf("插入失败: %v\n", err)
	}

	// 显示当前所有学生
	sm.GetAllStudents()

	log.Println("\n2. 查询年龄大于18岁的学生...")
	_, err = sm.GetStudentsByAge(18)
	if err != nil {
		log.Printf("查询失败: %v\n", err)
	}

	log.Println("\n3. 更新学生年级...")
	err = sm.UpdateStudentGrade("张三", "四年级")
	if err != nil {
		log.Printf("更新失败: %v\n", err)
	}

	// 显示更新后的状态
	sm.GetAllStudents()
	log.Println("\n4. 插入测试数据用于删除操作...")
	testStudents := []struct {
		name  string
		age   int
		grade string
	}{
		{"李四", 14, "二年级"},
		{"王五", 16, "三年级"},
		{"赵六", 13, "一年级"},
		{"孙七", 17, "三年级"},
	}

	for _, stu := range testStudents {
		err = sm.CreateStudent(stu.name, stu.age, stu.grade)
		if err != nil {
			log.Printf("插入测试数据失败: %v\n", err)
		}
		time.Sleep(100 * time.Millisecond)
	}

	log.Println("\n删除前的学生列表:")
	sm.GetAllStudents()
	log.Println("\n5. 删除年龄小于15岁的学生...")
	err = sm.DeleteStudentsByAge(15)
	if err != nil {
		log.Printf("删除失败: %v\n", err)
	}
	log.Println("\n🎊 最终学生列表:")
	sm.GetAllStudents()
}

func main() {
	// 创建学生管理器
	studentManager := NewStudentManager()

	// 执行演示
	demonstrateCRUDOperations(studentManager)

	// 程序结束提示
	fmt.Println("\n程序执行完毕！")

}
