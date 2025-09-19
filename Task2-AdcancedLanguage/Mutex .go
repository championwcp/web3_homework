package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

// ==================== 第一题：使用 sync.Mutex 保护共享计数器 ====================

func mutexCounterDemo() {
	fmt.Println("启动10个协程，每个协程对计数器进行1000次递增操作")
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
	fmt.Println("启动10个协程，每个协程对计数器进行1000次递增操作")
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


func main() {
	// 第一题：使用Mutex保护共享计数器
	mutexCounterDemo()

	// 第二题：使用原子操作实现无锁计数器
	atomicCounterDemo()

}
