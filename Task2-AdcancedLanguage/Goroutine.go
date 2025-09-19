package main

import (
	"fmt"
	"sync"
	"time"
)

// 第一题：奇数偶数打印协程
func printOddEven() {
	fmt.Println("=== 第一题：奇数偶数打印 ===")
	var wg sync.WaitGroup
	wg.Add(2)
	// 打印奇数的协程
	go func() {
		defer wg.Done()
		for i := 1; i <= 10; i += 2 {
			fmt.Printf("奇数: %d\n", i)
			time.Sleep(100 * time.Millisecond) // 稍微延迟
		}
	}()

	// 打印偶数的协程
	go func() {
		defer wg.Done()
		for i := 2; i <= 10; i += 2 {
			fmt.Printf("偶数: %d\n", i)
			time.Sleep(100 * time.Millisecond) // 延迟
		}
	}()

	wg.Wait()
	fmt.Println()
}

// 第二题：任务调度器
type Task func()
// 任务执行结果
type TaskResult struct {
	TaskName    string
	Duration    time.Duration
	StartTime   time.Time
	EndTime     time.Time
	Success     bool
	Error       error
}
// 任务调度器
type TaskScheduler struct {
	tasks     []func()
	taskNames []string
	results   []TaskResult
	wg        sync.WaitGroup
	mu        sync.Mutex
}
// 创建新的任务调度器
func NewTaskScheduler() *TaskScheduler {
	return &TaskScheduler{
		tasks:     make([]func(), 0),
		taskNames: make([]string, 0),
		results:   make([]TaskResult, 0),
	}
}
// 添加任务
func (ts *TaskScheduler) AddTask(name string, task Task) {
	ts.tasks = append(ts.tasks, task)
	ts.taskNames = append(ts.taskNames, name)
}
// 执行单个任务
func (ts *TaskScheduler) executeTask(index int) {
	defer ts.wg.Done()
	taskName := ts.taskNames[index]
	task := ts.tasks[index]
	startTime := time.Now()
	result := TaskResult{
		TaskName:  taskName,
		StartTime: startTime,
		Success:   true,
	}

	defer func() {
		// 捕获panic，确保一个任务的崩溃不会影响其他任务
		if r := recover(); r != nil {
			result.Success = false
			result.Error = fmt.Errorf("task panicked: %v", r)
		}

		result.EndTime = time.Now()
		result.Duration = result.EndTime.Sub(startTime)

		ts.mu.Lock()
		ts.results = append(ts.results, result)
		ts.mu.Unlock()
	}()
	// 执行任务
	task()
}
// 并发执行所有任务
func (ts *TaskScheduler) Run() []TaskResult {
	ts.results = make([]TaskResult, 0)
	ts.wg.Add(len(ts.tasks))
	fmt.Println("=== 第二题：任务调度器开始执行 ===")
	startTime := time.Now()
	// 启动协程并发执行所有任务
	for i := 0; i < len(ts.tasks); i++ {
		go ts.executeTask(i)
	}
	// 等待所有任务完成
	ts.wg.Wait()
	totalDuration := time.Since(startTime)
	fmt.Printf("所有任务执行完成，总耗时: %v\n", totalDuration)
	fmt.Println()
	return ts.results
}

// 打印任务执行结果
func (ts *TaskScheduler) PrintResults() {
	fmt.Println("=== 任务执行结果统计 ===")
	for _, result := range ts.results {
		status := "成功"
		if !result.Success {
			status = "失败"
		}
		fmt.Printf("任务: %-15s | 状态: %-4s | 耗时: %-10v | 开始时间: %v\n",
			result.TaskName, status, result.Duration, result.StartTime.Format("15:04:05.000"))
	}
}

// 模拟任务函数
func quickTask() {
	time.Sleep(200 * time.Millisecond)
	fmt.Println("快速任务执行完成")
}

func mediumTask() {
	time.Sleep(500 * time.Millisecond)
	fmt.Println("中等任务执行完成")
}

func longTask() {
	time.Sleep(1000 * time.Millisecond)
	fmt.Println("长时间任务执行完成")
}

func errorTask() {
	time.Sleep(300 * time.Millisecond)
	panic("模拟任务执行出错")
}

func calculationTask() {
	sum := 0
	for i := 0; i < 1000000; i++ {
		sum += i
	}
	fmt.Printf("计算任务完成，结果: %d\n", sum)
}

func main() {
	// 执行第一题：奇数偶数打印
	printOddEven()
	// 执行第二题：任务调度器
	scheduler := NewTaskScheduler()
	// 添加各种任务
	scheduler.AddTask("快速任务", quickTask)
	scheduler.AddTask("中等任务", mediumTask)
	scheduler.AddTask("长时间任务", longTask)
	scheduler.AddTask("计算任务", calculationTask)
	scheduler.AddTask("错误任务", errorTask) // 这个任务会panic
	// 并发执行所有任务并获取结果
	results := scheduler.Run()
	// 打印执行结果
	scheduler.PrintResults()
	// 统计
	successCount := 0
	for _, result := range results {
		if result.Success {
			successCount++
		}
	}
	fmt.Printf("\n任务执行统计: 成功 %d/%d\n", successCount, len(results))

}
