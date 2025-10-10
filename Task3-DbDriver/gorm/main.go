package main

import (
	"fmt"
	"log"
	"strings"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func main() {
	// 初始化数据库连接
	db := initDatabase()

	// 创建博客管理器
	blogManager := NewBlogManager(db)

	// 执行演示
	demonstrateGORMFeatures(blogManager)
}

// initDatabase 初始化数据库连接
func initDatabase() *gorm.DB {
	dsn := "root:st123456@tcp(127.0.0.1:3306)/gorm?charset=utf8mb4&parseTime=True&loc=Local"

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatalf("无法连接数据库: %v", err)
	}

	fmt.Println("数据库连接成功!")

	// 自动迁移表结构（题目1要求）
	err = db.AutoMigrate(&User{}, &Post{}, &Comment{})
	if err != nil {
		log.Fatalf("自动迁移失败: %v", err)
	}

	fmt.Println("表结构迁移完成!")
	return db
}

// demonstrateGORMFeatures 演示GORM功能
func demonstrateGORMFeatures(bm *BlogManager) {
	log.Println("开始演示GORM进阶功能...")

	// 1. 创建示例数据
	log.Println("\n" + strings.Repeat("=", 50))
	log.Println("1. 创建示例数据")
	log.Println(strings.Repeat("=", 50))

	err := bm.CreateSampleData()
	if err != nil {
		log.Printf("创建示例数据失败: %v\n", err)
		return
	}

	// 显示初始数据
	log.Println("\n" + strings.Repeat("=", 50))
	log.Println("2. 初始数据展示")
	log.Println(strings.Repeat("=", 50))

	displayInitialData(bm)

	// 3. 题目2要求1：查询用户的所有文章及其评论
	log.Println("\n" + strings.Repeat("=", 50))
	log.Println("3. 查询用户的所有文章及其评论（题目2要求1）")
	log.Println(strings.Repeat("=", 50))

	demonstrateUserPostsQuery(bm)

	// 4. 题目2要求2：查询评论数量最多的文章
	log.Println("\n" + strings.Repeat("=", 50))
	log.Println("4. 查询评论数量最多的文章（题目2要求2）")
	log.Println(strings.Repeat("=", 50))

	demonstrateMostCommentedPost(bm)

	// 5. 题目3：钩子函数演示
	log.Println("\n" + strings.Repeat("=", 50))
	log.Println("5. 钩子函数演示（题目3）")
	log.Println(strings.Repeat("=", 50))

	demonstrateHooks(bm)

	// 6. 最终数据状态
	log.Println("\n" + strings.Repeat("=", 50))
	log.Println("6. 最终数据状态")
	log.Println(strings.Repeat("=", 50))

	displayFinalData(bm)

	log.Println("GORM进阶功能演示完成!")
}

// displayInitialData 显示初始数据
func displayInitialData(bm *BlogManager) {
	users, err := bm.GetAllUsers()
	if err != nil {
		log.Printf("获取用户数据失败: %v\n", err)
		return
	}

	fmt.Println("\n👥 用户信息:")
	for _, user := range users {
		fmt.Printf("   用户: %s (ID: %d), 文章数量: %d\n", user.Name, user.ID, user.PostCount)
	}

	posts, err := bm.GetAllPosts()
	if err != nil {
		log.Printf("获取文章数据失败: %v\n", err)
		return
	}

	fmt.Println("\n文章信息:")
	for _, post := range posts {
		fmt.Printf("   文章: 《%s》- %s, 评论数量: %d, 状态: %s\n",
			post.Title, post.User.Name, post.CommentCount, post.CommentStatus)
	}
}

// demonstrateUserPostsQuery 演示用户文章查询
func demonstrateUserPostsQuery(bm *BlogManager) {
	// 查询用户1的所有文章及其评论
	userID := uint(1)
	posts, err := bm.GetUserPostsWithComments(userID)
	if err != nil {
		log.Printf("查询用户文章失败: %v\n", err)
		return
	}

	fmt.Printf("\n用户 %d 的所有文章及其评论:\n", userID)
	for _, post := range posts {
		fmt.Printf("\n   文章: 《%s》\n", post.Title)
		fmt.Printf("   内容: %.50s...\n", post.Content)
		fmt.Printf("   评论数量: %d\n", len(post.Comments))

		for i, comment := range post.Comments {
			fmt.Printf("     评论%d: %s - %.30s\n", i+1, comment.Author, comment.Content)
		}
	}
}

// demonstrateMostCommentedPost 演示最多评论文章查询
func demonstrateMostCommentedPost(bm *BlogManager) {
	post, err := bm.GetMostCommentedPost()
	if err != nil {
		log.Printf("查询最多评论文章失败: %v\n", err)
		return
	}

	fmt.Printf("\n评论数量最多的文章:\n")
	fmt.Printf("   标题: 《%s》\n", post.Title)
	fmt.Printf("   作者: %s\n", post.User.Name)
	fmt.Printf("   评论数量: %d\n", post.CommentCount)
	fmt.Printf("   评论状态: %s\n", post.CommentStatus)

	fmt.Println("   所有评论:")
	for i, comment := range post.Comments {
		fmt.Printf("     %d. %s: %s\n", i+1, comment.Author, comment.Content)
	}
}

// demonstrateHooks 演示钩子函数
func demonstrateHooks(bm *BlogManager) {
	// 创建新文章测试 AfterCreate 钩子
	fmt.Println("\n测试 Post AfterCreate 钩子:")
	newPost := Post{
		Title:   "测试钩子函数的文章",
		Content: "这篇文章用于测试AfterCreate钩子函数...",
		UserID:  2, // 李四
	}

	err := bm.db.Create(&newPost).Error
	if err != nil {
		log.Printf("创建测试文章失败: %v\n", err)
	} else {
		fmt.Printf("✅ 新文章创建成功，应该看到用户文章数量更新的消息\n")
	}

	// 删除评论测试 AfterDelete 钩子
	fmt.Println("\n🔔 测试 Comment AfterDelete 钩子:")

	// 先获取一个有评论的文章
	post, err := bm.GetPostWithComments(1)
	if err != nil {
		log.Printf("获取文章失败: %v\n", err)
		return
	}

	if len(post.Comments) > 0 {
		commentID := post.Comments[0].ID
		fmt.Printf("   删除评论 ID: %d\n", commentID)
		err = bm.DeleteComment(commentID)
		if err != nil {
			log.Printf("删除评论失败: %v\n", err)
		} else {
			fmt.Printf("评论删除成功，应该看到评论数量更新的消息\n")
		}
	}
}

// displayFinalData 显示最终数据
func displayFinalData(bm *BlogManager) {
	users, err := bm.GetAllUsers()
	if err != nil {
		log.Printf("获取用户数据失败: %v\n", err)
		return
	}

	fmt.Println("\n最终用户信息:")
	for _, user := range users {
		fmt.Printf("   用户: %s, 文章数量: %d\n", user.Name, user.PostCount)
	}

	posts, err := bm.GetAllPosts()
	if err != nil {
		log.Printf("获取文章数据失败: %v\n", err)
		return
	}

	fmt.Println("\n最终文章信息:")
	for _, post := range posts {
		fmt.Printf("   文章: 《%s》, 评论数量: %d, 状态: %s\n",
			post.Title, post.CommentCount, post.CommentStatus)
	}
}
