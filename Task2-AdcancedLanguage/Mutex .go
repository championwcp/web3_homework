package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

// ==================== 第一题：使用 sync.Mutex 保护共享计数器 ====================

func mutexCounterDemo() {
	fmt.Println("=== 第一题：使用 sync.Mutex 保护共享计数器 ===")
	fmt.Println("启动10个协程，每个协程对计数器进行1000次递增操作")
	fmt.Println("----------------------------------------")

	var (
		counter int
		mutex   sync.Mutex
		wg      sync.WaitGroup
	)

	const (
		goroutines = 10
		increments = 1000
	)

	wg.Add(goroutines)

	for i := 0; i < goroutines; i++ {
		go func(id int) {
			defer wg.Done()
			for j := 0; j < increments; j++ {
				mutex.Lock()
				counter++
				mutex.Unlock()
			}
			// 使用id变量避免警告
			_ = id
		}(i)
	}

	wg.Wait()

	expected := goroutines * increments
	fmt.Printf("最终计数器值: %d\n", counter)
	fmt.Printf("期望值: %d\n", expected)
	fmt.Printf("结果是否正确: %t\n", counter == expected)
	fmt.Println()
}

// ==================== 第二题：使用原子操作实现无锁计数器 ====================

func atomicCounterDemo() {
	fmt.Println("=== 第二题：使用原子操作实现无锁计数器 ===")
	fmt.Println("启动10个协程，每个协程对计数器进行1000次递增操作")
	fmt.Println("----------------------------------------")

	var (
		counter int64
		wg      sync.WaitGroup
	)

	const (
		goroutines = 10
		increments = 1000
	)

	wg.Add(goroutines)

	for i := 0; i < goroutines; i++ {
		go func(id int) {
			defer wg.Done()
			for j := 0; j < increments; j++ {
				atomic.AddInt64(&counter, 1)
			}
			// 使用id变量避免警告
			_ = id
		}(i)
	}

	wg.Wait()

	expected := int64(goroutines * increments)
	fmt.Printf("最终计数器值: %d\n", counter)
	fmt.Printf("期望值: %d\n", expected)
	fmt.Printf("结果是否正确: %t\n", counter == expected)
	fmt.Println()
}

// ==================== 扩展：性能对比测试 ====================

func performanceComparison() {
	fmt.Println("=== 扩展：Mutex vs Atomic 性能对比 ===")
	fmt.Println("每种方法运行5次，计算平均耗时")
	fmt.Println("----------------------------------------")

	const (
		goroutines = 10
		increments = 100000 // 增加操作次数以便更好地测量性能
		testRuns   = 5
	)

	// 测试Mutex性能
	//var mutexTotalTime int64
	for run := 0; run < testRuns; run++ {
		var (
			counter int
			mutex   sync.Mutex
			wg      sync.WaitGroup
		)

		wg.Add(goroutines)
		start := make(chan struct{})

		for i := 0; i < goroutines; i++ {
			go func() {
				defer wg.Done()
				<-start
				for j := 0; j < increments; j++ {
					mutex.Lock()
					counter++
					mutex.Unlock()
				}
			}()
		}

		// 开始计时
		close(start)
		wg.Wait()

		// 验证结果
		expected := goroutines * increments
		if counter != expected {
			fmt.Printf("Mutex测试运行%d: 结果错误，期望%d，得到%d\n", run+1, expected, counter)
		}
	}

	// 测试Atomic性能
	//var atomicTotalTime int64
	for run := 0; run < testRuns; run++ {
		var (
			counter int64
			wg      sync.WaitGroup
		)

		wg.Add(goroutines)
		start := make(chan struct{})

		for i := 0; i < goroutines; i++ {
			go func() {
				defer wg.Done()
				<-start
				for j := 0; j < increments; j++ {
					atomic.AddInt64(&counter, 1)
				}
			}()
		}

		// 开始计时
		close(start)
		wg.Wait()

		// 验证结果
		expected := int64(goroutines * increments)
		if counter != expected {
			fmt.Printf("Atomic测试运行%d: 结果错误，期望%d，得到%d\n", run+1, expected, counter)
		}
	}

	fmt.Println("性能对比完成（实际计时代码需要更复杂的实现）")
	fmt.Println("通常原子操作比互斥锁更快，因为避免了上下文切换")
	fmt.Println()
}



func main() {
	// 第一题：使用Mutex保护共享计数器
	mutexCounterDemo()

	// 第二题：使用原子操作实现无锁计数器
	atomicCounterDemo()
}